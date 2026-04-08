// Package strategy provides DSL-based automated trading and backtesting commands.
package strategy

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/deepcoinapi/agent-cli/pkg/client"
	"github.com/deepcoinapi/agent-cli/pkg/output"
)

// Cmd is the strategy command group.
var Cmd = &cobra.Command{
	Use:   "strategy",
	Short: "Strategy — DSL-based automated trading, backtesting",
}

func init() {
	Cmd.AddCommand(backtestCmd)
	Cmd.AddCommand(dslTriggerOrderCmd)
}

// loadDSL loads DSL from a JSON string or @filepath.
func loadDSL(dsl string) (map[string]any, error) {
	var raw string
	if strings.HasPrefix(dsl, "@") {
		data, err := os.ReadFile(dsl[1:])
		if err != nil {
			return nil, fmt.Errorf("reading DSL file: %w", err)
		}
		raw = string(data)
	} else {
		raw = dsl
	}
	var obj map[string]any
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		return nil, fmt.Errorf("invalid DSL JSON: %w", err)
	}
	return obj, nil
}

// ── backtest ────────────────────────────────────────────────────────

var backtestCmd = &cobra.Command{
	Use:   "backtest",
	Short: "Run a strategy backtest",
	Long: `Run a strategy backtest using the DSL engine.

The --dsl flag accepts either inline JSON or @filepath to read from a file.

Supported indicators: BOLL, MA, EMA, KDJ, RSI, WR
Condition operators: >=, <=, >, <, ==, cross_above, cross_below`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dslStr, _ := cmd.Flags().GetString("dsl")
		symbol, _ := cmd.Flags().GetString("symbol")
		fromTS, _ := cmd.Flags().GetString("from-ts")
		toTS, _ := cmd.Flags().GetString("to-ts")
		asJSON, _ := cmd.Flags().GetBool("json")

		dslObj, err := loadDSL(dslStr)
		if err != nil {
			return err
		}

		body := map[string]any{
			"dsl": dslObj,
			"data_source": map[string]any{
				"symbol":  symbol,
				"from_ts": fromTS,
				"to_ts":   toTS,
			},
		}
		c := client.New()
		resp, err := c.Post("/deepcoin/trade/backtest-run", body)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			data := client.GetDataMap(resp)
			if data == nil {
				output.JSON(resp)
				return nil
			}
			if summary, ok := data["summary"].(map[string]any); ok {
				fmt.Println("── Backtest Summary ──")
				fmt.Printf("  Symbol:       %v\n", summary["symbol"])
				fmt.Printf("  Realized PnL: %v\n", summary["realized_pnl"])
				fmt.Printf("  Total Trades: %v\n", summary["trades"])
				fmt.Printf("  Total Fee:    %v\n", summary["total_fee"])
			}
			if trades, ok := data["trades"].([]any); ok && len(trades) > 0 {
				fmt.Printf("\n  Trade details: %d trades (use --json for full output)\n", len(trades))
			}
		}
		return nil
	},
}

func init() {
	backtestCmd.Flags().String("dsl", "", "DSL JSON string or @filepath (required)")
	backtestCmd.Flags().String("symbol", "", "Symbol e.g. BTC-USDT-SWAP (required)")
	backtestCmd.Flags().String("from-ts", "", "Start time ISO 8601 (required)")
	backtestCmd.Flags().String("to-ts", "", "End time ISO 8601 (required)")
	backtestCmd.Flags().Bool("json", false, "Output raw JSON")
	backtestCmd.MarkFlagRequired("dsl")
	backtestCmd.MarkFlagRequired("symbol")
	backtestCmd.MarkFlagRequired("from-ts")
	backtestCmd.MarkFlagRequired("to-ts")
}

// ── dsl-trigger-order ───────────────────────────────────────────────

var dslTriggerOrderCmd = &cobra.Command{
	Use:   "dsl-trigger-order",
	Short: "Place a live DSL-driven trigger order",
	RunE: func(cmd *cobra.Command, args []string) error {
		dslStr, _ := cmd.Flags().GetString("dsl")
		symbol, _ := cmd.Flags().GetString("symbol")
		tradeMode, _ := cmd.Flags().GetString("trade-mode")
		mrgPos, _ := cmd.Flags().GetString("mrg-position")

		dslObj, err := loadDSL(dslStr)
		if err != nil {
			return err
		}

		body := map[string]any{
			"trade_info": map[string]any{
				"symbol":      symbol,
				"tradeMode":   tradeMode,
				"mrgPosition": mrgPos,
			},
			"dsl_json": dslObj,
		}
		c := client.New()
		resp, err := c.Post("/deepcoin/trade/dsl-trigger-order", body)
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	dslTriggerOrderCmd.Flags().String("dsl", "", "DSL JSON string or @filepath (required)")
	dslTriggerOrderCmd.Flags().String("symbol", "", "Symbol e.g. BTC-USDT-SWAP (required)")
	dslTriggerOrderCmd.Flags().String("trade-mode", "", "Trade mode: isolated/cross (required)")
	dslTriggerOrderCmd.Flags().String("mrg-position", "", "Position mode: merge/split (required)")
	dslTriggerOrderCmd.MarkFlagRequired("dsl")
	dslTriggerOrderCmd.MarkFlagRequired("symbol")
	dslTriggerOrderCmd.MarkFlagRequired("trade-mode")
	dslTriggerOrderCmd.MarkFlagRequired("mrg-position")
}
