package checkingrequiredfiles

import (
	"fmt"
	"github.com/Appliscale/perun/helpers"
	"github.com/Appliscale/tyr/configuration"
	"github.com/Appliscale/tyr/myuser"
	"os"
	"strings"
)

var regions = []string{
	"us-east-1",
	"us-east-2",
	"us-west-1",
	"us-west-2",
	"ca-central-1",
	"ca-central-1",
	"eu-west-1",
	"eu-west-2",
	"ap-northeast-1",
	"ap-northeast-2",
	"ap-southeast-1",
	"ap-southeast-2",
	"ap-south-1",
	"sa-east-1",
}

func CheckingRequiredFiles(config *configuration.Config) bool {
	homeDir, pathError := myuser.GetUserHomeDir()
	if pathError != nil {
		config.Logger.Error(pathError.Error())
		return false
	}

	configAWSExists, configError := isAWSConfigPresent(homeDir)
	if configError != nil {
		config.Logger.Error(configError.Error())
	}

	credentialsExists, credentialsError := isCredentialsPresent(homeDir)
	if credentialsError != nil {
		config.Logger.Error(credentialsError.Error())
	}

	profile := config.Profile

	if configAWSExists {
		profilesInConfig := getProfilesFromFile(homeDir+"/.aws/config", config.Logger)
		if !helpers.SliceContains(profilesInConfig, profile) {
			var ans string
			config.Logger.GetInput("You don't have the "+profile+" profile in your config file. Would you like to create one?", &ans)
			if strings.ToUpper(ans) == "Y" {
				region := GetRegion(config)
				CreateAWSConfigFile(config, profile, region, "")
			} else {
				config.Logger.Info("You can use another profile by setting the \"-p\" argument or specify a different default profile by setting the AWS_PROFILE variable")
				return false
			}
		}
		if credentialsExists {
			addProfileToCredentials(profile, homeDir, config)
		} else {
			CreateAWSCredentialsFile(config, profile)
		}
	} else {
		if credentialsExists {
			var ans string
			config.Logger.GetInput("File .aws/config does not exist, but .aws/credentials has been found. Do you want to create config file using one of the profiles in the .aws/credentials?", &ans)
			if strings.ToUpper(ans) == "Y" {
				profilesInCredentials := uniqueNonEmptyElementsOf(getProfilesFromFile(homeDir+"/.aws/credentials", config.Logger))
				config.Logger.Always("Available profile names are: " + fmt.Sprint("[ "+strings.Join(profilesInCredentials, ", ")+" ]"))
				config.Logger.GetInput("Profile", &profile)
				for !helpers.SliceContains(profilesInCredentials, profile) {
					config.Logger.Always("Invalid profile name. Available profile names are: " + fmt.Sprint("[ "+strings.Join(profilesInCredentials, ", ")+" ]"))
					config.Logger.GetInput("Profile", &profile)
				}

				region := GetRegion(config)
				CreateAWSConfigFile(config, profile, region, "")
				return true

			} else {
				profile = setProfileInfoAndCreateConfigFile(config)
				CreateAWSCredentialsFile(config, profile)
			}
		} else {
			profile = setProfileInfoAndCreateConfigFile(config)
			CreateAWSCredentialsFile(config, profile)
		}
	}
	return true

}

func uniqueNonEmptyElementsOf(s []string) []string {
	unique := make(map[string]bool, len(s))
	us := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}

	return us

}

func setProfileInfoAndCreateConfigFile(config *configuration.Config) (profile string) {
	var ans string
	config.Logger.GetInput("File .aws/config does not exist. Do you want to create a default profile with default values?", &ans)
	if strings.ToUpper(ans) == "Y" {
		region := "us-east-1"
		profile = "default"
		CreateAWSConfigFile(config, profile, region, "JSON")
	} else {
		profile = GetProfile(config)
		region := GetRegion(config)
		CreateAWSConfigFile(config, profile, region, "")
	}
	return profile
}

// Checking if profile is in .aws/credentials.
func addProfileToCredentials(profile string, homePath string, config *configuration.Config) {
	profilesInCredentials := getProfilesFromFile(homePath+"/.aws/credentials", config.Logger)
	if !helpers.SliceContains(profilesInCredentials, profile) {
		CreateAWSCredentialsFile(config, profile)
	}
}

func GetProfile(config *configuration.Config) string {
	profile, err := setProfile(config.Logger)
	for !err {
		config.Logger.Always("Try again, invalid profile")
		profile, err = setProfile(config.Logger)
	}
	config.Logger.Always("Your region is: " + profile)
	return profile
}

func GetRegion(config *configuration.Config) string {
	showRegions(config.Logger)
	region, err := setRegions(config.Logger)
	for err {
		config.Logger.Always("Try again, invalid region")
		region, err = setRegions(config.Logger)
	}
	config.Logger.Always("Your region is: " + region)
	return region
}

// Looking for .aws.config.
func isAWSConfigPresent(homePath string) (bool, error) {
	_, credentialsError := os.Open(homePath + "/.aws/config")
	if credentialsError != nil {
		return false, nil
	}
	return true, nil
}

// Looking for .aws/credentials.
func isCredentialsPresent(homePath string) (bool, error) {
	_, credentialsError := os.Open(homePath + "/.aws/credentials")
	if credentialsError != nil {
		return false, nil
	}
	return true, nil
}
