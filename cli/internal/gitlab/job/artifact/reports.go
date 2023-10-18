package artifact

type Reports struct {
	CoverageReport CoverageReport `yaml:"coverage_report,omitempty"`
	Junit          JunitReport    `yaml:"junit,omitempty"`
}

func (r *Reports) isValid() bool {
	return true
}

type ReportsYaml struct {
	CoverageReport *CoverageReportYaml `yaml:"coverage_report,omitempty"`
	Junit          string              `yaml:"junit,omitempty"`
}

func (r *Reports) Render() *ReportsYaml {
	if !r.isValid() {
		return nil
	}
	return &ReportsYaml{
		CoverageReport: r.CoverageReport.Render(),
		Junit:          r.Junit.Render(),
	}
}
