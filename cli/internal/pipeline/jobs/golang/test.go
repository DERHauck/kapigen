package golang

import (
	"errors"
	"fmt"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/job/artifact"
	"kapigen.kateops.com/internal/gitlab/job/artifact/reports"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/pipeline/types"
	"kapigen.kateops.com/internal/pipeline/wrapper"
	"strings"
)

func NewUnitTest(image string, path string, packages []string) (*types.Job, error) {
	if image == "" {
		return nil, errors.New("no image set, required")
	}
	if path == "" {
		return nil, errors.New("no path set, required")
	}

	return types.NewJob("Unit Test", image, func(ciJob *job.CiJob) {
		var reportPath = fmt.Sprintf("%s/report.xml", path)
		if path == "." {
			reportPath = "report.xml"
		}
		coveragePkg := strings.Join(packages, ",")
		testCmd := fmt.Sprintf("go test -json -cover ./... -coverpkg=%s -coverprofile=profile.cov", coveragePkg)
		ciJob.SetImageName(image).
			TagMediumPressure().
			SetStage(stages.TEST).
			AddBeforeScript(fmt.Sprintf("cd %s", path)).
			AddScript("go install github.com/jstemmer/go-junit-report/v2@latest").
			AddScript(fmt.Sprintf("%s 2>&1 | go-junit-report -parser gojson -iocopy -out report.xml", testCmd)).
			AddScript("go tool cover -func profile.cov").
			AddVariable("KTC_PATH", path).
			AddArtifact(job.Artifact{
				Name:  "report",
				Paths: *(wrapper.NewStringSlice().Add(reportPath)),
				Reports: artifact.NewReports().
					SetCoverageReport(artifact.NewCoverageReport(reports.Cobertura, reportPath)).
					SetJunit(artifact.NewJunitReport(reportPath)),
			}).
			SetCodeCoverage(`/\(statements\)(?:\s+)?(\d+(?:\.\d+)?%)/`)

		ciJob.Rules = *job.DefaultPipelineRules()
	}), nil
}
