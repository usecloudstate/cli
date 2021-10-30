package apps

import "github.com/usecloudstate/cli/pkg/client"

type Apps struct {
	client *client.Client
}

func Init(client *client.Client) *Apps {
	return &Apps{client: client}
}
