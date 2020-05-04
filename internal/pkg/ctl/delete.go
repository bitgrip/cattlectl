//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ctl

import (
	"fmt"

	"github.com/bitgrip/cattlectl/internal/pkg/config"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/client"
	"github.com/sirupsen/logrus"
)

var (
	deletableProjectResouceTypes = map[string]func(string, string, string, config.Config) (bool, error){
		"namespace":         deleteNamespace,
		"certificate":       deleteCertificate,
		"config-map":        deleteConfigMap,
		"docker-credential": deleteDockerCredential,
		"secret":            deleteSecret,
		"app":               deleteApp,
		"job":               deleteJob,
		"cron-job":          deleteCronJob,
		"deployment":        deleteDeployment,
		"daemon-set":        deleteDaemonSet,
		"stateful-set":      deleteStatefulSet,
	}
)

// DeleteProjectResouce is deleting one project resource from project
//
// * projectName: the project to delete the resource from
// * resourceType: the type of the resource to delete
// * name: the name of the resource to delete
func DeleteProjectResouce(projectName, namespace, resouceType, name string, config config.Config) (bool, error) {
	deleteFunc, supportedType := deletableProjectResouceTypes[resouceType]
	if !supportedType {
		return false, fmt.Errorf("Not supported resouce type [%s]", resouceType)
	}
	return deleteFunc(projectName, namespace, name, config)
}

func deleteNamespace(projectName, _namespace, name string, config config.Config) (deleted bool, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	namespace, err := projectClient.Namespace(name)
	if err != nil {
		return
	}

	return deleteProjectResouce(namespace, config.ClusterName(), projectName, "namespace", name)
}

func deleteCertificate(projectName, namespace, name string, config config.Config) (deleted bool, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	certificate, err := projectClient.Certificate(name, namespace)
	if err != nil {
		return
	}

	if namespace == "" {
		return deleteProjectResouce(certificate, config.ClusterName(), projectName, "certificate", name)
	}
	return deleteNamespaceResouce(certificate, config.ClusterName(), projectName, namespace, "certificate", name)
}

func deleteConfigMap(projectName, namespace, name string, config config.Config) (deleted bool, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	configMap, err := projectClient.ConfigMap(name, namespace)
	if err != nil {
		return
	}

	return deleteNamespaceResouce(configMap, config.ClusterName(), projectName, namespace, "config-map", name)
}

func deleteDockerCredential(projectName, namespace, name string, config config.Config) (deleted bool, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	dockerCredential, err := projectClient.DockerCredential(name, namespace)
	if err != nil {
		return
	}

	if namespace == "" {
		return deleteProjectResouce(dockerCredential, config.ClusterName(), projectName, "docker-credential", name)
	}
	return deleteNamespaceResouce(dockerCredential, config.ClusterName(), projectName, namespace, "docker-credential", name)
}

func deleteSecret(projectName, namespace, name string, config config.Config) (deleted bool, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	secret, err := projectClient.Secret(name, namespace)
	if err != nil {
		return
	}

	if namespace == "" {
		return deleteProjectResouce(secret, config.ClusterName(), projectName, "secret", name)
	}
	return deleteNamespaceResouce(secret, config.ClusterName(), projectName, namespace, "secret", name)
}

func deleteApp(projectName, namespace, name string, config config.Config) (deleted bool, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	app, err := projectClient.App(name)
	if err != nil {
		return
	}

	return deleteProjectResouce(app, config.ClusterName(), projectName, "app", name)
}

func deleteJob(projectName, namespace, name string, config config.Config) (deleted bool, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	job, err := projectClient.Job(name, namespace)
	if err != nil {
		return
	}

	return deleteNamespaceResouce(job, config.ClusterName(), projectName, namespace, "job", name)
}

func deleteCronJob(projectName, namespace, name string, config config.Config) (deleted bool, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	cronJob, err := projectClient.CronJob(name, namespace)
	if err != nil {
		return
	}

	return deleteNamespaceResouce(cronJob, config.ClusterName(), projectName, namespace, "cron-job", name)
}

func deleteDeployment(projectName, namespace, name string, config config.Config) (deleted bool, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	deployment, err := projectClient.Deployment(name, namespace)
	if err != nil {
		return
	}

	return deleteNamespaceResouce(deployment, config.ClusterName(), projectName, namespace, "deployment", name)
}

func deleteDaemonSet(projectName, namespace, name string, config config.Config) (deleted bool, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	daemonSet, err := projectClient.DaemonSet(name, namespace)
	if err != nil {
		return
	}

	return deleteNamespaceResouce(daemonSet, config.ClusterName(), projectName, namespace, "daemon-set", name)
}

func deleteStatefulSet(projectName, namespace, name string, config config.Config) (deleted bool, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	statefulSet, err := projectClient.StatefulSet(name, namespace)
	if err != nil {
		return
	}

	return deleteNamespaceResouce(statefulSet, config.ClusterName(), projectName, namespace, "stateful-set", name)
}

func deleteProjectResouce(resource client.ResourceClient, clusterName, projectName, resouceType, name string) (deleted bool, err error) {
	if exists, err := resource.Exists(); err != nil || !exists {
		if err != nil {
			return false, err
		}
		logrus.
			WithField("project-name", projectName).
			WithField("resouce-name", name).
			WithField("cluster-name", clusterName).
			Infof("No %s skip delete", resouceType)
		return false, nil
	}

	err = resource.Delete()
	deleted = err == nil
	return
}

func deleteNamespaceResouce(resource client.ResourceClient, clusterName, projectName, namespace, resouceType, name string) (deleted bool, err error) {
	if exists, err := resource.Exists(); err != nil || !exists {
		if err != nil {
			return false, err
		}
		logrus.
			WithField("project-name", projectName).
			WithField("namespace", namespace).
			WithField("resouce-name", name).
			WithField("cluster-name", clusterName).
			Infof("No %s skip delete", resouceType)
		return false, nil
	}

	err = resource.Delete()
	deleted = err == nil
	return
}
