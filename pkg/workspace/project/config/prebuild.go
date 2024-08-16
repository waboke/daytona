// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/docker/docker/pkg/stringid"
)

// PrebuildConfig holds configuration for the prebuild process
type PrebuildConfig struct {
	Id             string   `json:"id" validate:"required"`
	Branch         string   `json:"branch" validate:"required"`
	CommitInterval int      `json:"commitInterval" validate:"required"`
	TriggerFiles   []string `json:"triggerFiles" validate:"required"`
} // @name PrebuildConfig

func (p *PrebuildConfig) GenerateId() error {
	id := stringid.GenerateRandomID()
	id = stringid.TruncateID(id)

	p.Id = id
	return nil
}
