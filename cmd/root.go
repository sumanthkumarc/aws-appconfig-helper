package cmd

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/aws/aws-sdk-go/service/appconfigdata"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:              "aws-appconfig-helper",
		Short:            "AWS appconfig helper",
		Long:             "AWS appconfig helper",
		PersistentPreRun: initAWSConfig,
	}
)

var appId string
var envId string
var appConfigData *appconfigdata.AppConfigData
var appConfig *appconfig.AppConfig

func initAWSConfig(cmd *cobra.Command, args []string) {

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	})

	if err != nil {
		log.Fatal(err)
	}

	appConfigData = appconfigdata.New(awsSession)
	appConfig = appconfig.New(awsSession)
}

func Execute() error {
	return rootCmd.Execute()
}
