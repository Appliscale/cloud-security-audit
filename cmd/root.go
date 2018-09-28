package cmd

import (
	"os"

	"github.com/Appliscale/tyr/configuration"
	"github.com/Appliscale/tyr/resource"
	"github.com/Appliscale/tyr/tyrsession"

	"github.com/Appliscale/tyr/environment"
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
		if environment.CheckAWSConfigFiles(&config) {
			err := scanner.Run(&config)
			if err != nil {
				config.Logger.Error(err.Error())
			}
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

var (
	region  string
	service string
	profile string
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&region, "region", "r", "", "specify aws region to scan your account,e.g. --region us-east-1")

	rootCmd.Flags().StringVarP(&service, "service", "s", "", "specify aws service to scan in your account,e.g. --service [ec2:x,ec2:image]")

	rootCmd.Flags().StringVarP(&profile, "profile", "p", "", "specify aws profile e.g. --profile appliscale")
}

func getRegions() *[]string {
	if region != "" {
		return &[]string{region}
	}

	return tyrsession.GetAvailableRegions()
}

func getServices() *[]string {
	if service != "" {
		return &[]string{service}
	}

	return resource.GetAvailableServices()
}

func getProfile() string {
	if profile != "" {
		return profile
	}

	if profile, ok := os.LookupEnv("AWS_PROFILE"); ok {
		return profile
	}

	return "default"
}

func initConfig() {
	config.Regions = getRegions()
	config.Services = getServices()
	config.Profile = getProfile()
}
