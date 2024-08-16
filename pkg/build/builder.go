// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package build

import (
	"github.com/daytonaio/daytona/pkg/containerregistry"
	"github.com/daytonaio/daytona/pkg/logs"
)

type IBuilder interface {
	Build(build Build) (string, string, error)
	CleanUp() error
	Publish(build Build) error
}

type Builder struct {
	id                  string
	hash                string
	projectDir          string
	image               string
	containerRegistry   containerregistry.ContainerRegistry
	buildStore          Store
	buildImageNamespace string
	loggerFactory       logs.LoggerFactory
	defaultProjectImage string
	defaultProjectUser  string
}
