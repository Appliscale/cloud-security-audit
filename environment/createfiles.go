package environment

import (
	"fmt"
	"github.com/Appliscale/perun/helpers"
	"github.com/Appliscale/tyr/configuration"
	"os"
	"strings"
)

func CreateAWSCredentialsFile(config *configuration.Config, profile string) {
	if profile != "" {
		config.Logger.Always("You haven't got .aws/credentials file for profile " + profile)
		var awsAccessKeyID string
		var awsSecretAccessKey string
		var mfaSerial string

		config.Logger.GetInput("awsAccessKeyID", &awsAccessKeyID)
		config.Logger.GetInput("awsSecretAccessKey", &awsSecretAccessKey)
		config.Logger.GetInput("mfaSerial", &mfaSerial)

		homePath, pathError := GetUserHomeDir()
		if pathError != nil {
			config.Logger.Error(pathError.Error())
		}
		path := homePath + "/.aws/credentials"
		line := "\n[" + profile + "-long-term" + "]\n"
		appendStringToFile(path, line)
		line = "aws_access_key_id" + " = " + awsAccessKeyID + "\n"
		appendStringToFile(path, line)
		line = "aws_secret_access_key" + " = " + awsSecretAccessKey + "\n"
		appendStringToFile(path, line)
		line = "mfa_serial" + " = " + mfaSerial + "\n"
		appendStringToFile(path, line)
	}
}

func CreateAWSConfigFile(config *configuration.Config, profile string, region string, output string) {
	if output == "" {
		output = getUserOutput(config)
	}
	homePath, pathError := GetUserHomeDir()
	if pathError != nil {
		config.Logger.Error(pathError.Error())
	}
	path := homePath + "/.aws/config"
	line := "\n[" + profile + "]\n"
	appendStringToFile(path, line)
	line = "region" + " = " + region + "\n"
	appendStringToFile(path, line)
	line = "output" + " = " + output + "\n"
	appendStringToFile(path, line)
}

func createConfigProfileFromCredentials(homeDir string, config *configuration.Config, profile string) {
	profilesInCredentials := uniqueNonEmptyElementsOf(getProfilesFromFile(config, homeDir+"/.aws/credentials"))
	config.Logger.Always("Available profile names are: " + fmt.Sprint("[ "+strings.Join(profilesInCredentials, ", ")+" ]"))
	config.Logger.GetInput("Profile", &profile)
	for !helpers.SliceContains(profilesInCredentials, profile) {
		config.Logger.Always("Invalid profile name. Available profile names are: " + fmt.Sprint("[ "+strings.Join(profilesInCredentials, ", ")+" ]"))
		config.Logger.GetInput("Profile", &profile)
	}
	region := getUserRegion(config)
	CreateAWSConfigFile(config, profile, region, "")
}

func setProfileInfoAndCreateConfigFile(config *configuration.Config) (profile string) {
	var ans string
	config.Logger.GetInput("Do you want to create a default profile with default values? *Y* / *N*", &ans)
	if strings.ToUpper(ans) == "Y" {
		region := "us-east-1"
		profile = "default"
		CreateAWSConfigFile(config, profile, region, "JSON")
	} else {
		profile = getUserProfile(config)
		region := getUserRegion(config)
		CreateAWSConfigFile(config, profile, region, "")
	}
	return profile
}

func addProfileToCredentials(profile string, homePath string, config *configuration.Config) {
	profilesInCredentials := getProfilesFromFile(config, homePath+"/.aws/credentials")
	if !helpers.SliceContains(profilesInCredentials, profile) {
		CreateAWSCredentialsFile(config, profile)
	}
}

func appendStringToFile(path, text string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}
