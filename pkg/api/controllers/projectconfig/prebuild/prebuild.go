// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package prebuild

import (
	"fmt"
	"net/http"

	"github.com/daytonaio/daytona/pkg/server"
	"github.com/daytonaio/daytona/pkg/server/projectconfig/dto"
	"github.com/daytonaio/daytona/pkg/workspace/project/config"
	"github.com/gin-gonic/gin"
)

// GetPrebuild godoc
//
//	@Tags			prebuild
//	@Summary		Get prebuild
//	@Description	Get prebuild
//	@Accept			json
//	@Param			prebuildId	path		string	true	"Prebuild ID"
//	@Success		200			{object}	PrebuildDTO
//	@Router			/project-config/prebuild/{prebuildId} [get]
//
//	@id				GetPrebuild
func GetPrebuild(ctx *gin.Context) {
	prebuildId := ctx.Param("prebuildId")

	server := server.GetInstance(nil)
	res, err := server.ProjectConfigService.FindPrebuild(nil, &config.PrebuildFilter{
		Id: &prebuildId,
	})
	if err != nil {
		if config.IsPrebuildNotFound(err) {
			ctx.JSON(200, &config.PrebuildConfig{})
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get prebuild: %s", err.Error()))
		return
	}

	ctx.JSON(200, res)
}

// SetPrebuild godoc

// @Tags			prebuild
// @Summary		Set prebuild
// @Description	Set prebuild
// @Accept			json
// @Param			prebuild	body	CreatePrebuildDTO	true	"Prebuild"
// @Success		201
// @Router			/project-config/prebuild [put]
//
// @id				SetPrebuild
func SetPrebuild(ctx *gin.Context) {
	var req dto.CreatePrebuildDTO
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid request body: %s", err.Error()))
		return
	}

	server := server.GetInstance(nil)
	_, err = server.ProjectConfigService.SetPrebuild(req)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to set prebuild: %s", err.Error()))
		return
	}

	ctx.Status(201)
}

// ListPrebuilds godoc

// @Tags			prebuild
// @Summary		List prebuilds
// @Description	List prebuilds
// @Accept			json
// @Param			configName	path	string	true	"Config name"
//
// @Success		200			{array}	PrebuildDTO
//
// @Router			/project-config/prebuild [get]
//
// @id				ListPrebuilds
func ListPrebuilds(ctx *gin.Context) {
	configName := ctx.Param("configName")

	var projectConfigFilter *config.Filter

	if configName != "" {
		projectConfigFilter = &config.Filter{
			Name: &configName,
		}
	}

	server := server.GetInstance(nil)
	res, err := server.ProjectConfigService.ListPrebuilds(projectConfigFilter, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get prebuilds: %s", err.Error()))
		return
	}

	ctx.JSON(200, res)
}

// DeletePrebuild godoc
//
//	@Tags			prebuild
//	@Summary		Delete prebuild
//	@Description	Delete prebuild
//	@Accept			json
//	@Param			projectConfigName	path	string	true	"Project config name"
//	@Param			prebuildId			path	string	true	"Prebuild ID"
//	@Success		204
//	@Router			/project-config/prebuild/{projectConfigName}/{prebuildId} [delete]
//
//	@id				DeletePrebuild
func DeletePrebuild(ctx *gin.Context) {
	projectConfigName := ctx.Param("projectConfigName")
	prebuildId := ctx.Param("prebuildId")

	server := server.GetInstance(nil)
	err := server.ProjectConfigService.DeletePrebuild(projectConfigName, prebuildId)
	if err != nil {
		if config.IsPrebuildNotFound(err) {
			ctx.Status(204)
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get prebuild: %s", err.Error()))
		return
	}

	ctx.Status(204)
}
