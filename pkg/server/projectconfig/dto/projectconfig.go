// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package dto

import (
	"github.com/daytonaio/daytona/pkg/workspace/project/buildconfig"
)

type CreateProjectConfigDTO struct {
	Name          string                   `json:"name" validate:"required"`
	Image         *string                  `json:"image,omitempty" validate:"optional"`
	User          *string                  `json:"user,omitempty" validate:"optional"`
	BuildConfig   *buildconfig.BuildConfig `json:"buildConfig,omitempty" validate:"optional"`
	RepositoryUrl string                   `json:"repositoryUrl" validate:"required"`
	EnvVars       map[string]string        `json:"envVars" validate:"required"`
} // @name CreateProjectConfigDTO

type PrebuildDTO struct {
	Id                string   `json:"id" validate:"required"`
	ProjectConfigName string   `json:"projectConfigName" validate:"required"`
	Branch            string   `json:"branch" validate:"required"`
	CommitInterval    int      `json:"commitInterval" validate:"required"`
	TriggerFiles      []string `json:"triggerFiles" validate:"required"`
} // @name PrebuildDTO

// Todo - use PrebuildDTOs
type CreatePrebuildDTO struct {
	Id                string   `json:"id" validate:"optional"`
	ProjectConfigName string   `json:"projectConfigName" validate:"required"`
	Branch            string   `json:"branch" validate:"required"`
	CommitInterval    int      `json:"commitInterval" validate:"required"`
	TriggerFiles      []string `json:"triggerFiles" validate:"required"`
	RunAtInit         *bool    `json:"runAtInit" validate:"required"`
} // @name CreatePrebuildDTO
