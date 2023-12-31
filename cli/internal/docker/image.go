package docker

import (
	"fmt"
	"kapigen.kateops.com/internal/logger"
)

type Image int

const DEPENDENCY_PROXY = "gitlab.kateops.com/infrastructure/dependency_proxy/containers/"

const (
	Kapigen_Latest Image = iota
	Alpine_3_18
	Terraform_Base
	BUILDKIT
	BUILDKIT_DAEMON
	CRANE_DEBUG

	GOLANG_1_21
)

var values = map[Image]string{
	Kapigen_Latest:  "kapigen",
	Alpine_3_18:     DEPENDENCY_PROXY + "alpine:3.18",
	Terraform_Base:  "hub.kateops.com/base/terraform:latest",
	BUILDKIT:        DEPENDENCY_PROXY + "moby/buildkit:master",
	BUILDKIT_DAEMON: DEPENDENCY_PROXY + "moby/buildkit:master-rootless",
	CRANE_DEBUG:     "gcr.io/go-containerregistry/crane:debug",
	GOLANG_1_21:     DEPENDENCY_PROXY + "golang:1.21",
}

func (c Image) String() string {
	if value, ok := values[c]; ok {
		return value
	}
	logger.Error(fmt.Sprintf("Not found for id: '%d'", c))
	return values[Alpine_3_18]
}
