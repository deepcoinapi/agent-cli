// Package trade provides order management commands.
package trade

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/deepcoinapi/agent-cli/pkg/client"
	"github.com/deepcoinapi/agent-cli/pkg/output"
)

// Cmd is the trade command group.
var Cmd = &cobra.Command{
	Use:   "trade",
	Short: "Trading — place, cancel, amend orders, triggers, TP/SL",
}

func init() {
	Cmd.AddCommand(placeOrderCmd)
	Cmd.AddCommand(batchOrdersCmd)
	Cmd.AddCommand(cancelOrderCmd)
	Cmd.AddCommand(batchCancelCmd)
	Cmd.AddCommand(cancelAllCmd)
	Cmd.AddCommand(amendOrderCmd)
	Cmd.AddCommand(amendOrderSltpCmd)
	Cmd.AddCommand(getOrderCmd)
	Cmd.AddCommand(getHistoryOrderCmd)
	Cmd.AddCommand(pendingOrdersCmd)
	Cmd.AddCommand(orderHistoryCmd)
	Cmd.AddCommand(batchQueryCmd)
	Cmd.AddCommand(fillsCmd)
	Cmd.AddCommand(triggerOrderCmd)
	Cmd.AddCommand(cancelTriggerCmd)
	Cmd.AddCommand(cancelAllTriggersCmd)
	Cmd.AddCommand(triggerPendingCmd)
	Cmd.AddCommand(triggerHistoryCmd)
	Cmd.AddCommand(setPositionSltpCmd)
	Cmd.AddCommand(modifyPositionSltpCmd)
	Cmd.AddCommand(cancelPositionSltpCmd)
	Cmd.AddCommand(closePositionCmd)
	Cmd.AddCommand(batchClosePositionCmd)
	Cmd.AddCommand(traceOrderCmd)
	Cmd.AddCommand(traceOrdersCmd)
}

// helper to build body from flags, skipping empty values
func buildBody(pairs ...string) map[string]any {
	body := make(map[string]any)
	for i := 0; i+1 < len(pairs); i += 2 {
		if pairs[i+1] != "" {
			body[pairs[i]] = pairs[i+1]
		}
	}
	return body
}

// ── place-order ─────────────────────────────────────────────────────

var placeOrderCmd = &cobra.Command{
	Use:   "place-order",
	Short: "Place a new order",
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Flags()
		instID, _ := f.GetString("inst-id")
		tdMode, _ := f.GetString("td-mode")
		side, _ := f.GetString("side")
		ordType, _ := f.GetString("ord-type")
		sz, _ := f.GetString("sz")
		px, _ := f.GetString("px")
		posSide, _ := f.GetString("pos-side")
		mrgPos, _ := f.GetString("mrg-position")
		tpPx, _ := f.GetString("tp-trigger-px")
		slPx, _ := f.GetString("sl-trigger-px")
		clOrdID, _ := f.GetString("cl-ord-id")
		reduceOnly, _ := f.GetBool("reduce-only")
		tgtCcy, _ := f.GetString("tgt-ccy")
		asJSON, _ := f.GetBool("json")

		body := buildBody(
			"instId", instID, "tdMode", tdMode, "side", side,
			"ordType", ordType, "sz", sz, "px", px,
			"posSide", posSide, "mrgPosition", mrgPos,
			"tpTriggerPx", tpPx, "slTriggerPx", slPx,
			"clOrdId", clOrdID, "tgtCcy", tgtCcy,
		)
		if reduceOnly {
			body["reduceOnly"] = true
		}

		c := client.New()
		resp, err := c.Post("/deepcoin/trade/order", body)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			d := client.GetDataMap(resp)
			fmt.Printf("Order placed: ordId=%v  sCode=%v  sMsg=%v\n",
				d["ordId"], d["sCode"], d["sMsg"])
		}
		return nil
	},
}

func init() {
	f := placeOrderCmd.Flags()
	f.String("inst-id", "", "Instrument ID (required)")
	f.String("td-mode", "", "Trade mode: isolated/cross/cash (required)")
	f.String("side", "", "Order side: buy/sell (required)")
	f.String("ord-type", "", "Order type: market/limit/post_only/ioc (required)")
	f.String("sz", "", "Order size (required)")
	f.String("px", "", "Price (required for limit/post_only)")
	f.String("pos-side", "", "Position side: long/short (SWAP)")
	f.String("mrg-position", "", "Position mode: merge/split (SWAP)")
	f.String("tp-trigger-px", "", "Take profit trigger price")
	f.String("sl-trigger-px", "", "Stop loss trigger price")
	f.String("cl-ord-id", "", "Custom order ID")
	f.Bool("reduce-only", false, "Reduce only")
	f.String("tgt-ccy", "", "Target currency: base_ccy/quote_ccy")
	f.Bool("json", false, "Output raw JSON")
	placeOrderCmd.MarkFlagRequired("inst-id")
	placeOrderCmd.MarkFlagRequired("td-mode")
	placeOrderCmd.MarkFlagRequired("side")
	placeOrderCmd.MarkFlagRequired("ord-type")
	placeOrderCmd.MarkFlagRequired("sz")
}

// ── batch-orders ────────────────────────────────────────────────────

var batchOrdersCmd = &cobra.Command{
	Use:   "batch-orders",
	Short: "Place multiple orders at once (max 5)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ordersStr, _ := cmd.Flags().GetString("orders")
		var orders []any
		if err := json.Unmarshal([]byte(ordersStr), &orders); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}
		c := client.New()
		resp, err := c.Post("/deepcoin/trade/batch-orders", map[string]any{"orders": orders})
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	batchOrdersCmd.Flags().String("orders", "", "JSON array of order objects (required)")
	batchOrdersCmd.MarkFlagRequired("orders")
}

// ── cancel-order ────────────────────────────────────────────────────

var cancelOrderCmd = &cobra.Command{
	Use:   "cancel-order",
	Short: "Cancel an existing order",
	RunE: func(cmd *cobra.Command, args []string) error {
		instID, _ := cmd.Flags().GetString("inst-id")
		ordID, _ := cmd.Flags().GetString("ord-id")
		asJSON, _ := cmd.Flags().GetBool("json")

		c := client.New()
		resp, err := c.Post("/deepcoin/trade/cancel-order", map[string]any{
			"instId": instID, "ordId": ordID,
		})
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			d := client.GetDataMap(resp)
			fmt.Printf("Cancelled: ordId=%v  sCode=%v  sMsg=%v\n",
				d["ordId"], d["sCode"], d["sMsg"])
		}
		return nil
	},
}

func init() {
	cancelOrderCmd.Flags().String("inst-id", "", "Instrument ID (required)")
	cancelOrderCmd.Flags().String("ord-id", "", "Order ID (required)")
	cancelOrderCmd.Flags().Bool("json", false, "Output raw JSON")
	cancelOrderCmd.MarkFlagRequired("inst-id")
	cancelOrderCmd.MarkFlagRequired("ord-id")
}

// ── batch-cancel ────────────────────────────────────────────────────

var batchCancelCmd = &cobra.Command{
	Use:   "batch-cancel",
	Short: "Cancel multiple orders (max 50)",
	RunE: func(cmd *cobra.Command, args []string) error {
		idsStr, _ := cmd.Flags().GetString("order-ids")
		ids := strings.Split(idsStr, ",")
		for i := range ids {
			ids[i] = strings.TrimSpace(ids[i])
		}
		c := client.New()
		resp, err := c.Post("/deepcoin/trade/batch-cancel-order", map[string]any{"orderSysIDs": ids})
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	batchCancelCmd.Flags().String("order-ids", "", "Comma-separated order IDs (required)")
	batchCancelCmd.MarkFlagRequired("order-ids")
}

// ── cancel-all ──────────────────────────────────────────────────────

var cancelAllCmd = &cobra.Command{
	Use:   "cancel-all",
	Short: "Cancel all orders for a product group",
	RunE: func(cmd *cobra.Command, args []string) error {
		pg, _ := cmd.Flags().GetString("product-group")
		instID, _ := cmd.Flags().GetString("inst-id")
		cm, _ := cmd.Flags().GetString("cross-margin")
		mm, _ := cmd.Flags().GetString("merge-mode")

		body := buildBody("ProductGroup", pg, "IsCrossMargin", cm, "IsMergeMode", mm, "InstrumentID", instID)
		c := client.New()
		resp, err := c.Post("/deepcoin/trade/swap/cancel-all", body)
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	cancelAllCmd.Flags().String("product-group", "", "Product group: Swap/SwapU (required)")
	cancelAllCmd.Flags().String("inst-id", "", "Instrument ID")
	cancelAllCmd.Flags().String("cross-margin", "0", "Is cross margin: 0/1")
	cancelAllCmd.Flags().String("merge-mode", "0", "Is merge mode: 0/1")
	cancelAllCmd.MarkFlagRequired("product-group")
}

// ── amend-order ─────────────────────────────────────────────────────

var amendOrderCmd = &cobra.Command{
	Use:   "amend-order",
	Short: "Amend/modify an existing order",
	RunE: func(cmd *cobra.Command, args []string) error {
		orderID, _ := cmd.Flags().GetString("order-id")
		price, _ := cmd.Flags().GetString("price")
		volume, _ := cmd.Flags().GetString("volume")

		body := buildBody("OrderSysID", orderID, "price", price, "volume", volume)
		c := client.New()
		resp, err := c.Post("/deepcoin/trade/replace-order", body)
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	amendOrderCmd.Flags().String("order-id", "", "Order system ID (required)")
	amendOrderCmd.Flags().String("price", "", "New price")
	amendOrderCmd.Flags().String("volume", "", "New volume")
	amendOrderCmd.MarkFlagRequired("order-id")
}

// ── amend-order-sltp ────────────────────────────────────────────────

var amendOrderSltpCmd = &cobra.Command{
	Use:   "amend-order-sltp",
	Short: "Amend TP/SL on an existing order",
	RunE: func(cmd *cobra.Command, args []string) error {
		orderID, _ := cmd.Flags().GetString("order-id")
		tp, _ := cmd.Flags().GetString("tp-trigger-px")
		sl, _ := cmd.Flags().GetString("sl-trigger-px")

		body := buildBody("orderSysID", orderID, "tpTriggerPx", tp, "slTriggerPx", sl)
		c := client.New()
		resp, err := c.Post("/deepcoin/trade/replace-order-sltp", body)
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	amendOrderSltpCmd.Flags().String("order-id", "", "Order system ID (required)")
	amendOrderSltpCmd.Flags().String("tp-trigger-px", "", "New TP trigger price")
	amendOrderSltpCmd.Flags().String("sl-trigger-px", "", "New SL trigger price")
	amendOrderSltpCmd.MarkFlagRequired("order-id")
}

// ── get-order ───────────────────────────────────────────────────────

var getOrderCmd = &cobra.Command{
	Use:   "get-order",
	Short: "Get details of a specific order",
	RunE: func(cmd *cobra.Command, args []string) error {
		instID, _ := cmd.Flags().GetString("inst-id")
		ordID, _ := cmd.Flags().GetString("ord-id")

		c := client.New()
		resp, err := c.Get("/deepcoin/trade/orderByID", map[string]string{
			"instId": instID, "ordId": ordID,
		})
		if err != nil {
			return err
		}
		output.JSON(client.GetData(resp))
		return nil
	},
}

func init() {
	getOrderCmd.Flags().String("inst-id", "", "Instrument ID (required)")
	getOrderCmd.Flags().String("ord-id", "", "Order ID (required)")
	getOrderCmd.MarkFlagRequired("inst-id")
	getOrderCmd.MarkFlagRequired("ord-id")
}

// ── get-history-order ───────────────────────────────────────────────

var getHistoryOrderCmd = &cobra.Command{
	Use:   "get-history-order",
	Short: "Get a historical (finished) order by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		instID, _ := cmd.Flags().GetString("inst-id")
		ordID, _ := cmd.Flags().GetString("ord-id")

		c := client.New()
		resp, err := c.Get("/deepcoin/trade/finishOrderByID", map[string]string{
			"instId": instID, "ordId": ordID,
		})
		if err != nil {
			return err
		}
		output.JSON(client.GetData(resp))
		return nil
	},
}

func init() {
	getHistoryOrderCmd.Flags().String("inst-id", "", "Instrument ID (required)")
	getHistoryOrderCmd.Flags().String("ord-id", "", "Order ID (required)")
	getHistoryOrderCmd.MarkFlagRequired("inst-id")
	getHistoryOrderCmd.MarkFlagRequired("ord-id")
}

// ── pending-orders ──────────────────────────────────────────────────

var pendingOrdersCmd = &cobra.Command{
	Use:   "pending-orders",
	Short: "List pending (open) orders",
	RunE: func(cmd *cobra.Command, args []string) error {
		instID, _ := cmd.Flags().GetString("inst-id")
		limit, _ := cmd.Flags().GetString("limit")
		asJSON, _ := cmd.Flags().GetBool("json")

		params := map[string]string{"limit": limit}
		if instID != "" {
			params["instId"] = instID
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/trade/v2/orders-pending", params)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			rows := client.GetDataSlice(resp)
			output.Table(rows,
				[]string{"ordId", "instId", "side", "ordType", "px", "sz", "state"},
				[]string{"Order ID", "Instrument", "Side", "Type", "Price", "Size", "State"},
			)
		}
		return nil
	},
}

func init() {
	pendingOrdersCmd.Flags().String("inst-id", "", "Filter by instrument ID")
	pendingOrdersCmd.Flags().String("limit", "30", "Max results (max 100)")
	pendingOrdersCmd.Flags().Bool("json", false, "Output raw JSON")
}

// ── order-history ───────────────────────────────────────────────────

var orderHistoryCmd = &cobra.Command{
	Use:   "order-history",
	Short: "Get order history",
	RunE: func(cmd *cobra.Command, args []string) error {
		instType, _ := cmd.Flags().GetString("inst-type")
		instID, _ := cmd.Flags().GetString("inst-id")
		state, _ := cmd.Flags().GetString("state")
		ordType, _ := cmd.Flags().GetString("ord-type")
		limit, _ := cmd.Flags().GetString("limit")
		asJSON, _ := cmd.Flags().GetBool("json")

		params := map[string]string{"instType": instType, "limit": limit}
		if instID != "" {
			params["instId"] = instID
		}
		if state != "" {
			params["state"] = state
		}
		if ordType != "" {
			params["ordType"] = ordType
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/trade/orders-history", params)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			rows := client.GetDataSlice(resp)
			output.Table(rows,
				[]string{"ordId", "instId", "side", "ordType", "px", "sz", "avgPx", "state", "uTime"},
				[]string{"Order ID", "Instrument", "Side", "Type", "Price", "Size", "Avg Price", "State", "Time"},
			)
		}
		return nil
	},
}

func init() {
	orderHistoryCmd.Flags().String("inst-type", "", "Instrument type: SPOT/SWAP (required)")
	orderHistoryCmd.Flags().String("inst-id", "", "Filter by instrument")
	orderHistoryCmd.Flags().String("state", "", "Filter: canceled/filled")
	orderHistoryCmd.Flags().String("ord-type", "", "Filter by order type")
	orderHistoryCmd.Flags().String("limit", "20", "Max results (max 100)")
	orderHistoryCmd.Flags().Bool("json", false, "Output raw JSON")
	orderHistoryCmd.MarkFlagRequired("inst-type")
}

// ── batch-query ─────────────────────────────────────────────────────

var batchQueryCmd = &cobra.Command{
	Use:   "batch-query",
	Short: "Query multiple orders at once (max 5)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ordersStr, _ := cmd.Flags().GetString("orders")
		var orders []any
		if err := json.Unmarshal([]byte(ordersStr), &orders); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}
		c := client.New()
		resp, err := c.Post("/deepcoin/trade/batch-order-query", map[string]any{"orders": orders})
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	batchQueryCmd.Flags().String("orders", "", `JSON array: [{"instId":"...","ordId":"..."}] (required)`)
	batchQueryCmd.MarkFlagRequired("orders")
}

// ── fills ───────────────────────────────────────────────────────────

var fillsCmd = &cobra.Command{
	Use:   "fills",
	Short: "Get trade fill history",
	RunE: func(cmd *cobra.Command, args []string) error {
		instType, _ := cmd.Flags().GetString("inst-type")
		instID, _ := cmd.Flags().GetString("inst-id")
		ordID, _ := cmd.Flags().GetString("ord-id")
		limit, _ := cmd.Flags().GetString("limit")
		asJSON, _ := cmd.Flags().GetBool("json")

		params := map[string]string{"instType": instType, "limit": limit}
		if instID != "" {
			params["instId"] = instID
		}
		if ordID != "" {
			params["ordId"] = ordID
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/trade/fills", params)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			rows := client.GetDataSlice(resp)
			output.Table(rows,
				[]string{"tradeId", "instId", "fillPx", "fillSz", "side", "fee", "ts"},
				[]string{"Trade ID", "Instrument", "Price", "Size", "Side", "Fee", "Time"},
			)
		}
		return nil
	},
}

func init() {
	fillsCmd.Flags().String("inst-type", "", "Instrument type: SPOT/SWAP (required)")
	fillsCmd.Flags().String("inst-id", "", "Filter by instrument")
	fillsCmd.Flags().String("ord-id", "", "Filter by order ID")
	fillsCmd.Flags().String("limit", "20", "Max results (max 100)")
	fillsCmd.Flags().Bool("json", false, "Output raw JSON")
	fillsCmd.MarkFlagRequired("inst-type")
}

// ── trigger-order ───────────────────────────────────────────────────

var triggerOrderCmd = &cobra.Command{
	Use:   "trigger-order",
	Short: "Place a trigger (conditional) order",
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Flags()
		instID, _ := f.GetString("inst-id")
		pg, _ := f.GetString("product-group")
		side, _ := f.GetString("side")
		sz, _ := f.GetString("sz")
		trigPx, _ := f.GetString("trigger-price")
		trigType, _ := f.GetString("trigger-px-type")
		ordType, _ := f.GetString("order-type")
		price, _ := f.GetString("price")
		posSide, _ := f.GetString("pos-side")
		tdMode, _ := f.GetString("td-mode")
		cm, _ := f.GetString("cross-margin")
		mrgPos, _ := f.GetString("mrg-position")
		tp, _ := f.GetString("tp-trigger-px")
		sl, _ := f.GetString("sl-trigger-px")
		asJSON, _ := f.GetBool("json")

		body := buildBody(
			"instId", instID, "productGroup", pg, "side", side, "sz", sz,
			"triggerPrice", trigPx, "triggerPxType", trigType,
			"orderType", ordType, "price", price, "posSide", posSide,
			"tdMode", tdMode, "isCrossMargin", cm, "mrgPosition", mrgPos,
			"tpTriggerPx", tp, "slTriggerPx", sl,
		)
		c := client.New()
		resp, err := c.Post("/deepcoin/trade/trigger-order", body)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			fmt.Printf("Trigger order placed: %v\n", client.GetData(resp))
		}
		return nil
	},
}

func init() {
	f := triggerOrderCmd.Flags()
	f.String("inst-id", "", "Instrument ID (required)")
	f.String("product-group", "", "Product group: Swap/SwapU (required)")
	f.String("side", "", "Order side: buy/sell (required)")
	f.String("sz", "", "Order size (required)")
	f.String("trigger-price", "", "Trigger price (required)")
	f.String("trigger-px-type", "last", "Trigger price type: last/index/mark")
	f.String("order-type", "market", "Order type: market/limit")
	f.String("price", "", "Limit price")
	f.String("pos-side", "", "Position side: long/short")
	f.String("td-mode", "isolated", "Trade mode: isolated/cross")
	f.String("cross-margin", "", "Is cross margin")
	f.String("mrg-position", "", "Position mode: merge/split")
	f.String("tp-trigger-px", "", "Take profit trigger price")
	f.String("sl-trigger-px", "", "Stop loss trigger price")
	f.Bool("json", false, "Output raw JSON")
	triggerOrderCmd.MarkFlagRequired("inst-id")
	triggerOrderCmd.MarkFlagRequired("product-group")
	triggerOrderCmd.MarkFlagRequired("side")
	triggerOrderCmd.MarkFlagRequired("sz")
	triggerOrderCmd.MarkFlagRequired("trigger-price")
}

// ── cancel-trigger ──────────────────────────────────────────────────

var cancelTriggerCmd = &cobra.Command{
	Use:   "cancel-trigger",
	Short: "Cancel a trigger order",
	RunE: func(cmd *cobra.Command, args []string) error {
		instID, _ := cmd.Flags().GetString("inst-id")
		ordID, _ := cmd.Flags().GetString("ord-id")

		c := client.New()
		resp, err := c.Post("/deepcoin/trade/cancel-trigger-order", map[string]any{
			"instId": instID, "ordId": ordID,
		})
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	cancelTriggerCmd.Flags().String("inst-id", "", "Instrument ID (required)")
	cancelTriggerCmd.Flags().String("ord-id", "", "Trigger order ID (required)")
	cancelTriggerCmd.MarkFlagRequired("inst-id")
	cancelTriggerCmd.MarkFlagRequired("ord-id")
}

// ── cancel-all-triggers ─────────────────────────────────────────────

var cancelAllTriggersCmd = &cobra.Command{
	Use:   "cancel-all-triggers",
	Short: "Cancel all trigger orders",
	RunE: func(cmd *cobra.Command, args []string) error {
		pg, _ := cmd.Flags().GetString("product-group")
		instID, _ := cmd.Flags().GetString("inst-id")
		cm, _ := cmd.Flags().GetString("cross-margin")
		mm, _ := cmd.Flags().GetString("merge-mode")

		body := buildBody("ProductGroup", pg, "InstrumentID", instID, "IsCrossMargin", cm, "IsMergeMode", mm)
		c := client.New()
		resp, err := c.Post("/deepcoin/trade/swap/cancel-trigger-all", body)
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	cancelAllTriggersCmd.Flags().String("product-group", "", "Product group: Swap/SwapU (required)")
	cancelAllTriggersCmd.Flags().String("inst-id", "", "Instrument ID")
	cancelAllTriggersCmd.Flags().String("cross-margin", "-1", "Cross margin filter: -1/0/1")
	cancelAllTriggersCmd.Flags().String("merge-mode", "-1", "Merge mode filter: -1/0/1")
	cancelAllTriggersCmd.MarkFlagRequired("product-group")
}

// ── trigger-pending ─────────────────────────────────────────────────

var triggerPendingCmd = &cobra.Command{
	Use:   "trigger-pending",
	Short: "List pending trigger orders",
	RunE: func(cmd *cobra.Command, args []string) error {
		instType, _ := cmd.Flags().GetString("inst-type")
		instID, _ := cmd.Flags().GetString("inst-id")
		limit, _ := cmd.Flags().GetString("limit")
		asJSON, _ := cmd.Flags().GetBool("json")

		params := map[string]string{"instType": instType, "limit": limit}
		if instID != "" {
			params["instId"] = instID
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/trade/trigger-orders-pending", params)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			rows := client.GetDataSlice(resp)
			output.Table(rows, nil, nil)
		}
		return nil
	},
}

func init() {
	triggerPendingCmd.Flags().String("inst-type", "", "Instrument type: SPOT/SWAP (required)")
	triggerPendingCmd.Flags().String("inst-id", "", "Filter by instrument")
	triggerPendingCmd.Flags().String("limit", "20", "Max results (max 100)")
	triggerPendingCmd.Flags().Bool("json", false, "Output raw JSON")
	triggerPendingCmd.MarkFlagRequired("inst-type")
}

// ── trigger-history ─────────────────────────────────────────────────

var triggerHistoryCmd = &cobra.Command{
	Use:   "trigger-history",
	Short: "Get trigger order history",
	RunE: func(cmd *cobra.Command, args []string) error {
		instType, _ := cmd.Flags().GetString("inst-type")
		instID, _ := cmd.Flags().GetString("inst-id")
		limit, _ := cmd.Flags().GetString("limit")
		asJSON, _ := cmd.Flags().GetBool("json")

		params := map[string]string{"instType": instType, "limit": limit}
		if instID != "" {
			params["instId"] = instID
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/trade/trigger-orders-history", params)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			rows := client.GetDataSlice(resp)
			output.Table(rows, nil, nil)
		}
		return nil
	},
}

func init() {
	triggerHistoryCmd.Flags().String("inst-type", "", "Instrument type: SPOT/SWAP (required)")
	triggerHistoryCmd.Flags().String("inst-id", "", "Filter by instrument")
	triggerHistoryCmd.Flags().String("limit", "20", "Max results (max 100)")
	triggerHistoryCmd.Flags().Bool("json", false, "Output raw JSON")
	triggerHistoryCmd.MarkFlagRequired("inst-type")
}

// ── set-position-sltp ───────────────────────────────────────────────

var setPositionSltpCmd = &cobra.Command{
	Use:   "set-position-sltp",
	Short: "Set TP/SL on a position",
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Flags()
		instType, _ := f.GetString("inst-type")
		instID, _ := f.GetString("inst-id")
		posSide, _ := f.GetString("pos-side")
		posID, _ := f.GetString("pos-id")
		tdMode, _ := f.GetString("td-mode")
		mrgPos, _ := f.GetString("mrg-position")
		tpPx, _ := f.GetString("tp-trigger-px")
		tpType, _ := f.GetString("tp-trigger-px-type")
		tpOrd, _ := f.GetString("tp-ord-px")
		slPx, _ := f.GetString("sl-trigger-px")
		slType, _ := f.GetString("sl-trigger-px-type")
		slOrd, _ := f.GetString("sl-ord-px")
		sz, _ := f.GetString("sz")

		body := buildBody(
			"instType", instType, "instId", instID, "posSide", posSide,
			"tdMode", tdMode, "posId", posID, "mrgPosition", mrgPos,
			"tpTriggerPx", tpPx, "tpTriggerPxType", tpType, "tpOrdPx", tpOrd,
			"slTriggerPx", slPx, "slTriggerPxType", slType, "slOrdPx", slOrd, "sz", sz,
		)
		c := client.New()
		resp, err := c.Post("/deepcoin/trade/set-position-sltp", body)
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	f := setPositionSltpCmd.Flags()
	f.String("inst-type", "", "Instrument type (required)")
	f.String("inst-id", "", "Instrument ID (required)")
	f.String("pos-side", "", "Position side: long/short (required)")
	f.String("pos-id", "", "Position ID (split mode)")
	f.String("td-mode", "isolated", "Trade mode")
	f.String("mrg-position", "", "Position mode")
	f.String("tp-trigger-px", "", "TP trigger price")
	f.String("tp-trigger-px-type", "", "TP trigger price type")
	f.String("tp-ord-px", "", "TP order price (-1 for market)")
	f.String("sl-trigger-px", "", "SL trigger price")
	f.String("sl-trigger-px-type", "", "SL trigger price type")
	f.String("sl-ord-px", "", "SL order price (-1 for market)")
	f.String("sz", "", "Size")
	setPositionSltpCmd.MarkFlagRequired("inst-type")
	setPositionSltpCmd.MarkFlagRequired("inst-id")
	setPositionSltpCmd.MarkFlagRequired("pos-side")
}

// ── modify-position-sltp ────────────────────────────────────────────

var modifyPositionSltpCmd = &cobra.Command{
	Use:   "modify-position-sltp",
	Short: "Modify TP/SL on a position",
	RunE: func(cmd *cobra.Command, args []string) error {
		ordID, _ := cmd.Flags().GetString("ord-id")
		instID, _ := cmd.Flags().GetString("inst-id")
		tp, _ := cmd.Flags().GetString("tp-trigger-px")
		sl, _ := cmd.Flags().GetString("sl-trigger-px")

		body := buildBody("ordId", ordID, "instId", instID, "tpTriggerPx", tp, "slTriggerPx", sl)
		c := client.New()
		resp, err := c.Post("/deepcoin/trade/modify-position-sltp", body)
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	modifyPositionSltpCmd.Flags().String("ord-id", "", "SLTP order ID (required)")
	modifyPositionSltpCmd.Flags().String("inst-id", "", "Instrument ID (required)")
	modifyPositionSltpCmd.Flags().String("tp-trigger-px", "", "New TP trigger price")
	modifyPositionSltpCmd.Flags().String("sl-trigger-px", "", "New SL trigger price")
	modifyPositionSltpCmd.MarkFlagRequired("ord-id")
	modifyPositionSltpCmd.MarkFlagRequired("inst-id")
}

// ── cancel-position-sltp ────────────────────────────────────────────

var cancelPositionSltpCmd = &cobra.Command{
	Use:   "cancel-position-sltp",
	Short: "Cancel TP/SL on a position",
	RunE: func(cmd *cobra.Command, args []string) error {
		ordID, _ := cmd.Flags().GetString("ord-id")

		c := client.New()
		resp, err := c.Post("/deepcoin/trade/cancel-position-sltp", map[string]any{"ordId": ordID})
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	cancelPositionSltpCmd.Flags().String("ord-id", "", "SLTP order ID (required)")
	cancelPositionSltpCmd.MarkFlagRequired("ord-id")
}

// ── close-position ──────────────────────────────────────────────────

var closePositionCmd = &cobra.Command{
	Use:   "close-position",
	Short: "Close positions by IDs",
	RunE: func(cmd *cobra.Command, args []string) error {
		instID, _ := cmd.Flags().GetString("inst-id")
		pg, _ := cmd.Flags().GetString("product-group")
		idsStr, _ := cmd.Flags().GetString("position-ids")

		ids := strings.Split(idsStr, ",")
		for i := range ids {
			ids[i] = strings.TrimSpace(ids[i])
		}
		posIDs := make([]any, len(ids))
		for i, id := range ids {
			posIDs[i] = id
		}
		c := client.New()
		resp, err := c.Post("/deepcoin/trade/close-position-by-ids", map[string]any{
			"instId": instID, "productGroup": pg, "positionIds": posIDs,
		})
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	closePositionCmd.Flags().String("inst-id", "", "Instrument ID (required)")
	closePositionCmd.Flags().String("product-group", "", "Product group (required)")
	closePositionCmd.Flags().String("position-ids", "", "Comma-separated position IDs (required)")
	closePositionCmd.MarkFlagRequired("inst-id")
	closePositionCmd.MarkFlagRequired("product-group")
	closePositionCmd.MarkFlagRequired("position-ids")
}

// ── batch-close-position ────────────────────────────────────────────

var batchClosePositionCmd = &cobra.Command{
	Use:   "batch-close-position",
	Short: "Close all positions for an instrument",
	RunE: func(cmd *cobra.Command, args []string) error {
		instID, _ := cmd.Flags().GetString("inst-id")
		pg, _ := cmd.Flags().GetString("product-group")

		c := client.New()
		resp, err := c.Post("/deepcoin/trade/batch-close-position", map[string]any{
			"instId": instID, "productGroup": pg,
		})
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	batchClosePositionCmd.Flags().String("inst-id", "", "Instrument ID (required)")
	batchClosePositionCmd.Flags().String("product-group", "", "Product group (required)")
	batchClosePositionCmd.MarkFlagRequired("inst-id")
	batchClosePositionCmd.MarkFlagRequired("product-group")
}

// ── trace-order ─────────────────────────────────────────────────────

var traceOrderCmd = &cobra.Command{
	Use:   "trace-order",
	Short: "Place a trace (trailing) order",
	RunE: func(cmd *cobra.Command, args []string) error {
		instID, _ := cmd.Flags().GetString("inst-id")
		retrace, _ := cmd.Flags().GetString("retrace-point")
		trigPx, _ := cmd.Flags().GetString("trigger-price")
		posSide, _ := cmd.Flags().GetString("pos-side")

		c := client.New()
		resp, err := c.Post("/deepcoin/trade/trace-order", map[string]any{
			"instId": instID, "retracePoint": retrace,
			"triggerPrice": trigPx, "posSide": posSide,
		})
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	traceOrderCmd.Flags().String("inst-id", "", "Instrument ID (required)")
	traceOrderCmd.Flags().String("retrace-point", "", "Retrace point (required)")
	traceOrderCmd.Flags().String("trigger-price", "", "Trigger price (required)")
	traceOrderCmd.Flags().String("pos-side", "", "Position side: long/short (required)")
	traceOrderCmd.MarkFlagRequired("inst-id")
	traceOrderCmd.MarkFlagRequired("retrace-point")
	traceOrderCmd.MarkFlagRequired("trigger-price")
	traceOrderCmd.MarkFlagRequired("pos-side")
}

// ── trace-orders ────────────────────────────────────────────────────

var traceOrdersCmd = &cobra.Command{
	Use:   "trace-orders",
	Short: "List pending trace orders",
	RunE: func(cmd *cobra.Command, args []string) error {
		asJSON, _ := cmd.Flags().GetBool("json")

		c := client.New()
		resp, err := c.Get("/deepcoin/trade/trace-order-list", nil)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			rows := client.GetDataSlice(resp)
			output.Table(rows, nil, nil)
		}
		return nil
	},
}

func init() {
	traceOrdersCmd.Flags().Bool("json", false, "Output raw JSON")
}
