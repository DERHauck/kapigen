package infrastructure

import (
	"fmt"
	"kapigen.kateops.com/docker"
	"kapigen.kateops.com/internal/gitlab"
	"kapigen.kateops.com/internal/gitlab/environment"
	"kapigen.kateops.com/internal/gitlab/rules"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/types"
)

func NewTerraformInit(path string, state string, s3 bool) *types.Job {
	return types.NewJob("Init", docker.Terraform_Base, func(job *gitlab.CiJob) {
		job.BeforeScript.Value.AddSeveral([]string{
			"echo \"credentials \\\\\"${CI_SERVER_HOST}\\\\\" {\\n  token = \\\\\"${CI_PIPELINE_TOKEN}\\\\\"\\n}\" > gitlab.tfrc",
			"export TF_CLI_CONFIG_FILE=${PWD}/gitlab.tfrc",
		})
		project, err := environment.Lookup(environment.CI_PROJECT_ID)
		if err != nil {
			logger.ErrorE(err)
			project = "test"
		}

		job.Stage = stages.INIT
		job.Script.Value.Add(
			"terraform init \\\n" +
				" -backend-config=\"region=eu-central-1\" \\\n" +
				" -backend-config=\"access_key=${TF_STATE_ACCESS_KEY}\" \\\n" +
				" -backend-config=\"secret_key=${TF_STATE_SECRET_KEY}\" \\\n" +
				" -backend-config=\"bucket=${TF_STATE_BUCKET}\" \\\n" +
				fmt.Sprintf(` -backend-config="key=${TF_STATE_BUCKET}/states/%s/%s/terraform.tfstate"`, project, state),
		)

		job.Variables = map[string]string{}
		job.Variables["TF_STATE_PROJECT"] = project
		//job.AllowFailure.Failure = true
		job.Rules = *rules.DefaultPipelineRules()
		job.Cache.SetPolicy(gitlab.CachePolicyEnumPullPush).
			SetActive()
	})
}
