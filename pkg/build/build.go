// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package build

import (
	"time"

	"github.com/daytonaio/daytona/pkg/gitprovider"
	"github.com/daytonaio/daytona/pkg/workspace/project/buildconfig"
)

type BuildState string

const (
	BuildStatePending   BuildState = "pending"
	BuildStateRunning   BuildState = "running"
	BuildStateError     BuildState = "error"
	BuildStateSuccess   BuildState = "success"
	BuildStatePublished BuildState = "published"
)

type Build struct {
	Id          string                     `json:"id" validate:"required"`
	Hash        string                     `json:"hash" validate:"required"`
	State       BuildState                 `json:"state" validate:"required"`
	Image       string                     `json:"image" validate:"required"`
	User        string                     `json:"user" validate:"required"`
	BuildConfig *buildconfig.BuildConfig   `json:"buildConfig" validate:"optional"`
	Repository  *gitprovider.GitRepository `json:"repository" validate:"optional"`
	EnvVars     map[string]string          `json:"envVars" validate:"required"`
	PrebuildId  string                     `json:"prebuildId" validate:"required"`
	CreatedAt   time.Time                  `json:"createdAt" validate:"required"`
	UpdatedAt   time.Time                  `json:"updatedAt" validate:"required"`
} // @name Build

func FetchCachedBuild(build *Build, builds []*Build) *buildconfig.CachedBuild {
	var cachedBuild *Build

	for _, existingBuild := range builds {
		if existingBuild.PrebuildId != build.PrebuildId ||
			existingBuild.State != BuildStateSuccess && existingBuild.State != BuildStatePublished {
			continue
		}
		if cachedBuild == nil {
			cachedBuild = existingBuild
			continue
		} else {
			if existingBuild.CreatedAt.After(cachedBuild.CreatedAt) {
				cachedBuild = existingBuild
			}
		}
	}

	if cachedBuild != nil {
		return &buildconfig.CachedBuild{
			Image: cachedBuild.Image,
			User:  cachedBuild.User,
		}
	}

	return nil
}
