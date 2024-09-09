// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package list

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/views"
	"github.com/daytonaio/daytona/pkg/views/template/info"
	views_util "github.com/daytonaio/daytona/pkg/views/util"
	"golang.org/x/term"
)

type RowData struct {
	Name           string
	ProjectConfigs string
}

func ListTemplates(templateList []apiclient.Template, apiServerConfig *apiclient.ServerConfig, specifyGitProviders bool) {
	re := lipgloss.NewRenderer(os.Stdout)

	headers := []string{"Name", "ProjectConfigs"}

	data := [][]string{}

	for _, pc := range templateList {
		var rowData *RowData
		var row []string

		rowData = getTableRowData(pc, apiServerConfig, specifyGitProviders)
		row = getRowFromRowData(*rowData)
		data = append(data, row)
	}

	terminalWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println(data)
		return
	}

	breakpointWidth := views.GetContainerBreakpointWidth(terminalWidth)

	minWidth := views_util.GetTableMinimumWidth(data)

	if breakpointWidth == 0 || minWidth > breakpointWidth {
		renderUnstyledList(templateList, apiServerConfig)
		return
	}

	t := table.New().
		Headers(headers...).
		Rows(data...).
		BorderStyle(re.NewStyle().Foreground(views.LightGray)).
		BorderRow(false).BorderColumn(false).BorderLeft(false).BorderRight(false).BorderTop(false).BorderBottom(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return views.TableHeaderStyle
			}
			return views.BaseCellStyle
		}).Width(breakpointWidth - 2*views.BaseTableStyleHorizontalPadding - 1)

	fmt.Println(views.BaseTableStyle.Render(t.String()))
}

func renderUnstyledList(templateList []apiclient.Template, apiServerConfig *apiclient.ServerConfig) {
	for _, pc := range templateList {
		info.Render(&pc, apiServerConfig, true)

		if pc.Name != templateList[len(templateList)-1].Name {
			fmt.Printf("\n%s\n\n", views.SeparatorString)
		}
	}
}

func getRowFromRowData(rowData RowData) []string {
	row := []string{
		views.NameStyle.Render(rowData.Name),
		views.DefaultRowDataStyle.Render(rowData.ProjectConfigs),
	}

	return row
}

func getTableRowData(template apiclient.Template, apiServerConfig *apiclient.ServerConfig, specifyGitProviders bool) *RowData {
	rowData := RowData{"", ""}

	rowData.Name = template.Name + views_util.AdditionalPropertyPadding
	rowData.ProjectConfigs = template.ProjectConfigs[0]

	return &rowData
}
