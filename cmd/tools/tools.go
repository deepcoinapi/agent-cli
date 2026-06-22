// Package tools exposes machine-readable CLI capability discovery.
package tools

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

type commandSpec struct {
	Skill       string `json:"skill"`
	Group       string `json:"group"`
	Command     string `json:"command"`
	Description string `json:"description"`
	Auth        string `json:"auth"`
	Type        string `json:"type"`
}

// ListToolsCmd prints the stable command surface for agents.
var ListToolsCmd = &cobra.Command{
	Use:   "list-tools",
	Short: "List stable agent command entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)
		return enc.Encode(commandSpecs)
	},
}

var commandSpecs = []commandSpec{
	{"deepcoin-market", "market", "deepcoin-cli market instruments --inst-type <SPOT|SWAP> [--inst-id <id>] [--json]", "List tradeable instruments", "none", "READ"},
	{"deepcoin-market", "market", "deepcoin-cli market tickers --inst-type <SPOT|SWAP> [--json]", "Get tickers for an instrument type", "none", "READ"},
	{"deepcoin-market", "market", "deepcoin-cli market ticker <INST_ID> [--json]", "Get ticker for one instrument", "none", "READ"},
	{"deepcoin-market", "market", "deepcoin-cli market orderbook <INST_ID> [--sz <n>] [--json]", "Get order book depth", "none", "READ"},
	{"deepcoin-market", "market", "deepcoin-cli market candles <INST_ID> [--bar <bar>] [--limit <n>] [--after <ts>] [--json]", "Get candles", "none", "READ"},
	{"deepcoin-market", "market", "deepcoin-cli market trades <INST_ID> [--product-group <Spot|Swap|SwapU>] [--limit <n>] [--json]", "Get recent trades", "none", "READ"},
	{"deepcoin-market", "market", "deepcoin-cli market funding-rate --inst-type <SwapU|Swap> [--inst-id <id>] [--json]", "Get current funding rates", "none", "READ"},
	{"deepcoin-market", "market", "deepcoin-cli market funding-rate-history <INST_ID> [--page <n>] [--size <n>]", "Get funding rate history", "none", "READ"},
	{"deepcoin-market", "market", "deepcoin-cli market book-spread <INST_ID> [--value <value>] [--vtype <0|1>]", "Get bid-ask spread", "none", "READ"},
	{"deepcoin-market", "market", "deepcoin-cli market step-margin <INST_ID> [--json]", "Get margin tiers", "none", "READ"},
	{"deepcoin-market", "market", "deepcoin-cli market server-time", "Get server time", "none", "READ"},
	{"deepcoin-market", "market", "deepcoin-cli market ping", "Check connectivity", "none", "READ"},

	{"deepcoin-trade", "trade", "deepcoin-cli trade place-order --inst-id <id> --td-mode <mode> --side <buy|sell> --ord-type <type> --sz <size> [flags] [--json]", "Place order", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade batch-orders --orders '<json-array>'", "Place up to 5 orders", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade cancel-order --inst-id <id> --ord-id <id> [--json]", "Cancel order", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade batch-cancel --order-ids '<id,id>'", "Cancel up to 50 orders", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade cancel-all --product-group <Swap|SwapU> [--inst-id <id>] [--cross-margin <0|1>] [--merge-mode <0|1>]", "Cancel all swap orders", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade amend-order --order-id <id> [--price <px>] [--volume <sz>]", "Amend order", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade amend-order-sltp --order-id <id> [--tp-trigger-px <px>] [--sl-trigger-px <px>]", "Amend order TP/SL", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade get-order --inst-id <id> --ord-id <id>", "Get active/recent order", "required", "READ"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade get-history-order --inst-id <id> --ord-id <id>", "Get historical order", "required", "READ"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade pending-orders [--inst-id <id>] [--limit <n>] [--json]", "List pending orders", "required", "READ"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade order-history --inst-type <SPOT|SWAP> [flags] [--json]", "List order history", "required", "READ"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade batch-query --orders '<json-array>'", "Query up to 5 orders", "required", "READ"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade fills --inst-type <SPOT|SWAP> [flags] [--json]", "Get trade fills", "required", "READ"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade trigger-order --inst-id <id> --product-group <Swap|SwapU> --side <buy|sell> --sz <size> --trigger-price <px> [flags]", "Place trigger order", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade cancel-trigger --inst-id <id> --ord-id <id>", "Cancel trigger order", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade cancel-all-triggers --product-group <Swap|SwapU> [flags]", "Cancel all trigger orders", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade trigger-pending --inst-type <SPOT|SWAP> [--inst-id <id>] [--limit <n>] [--json]", "List pending trigger orders", "required", "READ"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade trigger-history --inst-type <SPOT|SWAP> [--inst-id <id>] [--limit <n>] [--json]", "List trigger history", "required", "READ"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade set-position-sltp --inst-type <SPOT|SWAP> --inst-id <id> --pos-side <side> [flags]", "Set position TP/SL", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade modify-position-sltp --ord-id <id> --inst-id <id> [flags]", "Modify position TP/SL", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade cancel-position-sltp --ord-id <id>", "Cancel position TP/SL", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade close-position --inst-id <id> --product-group <Swap|SwapU> --position-ids '<id,id>'", "Close positions by ID", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade batch-close-position --inst-id <id> --product-group <Swap|SwapU>", "Close all positions for an instrument", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade trace-order --inst-id <id> --retrace-point <value> --trigger-price <px> --pos-side <long|short>", "Place trace order", "required", "WRITE"},
	{"deepcoin-trade", "trade", "deepcoin-cli trade trace-orders [--json]", "List pending trace orders", "required", "READ"},

	{"deepcoin-portfolio", "account", "deepcoin-cli account balance [--inst-type <SPOT|SWAP>] [--ccy <ccy>] [--json]", "Get account balance", "required", "READ"},
	{"deepcoin-portfolio", "account", "deepcoin-cli account positions [--inst-type <SPOT|SWAP>] [--inst-id <id>] [--json]", "Get positions", "required", "READ"},
	{"deepcoin-portfolio", "account", "deepcoin-cli account bills --inst-type <SPOT|SWAP> [flags] [--json]", "Get account bills", "required", "READ"},
	{"deepcoin-portfolio", "account", "deepcoin-cli account set-leverage --inst-id <id> --lever <n> --mgn-mode <cross|isolated> [--mrg-position <merge|split>]", "Set leverage", "required", "WRITE"},
	{"deepcoin-portfolio", "account", "deepcoin-cli account uid", "Get account UID", "required", "READ"},
	{"deepcoin-portfolio", "account", "deepcoin-cli account sub-accounts [--json]", "List sub-accounts", "required", "READ"},
	{"deepcoin-portfolio", "account", "deepcoin-cli account sub-account-transfer [flags]", "Transfer between sub-accounts", "required", "WRITE"},
	{"deepcoin-portfolio", "account", "deepcoin-cli account transfer [flags]", "Transfer assets between accounts", "required", "WRITE"},
	{"deepcoin-portfolio", "account", "deepcoin-cli account internal-transfer [flags]", "Make internal transfer", "required", "WRITE"},

	{"deepcoin-withdrawal", "withdrawal", "deepcoin-cli withdrawal config [--ccy <ccy>] [--include-addresses true]", "Get aggregated withdrawal config", "required", "READ"},
	{"deepcoin-withdrawal", "withdrawal", "deepcoin-cli withdrawal assets [--ccy <ccy>]", "List withdrawable assets", "required", "READ"},
	{"deepcoin-withdrawal", "withdrawal", "deepcoin-cli withdrawal chains --ccy <ccy>", "List withdrawal chains", "required", "READ"},
	{"deepcoin-withdrawal", "withdrawal", "deepcoin-cli withdrawal addresses --ccy <ccy>", "List whitelist addresses", "required", "READ"},
	{"deepcoin-withdrawal", "withdrawal", "deepcoin-cli withdrawal records [flags]", "List withdrawal records", "required", "READ"},
	{"deepcoin-withdrawal", "withdrawal", "deepcoin-cli withdrawal status --wd-id <id> [--ccy <ccy>]", "Get withdrawal status", "required", "READ"},
	{"deepcoin-withdrawal", "withdrawal", "deepcoin-cli withdrawal create --ccy <ccy> --chain <chain> --amt <amount> --address-id <id> [flags]", "Create withdrawal", "required", "WRITE"},
	{"deepcoin-withdrawal", "withdrawal", "deepcoin-cli withdrawal cancel --wd-id <id> [--ccy <ccy>] [--client-id <id>]", "Cancel withdrawal", "required", "WRITE"},

	{"deepcoin-copytrade", "copytrade", "deepcoin-cli copytrade leader-settings --status <0|1> [flags]", "Update leader settings", "required", "WRITE"},
	{"deepcoin-copytrade", "copytrade", "deepcoin-cli copytrade support-contracts [--json]", "List supported contracts", "required", "READ"},
	{"deepcoin-copytrade", "copytrade", "deepcoin-cli copytrade set-contracts --contracts '<BTCUSDT,ETHUSDT>'", "Set leader contracts", "required", "WRITE"},
	{"deepcoin-copytrade", "copytrade", "deepcoin-cli copytrade followers --status <1|2> [--json]", "List followers", "required", "READ"},
	{"deepcoin-copytrade", "copytrade", "deepcoin-cli copytrade leader-positions [--page <n>] [--size <n>] [--json]", "List leader positions", "required", "READ"},
	{"deepcoin-copytrade", "copytrade", "deepcoin-cli copytrade position-type [--json]", "Get position type", "required", "READ"},
	{"deepcoin-copytrade", "copytrade", "deepcoin-cli copytrade set-position-type --type <1|2>", "Set position type", "required", "WRITE"},
	{"deepcoin-copytrade", "copytrade", "deepcoin-cli copytrade estimated-profit [--json]", "Get estimated profit", "required", "READ"},
	{"deepcoin-copytrade", "copytrade", "deepcoin-cli copytrade history-profit [--json]", "Get historical profit", "required", "READ"},

	{"deepcoin-strategy", "strategy", "deepcoin-cli strategy backtest --symbol <id> --from-ts <ts> --to-ts <ts> --dsl @strategy.json [--json]", "Run strategy backtest", "required", "READ"},
	{"deepcoin-strategy", "strategy", "deepcoin-cli strategy dsl-trigger-order --symbol <id> --trade-mode <cross|isolated> --mrg-position <merge|split> --dsl @strategy.json", "Place live DSL trigger order", "required", "WRITE"},
}
