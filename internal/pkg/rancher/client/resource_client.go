package client

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

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

type namespacedResourceClient struct {
	resourceClient
	namespaceID string
	namespace   string
	project     ProjectClient
}

func (client *namespacedResourceClient) NamespaceID() (string, error) {
	if client.namespace == "" {
		return "", nil
	}

	if client.namespaceID != "" {
		return client.namespaceID, nil
	}
	var namespace NamespaceClient
	var err error
	if namespace, err = client.project.Namespace(client.namespace); err != nil {
		client.logger.WithError(err).Error("Failed to read namespaceID")
		return "", fmt.Errorf("Failed to read namespaceID, %v", err)
	}
	if client.namespaceID, err = namespace.ID(); err != nil {
		client.logger.WithError(err).Error("Failed to read namespaceID")
		return "", fmt.Errorf("Failed to read namespaceID, %v", err)
	}
	return client.namespaceID, nil
}
func (client *namespacedResourceClient) Namespace() (string, error) {
	return client.namespace, nil
}
