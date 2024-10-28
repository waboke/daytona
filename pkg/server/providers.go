// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/daytonaio/daytona/pkg/provider/manager"
	log "github.com/sirupsen/logrus"
)

func (s *Server) downloadDefaultProviders() error {
	manifest, err := s.ProviderManager.GetProvidersManifest()
	if err != nil {
		return err
	}

	defaultProviders := manifest.GetDefaultProviders()

	log.Info("Downloading default providers")
	for providerName, provider := range defaultProviders {
		lockFilePath := filepath.Join(s.config.ProvidersDir, providerName, manager.INITIAL_SETUP_LOCK_FILE_NAME)

		_, err := os.Stat(lockFilePath)
		if err == nil {
			continue
		}

		_, err = s.ProviderManager.DownloadProvider(context.Background(), provider.DownloadUrls, providerName)
		if err != nil {
			if !manager.IsProviderAlreadyDownloaded(err, providerName) {
				log.Error(err)
			}
			continue
		}
	}

	log.Info("Default providers downloaded")

	return nil
}

func (s *Server) registerProviders() error {
	log.Info("Registering providers")

	manifest, err := s.ProviderManager.GetProvidersManifest()
	if err != nil {
		return err
	}

	files, err := os.ReadDir(s.config.ProvidersDir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Info("No providers found")
			return nil
		}
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			pluginPath, err := s.getPluginPath(filepath.Join(s.config.ProvidersDir, file.Name()))
			if err != nil {
				log.Error(err)
				continue
			}

			if strings.HasSuffix(pluginPath, manager.INITIAL_SETUP_LOCK_FILE_NAME) {
				continue
			}

			err = s.ProviderManager.RegisterProvider(pluginPath)
			if err != nil {
				log.Error(err)
				continue
			}

			// Check for updates
			provider, err := s.ProviderManager.GetProvider(file.Name())
			if err != nil {
				log.Error(err)
				continue
			}

			info, err := (*provider).GetInfo()
			if err != nil {
				log.Error(err)
				continue
			}

			if manifest.HasUpdateAvailable(info.Name, info.Version) {
				log.Infof("Update available for %s. Update with `daytona provider update`.", info.Name)
			}
		}
	}

	log.Info("Providers registered")

	return nil
}

func (s *Server) lockInitialProviderSetups() error {
	providers := s.ProviderManager.GetProviders()
	for providerName := range providers {
		lockFilePath := filepath.Join(s.config.ProvidersDir, providerName, manager.INITIAL_SETUP_LOCK_FILE_NAME)

		_, err := os.Stat(lockFilePath)
		if err == nil {
			continue
		}

		file, err := os.Create(lockFilePath)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	return nil
}

func (s *Server) getPluginPath(dir string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() {
			return filepath.Join(dir, file.Name()), nil
		}
	}

	return "", errors.New("no plugin found in " + dir)
}
