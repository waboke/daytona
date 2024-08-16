// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package add

import (
	"log"
	"strconv"
	"strings"

	"github.com/daytonaio/daytona/pkg/views"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

var DEFAULT_COMMIT_INTERVAL = "10"

type PrebuildAddView struct {
	ProjectConfigName string
	Branch            string
	CommitInterval    string
	TriggerFiles      []string
}

func PrebuildCreationView(prebuildAddView *PrebuildAddView, editing bool) {
	if prebuildAddView.CommitInterval == "" {
		prebuildAddView.CommitInterval = DEFAULT_COMMIT_INTERVAL
	}
	triggerFilesInput := ""

	for _, triggerFile := range prebuildAddView.TriggerFiles {
		triggerFilesInput += triggerFile + "\n"
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Commit interval").
				Description("Commit interval").
				Value(&prebuildAddView.CommitInterval).
				Validate(func(str string) error {
					_, err := strconv.Atoi(str)
					return err
				}),
			huh.NewText().
				Title("Trigger files").
				Description("Enter full paths for files whose changes you want to explicitly trigger a prebuild. Use newlines for multiple entries").
				Value(&triggerFilesInput),
		),
	).WithTheme(views.GetCustomTheme())

	keyMap := huh.NewDefaultKeyMap()
	keyMap.Text = huh.TextKeyMap{
		NewLine: key.NewBinding(key.WithKeys("alt+enter"), key.WithHelp("alt+enter", "new line")),
		Next:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "next")),
		Prev:    key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev")),
	}

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	prebuildAddView.TriggerFiles = []string{}
	lines := strings.Split(triggerFilesInput, "\n")

	for _, line := range lines {
		if line != "" {
			prebuildAddView.TriggerFiles = append(prebuildAddView.TriggerFiles, strings.TrimRight(line, " "))
		}
	}
}
