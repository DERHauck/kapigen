package types

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"kapigen.kateops.com/internal/logger"
	"os"
)

type PipelineTypeConfig struct {
	Type       PipelineType `yaml:"type"`
	Config     interface{}  `yaml:"config"`
	PipelineId string       `yaml:"id"`
}

func (p *PipelineTypeConfig) Decode(configTypes map[PipelineType]PipelineConfigInterface) (*Jobs, error) {
	logger.Info(fmt.Sprintf("decoding pipeline type %s", p.Type))
	var pipelineConfig = configTypes[p.Type].New()
	if pipelineConfig == nil {
		return nil, errors.New(fmt.Sprintf("no pipeline definition found for type: '%s'", p.Type))
	}
	err := mapstructure.Decode(p.Config, pipelineConfig)
	if err != nil {
		return nil, err
	}

	jobs, err := GetPipelineJobs(pipelineConfig, p.Type, p.PipelineId)
	if err != nil {
		return nil, err
	}
	for _, job := range jobs.GetJobs() {
		if p.PipelineId != "" {
			job.AddName(p.PipelineId)
		}
		job.AddName(string(p.Type))
		err = job.Render()
		if err != nil {
			return nil, err
		}
	}

	return jobs, nil
}

func (p *PipelineTypeConfig) GetType() PipelineType {
	return p.Type
}

type PipelineConfig struct {
	Pipelines []PipelineTypeConfig `yaml:"pipelines"`
}

func GetPipelineJobs(config PipelineConfigInterface, pipelineType PipelineType, pipelineId string) (*Jobs, error) {
	err := config.Validate()
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Pipeline type: %s, id: %s, encountered validation error: %s",
			pipelineType,
			pipelineId,
			err.Error(),
		))
	}

	jobs, err := config.Build(pipelineType, pipelineId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Pipeline type: %s, id: %s, encountered build error: %s",
			pipelineType,
			pipelineId,
			err.Error(),
		))
	}
	return jobs, nil
}

func LoadJobsFromPipelineConfig(configPath string, configTypes map[PipelineType]PipelineConfigInterface) (*Jobs, error) {
	body, err := os.ReadFile(configPath)

	if err != nil {
		return nil, err
	}

	var pipelineConfig PipelineConfig
	err = yaml.Unmarshal(body, &pipelineConfig)
	if err != nil {
		return nil, err
	}

	jobs, err := pipelineConfig.Decode(configTypes)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (p *PipelineConfig) Decode(configTypes map[PipelineType]PipelineConfigInterface) (*Jobs, error) {

	var pipelineJobs Jobs

	for i := 0; i < len(p.Pipelines); i++ {
		configuration := p.Pipelines[i]
		jobs, err := configuration.Decode(configTypes)
		if err != nil {
			return nil, err
		}
		pipelineJobs = append(pipelineJobs, jobs.GetJobs()...)
	}

	return &pipelineJobs, nil
}
