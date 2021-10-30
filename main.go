package main

import (
	"log"
	"os"

	"github.com/usecloudstate/cli/pkg/apps"
	"github.com/usecloudstate/cli/pkg/client"
	"github.com/usecloudstate/cli/pkg/session"
	"github.com/usecloudstate/cli/pkg/setup"
	"github.com/usecloudstate/cli/pkg/token"
)

type State int64

const (
	SInit = iota
	SClientInitiated
	SSessionFail
	SSessionSuccess
	SUnrecoverableError
	SAppCreateRequested
	SSetupRequested
	SSuccess
)

type StateMachineCtx struct {
	client  *client.Client
	sesh    *session.Session
	ac      *apps.Apps
	app     *apps.App
	cliArgs []string
}

func next(s State, ctx *StateMachineCtx) {
	switch s {
	case SInit:
		ctx.client = client.Init()
		next(SClientInitiated, ctx)

	case SClientInitiated:
		nsesh, err := session.Init(ctx.client)
		ac := apps.Init(ctx.client)
		ctx.sesh = nsesh
		ctx.ac = ac

		if err != nil {
			log.Println("Error initiating client:", err)
			next(SSessionFail, ctx)
		} else {
			next(SSessionSuccess, ctx)
		}

	case SSessionFail:
		err := token.RemoveMachine()

		if err != nil {
			next(SUnrecoverableError, ctx)
		} else {
			log.Println("❌ Error initiating new session.")
			next(SUnrecoverableError, ctx)
		}

	case SSessionSuccess:
		next(SSuccess, ctx)

	case SAppCreateRequested:
		app, err := ctx.ac.CreateNewApp()

		if err != nil {
			next(SUnrecoverableError, ctx)
		} else {
			ctx.app = app
			log.Println("💡 Hint to setup your app: $ cloudstate setup <path>")
			next(SSuccess, ctx)
		}

	case SSetupRequested:
		setup.RunSetup(ctx.ac, ctx.cliArgs[0])

	case SSuccess:
		return
	}
}

func main() {
	log.Println("🚀 Starting Cloudstate CLI")
	log.Println("Please report any issues at: https://github.com/usecloudstate/cli")
	log.Println("")

	argsWithoutProg := os.Args[1:]
	cmdVar := argsWithoutProg[0]
	argsWithoutCmd := argsWithoutProg[1:]

	ctx := StateMachineCtx{
		cliArgs: argsWithoutCmd,
	}
	next(SInit, &ctx)

	switch cmdVar {
	case "create":
		next(SAppCreateRequested, &ctx)
	case "setup":
		next(SSetupRequested, &ctx)
	default:
		log.Println("Unknown command:", cmdVar)
	}
}
