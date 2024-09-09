// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"fmt"
	"net/http"

	"github.com/daytonaio/daytona/pkg/server"
	"github.com/daytonaio/daytona/pkg/workspace/template"
	"github.com/gin-gonic/gin"
)

// GetTemplate godoc
//
//	@Tags			template
//	@Summary		Get template data
//	@Description	Get template data
//	@Accept			json
//	@Param			templateName	path		string	true	"Template name"
//	@Success		200				{object}	Template
//	@Router			/template/{templateName} [get]
//
//	@id				GetTemplate
func GetTemplate(ctx *gin.Context) {
	templateName := ctx.Param("configName")

	server := server.GetInstance(nil)

	template, err := server.TemplateService.Find(templateName)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get project config: %s", err.Error()))
		return
	}

	ctx.JSON(200, template)
}

// ListTemplates godoc
//
//	@Tags			template
//	@Summary		List templates
//	@Description	List templates
//	@Produce		json
//	@Success		200	{array}	Template
//	@Router			/template [get]
//
//	@id				ListTemplates
func ListTemplates(ctx *gin.Context) {
	server := server.GetInstance(nil)

	templates, err := server.TemplateService.List()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to list templates: %s", err.Error()))
		return
	}

	ctx.JSON(200, templates)
}

// SetTemplate godoc
//
//	@Tags			template
//	@Summary		Set template data
//	@Description	Set template data
//	@Accept			json
//	@Param			template	body	Template	true	"Template"
//	@Success		201
//	@Router			/template [put]
//
//	@id				SetTemplate
func SetTemplate(ctx *gin.Context) {
	var req template.Template
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid request body: %s", err.Error()))
		return
	}

	s := server.GetInstance(nil)

	err = s.TemplateService.Save(&req)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to save template: %s", err.Error()))
		return
	}

	ctx.Status(201)
}

// DeleteTemplate godoc
//
//	@Tags			template
//	@Summary		Delete template data
//	@Description	Delete template data
//	@Param			templateName	path	string	true	"Template name"
//	@Success		204
//	@Router			/template/{templateName} [delete]
//
//	@id				DeleteTemplate
func DeleteTemplate(ctx *gin.Context) {
	templateName := ctx.Param("templateName")

	server := server.GetInstance(nil)

	template, err := server.TemplateService.Find(templateName)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("failed to find project config: %s", err.Error()))
		return
	}

	err = server.TemplateService.Delete(template.Name)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(204)
}
