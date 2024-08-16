// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package builds_test

import (
	"testing"

	build_internal "github.com/daytonaio/daytona/internal/testing/build"
	"github.com/daytonaio/daytona/pkg/build"
	"github.com/daytonaio/daytona/pkg/gitprovider"
	"github.com/daytonaio/daytona/pkg/server/builds"
	"github.com/daytonaio/daytona/pkg/workspace/project/buildconfig"
	"github.com/stretchr/testify/suite"
)

var build1Image = "image1"
var build1User = "user1"

var build1 *build.Build = &build.Build{
	Id:          "id1",
	Image:       build1Image,
	User:        build1User,
	BuildConfig: &buildconfig.BuildConfig{},
	Repository: &gitprovider.GitRepository{
		Sha: "sha1",
	},
	State: build.BuildStatePublished,
}

var build2 *build.Build = &build.Build{
	Id:          "id2",
	Image:       "image2",
	User:        "user2",
	BuildConfig: nil,
	Repository: &gitprovider.GitRepository{
		Sha: "sha2",
	},
	State: build.BuildStatePublished,
}

var build3 *build.Build = &build.Build{
	Id:          "id3",
	Image:       "image3",
	User:        "user3",
	BuildConfig: nil,
	Repository: &gitprovider.GitRepository{
		Sha: "sha3",
	},
	State: build.BuildStatePending,
}

var build4 *build.Build = &build.Build{
	Id:          "id4",
	Image:       "image4",
	User:        "user4",
	BuildConfig: nil,
	Repository: &gitprovider.GitRepository{
		Sha: "sha4",
	},
	State: build.BuildStatePending,
}

var expectedBuilds []*build.Build
var expectedFilteredBuilds []*build.Build

var expectedBuildsMap map[string]*build.Build
var expectedFilteredBuildsMap map[string]*build.Build

type BuildServiceTestSuite struct {
	suite.Suite
	buildService builds.IBuildService
	buildStore   build.Store
}

func NewBuildServiceTestSuite() *BuildServiceTestSuite {
	return &BuildServiceTestSuite{}
}

func (s *BuildServiceTestSuite) SetupTest() {
	expectedBuilds = []*build.Build{
		build1, build2, build3,
	}

	expectedBuildsMap = map[string]*build.Build{
		build1.Id: build1,
		build2.Id: build2,
		build3.Id: build3,
	}

	expectedFilteredBuilds = []*build.Build{
		build1, build2,
	}

	expectedFilteredBuildsMap = map[string]*build.Build{
		build1.Id: build1,
		build2.Id: build2,
	}

	s.buildStore = build_internal.NewInMemoryBuildStore()
	s.buildService = builds.NewBuildService(builds.BuildServiceConfig{
		BuildStore: s.buildStore,
	})

	for _, b := range expectedBuilds {
		_ = s.buildStore.Save(b)
	}
}

func TestBuildService(t *testing.T) {
	suite.Run(t, NewBuildServiceTestSuite())
}

func (s *BuildServiceTestSuite) TestList() {
	require := s.Require()

	builds, err := s.buildService.List(nil)
	require.Nil(err)
	require.ElementsMatch(expectedBuilds, builds)
}

func (s *BuildServiceTestSuite) TestFind() {
	require := s.Require()

	build, err := s.buildService.Find(&build.Filter{
		Id: &build1.Id,
	})
	require.Nil(err)
	require.Equal(build1, build)
}

func (s *BuildServiceTestSuite) TestSave() {
	expectedBuilds = append(expectedBuilds, build4)

	require := s.Require()

	err := s.buildService.Create(build4)
	require.Nil(err)

	builds, err := s.buildService.List(nil)
	require.Nil(err)
	require.ElementsMatch(expectedBuilds, builds)
}

func (s *BuildServiceTestSuite) TestDelete() {
	expectedBuilds = expectedBuilds[:2]

	require := s.Require()

	err := s.buildService.Delete(build3.Id)
	require.Nil(err)

	builds, err := s.buildService.List(nil)
	require.Nil(err)
	require.ElementsMatch(expectedBuilds, builds)
}
