package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy files",
	Long:  `deploy files`,
	Run: func(cmd *cobra.Command, args []string) {
		validateDeployFlags(cmd, args)
		deployConfigs(cmd, args)
	},
}

var version string
var profileId string
var deploymentStrategyId string

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVarP(&appId, "app-id", "", "", "Application id")
	deployCmd.MarkFlagRequired("app-id")

	deployCmd.Flags().StringVarP(&envId, "env-id", "", "", "Environment id to deploy the config")
	deployCmd.MarkFlagRequired("env-id")

	deployCmd.Flags().StringVarP(&profileId, "profile-id", "", "all", "Single config profile id to deploy. default: all")
	deployCmd.MarkFlagRequired("profile-id")

	deployCmd.Flags().StringVarP(&deploymentStrategyId, "strategy", "", "AppConfig.AllAtOnce", "deployment stategy to use")
	deployCmd.MarkFlagRequired("strategy")

	deployCmd.Flags().StringVarP(&version, "version", "", "latest", "deploy the provided version of the given file. default: latest")
}

func validateDeployFlags(cmd *cobra.Command, args []string) {
	// @todo Validate flag values
	// check if provided app id, env id and file exists

	// enforce latest version if config name is not explicit.
	if profileId == "all" {
		version = "latest"
	}
}

func deployConfigs(cmd *cobra.Command, args []string) {

	if profileId != "all" {
		if err := deployConfig(profileId, version); err != nil {
			log.Fatalf("unable to deploy profile %s with error - %s", profileId, err)
		}

		fmt.Printf("deployed version %s of profile %s successfully", version, profileId)
	} else {

		configList, err := getConfigList(appId)
		if err != nil {
			log.Fatalf("error getting config list %s", err)
		}

		if len(configList.Items) > 0 {
			for _, configProfile := range configList.Items {
				// find latest version to deploy
				latestVersion, err := getLatestConfigVersion(*configProfile.Id)

				if err != nil {
					log.Fatalf("unable to fetch latest version for profile %s with error - %s", *configProfile.Id, err)
				}

				if err := deployConfig(*configProfile.Id, latestVersion); err != nil {
					log.Fatalf("unable to deploy profile %s with error - %s", profileId, err)
				}

			}

			fmt.Printf("deployed latest versions of all config profiles successfully")
		} else {
			fmt.Println("no config profiles found to deploy")
		}
	}

}

func deployConfig(profileId string, configVersion string) error {
	input := &appconfig.StartDeploymentInput{
		ApplicationId:          &appId,
		ConfigurationProfileId: &profileId,
		ConfigurationVersion:   &configVersion,
		DeploymentStrategyId:   &deploymentStrategyId,
		EnvironmentId:          &envId,
		// Tags: map[string]*string{
		// 	"deployed-by": "appconfig-helper",
		// },
	}

	_, err := appConfig.StartDeployment(input)
	if err != nil {
		return err
	}

	return nil
}

func getLatestConfigVersion(profileId string) (string, error) {
	var versions []int64

	versionsOutput, err := getConfigProfileVersions(profileId)

	if err != nil {
		return "", err
	}

	if len(versionsOutput.Items) > 0 {

		for _, version := range versionsOutput.Items {
			versions = append(versions, *version.VersionNumber)
		}
		return fmt.Sprintf("%d", max(versions)), nil
	} else {
		return "", fmt.Errorf("no versions found for profile %s", profileId)
	}
}

func max(n []int64) int64 {
	for i := 1; i < len(n); i++ {
		if n[0] < n[i] {
			n[0] = n[i]
		}
	}

	return n[0]
}
