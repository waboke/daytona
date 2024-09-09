// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package projectconfig

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"
	apiclient_util "github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/views"
	"github.com/daytonaio/daytona/pkg/views/workspace/selection"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var allFlag bool
var yesFlag bool
var forceFlag bool

var templateDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete a template",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var selectedTemplate *apiclient.Template
		var selectedTemplateName string

		apiClient, err := apiclient_util.GetApiClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		if allFlag {
			if !yesFlag {
				form := huh.NewForm(
					huh.NewGroup(
						huh.NewConfirm().
							Title("Delete all templates?").
							Description("Are you sure you want to delete all templates?").
							Value(&yesFlag),
					),
				).WithTheme(views.GetCustomTheme())

				err := form.Run()
				if err != nil {
					log.Fatal(err)
				}

				if !yesFlag {
					fmt.Println("Operation canceled.")
					return
				}
			}

			templates, res, err := apiClient.TemplateAPI.ListTemplates(context.Background()).Execute()
			if err != nil {
				log.Fatal(apiclient_util.HandleErrorResponse(res, err))
			}

			if len(templates) == 0 {
				views.RenderInfoMessage("No templates found")
				return
			}

			for _, template := range templates {
				selectedTemplateName = template.Name
				res, err := apiClient.TemplateAPI.DeleteTemplate(context.Background(), selectedTemplateName).Execute()
				if err != nil {
					log.Error(apiclient_util.HandleErrorResponse(res, err))
					continue
				}
				views.RenderInfoMessage("Deleted template: " + selectedTemplateName)
			}
			return
		}

		if len(args) == 0 {
			templates, res, err := apiClient.TemplateAPI.ListTemplates(context.Background()).Execute()
			if err != nil {
				log.Fatal(apiclient_util.HandleErrorResponse(res, err))
			}

			if len(templates) == 0 {
				views.RenderInfoMessage("No templates found")
				return
			}

			selectedTemplate = selection.GetTemplateFromPrompt(templates, "Delete")
			if selectedTemplate == nil {
				return
			}
			selectedTemplateName = selectedTemplate.Name
		} else {
			selectedTemplateName = args[0]
		}

		res, err := apiClient.TemplateAPI.DeleteTemplate(context.Background(), selectedTemplateName).Execute()
		if err != nil {
			log.Fatal(apiclient_util.HandleErrorResponse(res, err))
		}

		views.RenderInfoMessage("Template deleted successfully")
	},
}

func init() {
	templateDeleteCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Delete all project configs")
	templateDeleteCmd.Flags().BoolVarP(&yesFlag, "yes", "y", false, "Confirm deletion without prompt")
	templateDeleteCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force delete project configs")
}
