// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package build

import (
	"path/filepath"

	"github.com/daytonaio/daytona/pkg/containerregistry"
	"github.com/daytonaio/daytona/pkg/logs"
	"github.com/daytonaio/daytona/pkg/ports"
)

type IBuilderFactory interface {
	Create(build Build) (IBuilder, error)
	CheckExistingBuild(build Build) (*Build, error)
}

type BuilderFactory struct {
	containerRegistry   containerregistry.ContainerRegistry
	buildImageNamespace string
	buildStore          Store
	basePath            string
	loggerFactory       logs.LoggerFactory
	image               string
	defaultProjectImage string
	defaultProjectUser  string
}

type BuilderFactoryConfig struct {
	Image               string
	ContainerRegistry   containerregistry.ContainerRegistry
	BuildStore          Store
	BuildImageNamespace string // Namespace to be used when tagging and pushing the build image
	LoggerFactory       logs.LoggerFactory
	DefaultProjectImage string
	DefaultProjectUser  string
	BasePath            string
}

func NewBuilderFactory(config BuilderFactoryConfig) IBuilderFactory {
	return &BuilderFactory{
		image:               config.Image,
		containerRegistry:   config.ContainerRegistry,
		buildImageNamespace: config.BuildImageNamespace,
		buildStore:          config.BuildStore,
		loggerFactory:       config.LoggerFactory,
		defaultProjectImage: config.DefaultProjectImage,
		defaultProjectUser:  config.DefaultProjectUser,
		basePath:            config.BasePath,
	}
}

func (f *BuilderFactory) Create(build Build) (IBuilder, error) {
	// TODO: Implement factory logic after adding prebuilds and other builder types
	return f.newDevcontainerBuilder(build)
}

func (f *BuilderFactory) CheckExistingBuild(b Build) (*Build, error) {
	build, err := f.buildStore.Find(&Filter{
		Hash: &b.Hash,
	})
	if err != nil {
		return nil, err
	}

	return build, nil
}

func (f *BuilderFactory) newDevcontainerBuilder(build Build) (*DevcontainerBuilder, error) {
	builderDockerPort, err := ports.GetAvailableEphemeralPort()
	if err != nil {
		return nil, err
	}

	return &DevcontainerBuilder{
		Builder: &Builder{
			hash:                build.Hash,
			projectDir:          filepath.Join(f.basePath, build.Hash, "project"),
			image:               f.image,
			containerRegistry:   f.containerRegistry,
			buildImageNamespace: f.buildImageNamespace,
			buildStore:          f.buildStore,
			loggerFactory:       f.loggerFactory,
			defaultProjectImage: f.defaultProjectImage,
			defaultProjectUser:  f.defaultProjectUser,
		},
		builderDockerPort: builderDockerPort,
	}, nil
}
