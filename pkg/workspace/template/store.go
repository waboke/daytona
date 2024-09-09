// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import "errors"

type Store interface {
	List() ([]*Template, error)
	Find(name string) (*Template, error)
	Save(template *Template) error
	Delete(template *Template) error
}

var (
	ErrTemplateNotFound = errors.New("template not found")
)

func IsTemplateNotFound(err error) bool {
	return err.Error() == ErrTemplateNotFound.Error()
}
