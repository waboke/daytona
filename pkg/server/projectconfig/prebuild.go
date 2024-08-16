// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package projectconfig

import (
	"fmt"

	"github.com/daytonaio/daytona/internal/util"
	"github.com/daytonaio/daytona/pkg/build"
	"github.com/daytonaio/daytona/pkg/gitprovider"
	"github.com/daytonaio/daytona/pkg/server/projectconfig/dto"
	"github.com/daytonaio/daytona/pkg/workspace/project/config"
)

func (s *ProjectConfigService) SetPrebuild(createPrebuildDto dto.CreatePrebuildDTO) (*dto.PrebuildDTO, error) {
	projectConfig, err := s.Find(&config.Filter{
		Name: &createPrebuildDto.ProjectConfigName,
	})
	if err != nil {
		return nil, err
	}

	existingPrebuild, err := projectConfig.FindPrebuild(&config.PrebuildFilter{
		Branch: &createPrebuildDto.Branch,
	})
	if err != nil {
		return nil, err
	}

	if existingPrebuild != nil {
		return nil, fmt.Errorf("prebuild for the specified project config and branch already exists")
	}

	gitProvider, gitProviderId, err := s.gitProviderService.GetGitProviderForUrl(projectConfig.RepositoryUrl)
	if err != nil {
		return nil, err
	}

	repository, err := gitProvider.GetRepositoryFromUrl(projectConfig.RepositoryUrl)
	if err != nil {
		return nil, err
	}

	if projectConfig.WebhookId == nil {
		webhookId, err := s.gitProviderService.RegisterPrebuildWebhook(gitProviderId, repository, s.prebuildWebhookEndpoint)
		if err != nil {
			return nil, err
		}

		projectConfig.WebhookId = &webhookId

		err = s.configStore.Save(projectConfig)
		if err != nil {
			return nil, err
		}
	}

	prebuild := &config.PrebuildConfig{
		Branch:         createPrebuildDto.Branch,
		CommitInterval: createPrebuildDto.CommitInterval,
		TriggerFiles:   createPrebuildDto.TriggerFiles,
	}

	err = prebuild.GenerateId()
	if err != nil {
		return nil, err
	}

	// TODO: handle webhook registration removal if prebuild creation fails
	err = projectConfig.SetPrebuild(prebuild)
	if err != nil {
		return nil, err
	}

	err = s.configStore.Save(projectConfig)
	if err != nil {
		return nil, err
	}

	if createPrebuildDto.RunAtInit != nil && *createPrebuildDto.RunAtInit {
		err = s.buildService.Create(&build.Build{
			Image:       projectConfig.Image,
			User:        projectConfig.User,
			BuildConfig: projectConfig.BuildConfig,
			Repository:  repository,
			EnvVars:     projectConfig.EnvVars,
			PrebuildId:  prebuild.Id,
		})
		if err != nil {
			return nil, err
		}
	}

	return &dto.PrebuildDTO{
		Id:                prebuild.Id,
		ProjectConfigName: projectConfig.Name,
		Branch:            prebuild.Branch,
		CommitInterval:    prebuild.CommitInterval,
		TriggerFiles:      prebuild.TriggerFiles,
	}, nil
}

func (s *ProjectConfigService) DeletePrebuild(projectConfigName string, id string) error {
	projectConfig, err := s.Find(&config.Filter{
		Name: &projectConfigName,
	})
	if err != nil {
		return err
	}

	// If this is the last prebuild, unregister the Git provider webhook
	if len(projectConfig.Prebuilds) == 1 {
		gitProvider, gitProviderId, err := s.gitProviderService.GetGitProviderForUrl(projectConfig.RepositoryUrl)
		if err != nil {
			return err
		}

		repository, err := gitProvider.GetRepositoryFromUrl(projectConfig.RepositoryUrl)
		if err != nil {
			return err
		}

		if projectConfig.WebhookId != nil {
			err = s.gitProviderService.UnregisterPrebuildWebhook(gitProviderId, repository, *projectConfig.WebhookId)
			if err != nil {
				return err
			}
		}
	}

	err = projectConfig.RemovePrebuild(id)
	if err != nil {
		return err
	}

	return s.configStore.Save(projectConfig)
}

func (s *ProjectConfigService) FindPrebuild(projectConfigFilter *config.Filter, prebuildFilter *config.PrebuildFilter) (*dto.PrebuildDTO, error) {
	pc, err := s.configStore.Find(projectConfigFilter)
	if err != nil {
		return nil, config.ErrProjectConfigNotFound
	}

	prebuild, err := pc.FindPrebuild(prebuildFilter)
	if err != nil {
		return nil, err
	}

	return &dto.PrebuildDTO{
		Id:                prebuild.Id,
		ProjectConfigName: pc.Name,
		Branch:            prebuild.Branch,
		CommitInterval:    prebuild.CommitInterval,
		TriggerFiles:      prebuild.TriggerFiles,
	}, nil
}

func (s *ProjectConfigService) ListPrebuilds(projectConfigFilter *config.Filter, prebuildFilter *config.PrebuildFilter) ([]*dto.PrebuildDTO, error) {
	var result []*dto.PrebuildDTO
	pcs, err := s.configStore.List(projectConfigFilter)
	if err != nil {
		return nil, config.ErrProjectConfigNotFound
	}

	for _, pc := range pcs {
		for _, prebuild := range pc.Prebuilds {
			result = append(result, &dto.PrebuildDTO{
				Id:                prebuild.Id,
				ProjectConfigName: pc.Name,
				Branch:            prebuild.Branch,
				CommitInterval:    prebuild.CommitInterval,
				TriggerFiles:      prebuild.TriggerFiles,
			})
		}
	}

	return result, nil
}

func (s *ProjectConfigService) ProcessGitEvent(data gitprovider.GitEventData) error {
	var buildsToTrigger []build.Build

	projectConfigs, err := s.List(&config.Filter{
		Url: &data.Url,
	})
	if err != nil {
		return err
	}

	gitProvider, _, err := s.gitProviderService.GetGitProviderForUrl(data.Url)
	if err != nil {
		return err
	}

	repo, err := gitProvider.GetRepositoryFromUrl(data.Url)
	if err != nil {
		return err
	}

	for _, projectConfig := range projectConfigs {
		prebuild, err := projectConfig.FindPrebuild(&config.PrebuildFilter{
			Branch: &data.Branch,
		})
		if err != nil {
			return err
		}

		if prebuild == nil {
			continue
		}

		// Check if the commit's affected files and prebuild config's trigger files have any overlap
		if len(prebuild.TriggerFiles) > 0 {
			if SlicesHaveCommonEntry(prebuild.TriggerFiles, data.AffectedFiles) {
				buildsToTrigger = append(buildsToTrigger, build.Build{
					Image:       projectConfig.Image,
					User:        projectConfig.User,
					BuildConfig: projectConfig.BuildConfig,
					Repository:  repo,
					EnvVars:     projectConfig.EnvVars,
					PrebuildId:  prebuild.Id,
				})
				continue
			}
		}

		newestBuild, err := s.buildService.Find(&build.Filter{
			PrebuildIds: &[]string{prebuild.Id},
			GetNewest:   util.Pointer(true),
		})
		if err != nil {
			return err
		}

		commitsRange, err := gitProvider.GetCommitsRange(repo, data.Owner, newestBuild.Repository.Sha, data.Sha)
		if err != nil {
			return err
		}

		// Check if the commit interval has been reached
		if commitsRange > prebuild.CommitInterval {
			buildsToTrigger = append(buildsToTrigger, build.Build{
				Image:       projectConfig.Image,
				User:        projectConfig.User,
				BuildConfig: projectConfig.BuildConfig,
				Repository:  repo,
				EnvVars:     projectConfig.EnvVars,
				PrebuildId:  prebuild.Id,
			})
		}
	}

	for _, build := range buildsToTrigger {
		err = s.buildService.Create(&build)
		if err != nil {
			return err
		}
	}

	return nil
}

func SlicesHaveCommonEntry(slice1, slice2 []string) bool {
	entryMap := make(map[string]bool)

	for _, entry := range slice1 {
		entryMap[entry] = true
	}

	for _, entry := range slice2 {
		if entryMap[entry] {
			return true
		}
	}

	return false
}

// // TODO: check if multiple same branch prebuilds for single project config is feasible
// // TODO: check if "branch commit interval" and "branch trigger files" should be separate prebuild configs
// func getCommonEntries(slice1, slice2 []string) []string {
// 	// Create a map to store the entries of the first slice
// 	entryMap := make(map[string]bool)
// 	// Slice to hold common entries
// 	var commonEntries []string

// 	// Populate the map with entries from slice1
// 	for _, entry := range slice1 {
// 		entryMap[entry] = true
// 	}

// 	// Iterate over slice2 and check if any entry exists in the map
// 	for _, entry := range slice2 {
// 		if entryMap[entry] {
// 			commonEntries = append(commonEntries, entry)
// 		}
// 	}

// 	return commonEntries
// }
