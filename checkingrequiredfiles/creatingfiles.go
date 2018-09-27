package checkingrequiredfiles

import (
	"bufio"
	"github.com/Appliscale/perun/helpers"
	"github.com/Appliscale/tyr/configuration"
	"github.com/Appliscale/tyr/logger"
	"github.com/Appliscale/tyr/myuser"
	"os"
	"strconv"
	"strings"
)

// Creating .aws/credentials file.
func CreateAWSCredentialsFile(config *configuration.Config, profile string) {
	if profile != "" {
		config.Logger.Always("You haven't got .aws/credentials file for profile " + profile)
		var awsAccessKeyID string
		var awsSecretAccessKey string
		var mfaSerial string

		config.Logger.GetInput("awsAccessKeyID", &awsAccessKeyID)
		config.Logger.GetInput("awsSecretAccessKey", &awsSecretAccessKey)
		config.Logger.GetInput("mfaSerial", &mfaSerial)

		homePath, pathError := myuser.GetUserHomeDir()
		if pathError != nil {
			config.Logger.Error(pathError.Error())
		}
		path := homePath + "/.aws/credentials"
		line := "[" + profile + "-long-term" + "]\n"
		appendStringToFile(path, line)
		line = "aws_access_key_id" + " = " + awsAccessKeyID + "\n"
		appendStringToFile(path, line)
		line = "aws_secret_access_key" + " = " + awsSecretAccessKey + "\n"
		appendStringToFile(path, line)
		line = "mfa_serial" + " = " + mfaSerial + "\n"
		appendStringToFile(path, line)
	}
}

// Looking for [profiles] in credentials or config and return all.
func getProfilesFromFile(path string, mylogger *logger.Logger) []string {
	credentials, credentialsError := os.Open(path)
	if credentialsError != nil {
		mylogger.Error(credentialsError.Error())
		return []string{}
	}
	defer credentials.Close()
	profiles := make([]string, 0)
	scanner := bufio.NewScanner(credentials)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "[") {
			profile := strings.TrimPrefix(scanner.Text(), "[")
			profile = strings.TrimSuffix(profile, "]")
			if strings.Contains(profile, "profile ") {
				profile = strings.TrimPrefix(profile, "profile ")
			}
			if strings.Contains(profile, "-long-term") {
				profile = strings.TrimSuffix(profile, "-long-term")
			}
			profiles = append(profiles, profile)
		}
	}
	return profiles
}

//List of all available regions.
func showRegions(myLogger *logger.Logger) {
	myLogger.Always("Regions:")
	for i := 0; i < len(regions); i++ {
		pom := strconv.Itoa(i)
		myLogger.Always("Number " + pom + " region " + regions[i])
	}
}

// Choosing one region.
func setRegions(myLogger *logger.Logger) (region string, err bool) {
	var numberRegion int
	myLogger.GetInput("Choose region", &numberRegion)
	if numberRegion >= 0 && numberRegion < 14 {
		region = regions[numberRegion]
		err = false
	} else {
		err = true
	}
	return
}

// Choosing one profile.
func setProfile(myLogger *logger.Logger) (profile string, err bool) {
	myLogger.GetInput("Input name of profile", &profile)
	if profile != "" {
		err = true
	} else {
		err = false
	}
	return
}

func setOutput(config *configuration.Config) (output string, err bool) {
	config.Logger.GetInput("Input the output format [json, text, table]", &output)
	if !helpers.SliceContains([]string{"json", "text", "table"}, output) {
		err = true
	} else {
		err = false
	}
	return
}

func GetOutput(config *configuration.Config) string {
	output, err := setOutput(config)
	for err {
		config.Logger.Always("Try again, invalid output")
		output, err = setOutput(config)
	}
	config.Logger.Always("Your output is: " + output)
	return output
}

// Creating .aws/config file.
func CreateAWSConfigFile(config *configuration.Config, profile string, region string, output string) {
	if output == "" {
		output = GetOutput(config)
	}
	homePath, pathError := myuser.GetUserHomeDir()
	if pathError != nil {
		config.Logger.Error(pathError.Error())
	}
	path := homePath + "/.aws/config"
	line := "[" + profile + "]\n"
	appendStringToFile(path, line)
	line = "region" + " = " + region + "\n"
	appendStringToFile(path, line)
	line = "output" + " = " + output + "\n"
	appendStringToFile(path, line)
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
