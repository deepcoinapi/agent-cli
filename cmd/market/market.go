// Package market provides public market data commands.
package market

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/deepcoinapi/agent-cli/pkg/client"
	"github.com/deepcoinapi/agent-cli/pkg/output"
)

var jsonFlag bool

// Cmd is the market command group.
var Cmd = &cobra.Command{
	Use:   "market",
	Short: "Public market data — tickers, orderbook, candles, trades, funding rates",
}

func init() {
	Cmd.AddCommand(instrumentsCmd)
	Cmd.AddCommand(tickersCmd)
	Cmd.AddCommand(tickerCmd)
	Cmd.AddCommand(orderbookCmd)
	Cmd.AddCommand(candlesCmd)
	Cmd.AddCommand(tradesCmd)
	Cmd.AddCommand(fundingRateCmd)
	Cmd.AddCommand(fundingRateHistoryCmd)
	Cmd.AddCommand(bookSpreadCmd)
	Cmd.AddCommand(stepMarginCmd)
	Cmd.AddCommand(serverTimeCmd)
	Cmd.AddCommand(pingCmd)
}

// ── instruments ─────────────────────────────────────────────────────

var instrumentsCmd = &cobra.Command{
	Use:   "instruments",
	Short: "List tradeable instruments",
	RunE: func(cmd *cobra.Command, args []string) error {
		instType, _ := cmd.Flags().GetString("inst-type")
		instID, _ := cmd.Flags().GetString("inst-id")
		asJSON, _ := cmd.Flags().GetBool("json")

		c := client.New()
		params := map[string]string{"instType": instType}
		if instID != "" {
			params["instId"] = instID
		}
		resp, err := c.GetPublic("/deepcoin/market/instruments", params)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			rows := client.GetDataSlice(resp)
			output.Table(rows,
				[]string{"instId", "baseCcy", "quoteCcy", "tickSz", "lotSz", "minSz", "state"},
				[]string{"Instrument", "Base", "Quote", "Tick Size", "Lot Size", "Min Size", "State"},
			)
		}
		return nil
	},
}

func init() {
	instrumentsCmd.Flags().String("inst-type", "", "Instrument type: SPOT or SWAP (required)")
	instrumentsCmd.Flags().String("inst-id", "", "Filter by instrument ID")
	instrumentsCmd.Flags().Bool("json", false, "Output raw JSON")
	instrumentsCmd.MarkFlagRequired("inst-type")
}

// ── tickers ─────────────────────────────────────────────────────────

var tickersCmd = &cobra.Command{
	Use:   "tickers",
	Short: "Get market tickers for all instruments",
	RunE: func(cmd *cobra.Command, args []string) error {
		instType, _ := cmd.Flags().GetString("inst-type")
		asJSON, _ := cmd.Flags().GetBool("json")

		c := client.New()
		resp, err := c.GetPublic("/deepcoin/market/tickers", map[string]string{"instType": instType})
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			rows := client.GetDataSlice(resp)
			output.Table(rows,
				[]string{"instId", "last", "bidPx", "askPx", "high24h", "low24h", "vol24h"},
				[]string{"Instrument", "Last", "Bid", "Ask", "High 24h", "Low 24h", "Vol 24h"},
			)
		}
		return nil
	},
}

func init() {
	tickersCmd.Flags().String("inst-type", "", "Instrument type: SPOT or SWAP (required)")
	tickersCmd.Flags().Bool("json", false, "Output raw JSON")
	tickersCmd.MarkFlagRequired("inst-type")
}

// ── ticker ──────────────────────────────────────────────────────────

var tickerCmd = &cobra.Command{
	Use:   "ticker [INST_ID]",
	Short: "Get ticker for a single instrument",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		instID := args[0]
		asJSON, _ := cmd.Flags().GetBool("json")

		instType := "SPOT"
		if contains(instID, "SWAP") {
			instType = "SWAP"
		}
		c := client.New()
		resp, err := c.GetPublic("/deepcoin/market/tickers", map[string]string{"instType": instType})
		if err != nil {
			return err
		}
		rows := client.GetDataSlice(resp)
		var filtered []map[string]any
		for _, r := range rows {
			if fmt.Sprintf("%v", r["instId"]) == instID {
				filtered = append(filtered, r)
			}
		}
		if asJSON {
			output.JSON(map[string]any{"data": filtered})
		} else {
			output.Table(filtered,
				[]string{"instId", "last", "lastSz", "bidPx", "bidSz", "askPx", "askSz", "open24h", "high24h", "low24h", "vol24h", "volCcy24h"},
				[]string{"Instrument", "Last", "Last Size", "Bid", "Bid Size", "Ask", "Ask Size", "Open 24h", "High 24h", "Low 24h", "Vol 24h", "VolCcy 24h"},
			)
		}
		return nil
	},
}

func init() {
	tickerCmd.Flags().Bool("json", false, "Output raw JSON")
}

// ── orderbook ───────────────────────────────────────────────────────

var orderbookCmd = &cobra.Command{
	Use:   "orderbook [INST_ID]",
	Short: "Get order book depth",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		instID := args[0]
		sz, _ := cmd.Flags().GetString("sz")
		asJSON, _ := cmd.Flags().GetBool("json")

		c := client.New()
		resp, err := c.GetPublic("/deepcoin/market/books", map[string]string{"instId": instID, "sz": sz})
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
			return nil
		}

		data := client.GetDataMap(resp)
		if data == nil {
			output.JSON(resp)
			return nil
		}
		printBookSide("Asks (lowest first)", data["asks"], true)
		fmt.Println()
		printBookSide("Bids (highest first)", data["bids"], false)
		return nil
	},
}

func printBookSide(title string, raw any, reverse bool) {
	fmt.Printf("── %s ──\n", title)
	arr, ok := raw.([]any)
	if !ok || len(arr) == 0 {
		fmt.Println("  (empty)")
		return
	}
	limit := 10
	if len(arr) < limit {
		limit = len(arr)
	}
	var rows []map[string]any
	if reverse {
		for i := limit - 1; i >= 0; i-- {
			rows = append(rows, bookEntry(arr[i]))
		}
	} else {
		for i := 0; i < limit; i++ {
			rows = append(rows, bookEntry(arr[i]))
		}
	}
	output.Table(rows, []string{"price", "size"}, []string{"Price", "Size"})
}

func bookEntry(v any) map[string]any {
	if arr, ok := v.([]any); ok && len(arr) >= 2 {
		return map[string]any{"price": arr[0], "size": arr[1]}
	}
	return map[string]any{"price": v, "size": ""}
}

func init() {
	orderbookCmd.Flags().String("sz", "20", "Depth levels (max 400)")
	orderbookCmd.Flags().Bool("json", false, "Output raw JSON")
}

// ── candles ─────────────────────────────────────────────────────────

var candlesCmd = &cobra.Command{
	Use:   "candles [INST_ID]",
	Short: "Get K-line / candlestick data",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		instID := args[0]
		bar, _ := cmd.Flags().GetString("bar")
		limit, _ := cmd.Flags().GetString("limit")
		after, _ := cmd.Flags().GetString("after")
		asJSON, _ := cmd.Flags().GetBool("json")

		c := client.New()
		params := map[string]string{"instId": instID, "bar": bar, "limit": limit}
		if after != "" {
			params["after"] = after
		}
		resp, err := c.GetPublic("/deepcoin/market/candles", params)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
			return nil
		}

		// data is [][]any
		data, _ := resp["data"]
		arr, ok := data.([]any)
		if !ok {
			output.JSON(resp)
			return nil
		}
		var rows []map[string]any
		for _, item := range arr {
			if c, ok := item.([]any); ok && len(c) >= 6 {
				rows = append(rows, map[string]any{
					"ts": c[0], "open": c[1], "high": c[2], "low": c[3], "close": c[4], "vol": c[5],
				})
			}
		}
		output.Table(rows,
			[]string{"ts", "open", "high", "low", "close", "vol"},
			[]string{"Timestamp", "Open", "High", "Low", "Close", "Volume"},
		)
		return nil
	},
}

func init() {
	candlesCmd.Flags().String("bar", "1m", "Candle interval: 1m/5m/15m/30m/1H/4H/12H/1D/1W/1M/1Y")
	candlesCmd.Flags().String("limit", "100", "Number of candles (max 300)")
	candlesCmd.Flags().String("after", "", "Pagination: timestamp for older data")
	candlesCmd.Flags().Bool("json", false, "Output raw JSON")
}

// ── trades ──────────────────────────────────────────────────────────

var tradesCmd = &cobra.Command{
	Use:   "trades [INST_ID]",
	Short: "Get recent trades",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		instID := args[0]
		limit, _ := cmd.Flags().GetString("limit")
		pg, _ := cmd.Flags().GetString("product-group")
		asJSON, _ := cmd.Flags().GetBool("json")

		c := client.New()
		params := map[string]string{"instId": instID, "limit": limit}
		if pg != "" {
			params["productGroup"] = pg
		}
		resp, err := c.GetPublic("/deepcoin/market/trades", params)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			rows := client.GetDataSlice(resp)
			output.Table(rows,
				[]string{"instId", "tradeId", "px", "sz", "side", "ts"},
				[]string{"Instrument", "Trade ID", "Price", "Size", "Side", "Time"},
			)
		}
		return nil
	},
}

func init() {
	tradesCmd.Flags().String("limit", "50", "Number of trades (max 500)")
	tradesCmd.Flags().String("product-group", "", "Product group: Spot/Swap/SwapU")
	tradesCmd.Flags().Bool("json", false, "Output raw JSON")
}

// ── funding-rate ────────────────────────────────────────────────────

var fundingRateCmd = &cobra.Command{
	Use:   "funding-rate",
	Short: "Get current funding rates",
	RunE: func(cmd *cobra.Command, args []string) error {
		instType, _ := cmd.Flags().GetString("inst-type")
		instID, _ := cmd.Flags().GetString("inst-id")
		asJSON, _ := cmd.Flags().GetBool("json")

		c := client.New()
		params := map[string]string{"instType": instType}
		if instID != "" {
			params["instId"] = instID
		}
		resp, err := c.GetPublic("/deepcoin/trade/fund-rate/current-funding-rate", params)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			data := client.GetDataMap(resp)
			if rates, ok := data["current_fund_rates"]; ok {
				if arr, ok := rates.([]any); ok {
					var rows []map[string]any
					for _, item := range arr {
						if m, ok := item.(map[string]any); ok {
							rows = append(rows, m)
						}
					}
					output.Table(rows,
						[]string{"instrumentId", "fundingRate"},
						[]string{"Instrument", "Funding Rate"},
					)
					return nil
				}
			}
			output.JSON(data)
		}
		return nil
	},
}

func init() {
	fundingRateCmd.Flags().String("inst-type", "", "Instrument type: SwapU or Swap (required)")
	fundingRateCmd.Flags().String("inst-id", "", "Filter by instrument ID")
	fundingRateCmd.Flags().Bool("json", false, "Output raw JSON")
	fundingRateCmd.MarkFlagRequired("inst-type")
}

// ── funding-rate-history ────────────────────────────────────────────

var fundingRateHistoryCmd = &cobra.Command{
	Use:   "funding-rate-history [INST_ID]",
	Short: "Get historical funding rates",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		instID := args[0]
		page, _ := cmd.Flags().GetString("page")
		size, _ := cmd.Flags().GetString("size")

		c := client.New()
		resp, err := c.GetPublic("/deepcoin/trade/fund-rate/history", map[string]string{
			"instId": instID, "page": page, "size": size,
		})
		if err != nil {
			return err
		}
		output.JSON(client.GetData(resp))
		return nil
	},
}

func init() {
	fundingRateHistoryCmd.Flags().String("page", "1", "Page number")
	fundingRateHistoryCmd.Flags().String("size", "20", "Page size (max 100)")
}

// ── book-spread ─────────────────────────────────────────────────────

var bookSpreadCmd = &cobra.Command{
	Use:   "book-spread [INST_ID]",
	Short: "Get bid-ask spread information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		instID := args[0]
		value, _ := cmd.Flags().GetString("value")
		vType, _ := cmd.Flags().GetString("vtype")

		c := client.New()
		params := map[string]string{"instId": instID}
		if value != "" {
			params["value"] = value
		}
		if vType != "" {
			params["vType"] = vType
		}
		resp, err := c.GetPublic("/deepcoin/market/book-spread", params)
		if err != nil {
			return err
		}
		output.JSON(client.GetData(resp))
		return nil
	},
}

func init() {
	bookSpreadCmd.Flags().String("value", "", "Target value")
	bookSpreadCmd.Flags().String("vtype", "", "0=quoteCcy, 1=baseCcy")
}

// ── step-margin ─────────────────────────────────────────────────────

var stepMarginCmd = &cobra.Command{
	Use:   "step-margin [INST_ID]",
	Short: "Get margin tier info (SWAP only)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		instID := args[0]
		asJSON, _ := cmd.Flags().GetBool("json")

		c := client.New()
		resp, err := c.GetPublic("/deepcoin/market/step-margin", map[string]string{"instId": instID})
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			rows := client.GetDataSlice(resp)
			output.Table(rows,
				[]string{"grade", "leverage", "maxContractValue", "maintenanceMarginRate"},
				[]string{"Grade", "Leverage", "Max Contract Value", "Maintenance Margin Rate"},
			)
		}
		return nil
	},
}

func init() {
	stepMarginCmd.Flags().Bool("json", false, "Output raw JSON")
}

// ── server-time ─────────────────────────────────────────────────────

var serverTimeCmd = &cobra.Command{
	Use:   "server-time",
	Short: "Get server time",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New()
		resp, err := c.GetPublic("/deepcoin/market/time", nil)
		if err != nil {
			return err
		}
		output.JSON(client.GetData(resp))
		return nil
	},
}

// ── ping ────────────────────────────────────────────────────────────

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Check API connectivity",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New()
		_, err := c.GetPublic("/deepcoin/market/ping", nil)
		if err != nil {
			fmt.Println("FAIL:", err)
			return err
		}
		fmt.Println("OK")
		return nil
	},
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
