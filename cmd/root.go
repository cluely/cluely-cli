package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/cluely/cli/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cluely",
	Short: "Cluely CLI",
	Long:  "Command-line interface for the Cluely platform.",
	// No Run — root without subcommand prints help.
	SilenceUsage: true,
}

func init() {
	rootCmd.Version = version.Full()
	rootCmd.SetVersionTemplate("cluely {{.Version}}\n")
}

// Execute runs the root command with signal-aware context.
func Execute() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	return rootCmd.ExecuteContext(ctx)
}
