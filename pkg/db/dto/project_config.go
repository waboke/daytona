// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package dto

import (
	"github.com/daytonaio/daytona/pkg/workspace/project/config"
)

type ProjectConfigDTO struct {
	Name          string            `gorm:"primaryKey"`
	Image         string            `json:"image"`
	User          string            `json:"user"`
	Build         *ProjectBuildDTO  `json:"build,omitempty" gorm:"serializer:json"`
	RepositoryUrl string            `json:"repositoryUrl"`
	EnvVars       map[string]string `json:"envVars" gorm:"serializer:json"`
	Prebuilds     []PrebuildDTO     `gorm:"serializer:json"`
	IsDefault     bool              `json:"isDefault"`
	WebhookId     *string           `json:"webhookId,omitempty"`
}

type PrebuildDTO struct {
	Id             string   `json:"id"`
	Branch         string   `json:"branch"`
	CommitInterval int      `json:"commitInterval"`
	TriggerFiles   []string `json:"triggerFiles"`
}

func ToProjectConfigDTO(projectConfig *config.ProjectConfig) ProjectConfigDTO {
	prebuilds := []PrebuildDTO{}
	for _, prebuild := range projectConfig.Prebuilds {
		prebuilds = append(prebuilds, ToPrebuildDTO(prebuild))
	}

	return ProjectConfigDTO{
		Name:          projectConfig.Name,
		Image:         projectConfig.Image,
		User:          projectConfig.User,
		Build:         ToProjectBuildDTO(projectConfig.BuildConfig),
		RepositoryUrl: projectConfig.RepositoryUrl,
		EnvVars:       projectConfig.EnvVars,
		Prebuilds:     prebuilds,
		IsDefault:     projectConfig.IsDefault,
		WebhookId:     projectConfig.WebhookId,
	}
}

func ToProjectConfig(projectConfigDTO ProjectConfigDTO) *config.ProjectConfig {
	prebuilds := []*config.PrebuildConfig{}
	for _, prebuildDTO := range projectConfigDTO.Prebuilds {
		prebuilds = append(prebuilds, ToPrebuild(prebuildDTO))
	}

	return &config.ProjectConfig{
		Name:          projectConfigDTO.Name,
		Image:         projectConfigDTO.Image,
		User:          projectConfigDTO.User,
		BuildConfig:   ToProjectBuild(projectConfigDTO.Build),
		RepositoryUrl: projectConfigDTO.RepositoryUrl,
		EnvVars:       projectConfigDTO.EnvVars,
		Prebuilds:     prebuilds,
		IsDefault:     projectConfigDTO.IsDefault,
		WebhookId:     projectConfigDTO.WebhookId,
	}
}

func ToPrebuildDTO(prebuild *config.PrebuildConfig) PrebuildDTO {
	return PrebuildDTO{
		Id:             prebuild.Id,
		Branch:         prebuild.Branch,
		CommitInterval: prebuild.CommitInterval,
		TriggerFiles:   prebuild.TriggerFiles,
	}
}

func ToPrebuild(prebuildDTO PrebuildDTO) *config.PrebuildConfig {
	return &config.PrebuildConfig{
		Id:             prebuildDTO.Id,
		Branch:         prebuildDTO.Branch,
		CommitInterval: prebuildDTO.CommitInterval,
		TriggerFiles:   prebuildDTO.TriggerFiles,
	}
}
