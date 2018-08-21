package cmd

import (
	"fmt"
	"os"

	"github.com/Appliscale/tyr/configuration"
	"github.com/Appliscale/tyr/scanner"
	"github.com/Appliscale/tyr/tyrlogger"
	"github.com/Appliscale/tyr/tyrsession"
	"github.com/Appliscale/tyr/tyrsession/clientfactory"
	"github.com/Appliscale/tyr/tyrsession/sessionfactory"

	"github.com/spf13/cobra"
)

// var cfgFile string
var config configuration.Config
var logger = tyrlogger.GetInstance()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tyr",
	Short: "Scan for vulnerabilities in your AWS Account.",
	Long:  `Scan for vulnerabilities in your AWS Account.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := scanner.Run(&config)
		if err != nil {
			logger.Fatalln(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	region  string
	service string
	profile string
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&region, "region", "r", "", "specify aws region to scan your account,e.g. --region us-east-1")

	rootCmd.Flags().StringVarP(&service, "service", "s", "", "specify aws service to scan in your account,e.g. --service [ec2:x,ec2:image]")
	rootCmd.MarkFlagRequired("service")

	rootCmd.Flags().StringVarP(&profile, "profile", "p", "", "specify aws profile e.g. --profile appliscale")
}

func initConfig() {
	if region == "" {
		config.Regions = tyrsession.GetAvailableRegions()
	} else {
		config.Regions = &[]string{region}
	}

	config.Service = service

	config.Profile = profile

	config.SessionFactory = sessionfactory.New()

	config.ClientFactory = clientfactory.New(config.SessionFactory)
}
