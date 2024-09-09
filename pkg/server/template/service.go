// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"github.com/daytonaio/daytona/pkg/workspace/template"
)

type ITemplateService interface {
	Save(template *template.Template) error
	Find(name string) (*template.Template, error)
	List() ([]*template.Template, error)
	Delete(templateName string) error
}

type TemplateServiceConfig struct {
	TemplateStore template.Store
}

type TemplateService struct {
	templateStore template.Store
}

func NewTemplateService(config TemplateServiceConfig) ITemplateService {
	return &TemplateService{
		templateStore: config.TemplateStore,
	}
}

func (s *TemplateService) List() ([]*template.Template, error) {
	return s.templateStore.List()
}

func (s *TemplateService) Find(name string) (*template.Template, error) {
	return s.templateStore.Find(name)
}

func (s *TemplateService) Save(template *template.Template) error {
	return s.templateStore.Save(template)
}

func (s *TemplateService) Delete(name string) error {
	template, err := s.Find(name)
	if err != nil {
		return err
	}

	return s.Delete(template.Name)
}
