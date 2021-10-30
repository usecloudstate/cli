package setup

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/usecloudstate/cli/pkg/apps"
	"github.com/xeonx/timeago"
)

type State int64

const (
	SAppCreated = iota
	SCreateReactAppJs
	SCreateReactAppTs
	SDumpEnv
	SOpenAdminConsole
	SFinished
)

type StateMachineCtx struct {
	ac    *apps.Apps
	appId string
	path  string
}

func promptWhichApp(ac *apps.Apps) (string, error) {
	apps, err := ac.GetApps()

	if err != nil {
		return "", err
	}

	var appIds []string

	for i, j := 0, len(apps)-1; i < j; i, j = i+1, j-1 {
		apps[i], apps[j] = apps[j], apps[i]
	}

	for _, a := range apps {
		if a.CreatedAt.IsZero() {
			appIds = append(appIds, a.AppId)
		} else {
			s := timeago.English.Format(a.CreatedAt)
			appIds = append(appIds, fmt.Sprintf("%s (Created %s)", a.AppId, s))
		}
	}

	prompt := promptui.Select{
		Label: "What is your appId?",
		Items: appIds,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	return strings.Fields(result)[0], nil
}

func promptNextStep() (State, error) {
	prompt := promptui.Select{
		Label: "Do you want to...",
		Items: []string{"Run create-react-app (Javascript)", "Run create-react-app (Typescript)", "Print the AppId so I can set it up myself", "Go to admin console"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return SFinished, err
	}

	switch result {
	case "Run create-react-app (Javascript)":
		return SCreateReactAppJs, nil
	case "Run create-react-app (Typescript)":
		return SCreateReactAppTs, nil
	case "Go to admin console, so I can set it up myself":
		return SOpenAdminConsole, nil
	}

	return SFinished, nil
}

func next(state State, ctx *StateMachineCtx) error {
	switch state {
	case SAppCreated:
		id, err := promptWhichApp(ctx.ac)

		if err != nil {
			return err
		}

		ctx.appId = id

		nextS, err := promptNextStep()

		if err != nil {
			return err
		}

		next(nextS, ctx)

	case SCreateReactAppJs:
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd %s && npx create-react-app --template @usecloudstate/usecloudstate-js . && echo \"REACT_APP_CLOUD_STATE_APP_ID=%s\" >> .env", ctx.path, ctx.appId))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		next(SFinished, ctx)

	case SCreateReactAppTs:
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd %s && npx create-react-app --template @usecloudstate/usecloudstate-ts . && echo \"REACT_APP_CLOUD_STATE_APP_ID=%s\" >> .env", ctx.path, ctx.appId))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		next(SFinished, ctx)

	case SOpenAdminConsole:
		consoleUrl := fmt.Sprintf("https://usecloudstate.io/admin/apps/%s/settings", ctx.appId)

		log.Printf("Opening admin console at %s\n", consoleUrl)
		cmd := exec.Command("/bin/sh", "-c", "open", consoleUrl)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		next(SFinished, ctx)
	}

	return nil
}

func createPathIfNotExists(path string) error {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}

	return nil
}

func RunSetup(ac *apps.Apps, path string) error {
	if path == "" {
		return fmt.Errorf("path is required")
	}

	err := createPathIfNotExists(path)

	if err != nil {
		return err
	}

	ctx := StateMachineCtx{ac: ac, path: path}

	return next(SAppCreated, &ctx)
}
