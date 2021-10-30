package apps

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/manifoldco/promptui"
)

type App struct {
	AppId           string    `bson:"appId" json:"appId"`
	Name            string    `bson:"name" json:"name"`
	Description     string    `bson:"description" json:"description"`
	Origin          string    `bson:"origin" json:"origin"`
	AllowSelfSignup bool      `bson:"allowSelfSignup" json:"allowSelfSignup"`
	CreatedAt       time.Time `bson:"createdAt" json:"createdAt"`
}

func (a *Apps) CreateNewApp() (*App, error) {
	log.Println("Creating a new app...")

	wd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	prAppName := promptui.Prompt{
		Label: "App Name",
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("invalid input")
			}

			return nil
		},
		Default: filepath.Base(wd),
	}

	prAppDesc := promptui.Prompt{
		Label: "App Description",
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("invalid input")
			}

			return nil
		},
		Default: "My awesome app",
	}

	prOrigin := promptui.Prompt{
		Label: "Origin",
		Validate: func(input string) error {
			hasHttpOrHttps := regexp.MustCompile(`^(http|https)://`).MatchString(input)

			if !hasHttpOrHttps {
				return errors.New("invalid url")
			}

			return nil
		},
		Default: "http://localhost:3000",
	}

	prPubSignup := promptui.Prompt{
		Label:   "Allow public signup?",
		Default: "yes",
	}

	appName, err := prAppName.Run()

	if err != nil {
		return nil, err
	}

	appDesc, err := prAppDesc.Run()

	if err != nil {
		return nil, err
	}

	origin, err := prOrigin.Run()

	if err != nil {
		return nil, err
	}

	pubSignupStr, err := prPubSignup.Run()

	if err != nil {
		return nil, err
	}

	pubSignup := pubSignupStr[0] == 'y'

	return a.create(appName, appDesc, origin, pubSignup)
}

func (a *Apps) create(appName string, appDescription string, origin string, publicSignup bool) (*App, error) {
	payload := map[string]interface{}{
		"name":            appName,
		"description":     appDescription,
		"origin":          origin,
		"allowSelfSignup": publicSignup,
	}

	resp, err := a.client.Request("POST", "admin/apps", payload)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b := App{}
	err = json.NewDecoder(resp.Body).Decode(&b)

	if err != nil {
		return nil, err
	}

	log.Println("✅ Your new app is created.")
	log.Printf("App ID: %s\n", b.AppId)
	log.Printf("Name: %s\n", b.Name)
	log.Printf("Description: %s\n", b.Description)
	log.Printf("Origin: %s\n", b.Origin)
	log.Printf("Allow self signup: %t\n", b.AllowSelfSignup)

	return &b, nil
}
