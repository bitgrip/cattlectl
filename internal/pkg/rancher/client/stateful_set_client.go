package client

import (
	"fmt"

	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/rancher/norman/types"
	projectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

func newStatefulSetClient(
	name, namespace string,
	project ProjectClient,
	backendClient *projectClient.Client,
	logger *logrus.Entry,
) (StatefulSetClient, error) {
	return &statefulSetClient{
		resourceClient: resourceClient{
			name:   name,
			logger: logger.WithField("statefulSet_name", name).WithField("namespace", namespace),
		},
		namespace:     namespace,
		project:       project,
		backendClient: backendClient,
	}, nil
}

type statefulSetClient struct {
	resourceClient
	namespace     string
	namespaceID   string
	statefulSet   projectModel.StatefulSet
	project       ProjectClient
	backendClient *projectClient.Client
}

func (client *statefulSetClient) init() error {
	if client.namespaceID != "" {
		return nil
	}
	var namespace NamespaceClient
	var err error
	if namespace, err = client.project.Namespace(client.namespace); err != nil {
		client.logger.WithError(err).Error("Failed to read statefulSet list")
		return fmt.Errorf("Failed to read statefulSet list, %v", err)
	}
	if client.namespaceID, err = namespace.ID(); err != nil {
		client.logger.WithError(err).Error("Failed to read statefulSet list")
		return fmt.Errorf("Failed to read statefulSet list, %v", err)
	}
	return nil
}

func (client *statefulSetClient) Exists() (bool, error) {
	if err := client.init(); err != nil {
		return false, err
	}
	collection, err := client.backendClient.StatefulSet.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": client.namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read statefulSet list")
		return false, fmt.Errorf("Failed to read statefulSet list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name {
			return true, nil
		}
	}
	client.logger.Debug("StatefulSet not found")
	return false, nil
}

func (client *statefulSetClient) Create() error {
	if err := client.init(); err != nil {
		return err
	}
	client.logger.Info("Create new statefulSet")
	pattern, err := projectModel.ConvertStatefulSetToProjectAPI(client.statefulSet)
	if err != nil {
		return err
	}
	pattern.NamespaceId = client.namespaceID
	_, err = client.backendClient.StatefulSet.Create(&pattern)
	return err
}

func (client *statefulSetClient) Upgrade() error {
	return fmt.Errorf("upgrade statefulset not supported")
}

func (client *statefulSetClient) Data() (projectModel.StatefulSet, error) {
	return client.statefulSet, nil
}

func (client *statefulSetClient) SetData(statefulSet projectModel.StatefulSet) error {
	client.statefulSet = statefulSet
	return nil
}
