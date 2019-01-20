package plugin

import (
	"context"
	"errors"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
	"github.com/fbcbarbosa/drone-ignore-config/plugin/ignorer"
	"golang.org/x/oauth2"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-github/github"
)

// New returns a new ignore configuration plugin.
func New(server, token string) config.Plugin {
	return &plugin{
		server: server,
		token:  token,
	}
}

type plugin struct {
	server string
	token  string
}

func (p *plugin) Find(ctx context.Context, req *config.Request) (*drone.Config, error) {
	log.Debug("received build event")

	var client *github.Client

	// creates a github transport that authenticates
	// http requests using the github access token.
	trans := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: p.token},
	))

	// if a custom github endpoint is configured, for use
	// with github enterprise, we need to adjust the client
	// url accordingly.
	if p.server == "" {
		client = github.NewClient(trans)
	} else {
		var err error
		client, err = github.NewEnterpriseClient(p.server, p.server, trans)
		if err != nil {
			return nil, err
		}
	}

	// verify if we should ignore this build
	ignore, err := p.ShouldIgnore(ctx, client, req)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if ignore {
		log.Info("positive match: ignoring build")

		// this returns an error to force the build to be ignored
		// its cheap but prevents the build from being added to the pipeline
		return nil, errors.New(".droneignore match, skipping build")
	}

	log.Debug("negative match: will not ignore build")
	return nil, nil
}

func (p *plugin) ShouldIgnore(ctx context.Context, client *github.Client, req *config.Request) (bool, error) {
	path := ".droneignore"

	// if a .droneignore file exists, get it and apply the rules
	data, _, _, err := client.Repositories.GetContents(ctx, req.Repo.Namespace, req.Repo.Name, path, &github.RepositoryContentGetOptions{Ref: req.Build.After})

	if err != nil || data == nil {
		return false, err
	}

	ignoreExpr, err := data.GetContent()

	if err != nil || ignoreExpr == "" {
		return false, err
	}

	// build Ignorer object
	di, err := ignorer.New(strings.NewReader(ignoreExpr))

	if err != nil {
		return false, err
	}

	commit, _, err := client.Repositories.GetCommit(ctx, req.Repo.Namespace, req.Repo.Name, req.Build.After)

	for _, file := range commit.Files {
		if !di.ShouldIgnore(file.GetFilename()) {
			return false, nil
		}
	}

	return true, nil
}
