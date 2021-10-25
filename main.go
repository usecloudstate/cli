package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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
	SAppCreated
	SSetupRequested
	SSuccess
)

type StateMachineCtx struct {
	client *client.Client
	sesh   *session.Session
	ac     *apps.AppCreator
	app    *apps.App
}

func next(s State, ctx StateMachineCtx) {
	switch s {
	case SInit:
		ctx.client = client.Init()
		next(SClientInitiated, ctx)

	case SClientInitiated:
		nsesh, err := session.Init(ctx.client)
		ctx.sesh = nsesh

		if err != nil {
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
		ac := apps.Init(ctx.client)
		app, err := ac.CreateNewApp()

		if err != nil {
			next(SUnrecoverableError, ctx)
		} else {
			ctx.app = app
			ctx.ac = ac
			next(SAppCreated, ctx)
		}

	case SAppCreated:
		setup.PrintSetup(ctx.app.AppId)

	case SSuccess:
		os.Exit(0)
	}
}

func main() {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
		case syscall.SIGTERM:
			log.Println("Exiting...")
			os.Exit(1)
		}
	}()

	ctx := StateMachineCtx{}
	next(SInit, ctx)
}
