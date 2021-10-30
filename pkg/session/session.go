package session

import (
	"errors"
	"log"
	"net/mail"

	"github.com/manifoldco/promptui"
	"github.com/usecloudstate/cli/pkg/client"
	"github.com/usecloudstate/cli/pkg/token"
)

type Session struct {
	client *client.Client
	token  *token.Token
}

func validEmail() (string, error) {
	validate := func(email string) error {
		_, err := mail.ParseAddress(email)

		if err != nil {
			return errors.New("invalid email address")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Your e-mail address",
		Validate: validate,
	}

	return prompt.Run()
}

func Init(client *client.Client) (*Session, error) {
	log.Println("Checking for an existing session.")

	t, err := token.Get()

	if err != nil {
		log.Println("No existing session found.")
		return nil, err
	}

	notoken := t == nil

	expired := t != nil && t.Expired()

	s := &Session{
		client: client,
	}

	if notoken || expired {
		log.Println("No session found. Requesting a new one.")

		err := s.requestSession()

		if err != nil {
			return nil, err
		}

		newt, err := s.confirmCode()

		if err != nil {
			return nil, err
		}

		s.token = newt
	} else {
		log.Println("Found an existing session.")
		s.token = t
	}

	s.client.SetAuthToken(s.token.GetPass())

	log.Println("Verifying your session with Cloud State server...")
	err = s.VerifySession()

	if err != nil {
		return nil, err
	} else {
		log.Println("Session verified. ✅")
	}

	return s, nil
}

func (s *Session) VerifySession() error {
	// TODO: Implement a auth token verification later on
	resp, err := s.client.Request("GET", "admin/apps", nil)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("invalid session")
	}

	return nil
}

func (sesh *Session) requestSession() error {
	email, err := validEmail()

	if err != nil {
		return err
	}

	log.Println("Sending an e-mail to you...")

	_, err = sesh.client.Request("PUT", "apps/00000000000000000000000000000000/user_session_request", map[string]interface{}{
		"email":   email,
		"fromCli": true,
	})

	if err != nil {
		return err
	}

	log.Println("E-mail sent. Please click the login link and paste your token below.")

	return nil
}

func (sesh *Session) confirmCode() (*token.Token, error) {
	prompt := promptui.Prompt{
		Label:       "Confirmation token",
		HideEntered: true,
	}

	code, err := prompt.Run()

	if err != nil {
		return nil, err
	}

	t := token.New(code)
	err = t.Save()

	if err != nil {
		return t, err
	}

	log.Println("Saving your token.")

	return t, err
}
