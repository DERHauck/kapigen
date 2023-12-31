package config

import (
	"errors"
	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/jobs/golang"
	"kapigen.kateops.com/internal/pipeline/types"
)

type GolangCoverage struct {
	Packages []string `yaml:"packages"`
}

func (g *GolangCoverage) Validate() error {
	if len(g.Packages) == 0 {
		logger.Info("no package declared, using./...")
		g.Packages = []string{"./..."}
	}
	return nil
}

type Golang struct {
	ImageName string          `yaml:"imageName"`
	Path      string          `yaml:"path"`
	Docker    *Docker         `yaml:"docker,omitempty"`
	Coverage  *GolangCoverage `yaml:"coverage,omitempty"`
}

func (g *Golang) New() types.PipelineConfigInterface {
	return &Golang{}
}

func (g *Golang) Validate() error {
	if g.ImageName == "" {
		return errors.New("no imageName set, required")
	}

	if g.Path == "" {
		return errors.New("no path set, required")
	}

	if g.Coverage == nil {
		g.Coverage = &GolangCoverage{}
	}
	if err := g.Coverage.Validate(); err != nil {
		return err
	}

	return nil
}

func (g *Golang) Build(factory *factory.MainFactory, pipelineType types.PipelineType, Id string) (*types.Jobs, error) {
	var allJobs = types.Jobs{}
	test, err := golang.NewUnitTest(g.ImageName, g.Path, g.Coverage.Packages)
	if err != nil {
		return nil, err
	}
	docker := g.Docker
	if docker != nil {
		jobs, err := types.GetPipelineJobs(factory, docker, pipelineType, Id)
		if err != nil {
			return nil, err
		}
		for _, job := range jobs.GetJobs() {
			job.AddJobAsNeed(test)
		}
		allJobs = append(allJobs, jobs.GetJobs()...)

	}
	allJobs = append(allJobs, test)
	return &allJobs, nil
}
