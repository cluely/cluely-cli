package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/browser"
	"github.com/zalando/go-keyring"
)

const (
	keyringService = "cluely-cli"
	keyToken       = "cli_token"
	loginTimeout   = 5 * time.Minute
)

// StoreToken saves the CLI token to the OS keyring.
func StoreToken(token string) error {
	if err := keyring.Set(keyringService, keyToken, token); err != nil {
		return fmt.Errorf("store token: %w", err)
	}
	return nil
}

// LoadToken retrieves the CLI token from the OS keyring.
// Returns empty string (no error) if no token is stored.
func LoadToken() (string, error) {
	token, err := keyring.Get(keyringService, keyToken)
	if err == keyring.ErrNotFound {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("load token: %w", err)
	}
	return token, nil
}

// ClearToken removes the stored token from the OS keyring.
func ClearToken() error {
	if err := keyring.Delete(keyringService, keyToken); err != nil && err != keyring.ErrNotFound {
		return fmt.Errorf("clear token: %w", err)
	}
	return nil
}

// HasToken returns true if a token is stored.
func HasToken() bool {
	t, _ := keyring.Get(keyringService, keyToken)
	return t != ""
}

// callbackResult holds data received from the web app callback.
type callbackResult struct {
	token string
	err   error
}

// Login performs the browser-based OAuth flow via the Cluely web app.
// 1. Starts a localhost server
// 2. Opens browser to the web app's /cli-auth page
// 3. User signs in with Clerk on the web
// 4. Web app creates a long-lived CLI token and redirects to localhost callback
// 5. CLI stores the token in the OS keyring
func Login(ctx context.Context, webBaseURL string) error {
	state, err := generateState()
	if err != nil {
		return fmt.Errorf("generate state: %w", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("start callback server: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	callbackURL := fmt.Sprintf("http://127.0.0.1:%d/callback", port)

	result := make(chan callbackResult, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if q.Get("state") != state {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			result <- callbackResult{err: fmt.Errorf("state mismatch: possible CSRF attack")}
			return
		}

		if errMsg := q.Get("error"); errMsg != "" {
			http.Error(w, "Authentication failed", http.StatusBadRequest)
			result <- callbackResult{err: fmt.Errorf("auth error: %s", errMsg)}
			return
		}

		token := q.Get("token")
		if token == "" {
			http.Error(w, "Missing token", http.StatusBadRequest)
			result <- callbackResult{err: fmt.Errorf("no token in callback")}
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, successHTML)

		result <- callbackResult{token: token}
	})

	server := &http.Server{Handler: mux}
	go func() { _ = server.Serve(listener) }()
	defer server.Shutdown(context.Background())

	authURL := fmt.Sprintf("%s/cli-auth?callback_url=%s&state=%s",
		webBaseURL, url.QueryEscape(callbackURL), url.QueryEscape(state))

	fmt.Println("Opening browser for authentication...")
	if err := browser.OpenURL(authURL); err != nil {
		fmt.Println("Could not open browser automatically.")
		fmt.Printf("Open this URL in your browser:\n\n  %s\n\n", authURL)
	}
	fmt.Println("Waiting for authentication...")

	select {
	case res := <-result:
		if res.err != nil {
			return res.err
		}
		return StoreToken(res.token)
	case <-time.After(loginTimeout):
		return fmt.Errorf("login timed out after %s", loginTimeout)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

const successHTML = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Cluely CLI</title>
  <style>
    body {
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 100vh;
      margin: 0;
      background: #fafafa;
      color: #111;
    }
    .container { text-align: center; padding: 2rem; }
    h1 { font-size: 1.5rem; font-weight: 600; margin-bottom: 0.5rem; }
    p { color: #666; }
  </style>
</head>
<body>
  <div class="container">
    <h1>Authentication successful</h1>
    <p>You can close this tab and return to the terminal.</p>
  </div>
</body>
</html>`
