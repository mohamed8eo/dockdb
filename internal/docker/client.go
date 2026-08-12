package docker

import "github.com/moby/moby/client"

func NewClient() (Client *client.Client, err error) {
	return client.New(client.FromEnv)
}
