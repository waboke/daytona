// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package template_test

// import (
// 	"testing"

// 	git_provider_mock "github.com/daytonaio/daytona/internal/testing/gitprovider/mocks"
// 	projectconfig_internal "github.com/daytonaio/daytona/internal/testing/server/projectconfig"
// 	"github.com/daytonaio/daytona/internal/testing/server/workspaces/mocks"
// 	"github.com/daytonaio/daytona/internal/util"
// 	"github.com/daytonaio/daytona/pkg/server/projectconfig"
// 	"github.com/daytonaio/daytona/pkg/workspace/project/config"
// 	"github.com/stretchr/testify/suite"
// )

// var template1Image = "image1"
// var template1User = "user1"

// var template1 *config.Template = &config.Template{
// 	Name:          "pc1",
// 	Image:         template1Image,
// 	User:          template1User,
// 	BuildConfig:   nil,
// 	RepositoryUrl: repository1.Url,
// 	IsDefault:     true,
// 	Prebuilds: []*config.PrebuildConfig{
// 		prebuild1,
// 		prebuild2,
// 	},
// }

// var template2 *config.Template = &config.Template{
// 	Name:          "pc2",
// 	Image:         "image2",
// 	User:          "user2",
// 	BuildConfig:   nil,
// 	RepositoryUrl: "https://github.com/daytonaio/daytona.git",
// }

// var template3 *config.Template = &config.Template{
// 	Name:          "pc3",
// 	Image:         "image3",
// 	User:          "user3",
// 	BuildConfig:   nil,
// 	RepositoryUrl: "https://github.com/daytonaio/daytona3.git",
// }

// var template4 *config.Template = &config.Template{
// 	Name:          "pc4",
// 	Image:         "image4",
// 	User:          "user4",
// 	BuildConfig:   nil,
// 	RepositoryUrl: "https://github.com/daytonaio/daytona4.git",
// }

// var expectedTemplates []*config.Template
// var expectedFilteredTemplates []*config.Template

// var expectedTemplatesMap map[string]*config.Template
// var expectedFilteredTemplatesMap map[string]*config.Template

// type TemplateServiceTestSuite struct {
// 	suite.Suite
// 	templateService projectconfig.ITemplateService
// 	templateStore   config.Store
// 	gitProviderService   mocks.MockGitProviderService
// 	buildService         mocks.MockBuildService
// 	gitProvider          git_provider_mock.MockGitProvider
// }

// func NewConfigServiceTestSuite() *TemplateServiceTestSuite {
// 	return &TemplateServiceTestSuite{}
// }

// func (s *TemplateServiceTestSuite) SetupTest() {
// 	expectedTemplates = []*config.Template{
// 		template1, template2, template3,
// 	}

// 	expectedPrebuilds = []*config.PrebuildConfig{
// 		prebuild1, prebuild2,
// 	}

// 	expectedTemplatesMap = map[string]*config.Template{
// 		template1.Name: template1,
// 		template2.Name: template2,
// 		template3.Name: template3,
// 	}

// 	expectedPrebuildsMap = map[string]*config.PrebuildConfig{
// 		prebuild1.Id: prebuild1,
// 		prebuild2.Id: prebuild2,
// 	}

// 	expectedFilteredTemplates = []*config.Template{
// 		template1, template2,
// 	}

// 	expectedFilteredPrebuilds = []*config.PrebuildConfig{
// 		prebuild1,
// 	}

// 	expectedFilteredTemplatesMap = map[string]*config.Template{
// 		template1.Name: template1,
// 		template2.Name: template2,
// 	}

// 	expectedFilteredPrebuildsMap = map[string]*config.PrebuildConfig{
// 		prebuild1.Id: prebuild1,
// 	}

// 	s.templateStore = projectconfig_internal.NewInMemoryTemplateStore()
// 	s.templateService = projectconfig.NewTemplateService(projectconfig.TemplateServiceConfig{
// 		ConfigStore:        s.templateStore,
// 		GitProviderService: &s.gitProviderService,
// 		BuildService:       &s.buildService,
// 	})

// 	for _, pc := range expectedTemplates {
// 		_ = s.templateStore.Save(pc)
// 	}
// }

// func TestTemplateService(t *testing.T) {
// 	suite.Run(t, NewConfigServiceTestSuite())
// }

// func (s *TemplateServiceTestSuite) TestList() {
// 	require := s.Require()

// 	templates, err := s.templateService.List(nil)
// 	require.Nil(err)
// 	require.ElementsMatch(expectedTemplates, templates)
// }

// func (s *TemplateServiceTestSuite) TestFind() {
// 	require := s.Require()

// 	template, err := s.templateService.Find(&config.TemplateFilter{
// 		Name: &template1.Name,
// 	})
// 	require.Nil(err)
// 	require.Equal(template1, template)
// }
// func (s *TemplateServiceTestSuite) TestSetDefault() {
// 	require := s.Require()

// 	err := s.templateService.SetDefault(template2.Name)
// 	require.Nil(err)

// 	template, err := s.templateService.Find(&config.TemplateFilter{
// 		Url:     util.Pointer(template1.RepositoryUrl),
// 		Default: util.Pointer(true),
// 	})
// 	require.Nil(err)

// 	require.Equal(template2, template)
// }

// func (s *TemplateServiceTestSuite) TestSave() {
// 	expectedTemplates = append(expectedTemplates, template4)

// 	require := s.Require()

// 	err := s.templateService.Save(template4)
// 	require.Nil(err)

// 	templates, err := s.templateService.List(nil)
// 	require.Nil(err)
// 	require.ElementsMatch(expectedTemplates, templates)
// }

// func (s *TemplateServiceTestSuite) TestDelete() {
// 	expectedTemplates = expectedTemplates[:2]

// 	require := s.Require()

// 	err := s.templateService.Delete(template3.Name, false)
// 	require.Nil(err)

// 	templates, errs := s.templateService.List(nil)
// 	require.Nil(errs)
// 	require.ElementsMatch(expectedTemplates, templates)
// }

// func (s *TemplateServiceTestSuite) AfterTest(_, _ string) {
// 	s.gitProviderService.AssertExpectations(s.T())
// 	s.gitProviderService.ExpectedCalls = nil
// 	s.buildService.AssertExpectations(s.T())
// 	s.buildService.ExpectedCalls = nil
// 	s.gitProvider.AssertExpectations(s.T())
// 	s.gitProvider.ExpectedCalls = nil
// }
