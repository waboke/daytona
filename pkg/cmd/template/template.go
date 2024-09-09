// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package projectconfig

import (
	"github.com/daytonaio/daytona/internal/util"
	"github.com/spf13/cobra"
)

var TemplateCmd = &cobra.Command{
	Use:     "template",
	Short:   "Manage workspace templates",
	Aliases: []string{"templates"},
	GroupID: util.WORKSPACE_GROUP,
}

func init() {
	TemplateCmd.AddCommand(templateListCmd)
	TemplateCmd.AddCommand(templateInfoCmd)
	TemplateCmd.AddCommand(templateAddCmd)
	TemplateCmd.AddCommand(templateUpdateCmd)
	TemplateCmd.AddCommand(templateDeleteCmd)
}
