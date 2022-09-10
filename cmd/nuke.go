package cmd

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/spf13/cobra"
)

var nukeCmd = &cobra.Command{
	Use:   "nuke",
	Short: "Wholesale deletion",
	Long:  `Wholesale deletion`,
	Run:   deletionHandler,
	Args: func(cmd *cobra.Command, args []string) error {

		if len(args) < 1 {
			return errors.New("requires at least one arg. either of app or env")
		}

		if !isValidArg(args[0]) {
			return fmt.Errorf("only either of %s is valid", stringify(validArgs))
		} else {
			if args[0] == "app" && appId == "" {
				return fmt.Errorf("please set --app-id value to delete")
			}

			if args[0] == "env" && (appId == "" || envId == "") {
				return fmt.Errorf("please set both --app-id and --env-id value to delete env")
			}
		}

		return nil
	},
}

var validArgs = []string{"app", "env"}

func isValidArg(arg string) bool {
	for _, v := range validArgs {
		if v == arg {
			return true
		}
	}
	return false
}

func stringify(stringList []string) string {
	return strings.Join(stringList, ", ")
}

func init() {
	rootCmd.AddCommand(nukeCmd)
	nukeCmd.Flags().StringVarP(&appId, "app-id", "", "", "Application id")
	nukeCmd.Flags().StringVarP(&envId, "env-id", "", "", "Environment id")
}

func deletionHandler(cmd *cobra.Command, args []string) {

	if args[0] == "env" {
		// delete env with --env-id
		if err := deleteEnv(envId); err != nil {
			log.Fatalf("unable to delete the env with id - %s with error - %s", envId, err)
		} else {
			log.Printf("env %s has been deleted sucessfully", envId)
		}
	}

	if args[0] == "app" {
		//  get list of env, and pass to env delete function
		envList, err := getEnvironmentList()
		if err != nil {
			log.Fatalf("unable to get env list. got error - %s", err)
		}

		if len(envList.Items) > 0 {
			// loop and delete the env
			for _, env := range envList.Items {
				if err := deleteEnv(*env.Id); err != nil {
					log.Fatalf("unable to delete env. got error - %s \n", err)
				} else {
					fmt.Printf("deleted the env %s \n", *env.Name)
				}
			}
		} else {
			fmt.Println("no envs found to delete")
		}

		//  get list of env config and their version and delete all the versions
		configList, err := getConfigList(appId)
		if err != nil {
			log.Fatalf("error getting config list %s", err)
		}

		if len(configList.Items) > 0 {
			for _, configProfile := range configList.Items {
				// delete all versions of config first
				if err := deleteConfigVersion(*configProfile.Id, -1); err != nil {
					log.Fatalf("unable to delete versions for config profile %s with error %s", *configProfile.Id, err)
				}

				if err := deleteConfigProfile(*configProfile.Id); err != nil {
					log.Fatalf("unable to delete profile with error - %s", err)
				}
			}
		} else {
			fmt.Println("no config profiles found to delete")
		}

		input := &appconfig.DeleteApplicationInput{
			ApplicationId: &appId,
		}

		_, err = appConfig.DeleteApplication(input)
		if err != nil {
			log.Fatalf("unable to delete application. error - %s", err)
		}

		fmt.Printf("deleted the application %s successfully", appId)

	}

}

func deleteEnv(env string) error {

	deleteEnvRequest := &appconfig.DeleteEnvironmentInput{
		ApplicationId: &appId,
		EnvironmentId: &env,
	}

	_, err := appConfig.DeleteEnvironment(deleteEnvRequest)

	if err != nil {
		return err
	}

	return nil
}

func getEnvironmentList() (appconfig.ListEnvironmentsOutput, error) {
	input := &appconfig.ListEnvironmentsInput{
		ApplicationId: &appId,
	}

	result, err := appConfig.ListEnvironments(input)
	if err != nil {
		return appconfig.ListEnvironmentsOutput{}, err
	}

	return *result, nil
}

func deleteConfigVersion(profileId string, version int64) error {
	if version == -1 {
		// delete all versions

		versions, err := getConfigProfileVersions(profileId)
		if err != nil {
			return err
		}

		if len(versions.Items) > 0 {
			for _, version := range versions.Items {
				input := &appconfig.DeleteHostedConfigurationVersionInput{
					ApplicationId:          &appId,
					ConfigurationProfileId: &profileId,
					VersionNumber:          version.VersionNumber,
				}

				_, err := appConfig.DeleteHostedConfigurationVersion(input)
				if err != nil {
					return err
				}
			}

			fmt.Printf("deleted the all versions for profile %s \n", profileId)
		} else {
			fmt.Printf("no versions found to delete for profile %s \n", profileId)
		}

	} else {
		input := &appconfig.DeleteHostedConfigurationVersionInput{
			ApplicationId:          &appId,
			ConfigurationProfileId: &profileId,
			VersionNumber:          &version,
		}

		_, err := appConfig.DeleteHostedConfigurationVersion(input)
		if err != nil {
			return err
		}

		fmt.Printf("deleted the version %d for profile %s \n", version, profileId)
	}

	return nil
}

func getConfigProfileVersions(profileId string) (appconfig.ListHostedConfigurationVersionsOutput, error) {
	input := &appconfig.ListHostedConfigurationVersionsInput{
		ApplicationId:          &appId,
		ConfigurationProfileId: &profileId,
	}

	result, err := appConfig.ListHostedConfigurationVersions(input)
	if err != nil {
		return appconfig.ListHostedConfigurationVersionsOutput{}, err
	}

	return *result, nil
}

func deleteConfigProfile(profileId string) error {
	input := &appconfig.DeleteConfigurationProfileInput{
		ApplicationId:          &appId,
		ConfigurationProfileId: &profileId,
	}

	_, err := appConfig.DeleteConfigurationProfile(input)
	if err != nil {
		return err
	}

	fmt.Printf("deleted the profile %s \n", profileId)

	return nil
}
