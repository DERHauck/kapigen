package config

import (
	"errors"
	"fmt"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/los"
	"kapigen.kateops.com/internal/pipeline/jobs/docker"
	"kapigen.kateops.com/internal/pipeline/types"
)

type Docker struct {
	Path       string `yaml:"path"`
	Context    string `yaml:"context"`
	Name       string `yaml:"name"`
	Dockerfile string `yaml:"dockerfile"`
}

func (d *Docker) New() types.PipelineConfigInterface {
	return &Docker{}
}

func (d *Docker) Validate() error {
	if d.Path == "" {
		return errors.New("need path set")
	}

	if d.Dockerfile == "" {
		d.Dockerfile = "Dockerfile"
	}

	if d.Context == "" {
		d.Context = d.Path
	}

	if d.Name == "" {

	}

	return nil
}

func (d *Docker) Build(pipelineType types.PipelineType, Id string) (*types.Jobs, error) {
	tag := los.GetVersion(environment.CI_PROJECT_ID.Get(), d.Path)
	//environment.GetNewVersion(tag)
	build := docker.NewBuildkitBuild(
		d.Path,
		d.Context,
		d.Dockerfile,
		d.DefaultRegistry(environment.GetFeatureBranchVersion(tag)),
	)
	return &types.Jobs{build}, nil

}

func (d *Docker) DefaultRegistry(tag string) string {
	if tag == "" {
		tag = "latest"
	}
	if d.Name != "" {
		return fmt.Sprintf("${CI_REGISTRY_IMAGE}/%s:%s", d.Name, tag)
	}
	return fmt.Sprintf("${CI_REGISTRY_IMAGE}:%s", tag)

}
