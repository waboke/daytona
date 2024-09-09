// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"github.com/daytonaio/daytona/pkg/workspace/template"

	"gorm.io/gorm"

	. "github.com/daytonaio/daytona/pkg/db/dto"
	"github.com/daytonaio/daytona/pkg/provider"
)

type TemplateStore struct {
	db *gorm.DB
}

func NewTemplateStore(db *gorm.DB) (*TemplateStore, error) {
	err := db.AutoMigrate(&TemplateDTO{})
	if err != nil {
		return nil, err
	}

	return &TemplateStore{db: db}, nil
}

func (s *TemplateStore) List() ([]*template.Template, error) {
	templatesDTOs := []TemplateDTO{}
	tx := s.db.Find(&templatesDTOs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	templates := []*template.Template{}
	for _, templateDTO := range templatesDTOs {
		templates = append(templates, ToTemplate(templateDTO))
	}

	return templates, nil
}

func (s *TemplateStore) Find(name string) (*template.Template, error) {
	templateDTO := TemplateDTO{}
	tx := s.db.Where("name = ?", name).First(&templateDTO)
	if tx.Error != nil {
		if IsRecordNotFound(tx.Error) {
			return nil, provider.ErrTargetNotFound
		}
		return nil, tx.Error
	}

	return ToTemplate(templateDTO), nil
}

func (s *TemplateStore) Save(template *template.Template) error {
	tx := s.db.Save(ToTemplateDTO(template))
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *TemplateStore) Delete(template *template.Template) error {
	tx := s.db.Delete(ToTemplateDTO(template))
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return provider.ErrTargetNotFound
	}

	return nil
}
