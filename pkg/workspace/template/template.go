// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package template

type Template struct {
	Name           string   `json:"name" validate:"required"`
	ProjectConfigs []string `json:"projectConfigs" validate:"required"`
} // @name Template
