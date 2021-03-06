package cmd

import (
	"fmt"
	"time"

	"github.com/mvisonneau/gitlab-merger/logger"
	"github.com/mvisonneau/go-gitlab"
	"github.com/nlopes/slack"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

type client struct {
	gitlab      *gitlab.Client
	gitlabAdmin *gitlab.Client
	slack       *slack.Client
}

type EmailMappings map[string]*Mapping

type Mapping struct {
	GitlabUserID int
	SlackUserID  string
}

var start time.Time
var c *client

func configure(ctx *cli.Context) error {
	start = ctx.App.Metadata["startTime"].(time.Time)

	lc := &logger.Config{
		Level:  ctx.GlobalString("log-level"),
		Format: ctx.GlobalString("log-format"),
	}

	if err := lc.Configure(); err != nil {
		return err
	}

	requiredFlags := []string{
		"gitlab-url",
		"gitlab-token",
	}

	if err := mandatoryStringOptions(ctx, requiredFlags); err != nil {
		return err
	}

	c = &client{
		gitlab: gitlab.NewClient(nil, ctx.GlobalString("gitlab-token")),
	}
	c.gitlab.SetBaseURL(ctx.GlobalString("gitlab-url"))

	if ctx.String("gitlab-admin-token") != "" {
		c.gitlabAdmin = gitlab.NewClient(nil, ctx.String("gitlab-admin-token"))
		c.gitlabAdmin.SetBaseURL(ctx.GlobalString("gitlab-url"))
	} else {
		c.gitlabAdmin = c.gitlab
	}

	if ctx.String("slack-token") != "" {
		c.slack = slack.New(ctx.String("slack-token"))
	}

	return nil
}

func mandatoryStringOptions(ctx *cli.Context, opts []string) (err error) {
	for _, o := range opts {
		if ctx.GlobalString(o) == "" && ctx.String(o) == "" {
			return fmt.Errorf("%s is required", o)
		}
	}
	return nil
}

func exit(err error, exitCode int) *cli.ExitError {
	defer log.Debugf("Executed in %s, exiting..", time.Since(start))
	if err != nil {
		log.Error(err.Error())
		return cli.NewExitError("", exitCode)
	}

	return cli.NewExitError("", 0)
}
