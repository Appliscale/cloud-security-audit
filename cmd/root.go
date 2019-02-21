package cmd

import (
	"os"

	"github.com/Appliscale/cloud-security-audit/configuration"
	"github.com/Appliscale/cloud-security-audit/csasession"
	"github.com/Appliscale/cloud-security-audit/environment"
	"github.com/Appliscale/cloud-security-audit/resource"
	"github.com/Appliscale/cloud-security-audit/scanner"
	"github.com/spf13/cobra"
)

// var cfgFile string
var config = configuration.GetConfig()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloud-security-audit",
	Short: "Scan for vulnerabilities in your AWS Account.",
	Long:  `Scan for vulnerabilities in your AWS Account.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, ok := os.LookupEnv("AWS_ACCESS_KEY_ID") // If csa is running on lambda then env will be available. In other case csa needs config files.
		if !ok {
			if environment.CheckAWSConfigFiles(&config) {
				err := scanner.Run(&config)
				if err != nil {
					config.Logger.Error(err.Error())
				}
			}
		} else {
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
	region      string
	service     string
	profile     string
	mfa         bool
	mfaDuration int64
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&region, "region", "r", "", "specify aws region to scan your account,e.g. --region us-east-1")

	rootCmd.Flags().StringVarP(&service, "service", "s", "", "specify aws service to scan in your account,e.g. --service [ec2:x,ec2:image]")

	rootCmd.Flags().StringVarP(&profile, "profile", "p", "", "specify aws profile e.g. --profile appliscale")

	rootCmd.Flags().BoolVarP(&mfa, "mfa", "m", false, "indicates usage of Multi Factor Authentication")
	rootCmd.Flags().Int64VarP(&mfaDuration, "mfa-duration", "d", 0, "sets the duration of the MFA session")
}

func getRegions() *[]string {
	if region != "" {
		return &[]string{region}
	}

	return csasession.GetAvailableRegions()
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
	config.Mfa = mfa
	config.MfaDuration = mfaDuration
	configuration.InitialiseMFA(config)
}
