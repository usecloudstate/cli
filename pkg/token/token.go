package token

import (
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jdxcode/netrc"
)

type Token struct {
	user string
	pass string
}

var (
	machine = "api.usecloudstate.io"
)

func New(code string) *Token {
	parsed, err := parse(code)

	if err != nil {
		return nil
	}

	return &Token{
		user: parsed.Subject,
		pass: code,
	}
}

func parse(token string) (*jwt.StandardClaims, error) {
	parser := jwt.Parser{}

	t, _, err := parser.ParseUnverified(token, &jwt.StandardClaims{})

	if err != nil {
		return nil, err
	}

	return t.Claims.(*jwt.StandardClaims), nil
}

func (t *Token) GetPass() string {
	return t.pass
}

func (t *Token) GetUser() string {
	return t.user
}

func (t *Token) Expired() bool {
	parsed, err := parse(t.pass)

	if err != nil {
		return true
	}

	return parsed.ExpiresAt < time.Now().Unix()
}

func (t *Token) Save() error {
	if t == nil || t.user == "" || t.pass == "" {
		return fmt.Errorf("user and pass must be set")
	}

	home, err := user.Current()
	if err != nil {
		return err
	}

	path := filepath.Join(home.HomeDir, ".netrc")
	f, err := netrc.Parse(path)

	if err != nil {
		return err
	}

	f.RemoveMachine(machine)
	f.AddMachine(machine, t.user, t.pass)

	return f.Save()
}

func Get() (*Token, error) {
	home, err := user.Current()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(home.HomeDir, ".netrc")
	f, err := netrc.Parse(path)

	if err != nil {
		log.Println("Error opening your .netrc file.")
		log.Printf("\t Please make sure a .netrc file exists at %s \n", path)
		log.Printf("\t $ touch %s # creates a new file", path)
		return nil, err
	}

	m := f.Machine(machine)

	if m == nil {
		return nil, nil
	}

	return &Token{
		user: m.Get("login"),
		pass: m.Get("password"),
	}, nil
}

func RemoveMachine() error {
	home, err := user.Current()
	if err != nil {
		return err
	}

	path := filepath.Join(home.HomeDir, ".netrc")
	f, err := netrc.Parse(path)

	if err != nil {
		return err
	}

	f.RemoveMachine(machine)
	return f.Save()
}
