// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const WINDOWS_USER_HOME_ENV = "USERPROFILE"

var userHomeDirectory string
var wslWindowsHomeDir string

func ensureSshFilesLinked(userHomeDir string) error {
	// Make sure ~/.ssh/config file exists if not create it
	sshDir := filepath.Join(userHomeDir, ".ssh")
	configPath := filepath.Join(sshDir, "config")

	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(sshDir, 0700)
		if err != nil {
			return err
		}
		err = os.WriteFile(configPath, []byte{}, 0600)
		if err != nil {
			return err
		}
	}

	// Make sure daytona_config file exists
	daytonaConfigPath := filepath.Join(sshDir, "daytona_config")

	_, err = os.Stat(daytonaConfigPath)
	if os.IsNotExist(err) {
		err := os.WriteFile(daytonaConfigPath, []byte{}, 0600)
		if err != nil {
			return err
		}
	}

	// Make sure daytona_config is included
	configFile := filepath.Join(sshDir, "config")
	_, err = os.Stat(configFile)
	if os.IsNotExist(err) {
		err := os.WriteFile(configFile, []byte("Include daytona_config\n\n"), 0600)
		if err != nil {
			return err
		}
	} else {
		content, err := os.ReadFile(configFile)
		if err != nil {
			return err
		}

		newContent := strings.ReplaceAll(string(content), "Include daytona_config\n\n", "")
		newContent = strings.ReplaceAll(string(newContent), "Include daytona_config\n", "")
		newContent = strings.ReplaceAll(string(newContent), "Include daytona_config", "")
		newContent = "Include daytona_config\n\n" + newContent
		err = os.WriteFile(configFile, []byte(newContent), 0600)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnlinkSshFiles() error {
	sshDirPath := filepath.Join(userHomeDirectory, ".ssh")
	sshConfigPath := filepath.Join(sshDirPath, "config")
	daytonaConfigPath := filepath.Join(sshDirPath, "daytona_config")

	// Remove the include line from the config file
	_, err := os.Stat(sshConfigPath)
	if os.IsExist(err) {
		content, err := os.ReadFile(sshConfigPath)
		if err != nil {
			return err
		}

		newContent := strings.ReplaceAll(string(content), "Include daytona_config\n\n", "")
		newContent = strings.ReplaceAll(string(newContent), "Include daytona_config", "")
		err = os.WriteFile(sshConfigPath, []byte(newContent), 0600)
		if err != nil {
			return err
		}
	}

	// Remove the daytona_config file
	_, err = os.Stat(daytonaConfigPath)
	if os.IsExist(err) {
		err = os.Remove(daytonaConfigPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// Add ssh entry

func generateSshConfigEntry(profileId, workspaceId, projectName, knownHostsPath string, addWslPath bool) (string, error) {
	daytonaPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	tab := "\t"
	projectHostname := GetProjectHostname(profileId, workspaceId, projectName)

	if addWslPath {
		daytonaPath = fmt.Sprintf("\"C:\\Windows\\system32\\wsl.exe\" \"%s\"", daytonaPath)
	}

	config := fmt.Sprintf("Host %s\n"+
		tab+"User daytona\n"+
		tab+"StrictHostKeyChecking no\n"+
		tab+"UserKnownHostsFile %s\n"+
		tab+"ProxyCommand %s ssh-proxy %s %s %s\n"+
		tab+"ForwardAgent yes\n\n",
		projectHostname, knownHostsPath, daytonaPath, profileId, workspaceId, projectName)

	return config, nil
}

func EnsureSshConfigEntryAdded(profileId, workspaceName, projectName string) error {
	// Check if being run on WSL
	// if yes, do this for windows
	cmd := exec.Command("uname", "-a")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	if strings.Contains(string(output), "WSL") {
		getWindowsPathCmd := exec.Command("bash", "-c", `echo /mnt/c/Users/$(cmd.exe /C "echo %USERNAME%" | tr -d '\r')`)

		windowsPathOutput, err := getWindowsPathCmd.Output()
		if err != nil {
			return err
		}
		wslWindowsHomeDir = strings.TrimSpace(string(windowsPathOutput))

		err = ensureSshConfigEntryAdded(wslWindowsHomeDir, profileId, workspaceName, projectName, true)
		if err != nil {
			return err
		}
	}

	return ensureSshConfigEntryAdded(userHomeDirectory, profileId, workspaceName, projectName, false)
}

func ensureSshConfigEntryAdded(userHomeDir, profileId, workspaceName, projectName string, addWslPath bool) error {
	err := ensureSshFilesLinked(userHomeDir)
	if err != nil {
		return err
	}

	knownHostsFile := "/dev/null"
	if runtime.GOOS == "windows" {
		knownHostsFile = "NUL"
	}

	data, err := generateSshConfigEntry(profileId, workspaceName, projectName, knownHostsFile, addWslPath)
	if err != nil {
		return err
	}

	sshDir := filepath.Join(userHomeDir, ".ssh")
	configPath := filepath.Join(sshDir, "daytona_config")

	// Read existing content from the file
	existingContent, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if strings.Contains(string(existingContent), data) {
		return nil
	}

	// Combine the new data with existing content
	newData := data + string(existingContent)

	// Open the file for writing
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(newData)
	if err != nil {
		return err
	}

	return nil
}

func RemoveWorkspaceSshEntries(profileId, workspaceId string) error {
	sshDir := filepath.Join(userHomeDirectory, ".ssh")
	configPath := filepath.Join(sshDir, "daytona_config")

	// Read existing content from the file
	existingContent, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return nil
	}

	regex := regexp.MustCompile(fmt.Sprintf(`Host %s-%s-\w+\n(?:\t.*\n?)*`, profileId, workspaceId))
	newContent := regex.ReplaceAllString(string(existingContent), "")

	newContent = strings.Trim(newContent, "\n")

	// Open the file for writing
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(newContent)
	if err != nil {
		return err
	}

	return nil
}

func GetProjectHostname(profileId, workspaceId, projectName string) string {
	return fmt.Sprintf("%s-%s-%s", profileId, workspaceId, projectName)
}

func init() {
	if runtime.GOOS == "windows" {
		userHomeDirectory = os.Getenv(WINDOWS_USER_HOME_ENV)
	} else {
		userHomeDirectory = os.Getenv("HOME")
	}
}
