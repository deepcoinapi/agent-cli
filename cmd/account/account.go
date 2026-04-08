// Package account provides account, portfolio, asset, and sub-account commands.
package account

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/deepcoinapi/agent-cli/pkg/client"
	"github.com/deepcoinapi/agent-cli/pkg/output"
)

// Cmd is the account command group.
var Cmd = &cobra.Command{
	Use:   "account",
	Short: "Account — balances, positions, leverage, sub-accounts, assets, transfers",
}

func init() {
	Cmd.AddCommand(balanceCmd)
	Cmd.AddCommand(positionsCmd)
	Cmd.AddCommand(billsCmd)
	Cmd.AddCommand(setLeverageCmd)
	Cmd.AddCommand(uidCmd)
	Cmd.AddCommand(subAccountsCmd)
	Cmd.AddCommand(subAccountBalanceCmd)
	Cmd.AddCommand(subAccountTransferCmd)
	Cmd.AddCommand(subAccountTransferRecordsCmd)
	Cmd.AddCommand(depositListCmd)
	Cmd.AddCommand(withdrawListCmd)
	Cmd.AddCommand(transferCmd)
	Cmd.AddCommand(rechargeChainsCmd)
	Cmd.AddCommand(internalTransferSupportCmd)
	Cmd.AddCommand(internalTransferCmd)
	Cmd.AddCommand(internalTransferHistoryCmd)
	Cmd.AddCommand(rebateSummaryCmd)
	Cmd.AddCommand(affiliatesCmd)
	Cmd.AddCommand(tradeStatsDailyCmd)
	Cmd.AddCommand(tradeStatsTotalCmd)
}

// helper
func buildBody(pairs ...string) map[string]any {
	body := make(map[string]any)
	for i := 0; i+1 < len(pairs); i += 2 {
		if pairs[i+1] != "" {
			body[pairs[i]] = pairs[i+1]
		}
	}
	return body
}

// ═══════════════════════════════════════════════════════════════════
// Account
// ═══════════════════════════════════════════════════════════════════

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Get account balance",
	RunE: func(cmd *cobra.Command, args []string) error {
		instType, _ := cmd.Flags().GetString("inst-type")
		ccy, _ := cmd.Flags().GetString("ccy")
		asJSON, _ := cmd.Flags().GetBool("json")

		params := map[string]string{}
		if instType != "" {
			params["instType"] = instType
		}
		if ccy != "" {
			params["ccy"] = ccy
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/account/balances", params)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			rows := client.GetDataSlice(resp)
			output.Table(rows,
				[]string{"ccy", "bal", "availBal", "frozenBal", "unrealizedProfit", "equity"},
				[]string{"Currency", "Balance", "Available", "Frozen", "Unrealized P&L", "Equity"},
			)
		}
		return nil
	},
}

func init() {
	balanceCmd.Flags().String("inst-type", "", "Instrument type: SPOT/SWAP")
	balanceCmd.Flags().String("ccy", "", "Currency filter (e.g. USDT)")
	balanceCmd.Flags().Bool("json", false, "Output raw JSON")
}

var positionsCmd = &cobra.Command{
	Use:   "positions",
	Short: "Get open positions",
	RunE: func(cmd *cobra.Command, args []string) error {
		instType, _ := cmd.Flags().GetString("inst-type")
		instID, _ := cmd.Flags().GetString("inst-id")
		asJSON, _ := cmd.Flags().GetBool("json")

		params := map[string]string{}
		if instType != "" {
			params["instType"] = instType
		}
		if instID != "" {
			params["instId"] = instID
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/account/positions", params)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			rows := client.GetDataSlice(resp)
			output.Table(rows,
				[]string{"posId", "instId", "posSide", "pos", "avgPx", "lastPx", "lever", "unrealizedProfit", "liqPx", "mgnMode"},
				[]string{"Pos ID", "Instrument", "Side", "Size", "Avg Price", "Last Price", "Lever", "Unrealized P&L", "Liq Price", "Margin"},
			)
		}
		return nil
	},
}

func init() {
	positionsCmd.Flags().String("inst-type", "", "Instrument type: SPOT/SWAP")
	positionsCmd.Flags().String("inst-id", "", "Filter by instrument ID")
	positionsCmd.Flags().Bool("json", false, "Output raw JSON")
}

var billsCmd = &cobra.Command{
	Use:   "bills",
	Short: "Get account bill history",
	RunE: func(cmd *cobra.Command, args []string) error {
		instType, _ := cmd.Flags().GetString("inst-type")
		ccy, _ := cmd.Flags().GetString("ccy")
		billType, _ := cmd.Flags().GetString("type")
		limit, _ := cmd.Flags().GetString("limit")
		asJSON, _ := cmd.Flags().GetBool("json")

		params := map[string]string{"instType": instType, "limit": limit}
		if ccy != "" {
			params["ccy"] = ccy
		}
		if billType != "" {
			params["type"] = billType
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/account/bills", params)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			rows := client.GetDataSlice(resp)
			output.Table(rows,
				[]string{"billId", "ccy", "balChg", "bal", "type", "ts"},
				[]string{"Bill ID", "Currency", "Change", "Balance", "Type", "Time"},
			)
		}
		return nil
	},
}

func init() {
	billsCmd.Flags().String("inst-type", "", "Instrument type: SPOT/SWAP (required)")
	billsCmd.Flags().String("ccy", "", "Currency filter")
	billsCmd.Flags().String("type", "", "Bill type: 2=income, 3=expense, 4=transfer, 5=fee")
	billsCmd.Flags().String("limit", "20", "Max results (max 100)")
	billsCmd.Flags().Bool("json", false, "Output raw JSON")
	billsCmd.MarkFlagRequired("inst-type")
}

var setLeverageCmd = &cobra.Command{
	Use:   "set-leverage",
	Short: "Set leverage for an instrument",
	RunE: func(cmd *cobra.Command, args []string) error {
		instID, _ := cmd.Flags().GetString("inst-id")
		lever, _ := cmd.Flags().GetString("lever")
		mgnMode, _ := cmd.Flags().GetString("mgn-mode")
		mrgPos, _ := cmd.Flags().GetString("mrg-position")

		body := buildBody("instId", instID, "lever", lever, "mgnMode", mgnMode, "mrgPosition", mrgPos)
		c := client.New()
		resp, err := c.Post("/deepcoin/account/set-leverage", body)
		if err != nil {
			return err
		}
		fmt.Printf("Leverage set: %v\n", client.GetData(resp))
		return nil
	},
}

func init() {
	setLeverageCmd.Flags().String("inst-id", "", "Instrument ID (required)")
	setLeverageCmd.Flags().String("lever", "", "Leverage 0.01-125 (required)")
	setLeverageCmd.Flags().String("mgn-mode", "", "Margin mode: cross/isolated (required)")
	setLeverageCmd.Flags().String("mrg-position", "", "Position mode: merge/split")
	setLeverageCmd.MarkFlagRequired("inst-id")
	setLeverageCmd.MarkFlagRequired("lever")
	setLeverageCmd.MarkFlagRequired("mgn-mode")
}

var uidCmd = &cobra.Command{
	Use:   "uid",
	Short: "Get current account UID",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New()
		resp, err := c.Get("/deepcoin/account/uid", nil)
		if err != nil {
			return err
		}
		data := client.GetDataMap(resp)
		if data != nil {
			fmt.Printf("UID: %v\n", data["uid"])
		} else {
			output.JSON(client.GetData(resp))
		}
		return nil
	},
}

// ═══════════════════════════════════════════════════════════════════
// Sub-Accounts
// ═══════════════════════════════════════════════════════════════════

var subAccountsCmd = &cobra.Command{
	Use:   "sub-accounts",
	Short: "List sub-accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		asJSON, _ := cmd.Flags().GetBool("json")
		c := client.New()
		resp, err := c.Get("/deepcoin/sub-account/sub-account-list", nil)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			rows := client.GetDataSlice(resp)
			output.Table(rows, []string{"subUid", "subNickname", "subAccount"}, []string{"Sub UID", "Nickname", "Account"})
		}
		return nil
	},
}

func init() {
	subAccountsCmd.Flags().Bool("json", false, "Output raw JSON")
}

var subAccountBalanceCmd = &cobra.Command{
	Use:   "sub-account-balance",
	Short: "Get total balance across all sub-accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New()
		resp, err := c.Get("/deepcoin/sub-account/sub-account-balance-total", nil)
		if err != nil {
			return err
		}
		output.JSON(client.GetData(resp))
		return nil
	},
}

var subAccountTransferCmd = &cobra.Command{
	Use:   "sub-account-transfer",
	Short: "Transfer between sub-accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		fromUID, _ := cmd.Flags().GetString("from-uid")
		toUID, _ := cmd.Flags().GetString("to-uid")
		fromID, _ := cmd.Flags().GetString("from-id")
		toID, _ := cmd.Flags().GetString("to-id")
		amount, _ := cmd.Flags().GetString("amount")
		coin, _ := cmd.Flags().GetString("coin")

		c := client.New()
		resp, err := c.Post("/deepcoin/sub-account/sub-account-transfer", map[string]any{
			"fromUid": fromUID, "toUid": toUID, "fromId": fromID,
			"toId": toID, "amount": amount, "coin": coin,
		})
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	f := subAccountTransferCmd.Flags()
	f.String("from-uid", "", "Source UID (required)")
	f.String("to-uid", "", "Destination UID (required)")
	f.String("from-id", "", "Source account type: 1/2/5/7 (required)")
	f.String("to-id", "", "Destination account type: 1/2/5/7 (required)")
	f.String("amount", "", "Transfer amount (required)")
	f.String("coin", "", "Currency e.g. USDT (required)")
	subAccountTransferCmd.MarkFlagRequired("from-uid")
	subAccountTransferCmd.MarkFlagRequired("to-uid")
	subAccountTransferCmd.MarkFlagRequired("from-id")
	subAccountTransferCmd.MarkFlagRequired("to-id")
	subAccountTransferCmd.MarkFlagRequired("amount")
	subAccountTransferCmd.MarkFlagRequired("coin")
}

var subAccountTransferRecordsCmd = &cobra.Command{
	Use:   "sub-account-transfer-records",
	Short: "Get sub-account transfer records",
	RunE: func(cmd *cobra.Command, args []string) error {
		coin, _ := cmd.Flags().GetString("coin")
		page, _ := cmd.Flags().GetString("page")
		size, _ := cmd.Flags().GetString("size")
		asJSON, _ := cmd.Flags().GetBool("json")

		params := map[string]string{"page": page, "size": size}
		if coin != "" {
			params["coin"] = coin
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/sub-account/sub-account-transfer-record", params)
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
	subAccountTransferRecordsCmd.Flags().String("coin", "", "Currency filter")
	subAccountTransferRecordsCmd.Flags().String("page", "1", "Page number")
	subAccountTransferRecordsCmd.Flags().String("size", "20", "Page size (max 100)")
	subAccountTransferRecordsCmd.Flags().Bool("json", false, "Output raw JSON")
}

// ═══════════════════════════════════════════════════════════════════
// Assets
// ═══════════════════════════════════════════════════════════════════

var depositListCmd = &cobra.Command{
	Use:   "deposit-list",
	Short: "Get deposit history",
	RunE: func(cmd *cobra.Command, args []string) error {
		coin, _ := cmd.Flags().GetString("coin")
		page, _ := cmd.Flags().GetString("page")
		size, _ := cmd.Flags().GetString("size")
		asJSON, _ := cmd.Flags().GetBool("json")

		params := map[string]string{"page": page, "size": size}
		if coin != "" {
			params["coin"] = coin
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/asset/deposit-list", params)
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
	depositListCmd.Flags().String("coin", "", "Currency filter")
	depositListCmd.Flags().String("page", "1", "Page number")
	depositListCmd.Flags().String("size", "20", "Page size")
	depositListCmd.Flags().Bool("json", false, "Output raw JSON")
}

var withdrawListCmd = &cobra.Command{
	Use:   "withdraw-list",
	Short: "Get withdrawal history",
	RunE: func(cmd *cobra.Command, args []string) error {
		coin, _ := cmd.Flags().GetString("coin")
		page, _ := cmd.Flags().GetString("page")
		size, _ := cmd.Flags().GetString("size")
		asJSON, _ := cmd.Flags().GetBool("json")

		params := map[string]string{"page": page, "size": size}
		if coin != "" {
			params["coin"] = coin
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/asset/withdraw-list", params)
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
	withdrawListCmd.Flags().String("coin", "", "Currency filter")
	withdrawListCmd.Flags().String("page", "1", "Page number")
	withdrawListCmd.Flags().String("size", "20", "Page size")
	withdrawListCmd.Flags().Bool("json", false, "Output raw JSON")
}

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer assets between accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		curID, _ := cmd.Flags().GetString("currency-id")
		amount, _ := cmd.Flags().GetString("amount")
		fromID, _ := cmd.Flags().GetString("from-id")
		toID, _ := cmd.Flags().GetString("to-id")
		uid, _ := cmd.Flags().GetString("uid")

		body := buildBody("currency_id", curID, "amount", amount, "from_id", fromID, "to_id", toID, "uid", uid)
		c := client.New()
		resp, err := c.Post("/deepcoin/asset/transfer", body)
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	transferCmd.Flags().String("currency-id", "", "Currency ID (required)")
	transferCmd.Flags().String("amount", "", "Amount (required)")
	transferCmd.Flags().String("from-id", "", "From account: 1/2/3/5/7/10 (required)")
	transferCmd.Flags().String("to-id", "", "To account: 1/2/3/5/7/10 (required)")
	transferCmd.Flags().String("uid", "", "Target UID (optional)")
	transferCmd.MarkFlagRequired("currency-id")
	transferCmd.MarkFlagRequired("amount")
	transferCmd.MarkFlagRequired("from-id")
	transferCmd.MarkFlagRequired("to-id")
}

var rechargeChainsCmd = &cobra.Command{
	Use:   "recharge-chains",
	Short: "Get supported deposit chains for a currency",
	RunE: func(cmd *cobra.Command, args []string) error {
		curID, _ := cmd.Flags().GetString("currency-id")
		asJSON, _ := cmd.Flags().GetBool("json")

		c := client.New()
		resp, err := c.Get("/deepcoin/asset/recharge-chain-list", map[string]string{"currency_id": curID})
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
	rechargeChainsCmd.Flags().String("currency-id", "", "Currency ID (required)")
	rechargeChainsCmd.Flags().Bool("json", false, "Output raw JSON")
	rechargeChainsCmd.MarkFlagRequired("currency-id")
}

// ═══════════════════════════════════════════════════════════════════
// Internal Transfer
// ═══════════════════════════════════════════════════════════════════

var internalTransferSupportCmd = &cobra.Command{
	Use:   "internal-transfer-support",
	Short: "Get supported coins for internal transfer",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New()
		resp, err := c.Get("/deepcoin/internal-transfer/support", nil)
		if err != nil {
			return err
		}
		output.JSON(client.GetData(resp))
		return nil
	},
}

var internalTransferCmd = &cobra.Command{
	Use:   "internal-transfer",
	Short: "Make an internal transfer",
	RunE: func(cmd *cobra.Command, args []string) error {
		amount, _ := cmd.Flags().GetString("amount")
		coin, _ := cmd.Flags().GetString("coin")
		recv, _ := cmd.Flags().GetString("receiver-account")
		accType, _ := cmd.Flags().GetString("account-type")
		recvUID, _ := cmd.Flags().GetString("receiver-uid")

		body := buildBody("amount", amount, "coin", coin, "receiverAccount", recv,
			"accountType", accType, "receiverUID", recvUID)
		c := client.New()
		resp, err := c.Post("/deepcoin/internal-transfer", body)
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	f := internalTransferCmd.Flags()
	f.String("amount", "", "Transfer amount (required)")
	f.String("coin", "", "Currency (required)")
	f.String("receiver-account", "", "Receiver account (required)")
	f.String("account-type", "", "Account type (required)")
	f.String("receiver-uid", "", "Receiver UID")
	internalTransferCmd.MarkFlagRequired("amount")
	internalTransferCmd.MarkFlagRequired("coin")
	internalTransferCmd.MarkFlagRequired("receiver-account")
	internalTransferCmd.MarkFlagRequired("account-type")
}

var internalTransferHistoryCmd = &cobra.Command{
	Use:   "internal-transfer-history",
	Short: "Get internal transfer history",
	RunE: func(cmd *cobra.Command, args []string) error {
		coin, _ := cmd.Flags().GetString("coin")
		status, _ := cmd.Flags().GetString("status")
		page, _ := cmd.Flags().GetString("page")
		size, _ := cmd.Flags().GetString("size")
		asJSON, _ := cmd.Flags().GetBool("json")

		params := map[string]string{"page": page, "size": size}
		if coin != "" {
			params["coin"] = coin
		}
		if status != "" {
			params["status"] = status
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/internal-transfer/history-order", params)
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
	internalTransferHistoryCmd.Flags().String("coin", "", "Currency filter")
	internalTransferHistoryCmd.Flags().String("status", "", "Status: 1/2/3")
	internalTransferHistoryCmd.Flags().String("page", "1", "Page number")
	internalTransferHistoryCmd.Flags().String("size", "20", "Page size")
	internalTransferHistoryCmd.Flags().Bool("json", false, "Output raw JSON")
}

// ═══════════════════════════════════════════════════════════════════
// Rebate / Affiliate
// ═══════════════════════════════════════════════════════════════════

var rebateSummaryCmd = &cobra.Command{
	Use:   "rebate-summary",
	Short: "Get rebate summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		uid, _ := cmd.Flags().GetString("uid")
		rType, _ := cmd.Flags().GetString("type")
		st, _ := cmd.Flags().GetString("start-time")
		et, _ := cmd.Flags().GetString("end-time")

		params := map[string]string{"uid": uid, "type": rType}
		if st != "" {
			params["startTime"] = st
		}
		if et != "" {
			params["endTime"] = et
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/agents/users/rebates", params)
		if err != nil {
			return err
		}
		output.JSON(client.GetData(resp))
		return nil
	},
}

func init() {
	rebateSummaryCmd.Flags().String("uid", "", "UID (required)")
	rebateSummaryCmd.Flags().String("type", "0", "Type: 0=all, 1=spot, 2=swap")
	rebateSummaryCmd.Flags().String("start-time", "", "Start time")
	rebateSummaryCmd.Flags().String("end-time", "", "End time")
	rebateSummaryCmd.MarkFlagRequired("uid")
}

var affiliatesCmd = &cobra.Command{
	Use:   "affiliates",
	Short: "Get affiliate list",
	RunE: func(cmd *cobra.Command, args []string) error {
		uid, _ := cmd.Flags().GetString("uid")
		st, _ := cmd.Flags().GetString("start-time")
		et, _ := cmd.Flags().GetString("end-time")

		params := map[string]string{"uid": uid}
		if st != "" {
			params["startTime"] = st
		}
		if et != "" {
			params["endTime"] = et
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/agents/users", params)
		if err != nil {
			return err
		}
		output.JSON(client.GetData(resp))
		return nil
	},
}

func init() {
	affiliatesCmd.Flags().String("uid", "", "UID (required)")
	affiliatesCmd.Flags().String("start-time", "", "Start time")
	affiliatesCmd.Flags().String("end-time", "", "End time")
	affiliatesCmd.MarkFlagRequired("uid")
}

// ═══════════════════════════════════════════════════════════════════
// Trade Statistics
// ═══════════════════════════════════════════════════════════════════

var tradeStatsDailyCmd = &cobra.Command{
	Use:   "trade-stats-daily",
	Short: "Get daily trade statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		appid, _ := cmd.Flags().GetString("appid")
		uid, _ := cmd.Flags().GetString("uid")
		st, _ := cmd.Flags().GetString("start-time")
		et, _ := cmd.Flags().GetString("end-time")
		instIDs, _ := cmd.Flags().GetString("instrument-ids")

		params := map[string]string{"appid": appid, "uid": uid}
		if st != "" {
			params["startTime"] = st
		}
		if et != "" {
			params["endTime"] = et
		}
		if instIDs != "" {
			params["instrumentIds"] = instIDs
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/apiUserTradeStats/daily", params)
		if err != nil {
			return err
		}
		output.JSON(client.GetData(resp))
		return nil
	},
}

func init() {
	tradeStatsDailyCmd.Flags().String("appid", "", "App ID (required)")
	tradeStatsDailyCmd.Flags().String("uid", "", "UID (required)")
	tradeStatsDailyCmd.Flags().String("start-time", "", "Start time")
	tradeStatsDailyCmd.Flags().String("end-time", "", "End time")
	tradeStatsDailyCmd.Flags().String("instrument-ids", "", "Instrument IDs (comma-separated)")
	tradeStatsDailyCmd.MarkFlagRequired("appid")
	tradeStatsDailyCmd.MarkFlagRequired("uid")
}

var tradeStatsTotalCmd = &cobra.Command{
	Use:   "trade-stats-total",
	Short: "Get total trade statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		appid, _ := cmd.Flags().GetString("appid")
		uid, _ := cmd.Flags().GetString("uid")
		st, _ := cmd.Flags().GetString("start-time")
		et, _ := cmd.Flags().GetString("end-time")

		params := map[string]string{"appid": appid, "uid": uid}
		if st != "" {
			params["startTime"] = st
		}
		if et != "" {
			params["endTime"] = et
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/apiUserTradeStats/total", params)
		if err != nil {
			return err
		}
		output.JSON(client.GetData(resp))
		return nil
	},
}

func init() {
	tradeStatsTotalCmd.Flags().String("appid", "", "App ID (required)")
	tradeStatsTotalCmd.Flags().String("uid", "", "UID (required)")
	tradeStatsTotalCmd.Flags().String("start-time", "", "Start time")
	tradeStatsTotalCmd.Flags().String("end-time", "", "End time")
	tradeStatsTotalCmd.MarkFlagRequired("appid")
	tradeStatsTotalCmd.MarkFlagRequired("uid")
}
