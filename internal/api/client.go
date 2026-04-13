package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cluely/cli/internal/auth"
	"github.com/cluely/cli/internal/config"
)

// orpcRequest is the wire format ORPC expects: {"json": input, "meta": [...]}
type orpcRequest struct {
	JSON interface{} `json:"json"`
}

// orpcResponse is the wire format ORPC returns: {"json": output, "meta": [...]}
type orpcResponse struct {
	JSON json.RawMessage `json:"json"`
	Meta json.RawMessage `json:"meta,omitempty"`
}

// CallRaw makes an authenticated ORPC call and returns the raw JSON output.
func CallRaw(procedure string, input interface{}) (json.RawMessage, error) {
	var envelope orpcResponse
	if err := call(procedure, input, &envelope); err != nil {
		return nil, err
	}
	return envelope.JSON, nil
}

// Call makes an authenticated ORPC call and decodes the output.
func Call(procedure string, input interface{}, output interface{}) error {
	raw, err := CallRaw(procedure, input)
	if err != nil {
		return err
	}
	if output != nil {
		if err := json.Unmarshal(raw, output); err != nil {
			return fmt.Errorf("parse response data: %w", err)
		}
	}
	return nil
}

func call(procedure string, input interface{}, envelope *orpcResponse) error {
	token, err := auth.LoadToken()
	if err != nil {
		return fmt.Errorf("failed to read credentials: %w", err)
	}
	if token == "" {
		return fmt.Errorf("not logged in — run 'cluely auth login'")
	}

	url := config.APIURL + "/rpc/" + procedure

	// ORPC wire format wraps input in {"json": input}
	reqEnvelope := orpcRequest{JSON: input}
	data, err := json.Marshal(reqEnvelope)
	if err != nil {
		return fmt.Errorf("marshal input: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("session expired — run 'cluely auth login'")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, respBody)
	}

	if err := json.Unmarshal(respBody, envelope); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}
	return nil
}
