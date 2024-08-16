// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package projectconfig

import (
	"context"
	"net/http"

	apiclient_util "github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/pkg/apiclient"
	workspace_util "github.com/daytonaio/daytona/pkg/cmd/workspace/util"
	"github.com/daytonaio/daytona/pkg/views"
	"github.com/daytonaio/daytona/pkg/views/workspace/create"
	"github.com/daytonaio/daytona/pkg/views/workspace/selection"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var projectConfigUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a project config",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var projectConfig *apiclient.ProjectConfig
		var projectDtos []apiclient.CreateProjectDTO
		var res *http.Response
		ctx := context.Background()

		apiClient, err := apiclient_util.GetApiClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		if len(args) == 0 {
			projectConfigList, res, err := apiClient.ProjectConfigAPI.ListProjectConfigs(ctx).Execute()
			if err != nil {
				log.Fatal(apiclient_util.HandleErrorResponse(res, err))
			}

			projectConfig = selection.GetProjectConfigFromPrompt(projectConfigList, 0, false, "Update")
			if projectConfig == nil {
				return
			}
		} else {
			projectConfig, res, err = apiClient.ProjectConfigAPI.GetProjectConfig(ctx, args[0]).Execute()
			if err != nil {
				log.Fatal(apiclient_util.HandleErrorResponse(res, err))
			}
		}

		if projectConfig == nil {
			return
		}

		projectDtos = append(projectDtos, apiclient.CreateProjectDTO{
			Name: projectConfig.Name,
			Source: apiclient.CreateProjectSourceDTO{
				Repository: apiclient.GitRepository{
					Url: projectConfig.RepositoryUrl,
				},
			},
			BuildConfig: projectConfig.BuildConfig,
			EnvVars:     projectConfig.EnvVars,
		})

		projectDefaults := &create.ProjectConfigDefaults{
			BuildChoice: create.AUTOMATIC,
			Image:       &projectConfig.Image,
			ImageUser:   &projectConfig.User,
		}

		if projectConfig.BuildConfig != nil && projectConfig.BuildConfig.Devcontainer != nil {
			projectDefaults.DevcontainerFilePath = projectConfig.BuildConfig.Devcontainer.FilePath
		}

		create.ProjectsConfigurationChanged, err = create.RunProjectConfiguration(&projectDtos, *projectDefaults)
		if err != nil {
			log.Fatal(err)
		}

		newProjectConfig := apiclient.CreateProjectConfigDTO{
			Name:          projectConfig.Name,
			BuildConfig:   projectDtos[0].BuildConfig,
			Image:         projectDtos[0].Image,
			User:          projectDtos[0].User,
			RepositoryUrl: projectDtos[0].Source.Repository.Url,
		}

		newProjectConfig.EnvVars = *workspace_util.GetEnvVariables(projectDtos[0].EnvVars, nil)

		res, err = apiClient.ProjectConfigAPI.SetProjectConfig(ctx).ProjectConfig(newProjectConfig).Execute()
		if err != nil {
			log.Fatal(apiclient_util.HandleErrorResponse(res, err))
		}

		views.RenderInfoMessage("Project config updated successfully")
	},
}
