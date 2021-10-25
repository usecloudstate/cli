package main

import (
	"os"

	"github.com/usecloudstate/cli/pkg/apps"
	"github.com/usecloudstate/cli/pkg/client"
	"github.com/usecloudstate/cli/pkg/session"
	"github.com/usecloudstate/cli/pkg/setup"
)

type State int64

const (
	SInit = iota
	SClientInitiated
	SSessionFail
	SSessionSuccess
	SUnrecoverableError
	SAppCreated
	SSetupRequested
	SSuccess
)

func next(s State, c *client.Client, sesh *session.Session, ac *apps.AppCreator, app *apps.App) {
	switch s {
	case SInit:
		next(SClientInitiated, client.Init(), nil, nil, nil)

	case SClientInitiated:
		nsesh, err := session.Init(c)

		if err != nil {
			next(SSessionFail, c, nil, nil, nil)
		} else {
			next(SSessionSuccess, c, nsesh, nil, nil)
		}

	case SSessionFail:
		err := sesh.PurgeSession()

		if err != nil {
			next(SUnrecoverableError, c, sesh, ac, nil)
		} else {
			next(SClientInitiated, c, nil, nil, nil)
		}

	case SSessionSuccess:
		ac = apps.Init(c)
		app, err := ac.CreateNewApp()

		if err != nil {
			next(SUnrecoverableError, c, sesh, ac, nil)
		} else {
			next(SAppCreated, c, sesh, ac, app)
		}

	case SAppCreated:
		setup.PrintSetup(app.AppId)

	case SSuccess:
		os.Exit(0)
	}
}

func main() {
	next(SInit, nil, nil, nil, nil)
}
