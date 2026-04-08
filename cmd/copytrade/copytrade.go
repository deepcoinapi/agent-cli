// Package copytrade provides copy trading commands.
package copytrade

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/deepcoinapi/agent-cli/pkg/client"
	"github.com/deepcoinapi/agent-cli/pkg/output"
)

// Cmd is the copytrade command group.
var Cmd = &cobra.Command{
	Use:   "copytrade",
	Short: "Copy trading — leader/follower management, positions, profit tracking",
}

func init() {
	Cmd.AddCommand(leaderSettingsCmd)
	Cmd.AddCommand(supportContractsCmd)
	Cmd.AddCommand(setContractsCmd)
	Cmd.AddCommand(followersCmd)
	Cmd.AddCommand(leaderPositionsCmd)
	Cmd.AddCommand(positionTypeCmd)
	Cmd.AddCommand(setPositionTypeCmd)
	Cmd.AddCommand(estimatedProfitCmd)
	Cmd.AddCommand(historyProfitCmd)
}

var leaderSettingsCmd = &cobra.Command{
	Use:   "leader-settings",
	Short: "Update leader settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		homeMode, _ := cmd.Flags().GetString("home-mode")
		closedCode, _ := cmd.Flags().GetString("is-closed-copy-code")
		copyCode, _ := cmd.Flags().GetString("copy-code")

		body := map[string]any{"status": status}
		if homeMode != "" {
			body["homeMode"] = homeMode
		}
		if closedCode != "" {
			body["isClosedCopyCode"] = closedCode
		}
		if copyCode != "" {
			body["copyCode"] = copyCode
		}
		c := client.New()
		resp, err := c.Post("/deepcoin/copytrading/leader-settings", body)
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	leaderSettingsCmd.Flags().String("status", "", "0=disable, 1=enable (required)")
	leaderSettingsCmd.Flags().String("home-mode", "", "Home page mode: 1/3")
	leaderSettingsCmd.Flags().String("is-closed-copy-code", "", "Close copy code")
	leaderSettingsCmd.Flags().String("copy-code", "", "Copy code")
	leaderSettingsCmd.MarkFlagRequired("status")
}

var supportContractsCmd = &cobra.Command{
	Use:   "support-contracts",
	Short: "Get supported copy trading contracts",
	RunE: func(cmd *cobra.Command, args []string) error {
		asJSON, _ := cmd.Flags().GetBool("json")
		c := client.New()
		resp, err := c.Get("/deepcoin/copytrading/support-contracts", nil)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			data := client.GetDataMap(resp)
			if list, ok := data["List"]; ok {
				if arr, ok := list.([]any); ok {
					for _, v := range arr {
						fmt.Println(v)
					}
					return nil
				}
			}
			output.JSON(client.GetData(resp))
		}
		return nil
	},
}

func init() {
	supportContractsCmd.Flags().Bool("json", false, "Output raw JSON")
}

var setContractsCmd = &cobra.Command{
	Use:   "set-contracts",
	Short: "Set copy trading contracts",
	RunE: func(cmd *cobra.Command, args []string) error {
		contractsStr, _ := cmd.Flags().GetString("contracts")
		parts := strings.Split(contractsStr, ",")
		contracts := make([]any, len(parts))
		for i, p := range parts {
			contracts[i] = strings.TrimSpace(p)
		}
		c := client.New()
		resp, err := c.Post("/deepcoin/copytrading/set-contracts", map[string]any{"contracts": contracts})
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	setContractsCmd.Flags().String("contracts", "", "Comma-separated contract symbols (required)")
	setContractsCmd.MarkFlagRequired("contracts")
}

var followersCmd = &cobra.Command{
	Use:   "followers",
	Short: "Get follower list and stats",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		asJSON, _ := cmd.Flags().GetBool("json")

		c := client.New()
		resp, err := c.Get("/deepcoin/copytrading/follower-rank", map[string]string{"status": status})
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			data := client.GetDataMap(resp)
			fmt.Printf("Followers: %v / %v\n", data["followerNum"], data["maxFollowerNum"])
			if list, ok := data["list"].([]any); ok {
				var rows []map[string]any
				for _, item := range list {
					if m, ok := item.(map[string]any); ok {
						rows = append(rows, m)
					}
				}
				output.Table(rows,
					[]string{"userId", "nickName", "totalProfit"},
					[]string{"User ID", "Nickname", "Total Profit"},
				)
			}
		}
		return nil
	},
}

func init() {
	followersCmd.Flags().String("status", "", "1=active, 2=inactive (required)")
	followersCmd.Flags().Bool("json", false, "Output raw JSON")
	followersCmd.MarkFlagRequired("status")
}

var leaderPositionsCmd = &cobra.Command{
	Use:   "leader-positions",
	Short: "Get leader's current positions",
	RunE: func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetString("page")
		size, _ := cmd.Flags().GetString("size")
		asJSON, _ := cmd.Flags().GetBool("json")

		c := client.New()
		resp, err := c.Get("/deepcoin/copytrading/leader-position", map[string]string{
			"pageNum": page, "pageSize": size,
		})
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			rows := client.GetDataSlice(resp)
			output.Table(rows,
				[]string{"instrumentId", "positionDirection", "position", "openPrice", "leverage", "isCrossMargin"},
				[]string{"Instrument", "Direction", "Size", "Open Price", "Leverage", "Cross"},
			)
		}
		return nil
	},
}

func init() {
	leaderPositionsCmd.Flags().String("page", "1", "Page number")
	leaderPositionsCmd.Flags().String("size", "20", "Page size")
	leaderPositionsCmd.Flags().Bool("json", false, "Output raw JSON")
}

var positionTypeCmd = &cobra.Command{
	Use:   "position-type",
	Short: "Get current position type (hedge/one-way)",
	RunE: func(cmd *cobra.Command, args []string) error {
		asJSON, _ := cmd.Flags().GetBool("json")
		c := client.New()
		resp, err := c.Get("/deepcoin/copytrading/position-type", nil)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			data := client.GetDataMap(resp)
			pt := fmt.Sprintf("%v", data["positionType"])
			label := pt
			if pt == "1" {
				label = "Hedge"
			} else if pt == "2" {
				label = "One-way"
			}
			fmt.Printf("Position type: %s (%s)\n", label, pt)
		}
		return nil
	},
}

func init() {
	positionTypeCmd.Flags().Bool("json", false, "Output raw JSON")
}

var setPositionTypeCmd = &cobra.Command{
	Use:   "set-position-type",
	Short: "Update position type",
	RunE: func(cmd *cobra.Command, args []string) error {
		posType, _ := cmd.Flags().GetString("type")
		c := client.New()
		resp, err := c.Post("/deepcoin/copytrading/position-type", map[string]any{"positionType": posType})
		if err != nil {
			return err
		}
		output.JSON(resp)
		return nil
	},
}

func init() {
	setPositionTypeCmd.Flags().String("type", "", "1=Hedge, 2=One-way (required)")
	setPositionTypeCmd.MarkFlagRequired("type")
}

var estimatedProfitCmd = &cobra.Command{
	Use:   "estimated-profit",
	Short: "Get estimated profit from followers",
	RunE: func(cmd *cobra.Command, args []string) error {
		asJSON, _ := cmd.Flags().GetBool("json")
		c := client.New()
		resp, err := c.Get("/deepcoin/copytrading/estimate-profit", nil)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			data := client.GetDataMap(resp)
			if list, ok := data["list"].([]any); ok {
				var rows []map[string]any
				for _, item := range list {
					if m, ok := item.(map[string]any); ok {
						rows = append(rows, m)
					}
				}
				output.Table(rows, []string{"userID", "nickName", "estimateProfit"}, []string{"User ID", "Nickname", "Est. Profit"})
			} else {
				output.JSON(client.GetData(resp))
			}
		}
		return nil
	},
}

func init() {
	estimatedProfitCmd.Flags().Bool("json", false, "Output raw JSON")
}

var historyProfitCmd = &cobra.Command{
	Use:   "history-profit",
	Short: "Get historical profit from copy trading",
	RunE: func(cmd *cobra.Command, args []string) error {
		asJSON, _ := cmd.Flags().GetBool("json")
		c := client.New()
		resp, err := c.Get("/deepcoin/copytrading/history-profit", nil)
		if err != nil {
			return err
		}
		if asJSON {
			output.JSON(resp)
		} else {
			data := client.GetDataMap(resp)
			if list, ok := data["list"].([]any); ok {
				var rows []map[string]any
				for _, item := range list {
					if m, ok := item.(map[string]any); ok {
						rows = append(rows, m)
					}
				}
				output.Table(rows, []string{"settlementTime", "profit"}, []string{"Settlement Time", "Profit"})
			} else {
				output.JSON(client.GetData(resp))
			}
		}
		return nil
	},
}

func init() {
	historyProfitCmd.Flags().Bool("json", false, "Output raw JSON")
}
