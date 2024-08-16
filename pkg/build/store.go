// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package build

import (
	"errors"
)

type Filter struct {
	Id          *string
	Hash        *string
	States      []*BuildState
	PrebuildIds *[]string
	GetNewest   *bool
}

type Store interface {
	Find(filter *Filter) (*Build, error)
	List(filter *Filter) ([]*Build, error)
	Save(build *Build) error
	Delete(id string) error
}

var (
	ErrBuildNotFound = errors.New("build not found")
)

func IsBuildNotFound(err error) bool {
	return err.Error() == ErrBuildNotFound.Error()
}
