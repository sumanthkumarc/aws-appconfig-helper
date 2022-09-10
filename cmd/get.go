package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/aws/aws-sdk-go/service/appconfigdata"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetch files",
	Long:  `Fetch files`,
	Run: func(cmd *cobra.Command, args []string) {
		validateFlags(cmd, args)
		getFiles(cmd, args)
	},
}

var watch bool
var files string
var srcFiles = make(map[string]string)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&appId, "app-id", "", "", "Application id")
	getCmd.MarkFlagRequired("app-id")

	getCmd.Flags().StringVarP(&envId, "env-id", "", "", "Environment id")
	getCmd.MarkFlagRequired("env-id")

	getCmd.Flags().BoolVarP(&watch, "watch", "w", false, "watch files. Default false.")
	getCmd.Flags().StringVarP(&files, "files", "f", "all", "File list in format: <src_file:dest_path>,<src_file:dest_path>")
}

func validateFlags(cmd *cobra.Command, args []string) {
	// @todo Validate flag values
	// check if provided app id and env id exists
	if files == "all" {

		configProfilesOutput, err := getConfigList(appId)
		if err != nil {
			log.Fatalf("error getting config list %s", err)
		}

		if configProfilesOutput.Items == nil {
			log.Println("No config profiles found")
		}

		for _, c := range configProfilesOutput.Items {
			srcFiles[*c.Name] = *c.Name
		}
	} else {
		for _, fileMapping := range strings.Split(files, ",") {
			files := strings.Split(fileMapping, ":")
			srcFiles[files[0]] = files[1]
		}
	}
}

func getFiles(cmd *cobra.Command, args []string) {

	for srcFile, destPath := range srcFiles {

		appConfigDetails := &appconfigdata.StartConfigurationSessionInput{
			ApplicationIdentifier:          &appId,
			ConfigurationProfileIdentifier: &srcFile,
			EnvironmentIdentifier:          &envId,
		}

		sessOutput, err := appConfigData.StartConfigurationSession(appConfigDetails)

		if err != nil {
			log.Fatal(err)
		}

		configInput := appconfigdata.GetLatestConfigurationInput{
			ConfigurationToken: sessOutput.InitialConfigurationToken,
		}

		// @todo implement watch
		if !watch {
			configData, err := appConfigData.GetLatestConfiguration(&configInput)

			if err != nil {
				log.Fatal(err)
			}

			if configData.Configuration != nil {
				// fmt.Println(configData.Configuration)
				writeData(configData.Configuration, destPath)
			}

		}
	}

	fmt.Println("successfully fetched files from given environment")

}

func writeData(data []byte, path string) {
	// @todo check if path is full fqdn to file else throw error
	//  We need file path to write data
	file, err := os.Create(path)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	file.Write(data)
}

func getConfigList(app string) (appconfig.ListConfigurationProfilesOutput, error) {
	configProfilesInput := &appconfig.ListConfigurationProfilesInput{
		ApplicationId: &app,
	}

	configProfilesOutput, err := appConfig.ListConfigurationProfiles(configProfilesInput)
	if err != nil {
		return appconfig.ListConfigurationProfilesOutput{}, err
	}

	return *configProfilesOutput, nil
}
