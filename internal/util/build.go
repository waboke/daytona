// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/daytonaio/daytona/pkg/workspace/project/buildconfig"
)

// GetBuildHash returns a SHA-256 hash of the build's configuration, repository SHA and environment variables.
func GetBuildHash(buildConfig *buildconfig.BuildConfig, sha string, envVars map[string]string) (string, error) {
	var buildJson []byte
	var err error

	if buildConfig != nil && buildConfig.Devcontainer != nil {
		buildJson, err = json.Marshal(buildConfig.Devcontainer)
		if err != nil {
			return "", err
		}
	}

	envVarsJson, err := json.Marshal(envVars)
	if err != nil {
		return "", err
	}

	data := string(buildJson) + sha + string(envVarsJson)
	hash := sha256.Sum256([]byte(data))
	hashStr := hex.EncodeToString(hash[:])

	return hashStr, nil
}
