// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package info

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/daytonaio/daytona/internal/util"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/views"
	"golang.org/x/term"
)

const propertyNameWidth = 20

var propertyNameStyle = lipgloss.NewStyle().
	Foreground(views.LightGray)

var propertyValueStyle = lipgloss.NewStyle().
	Foreground(views.Light).
	Bold(true)

func Render(template *apiclient.Template, apiServerConfig *apiclient.ServerConfig, forceUnstyled bool) {
	var output string
	output += "\n\n"

	output += views.GetStyledMainTitle("Template Info") + "\n\n"

	output += getInfoLine("Name", template.Name) + "\n"

	projectConfigCount := len(template.ProjectConfigs)

	if projectConfigCount > 0 {
		if projectConfigCount == 1 {
			output += getInfoLine("Prebuild: ", getProjectConfigLine(template.ProjectConfigs[0], nil)) + "\n"
		} else {
			output += getInfoLine("Prebuilds: ", "") + "\n"
			for i, pc := range template.ProjectConfigs {
				if len(template.ProjectConfigs) != 1 {
					output += getProjectConfigLine(pc, util.Pointer(i+1)) + "\n"
				} else {
					output += getProjectConfigLine(pc, nil) + "\n"
				}
			}
		}
	}

	terminalWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println(output)
		return
	}
	if terminalWidth < views.TUITableMinimumWidth || forceUnstyled {
		renderUnstyledInfo(output)
		return
	}

	renderTUIView(output, views.GetContainerBreakpointWidth(terminalWidth))
}

func renderUnstyledInfo(output string) {
	fmt.Println(output)
}

func renderTUIView(output string, width int) {
	output = lipgloss.NewStyle().PaddingLeft(3).Render(output)

	content := lipgloss.
		NewStyle().Width(width).
		Render(output)

	fmt.Println(content)
}

func getInfoLine(key, value string) string {
	return propertyNameStyle.Render(fmt.Sprintf("%-*s", propertyNameWidth, key)) + propertyValueStyle.Render(value) + "\n"
}

func getProjectConfigLine(projectConfigName string, order *int) string {
	var line string
	if order != nil {
		line += propertyNameStyle.Render(fmt.Sprintf("%s#%d%s", strings.Repeat(" ", 3), *order, strings.Repeat(" ", 2)))
	}

	line += propertyValueStyle.Render(projectConfigName)

	if order != nil {
		line += "\n"
	}

	return line
}
