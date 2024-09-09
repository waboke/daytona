// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package selection

import (
	"fmt"
	"os"

	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/views"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var BlankProjectIdentifier = "<BLANK_PROJECT>"
var NewProjectConfigIdentifier = "<NEW_PROJECT_CONFIG>"
var DoneSelecting = "<DONE_SELECTING>"

type ProjectConfigPromptConfig struct {
	ProjectConfigs  []apiclient.ProjectConfig
	ProjectOrder    int
	ShowBlankOption bool
	ShowNewOption   bool
	ShowDoneOption  bool
	ActionVerb      string
}

func GetProjectConfigFromPrompt(config ProjectConfigPromptConfig) *apiclient.ProjectConfig {
	choiceChan := make(chan *apiclient.ProjectConfig)
	go selectProjectConfigPrompt(config, choiceChan)
	return <-choiceChan
}

func selectProjectConfigPrompt(config ProjectConfigPromptConfig, choiceChan chan<- *apiclient.ProjectConfig) {
	items := []list.Item{}

	if config.ShowDoneOption {
		newItem := item[apiclient.ProjectConfig]{title: "Done selecting project configurations", desc: "", choiceProperty: apiclient.ProjectConfig{
			Name: DoneSelecting,
		}}
		items = append(items, newItem)
	}

	if config.ShowBlankOption {
		newItem := item[apiclient.ProjectConfig]{title: "Make a blank project", desc: "(default project configuration)", choiceProperty: apiclient.ProjectConfig{
			Name: BlankProjectIdentifier,
		}}
		items = append(items, newItem)
	}

	for _, pc := range config.ProjectConfigs {
		projectConfigName := pc.Name
		if pc.Name == "" {
			projectConfigName = "Unnamed Project Config"
		}

		newItem := item[apiclient.ProjectConfig]{title: projectConfigName, desc: pc.RepositoryUrl, choiceProperty: pc}
		items = append(items, newItem)
	}

	if config.ShowNewOption {
		newItem := item[apiclient.ProjectConfig]{title: "+ Create a new project configuration", desc: "", choiceProperty: apiclient.ProjectConfig{
			Name: NewProjectConfigIdentifier,
		}}
		items = append(items, newItem)
	}

	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(views.Green).
		Foreground(views.Green).
		Bold(true).
		Padding(0, 0, 0, 1)

	d.Styles.SelectedDesc = d.Styles.SelectedTitle.Copy().Foreground(views.DimmedGreen)

	l := list.New(items, d, 0, 0)

	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(views.Green)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(views.Green)

	l.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(views.Green)
	l.FilterInput.TextStyle = lipgloss.NewStyle().Foreground(views.Green)

	title := "Select a Project Config To " + config.ActionVerb
	if config.ProjectOrder > 1 {
		title += fmt.Sprintf(" (Project #%d)", config.ProjectOrder)
	}
	l.Title = views.GetStyledMainTitle(title)
	l.Styles.Title = titleStyle

	m := model[apiclient.ProjectConfig]{list: l}

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := p.(model[apiclient.ProjectConfig]); ok && m.choice != nil {
		choiceChan <- m.choice
	} else {
		choiceChan <- nil
	}
}
