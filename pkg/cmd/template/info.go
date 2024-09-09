// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package projectconfig

import (
	"context"
	"net/http"

	apiclient_util "github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/cmd/format"
	"github.com/daytonaio/daytona/pkg/views/template/info"
	"github.com/daytonaio/daytona/pkg/views/workspace/selection"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var templateInfoCmd = &cobra.Command{
	Use:     "info",
	Short:   "Show template info",
	Aliases: []string{"view", "inspect"},
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		apiClient, err := apiclient_util.GetApiClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		apiServerConfig, res, err := apiClient.ServerAPI.GetConfig(context.Background()).Execute()
		if err != nil {
			log.Fatal(apiclient_util.HandleErrorResponse(res, err))
		}

		var template *apiclient.Template

		if len(args) == 0 {
			templateList, res, err := apiClient.TemplateAPI.ListTemplates(ctx).Execute()
			if err != nil {
				log.Fatal(apiclient_util.HandleErrorResponse(res, err))
			}

			if format.FormatFlag != "" {
				format.UnblockStdOut()
			}

			template = selection.GetTemplateFromPrompt(templateList, "View")
			if format.FormatFlag != "" {
				format.BlockStdOut()
			}

		} else {
			var res *http.Response
			template, res, err = apiClient.TemplateAPI.GetTemplate(ctx, args[0]).Execute()
			if err != nil {
				log.Fatal(apiclient_util.HandleErrorResponse(res, err))
			}
		}

		if template == nil {
			return
		}

		if format.FormatFlag != "" {
			formattedData := format.NewFormatter(template)
			formattedData.Print()
			return
		}

		info.Render(template, apiServerConfig, false)
	},
}

func init() {
	format.RegisterFormatFlag(templateInfoCmd)
}
