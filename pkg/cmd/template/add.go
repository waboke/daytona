// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package projectconfig

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	apiclient_util "github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/cmd/projectconfig"
	"github.com/daytonaio/daytona/pkg/views"
	"github.com/daytonaio/daytona/pkg/views/workspace/create"
	"github.com/daytonaio/daytona/pkg/views/workspace/selection"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var templateAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"new", "create"},
	Short:   "Add a workspace template",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var projectConfig *apiclient.ProjectConfig
		var newTemplate apiclient.Template
		ctx := context.Background()

		apiClient, err := apiclient_util.GetApiClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		gitProviders, res, err := apiClient.GitProviderAPI.ListGitProviders(ctx).Execute()
		if err != nil {
			log.Fatal(apiclient_util.HandleErrorResponse(res, err))
		}

		projectConfigList, res, err := apiClient.ProjectConfigAPI.ListProjectConfigs(ctx).Execute()
		if err != nil {
			log.Fatal(apiclient_util.HandleErrorResponse(res, err))
		}

		for {
			projectConfig = selection.GetProjectConfigFromPrompt(selection.ProjectConfigPromptConfig{
				ProjectConfigs: projectConfigList,
				ProjectOrder:   0,
				ShowNewOption:  true,
				ShowDoneOption: true,
				ActionVerb:     "Add To Template",
			})
			if projectConfig == nil {
				log.Fatal("No project config selected")
			}

			if projectConfig.Name == selection.DoneSelecting {
				break
			}

			if projectConfig.Name == selection.NewProjectConfigIdentifier {
				projectConfig, err = projectconfig.RunProjectConfigAddFlow(apiClient, gitProviders, ctx)
				if err != nil {
					log.Fatal(err)
				}
				if projectConfig == nil {
					return
				}
			}

			newTemplate.ProjectConfigs = append(newTemplate.ProjectConfigs, projectConfig.Name)
		}

		m := create.NewSummaryModel(create.SubmissionFormConfig{
			ChosenName:    &newTemplate.Name,
			SuggestedName: "test",
			ExistingNames: []string{},
			NameLabel:     "Workspace Template",
			Defaults:      nil,
		})

		if _, err := tea.NewProgram(m).Run(); err != nil {
			log.Fatal(err)
		}

		res, err = apiClient.TemplateAPI.SetTemplate(ctx).Template(apiclient.Template{
			Name:           "test",
			ProjectConfigs: newTemplate.ProjectConfigs,
		}).Execute()
		if err != nil {
			log.Fatal(apiclient_util.HandleErrorResponse(res, err))
		}

		views.RenderInfoMessage(fmt.Sprintf("Template %s added successfully", "test"))
	},
}

var manualFlag bool

func init() {
	templateAddCmd.Flags().BoolVar(&manualFlag, "manual", false, "Manually enter the Git repository")
}
