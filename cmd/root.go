package cmd

import (
	"os"

	"github.com/Appliscale/tyr/configuration"
	"github.com/Appliscale/tyr/scanner"

	"github.com/spf13/cobra"
)

// var cfgFile string
var config = configuration.GetConfig()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tyr",
	Short: "Scan for vulnerabilities in your AWS Account.",
	Long:  `Scan for vulnerabilities in your AWS Account.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := scanner.Run(&config)
		if err != nil {
			config.Logger.Error(err.Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		config.Logger.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&config.Region, "region", "r", "", "specify aws region to scan your account,e.g. --region us-east-1")
	rootCmd.MarkFlagRequired("region")

	rootCmd.Flags().StringVarP(&config.Service, "service", "s", "", "specify aws service to scan in your account,e.g. --service [ec2:x,ec2:image]")
	rootCmd.MarkFlagRequired("service")

	rootCmd.Flags().StringVarP(&config.Profile, "profile", "p", "", "specify aws profile e.g. --profile appliscale")
}
