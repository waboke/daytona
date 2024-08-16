// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"fmt"
	"sync"

	"github.com/daytonaio/daytona/internal/util"
	"github.com/daytonaio/daytona/pkg/logs"
	"github.com/daytonaio/daytona/pkg/scheduler"
	"github.com/daytonaio/daytona/pkg/telemetry"
	log "github.com/sirupsen/logrus"
)

type BuildRunnerInstanceConfig struct {
	Interval         string
	Scheduler        scheduler.IScheduler
	BuildRunnerId    string
	BuildStore       Store
	BuilderFactory   IBuilderFactory
	LoggerFactory    logs.LoggerFactory
	TelemetryEnabled bool
	TelemetryService telemetry.TelemetryService
}

type BuildRunner struct {
	Id               string
	scheduler        scheduler.IScheduler
	interval         string
	buildStore       Store
	builderFactory   IBuilderFactory
	loggerFactory    logs.LoggerFactory
	telemetryEnabled bool
	telemetryService telemetry.TelemetryService
}

func NewBuildRunner(config BuildRunnerInstanceConfig) *BuildRunner {
	runner := &BuildRunner{
		Id:               config.BuildRunnerId,
		scheduler:        config.Scheduler,
		interval:         config.Interval,
		buildStore:       config.BuildStore,
		builderFactory:   config.BuilderFactory,
		loggerFactory:    config.LoggerFactory,
		telemetryEnabled: config.TelemetryEnabled,
		telemetryService: config.TelemetryService,
	}

	return runner
}

func (r *BuildRunner) Start() error {
	err := r.scheduler.AddFunc(r.interval, func() { r.Run() })
	if err != nil {
		return err
	}

	r.scheduler.Start()

	return nil
}

func (r *BuildRunner) Stop() {
	r.scheduler.Stop()
}

func (r *BuildRunner) Run() {
	builds, err := r.buildStore.List(&Filter{
		States: []*BuildState{util.Pointer(BuildStatePending), util.Pointer(BuildStateSuccess), util.Pointer(BuildStatePublished)},
	})
	if err != nil {
		log.Error(err)
		return
	}

	var wg sync.WaitGroup
	for _, currentBuild := range builds {
		if currentBuild.State == BuildStatePending {
			wg.Add(1)
			if currentBuild.BuildConfig != nil {
				currentBuild.BuildConfig.CachedBuild = GetCachedBuild(currentBuild, builds)
			}
			go r.RunBuildProcess(currentBuild, &wg)
		}
	}

	wg.Wait()
}

func (r *BuildRunner) RunBuildProcess(b *Build, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	if b.BuildConfig == nil {
		return
	}

	buildLogger := r.loggerFactory.CreateBuildLogger(b.Id, logs.LogSourceBuilder)
	defer buildLogger.Close()

	builder, err := r.builderFactory.Create(*b)
	if err != nil {
		r.handleBuildError(*b, builder, err, buildLogger)
		return
	}

	b.State = BuildStateRunning
	err = r.buildStore.Save(b)
	if err != nil {
		r.handleBuildError(*b, builder, err, buildLogger)
		return
	}

	image, user, err := builder.Build(*b)
	if err != nil {
		r.handleBuildError(*b, builder, err, buildLogger)
		return
	}

	b.Image = image
	b.User = user
	b.State = BuildStateSuccess
	err = r.buildStore.Save(b)
	if err != nil {
		r.handleBuildError(*b, builder, err, buildLogger)
		return
	}

	err = builder.Publish(*b)
	if err != nil {
		r.handleBuildError(*b, builder, err, buildLogger)
		return
	}

	b.State = BuildStatePublished
	err = r.buildStore.Save(b)
	if err != nil {
		r.handleBuildError(*b, builder, err, buildLogger)
		return
	}

	err = builder.CleanUp()
	if err != nil {
		errMsg := fmt.Sprintf("Error cleaning up build: %s\n", err.Error())
		buildLogger.Write([]byte(errMsg + "\n"))
	}

	if r.telemetryEnabled {
		r.logTelemetry(context.Background(), *b, err)
	}
}

func (r *BuildRunner) handleBuildError(b Build, builder IBuilder, err error, buildLogger logs.Logger) {
	var errMsg string
	errMsg += "################################################\n"
	errMsg += fmt.Sprintf("#### BUILD FAILED FOR %s: %s\n", b.Id, err.Error())
	errMsg += "################################################\n"

	b.State = BuildStateError
	err = r.buildStore.Save(&b)
	if err != nil {
		errMsg += fmt.Sprintf("Error saving build: %s\n", err.Error())
	}

	cleanupErr := builder.CleanUp()
	if cleanupErr != nil {
		errMsg += fmt.Sprintf("Error cleaning up build: %s\n", cleanupErr.Error())
	}

	buildLogger.Write([]byte(errMsg + "\n"))

	if r.telemetryEnabled {
		r.logTelemetry(context.Background(), b, err)
	}
}

func (r *BuildRunner) logTelemetry(ctx context.Context, b Build, err error) {
	telemetryProps := telemetry.NewBuildRunnerEventProps(ctx, b.Id, string(b.State))
	event := telemetry.BuildEventRunBuild
	if err != nil {
		telemetryProps["error"] = err.Error()
		event = telemetry.BuildEventRunBuildError
	}
	telemetryError := r.telemetryService.TrackBuildEvent(event, r.Id, telemetryProps)
	if telemetryError != nil {
		log.Trace(telemetryError)
	}
}
