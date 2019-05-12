package client

import "github.com/sirupsen/logrus"

type resourceClient struct {
	id     string
	name   string
	logger *logrus.Entry
}

func (client *resourceClient) ID() (string, error) {
	return client.id, nil
}
func (client *resourceClient) Name() (string, error) {
	return client.name, nil
}
