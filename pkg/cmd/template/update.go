// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package projectconfig

import (
	"context"
	"net/http"

	apiclient_util "github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/views"
	"github.com/daytonaio/daytona/pkg/views/workspace/selection"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var templateUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a template",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var template *apiclient.Template
		var res *http.Response
		ctx := context.Background()

		apiClient, err := apiclient_util.GetApiClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		if len(args) == 0 {
			templateList, res, err := apiClient.TemplateAPI.ListTemplates(ctx).Execute()
			if err != nil {
				log.Fatal(apiclient_util.HandleErrorResponse(res, err))
			}

			template = selection.GetTemplateFromPrompt(templateList, "Update")
			if template == nil {
				return
			}
		} else {
			template, res, err = apiClient.TemplateAPI.GetTemplate(ctx, args[0]).Execute()
			if err != nil {
				log.Fatal(apiclient_util.HandleErrorResponse(res, err))
			}
		}

		if template == nil {
			return
		}

		var newTemplate apiclient.Template

		res, err = apiClient.TemplateAPI.SetTemplate(ctx).Template(newTemplate).Execute()
		if err != nil {
			log.Fatal(apiclient_util.HandleErrorResponse(res, err))
		}

		views.RenderInfoMessage("Project config updated successfully")
	},
}
