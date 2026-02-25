package main

import (
	"biz-auto-api/cmd/auth"
	"biz-auto-api/cmd/billing"
	"biz-auto-api/cmd/business"
	"biz-auto-api/cmd/cronjob"
	"biz-auto-api/cmd/migrate"
	"biz-auto-api/cmd/network"
	"biz-auto-api/cmd/ops"
	"biz-auto-api/cmd/price"
	"biz-auto-api/cmd/system"
	"biz-auto-api/cmd/vpn"
	"biz-auto-api/pkg/logger"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "biz-auto-api",
	Short: "biz-auto-api",
	//SilenceUsage: true,
	Long: `biz-auto-api`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(migrate.Cmd)
	rootCmd.AddCommand(system.Cmd)
	rootCmd.AddCommand(cronjob.Cmd)
	rootCmd.AddCommand(price.Cmd)
	rootCmd.AddCommand(business.Cmd)
	rootCmd.AddCommand(billing.Cmd)
	rootCmd.AddCommand(network.Cmd)
	rootCmd.AddCommand(vpn.Cmd)
	rootCmd.AddCommand(auth.Cmd)
	rootCmd.AddCommand(ops.Cmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.NewLogger("fatal").Fatal(err)
	}
}
