package version

import (
	"github.com/xanzy/go-gitlab"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/los"
	"strings"
)

const emptyTag = "0.0.0"

type Mode int

const (
	Gitlab Mode = iota
	Los
	None
)

func (m Mode) getTag() string {
	switch m {
	case Gitlab:
		return Gitlab.Name()
	case Los:
		return Los.Name()
	default:
		return None.Name()

	}
}

var values = map[Mode]string{
	Gitlab: "Gitlab",
	Los:    "Logic Operator Server",
	None:   "No versioning",
}

func (m Mode) Name() string {
	if value, ok := values[m]; ok {
		return value
	}
	logger.Error("no mode found, use default")
	m = Gitlab
	return values[m]
}

type Controller struct {
	current      string
	intermediate string
	new          string
	mode         Mode
	gitlabClient *gitlab.Client
	losClient    *los.Client
	refresh      bool
}

func NewController(mode Mode, gitlabClient *gitlab.Client, losClient *los.Client) *Controller {
	return &Controller{
		"",
		"",
		"",
		mode,
		gitlabClient,
		losClient,
		false,
	}
}

func (c *Controller) getTagFromGitlab() string {
	if c.refresh == false && c.current != "" {
		return c.current
	}
	oderBy := "updated"
	sort := "desc"
	if c.gitlabClient == nil {
		logger.Error("no gitlab client configured")
		return emptyTag
	}
	tags, _, err := c.gitlabClient.Tags.ListTags(environment.CI_PROJECT_ID.Get(), &gitlab.ListTagsOptions{OrderBy: &oderBy, Sort: &sort})
	if err != nil {
		logger.ErrorE(err)
		return emptyTag
	}
	logger.DebugAny(tags)
	return tags[0].Name
}

func (c *Controller) getTagFromLos(path string) string {
	if c.refresh == false && c.current != "" {
		return c.current
	}
	if c.losClient == nil {
		return emptyTag
	}
	return c.losClient.GetLatestVersion(environment.CI_PROJECT_ID.Get(), path)
}
func (c *Controller) Refresh() *Controller {
	c.refresh = true
	return c
}
func (c *Controller) GetCurrentTag(path string) string {
	if c.current == "" || c.refresh {
		if c.mode == Gitlab {
			c.current = c.getTagFromGitlab()
		}
		if c.mode == Los {
			c.current = c.getTagFromLos(path)
		}
	}

	c.refresh = false
	return c.current
}

func (c *Controller) GetCurrentPipelineTag(path string) string {
	if environment.IsRelease() {
		return c.GetNewTag(path)
	}

	return c.GetIntermediateTag(path)
}

func (c *Controller) GetNewTag(path string) string {
	if c.new == "" || c.refresh {

		if c.mode == Gitlab {
			c.new = getNewVersion(
				c.GetCurrentTag(path),
				c.getVersionIncrease(environment.CI_PROJECT_ID.Get(), environment.GetMergeRequestId()),
			)
		}
		if c.mode == Los {
			c.current = "0.0.0"
		}

	}
	c.refresh = false
	return c.new
}

func (c *Controller) GetIntermediateTag(path string) string {
	if c.intermediate == "" || c.refresh {
		if c.mode == Gitlab {
			c.intermediate = GetFeatureBranchVersion(c.GetCurrentTag(path), environment.GetBranchName())
		}
		if c.mode == Los {
			c.current = "0.0.0"
		}
	}
	c.refresh = false
	return c.intermediate
}

func (c *Controller) getVersionIncrease(projectId string, mrId int) string {
	if environment.IsRelease() {
		return getVersionIncreaseFromLabels(
			c.getMrLabelsFromApi(projectId, mrId),
		)
	}
	return getVersionIncreaseFromLabels(environment.CI_MERGE_REQUEST_LABELS.Get())
}

func (c *Controller) getMrLabelsFromApi(projectId string, mrId int) string {
	mr, _, err := c.gitlabClient.MergeRequests.GetMergeRequest(projectId, mrId, nil)
	if err != nil {
		logger.ErrorE(err)
		return "none"
	}
	return strings.Join(mr.Labels, ",")
}