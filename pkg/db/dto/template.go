// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package dto

import "github.com/daytonaio/daytona/pkg/workspace/template"

type TemplateDTO struct {
	Name           string   `json:"name" gorm:"primaryKey"`
	ProjectConfigs []string `json:"projectConfigs" gorm:"serializer:json"`
}

func ToTemplateDTO(template *template.Template) TemplateDTO {
	return TemplateDTO{
		Name:           template.Name,
		ProjectConfigs: template.ProjectConfigs,
	}
}

func ToTemplate(templateDTO TemplateDTO) *template.Template {
	return &template.Template{
		Name:           templateDTO.Name,
		ProjectConfigs: templateDTO.ProjectConfigs,
	}
}
