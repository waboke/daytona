// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package builds

import (
	"github.com/daytonaio/daytona/internal/util"
	"github.com/daytonaio/daytona/pkg/build"
	"github.com/docker/docker/pkg/stringid"
)

type IBuildService interface {
	Create(*build.Build) error
	Find(filter *build.Filter) (*build.Build, error)
	List(filter *build.Filter) ([]*build.Build, error)
	Delete(id string) error
}

type BuildServiceConfig struct {
	BuildStore build.Store
}

type BuildService struct {
	buildStore build.Store
}

func NewBuildService(config BuildServiceConfig) IBuildService {
	return &BuildService{
		buildStore: config.BuildStore,
	}
}

func (s *BuildService) Create(b *build.Build) error {
	id := stringid.GenerateRandomID()
	id = stringid.TruncateID(id)

	b.Id = id

	hash, err := util.GetBuildHash(b.BuildConfig, b.Repository.Sha, b.EnvVars)
	if err != nil {
		return err
	}

	b.Hash = hash
	b.State = build.BuildStatePending

	return s.buildStore.Save(b)
}

func (s *BuildService) Find(filter *build.Filter) (*build.Build, error) {
	result, err := s.buildStore.Find(filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *BuildService) List(filter *build.Filter) ([]*build.Build, error) {
	result, err := s.buildStore.List(filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *BuildService) Delete(id string) error {
	return s.buildStore.Delete(id)
}
