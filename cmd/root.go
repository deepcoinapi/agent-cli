// Package cmd defines the CLI command tree.
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/deepcoinapi/agent-cli/cmd/account"
	"github.com/deepcoinapi/agent-cli/cmd/copytrade"
	"github.com/deepcoinapi/agent-cli/cmd/market"
	"github.com/deepcoinapi/agent-cli/cmd/strategy"
	"github.com/deepcoinapi/agent-cli/cmd/tools"
	"github.com/deepcoinapi/agent-cli/cmd/trade"
	"github.com/deepcoinapi/agent-cli/cmd/withdrawal"
)

var rootCmd = &cobra.Command{
	Use:     "dcli",
	Short:   "DeepCoin Agent CLI — interact with DeepCoin exchange",
	Version: "0.1.0",
}

func init() {
	rootCmd.AddCommand(market.Cmd)
	rootCmd.AddCommand(trade.Cmd)
	rootCmd.AddCommand(account.Cmd)
	rootCmd.AddCommand(withdrawal.Cmd)
	rootCmd.AddCommand(copytrade.Cmd)
	rootCmd.AddCommand(strategy.Cmd)
	rootCmd.AddCommand(tools.ListToolsCmd)
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
