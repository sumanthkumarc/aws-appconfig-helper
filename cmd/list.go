package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list given type",
	Long:  `list given type`,
	Run:   listTypes,
}

var entity string

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&appId, "app-id", "", "", "Application id")
	listCmd.MarkFlagRequired("app-id")

	listCmd.Flags().StringVarP(&entity, "entity", "e", "", "Entity type to list. valid values - deployment-strategy")
	listCmd.MarkFlagRequired("entity")
}

func listTypes(cmd *cobra.Command, args []string) {

	switch entity {
	case "deployment-strategy":
		output, err := listDeploymentStrategies()
		if err != nil {
			log.Fatalf("error when getting strategies %s", err)
		}

		fmt.Print(output)
	}
}

func listDeploymentStrategies() (appconfig.ListDeploymentStrategiesOutput, error) {
	input := &appconfig.ListDeploymentStrategiesInput{}

	result, err := appConfig.ListDeploymentStrategies(input)

	if err != nil {
		return appconfig.ListDeploymentStrategiesOutput{}, err
	}

	return *result, nil
}
