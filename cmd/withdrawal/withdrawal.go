// Package withdrawal provides on-chain withdrawal commands.
package withdrawal

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/deepcoinapi/agent-cli/pkg/client"
	"github.com/deepcoinapi/agent-cli/pkg/output"
)

// Cmd is the withdrawal command group.
var Cmd = &cobra.Command{
	Use:   "withdrawal",
	Short: "On-chain withdrawals — config, whitelist, create, cancel, status, records",
}

func init() {
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(cancelCmd)
	Cmd.AddCommand(recordsCmd)
	Cmd.AddCommand(statusCmd)
	Cmd.AddCommand(assetsCmd)
	Cmd.AddCommand(chainsCmd)
	Cmd.AddCommand(addressesCmd)
	Cmd.AddCommand(configCmd)
}

func buildBody(pairs ...string) map[string]any {
	body := make(map[string]any)
	for i := 0; i+1 < len(pairs); i += 2 {
		if pairs[i+1] != "" {
			body[pairs[i]] = pairs[i+1]
		}
	}
	return body
}

func splitCSV(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an on-chain withdrawal",
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Flags()
		ccy, _ := f.GetString("ccy")
		chain, _ := f.GetString("chain")
		amt, _ := f.GetString("amt")
		addressID, _ := f.GetString("address-id")
		toAddr, _ := f.GetString("to-addr")
		memo, _ := f.GetString("memo")
		accountTypes, _ := f.GetString("account-types")
		clientID, _ := f.GetString("client-id")
		remark, _ := f.GetString("remark")

		body := buildBody(
			"ccy", ccy, "chain", chain, "amt", amt, "addressId", addressID,
			"toAddr", toAddr, "memo", memo, "clientId", clientID, "remark", remark,
		)
		if accounts := splitCSV(accountTypes); len(accounts) > 0 {
			body["accountTypes"] = accounts
		}

		c := client.New()
		resp, err := c.Post("/deepcoin/asset/withdrawal", body)
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	f := createCmd.Flags()
	f.String("ccy", "", "Coin, e.g. USDT (required)")
	f.String("chain", "", "Withdrawal chain, e.g. USDT-TRC20 (required)")
	f.String("amt", "", "Withdrawal amount (required)")
	f.String("address-id", "", "Whitelist address ID (required)")
	f.String("to-addr", "", "Optional address consistency check")
	f.String("memo", "", "Memo/tag/payment ID when required")
	f.String("account-types", "", "Comma-separated account types; API accepts at most one: funding, spot, swap")
	f.String("client-id", "", "Client request ID")
	f.String("remark", "", "Remark")
	createCmd.MarkFlagRequired("ccy")
	createCmd.MarkFlagRequired("chain")
	createCmd.MarkFlagRequired("amt")
	createCmd.MarkFlagRequired("address-id")
}

var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel an on-chain withdrawal",
	RunE: func(cmd *cobra.Command, args []string) error {
		wdID, _ := cmd.Flags().GetString("wd-id")
		ccy, _ := cmd.Flags().GetString("ccy")
		clientID, _ := cmd.Flags().GetString("client-id")

		c := client.New()
		resp, err := c.Post("/deepcoin/asset/cancel-withdrawal", buildBody(
			"wdId", wdID, "ccy", ccy, "clientId", clientID,
		))
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	cancelCmd.Flags().String("wd-id", "", "Withdrawal ID (required)")
	cancelCmd.Flags().String("ccy", "", "Coin")
	cancelCmd.Flags().String("client-id", "", "Client request ID")
	cancelCmd.MarkFlagRequired("wd-id")
}

var recordsCmd = &cobra.Command{
	Use:   "records",
	Short: "List withdrawal records",
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Flags()
		params := map[string]string{}
		for _, name := range []string{"coin", "ccy", "chain", "tx-hash", "tx-id", "wd-id", "state", "start-time", "end-time", "page", "size"} {
			value, _ := f.GetString(name)
			if value != "" {
				params[toAPIParam(name)] = value
			}
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/asset/withdraw-list", params)
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func toAPIParam(name string) string {
	switch name {
	case "tx-hash":
		return "txHash"
	case "tx-id":
		return "txId"
	case "wd-id":
		return "wdId"
	case "start-time":
		return "startTime"
	case "end-time":
		return "endTime"
	default:
		return name
	}
}

func init() {
	f := recordsCmd.Flags()
	f.String("coin", "", "Coin filter")
	f.String("ccy", "", "Coin filter alias")
	f.String("chain", "", "Chain filter")
	f.String("tx-hash", "", "Transaction hash filter")
	f.String("tx-id", "", "Transaction ID filter")
	f.String("wd-id", "", "Withdrawal ID filter")
	f.String("state", "", "State filter")
	f.String("start-time", "", "Start time in milliseconds")
	f.String("end-time", "", "End time in milliseconds")
	f.String("page", "1", "Page number")
	f.String("size", "20", "Page size, max 100")
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get a single withdrawal status",
	RunE: func(cmd *cobra.Command, args []string) error {
		wdID, _ := cmd.Flags().GetString("wd-id")
		ccy, _ := cmd.Flags().GetString("ccy")

		params := map[string]string{"wdId": wdID}
		if ccy != "" {
			params["ccy"] = ccy
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/asset/withdrawal-status", params)
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	statusCmd.Flags().String("wd-id", "", "Withdrawal ID (required)")
	statusCmd.Flags().String("ccy", "", "Coin")
	statusCmd.MarkFlagRequired("wd-id")
}

var assetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "List withdrawable assets",
	RunE: func(cmd *cobra.Command, args []string) error {
		ccy, _ := cmd.Flags().GetString("ccy")
		params := map[string]string{}
		if ccy != "" {
			params["ccy"] = ccy
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/asset/withdraw-assets", params)
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	assetsCmd.Flags().String("ccy", "", "Coin filter")
}

var chainsCmd = &cobra.Command{
	Use:   "chains",
	Short: "List withdrawal chains for a coin",
	RunE: func(cmd *cobra.Command, args []string) error {
		ccy, _ := cmd.Flags().GetString("ccy")
		c := client.New()
		resp, err := c.Get("/deepcoin/asset/withdraw-chains", map[string]string{"ccy": ccy})
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	chainsCmd.Flags().String("ccy", "", "Coin, e.g. USDT (required)")
	chainsCmd.MarkFlagRequired("ccy")
}

var addressesCmd = &cobra.Command{
	Use:   "addresses",
	Short: "List withdrawal whitelist addresses for a coin",
	RunE: func(cmd *cobra.Command, args []string) error {
		ccy, _ := cmd.Flags().GetString("ccy")
		c := client.New()
		resp, err := c.Get("/deepcoin/asset/withdraw-addresses", map[string]string{"ccy": ccy})
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	addressesCmd.Flags().String("ccy", "", "Coin, e.g. USDT (required)")
	addressesCmd.MarkFlagRequired("ccy")
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get aggregated withdrawal config",
	RunE: func(cmd *cobra.Command, args []string) error {
		ccy, _ := cmd.Flags().GetString("ccy")
		includeAddresses, _ := cmd.Flags().GetString("include-addresses")
		params := map[string]string{}
		if ccy != "" {
			params["ccy"] = ccy
		}
		if includeAddresses != "" {
			params["includeAddresses"] = includeAddresses
		}
		c := client.New()
		resp, err := c.Get("/deepcoin/asset/withdraw-config", params)
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	configCmd.Flags().String("ccy", "", "Coin filter")
	configCmd.Flags().String("include-addresses", "", "true to include whitelist addresses")
}
