// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package prebuild

import (
	"context"
	"log"
	"net/http"

	apiclient_util "github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/views"
	"github.com/daytonaio/daytona/pkg/views/workspace/selection"
	"github.com/spf13/cobra"
)

var prebuildDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a prebuild configuration",
	Aliases: []string{"remove", "rm"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var selectedPrebuild *apiclient.PrebuildDTO
		var selectedPrebuildId string

		apiClient, err := apiclient_util.GetApiClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		if len(args) == 0 {
			prebuilds, res, err := apiClient.PrebuildAPI.ListPrebuilds(context.Background(), "").Execute()
			if err != nil {
				log.Fatal(apiclient_util.HandleErrorResponse(res, err))
			}

			if len(prebuilds) == 0 {
				views.RenderInfoMessage("No prebuilds found")
				return
			}

			selectedPrebuild = selection.GetPrebuildFromPrompt(prebuilds, "Delete")
			selectedPrebuildId = selectedPrebuild.Id
		} else {
			var res *http.Response

			selectedPrebuild, res, err = apiClient.PrebuildAPI.GetPrebuild(context.Background(), args[0]).Execute()
			if err != nil {
				log.Fatal(apiclient_util.HandleErrorResponse(res, err))
			}
		}

		res, err := apiClient.PrebuildAPI.DeletePrebuild(context.Background(), selectedPrebuild.ProjectConfigName, selectedPrebuildId).Execute()
		if err != nil {
			log.Fatal(apiclient_util.HandleErrorResponse(res, err))
		}

		views.RenderInfoMessage("Project config deleted successfully")
	},
}
