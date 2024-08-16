// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package prebuild

import (
	"context"
	"log"
	"strconv"

	apiclient_util "github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/views"
	"github.com/daytonaio/daytona/pkg/views/prebuild/add"
	"github.com/daytonaio/daytona/pkg/views/workspace/selection"
	"github.com/spf13/cobra"
)

var prebuildUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a prebuild configuration",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var prebuildAddView add.PrebuildAddView
		var prebuild *apiclient.PrebuildDTO
		ctx := context.Background()

		apiClient, err := apiclient_util.GetApiClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		userGitProviders, res, err := apiClient.GitProviderAPI.ListGitProviders(ctx).Execute()
		if err != nil {
			log.Fatal(apiclient_util.HandleErrorResponse(res, err))
		}

		if len(userGitProviders) == 0 {
			views.RenderInfoMessage("No registered Git providers have been found - please register a Git provider using 'daytona git-provider add' in order to start using prebuilds.")
			return
		}

		if len(args) == 0 {
			prebuilds, res, err := apiClient.PrebuildAPI.ListPrebuilds(context.Background(), "").Execute()
			if err != nil {
				log.Fatal(apiclient_util.HandleErrorResponse(res, err))
			}

			if len(prebuilds) == 0 {
				views.RenderInfoMessage("No project configs found")
				return
			}

			prebuild = selection.GetPrebuildFromPrompt(prebuilds, "Update")
			if prebuild == nil {
				return
			}
		} else {
			prebuild, res, err = apiClient.PrebuildAPI.GetPrebuild(ctx, args[0]).Execute()
			if err != nil {
				log.Fatal(apiclient_util.HandleErrorResponse(res, err))
			}
		}

		prebuildAddView.Branch = prebuild.Branch
		prebuildAddView.CommitInterval = strconv.Itoa(int(prebuild.CommitInterval))
		prebuildAddView.TriggerFiles = prebuild.TriggerFiles
		prebuildAddView.ProjectConfigName = prebuild.ProjectConfigName

		add.PrebuildCreationView(&prebuildAddView, false)

		commitInterval, err := strconv.Atoi(prebuildAddView.CommitInterval)
		if err != nil {
			log.Fatal("commit interval must be a number")
		}

		newPrebuild := apiclient.CreatePrebuildDTO{
			Id:                &prebuild.Id,
			ProjectConfigName: prebuildAddView.ProjectConfigName,
			Branch:            prebuildAddView.Branch,
			CommitInterval:    int32(commitInterval),
			TriggerFiles:      prebuildAddView.TriggerFiles,
			RunAtInit:         runFlag,
		}

		res, err = apiClient.PrebuildAPI.SetPrebuild(ctx).Prebuild(newPrebuild).Execute()
		if err != nil {
			log.Fatal(apiclient_util.HandleErrorResponse(res, err))
		}

		views.RenderInfoMessage("Prebuild updated successfully")
	},
}

func init() {
	prebuildUpdateCmd.Flags().BoolVar(&runFlag, "run", false, "Run the prebuild once after updating it")
}
