package configuration

import (
	"errors"
	"github.com/Appliscale/perun/utilities"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
	"os"
	"os/user"
	"time"
)

const dateFormat = "2006-01-02 15:04:05"

func InitialiseMFA(config Config) {
	tokenError := UpdateSessionToken(config, config.Profile, (*config.Regions)[0])
	if tokenError != nil {
		config.Logger.Error(tokenError.Error())
		os.Exit(1)
	}
}

func UpdateSessionToken(config Config, profile string, region string) error {
	if config.Mfa {
		currentUser, userError := user.Current()
		if userError != nil {
			return userError
		}

		credentialsFilePath := currentUser.HomeDir + "/.aws/credentials"
		configuration, loadCredentialsError := ini.Load(credentialsFilePath)
		if loadCredentialsError != nil {
			return loadCredentialsError
		}

		section, sectionError := configuration.GetSection(profile)
		if sectionError != nil {
			section, sectionError = configuration.NewSection(profile)
			if sectionError != nil {
				return sectionError
			}
		}

		profileLongTerm := profile + "-long-term"
		sectionLongTerm, profileLongTermError := configuration.GetSection(profileLongTerm)
		if profileLongTermError != nil {
			return profileLongTermError
		}

		sessionToken := section.Key("aws_session_token")
		expiration := section.Key("expiration")

		expirationDate, dataError := time.Parse(dateFormat, section.Key("expiration").Value())
		if dataError == nil {
			config.Logger.Info("Session token will expire in " + utilities.TruncateDuration(time.Since(expirationDate)).String() + " (" + expirationDate.Format(dateFormat) + ")")
		}

		mfaDevice := sectionLongTerm.Key("mfa_serial").Value()
		if mfaDevice == "" {
			return errors.New("There is no mfa_serial for the profile " + profileLongTerm + ". If you haven't used --mfa option you can change the default decision for MFA in the configuration file")
		}

		if sessionToken.Value() == "" || expiration.Value() == "" || time.Since(expirationDate).Nanoseconds() > 0 {
			currentSession, sessionError := session.NewSessionWithOptions(
				session.Options{
					Config: aws.Config{
						Region: &region,
					},
					Profile: profileLongTerm,
				})
			if sessionError != nil {
				return sessionError
			}

			var duration int64
			if config.MfaDuration == 0 {
				sessionError = config.Logger.GetInput("Duration", &duration)
				if sessionError != nil {
					return sessionError
				}
			} else {
				duration = config.MfaDuration
			}

			var tokenCode string
			sessionError = config.Logger.GetInput("MFA token code", &tokenCode)
			if sessionError != nil {
				return sessionError
			}

			stsSession := sts.New(currentSession)
			newToken, tokenError := stsSession.GetSessionToken(&sts.GetSessionTokenInput{
				DurationSeconds: &duration,
				SerialNumber:    aws.String(mfaDevice),
				TokenCode:       &tokenCode,
			})
			if tokenError != nil {
				return tokenError
			}

			section.Key("aws_access_key_id").SetValue(*newToken.Credentials.AccessKeyId)
			section.Key("aws_secret_access_key").SetValue(*newToken.Credentials.SecretAccessKey)
			sessionToken.SetValue(*newToken.Credentials.SessionToken)
			section.Key("expiration").SetValue(newToken.Credentials.Expiration.Format(dateFormat))

			configuration.SaveTo(credentialsFilePath)
		}
	}
	return nil
}
