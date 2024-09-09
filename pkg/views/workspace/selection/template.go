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

var MAX_DESCRIPTION_LENGTH = 30

func GetTemplateFromPrompt(templates []apiclient.Template, actionVerb string) *apiclient.Template {
	choiceChan := make(chan *apiclient.Template)
	go selectTemplatePrompt(templates, actionVerb, choiceChan)
	return <-choiceChan
}

func selectTemplatePrompt(templates []apiclient.Template, actionVerb string, choiceChan chan<- *apiclient.Template) {
	items := []list.Item{}

	for _, t := range templates {
		templateName := t.Name
		if t.Name == "" {
			templateName = "Unnamed Project Config"
		}

		var description string

		for _, projectConfig := range t.ProjectConfigs {
			description += projectConfig + ", "
		}

		if len(description) > 1 {
			description = description[:len(description)-2]
		}

		if len(description) > MAX_DESCRIPTION_LENGTH {
			description = description[:MAX_DESCRIPTION_LENGTH] + " ..."
		}

		newItem := item[apiclient.Template]{title: templateName, desc: description, choiceProperty: t}
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

	title := "Select a Template To " + actionVerb
	l.Title = views.GetStyledMainTitle(title)
	l.Styles.Title = titleStyle

	m := model[apiclient.Template]{list: l}

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := p.(model[apiclient.Template]); ok && m.choice != nil {
		choiceChan <- m.choice
	} else {
		choiceChan <- nil
	}
}
