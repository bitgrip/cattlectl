//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ctl

import (
	"fmt"
	"regexp"

	"github.com/bitgrip/cattlectl/internal/pkg/config"
)

var (
	listableProjectResouceTypes = map[string]func(string, string, config.Config) ([]string, error){
		"namespace":          listNamespaces,
		"namespaces":         listNamespaces,
		"certificate":        listCertificates,
		"certificates":       listCertificates,
		"config-map":         listConfigMaps,
		"config-maps":        listConfigMaps,
		"docker-credential":  listDockerCredentials,
		"docker-credentials": listDockerCredentials,
		"secret":             listSecrets,
		"secrets":            listSecrets,
		"app":                listApps,
		"apps":               listApps,
		"job":                listJobs,
		"jobs":               listJobs,
		"cron-job":           listCronJobs,
		"cron-jobs":          listCronJobs,
		"deployment":         listDeployments,
		"deployments":        listDeployments,
		"daemon-set":         listDaemonSets,
		"daemon-sets":        listDaemonSets,
		"stateful-set":       listStatefulSets,
		"stateful-sets":      listStatefulSets,
	}
)

// ListProjectResouces list all resources of a project to stdout
//
// * projectName: the project to list the resources from
// * namespace: the namespace to list the resources from
// * resourceType: the type of the resources to list
func ListProjectResouces(projectName, namespace, resouceType, pattern string, config config.Config) (err error) {
	listFunc, supportedType := listableProjectResouceTypes[resouceType]
	if !supportedType {
		return fmt.Errorf("Not supported resouce type [%s]", resouceType)
	}
	names, err := listFunc(projectName, namespace, config)
	if err != nil {
		return
	}
	for _, name := range names {
		matched, _ := regexp.MatchString(pattern, name)
		if !matched {
			continue
		}
		fmt.Println(name)
	}
	return
}

func listNamespaces(projectName, namespace string, config config.Config) (names []string, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	namespaces, err := projectClient.Namespaces()
	if err != nil {
		return
	}

	for _, namespace := range namespaces {
		name, err := namespace.Name()
		if err != nil {
			return names, err
		}
		names = append(names, name)
	}

	return
}

func listCertificates(projectName, namespace string, config config.Config) (names []string, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	certificates, err := projectClient.Certificates(namespace)
	if err != nil {
		return
	}

	for _, certificate := range certificates {
		name, err := certificate.Name()
		if err != nil {
			return names, err
		}
		names = append(names, name)
	}

	return
}

func listConfigMaps(projectName, namespace string, config config.Config) (names []string, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	configMaps, err := projectClient.ConfigMaps(namespace)
	if err != nil {
		return
	}

	for _, configMap := range configMaps {
		name, err := configMap.Name()
		if err != nil {
			return names, err
		}
		names = append(names, name)
	}

	return
}

func listDockerCredentials(projectName, namespace string, config config.Config) (names []string, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	dockerCredentials, err := projectClient.DockerCredentials(namespace)
	if err != nil {
		return
	}

	for _, dockerCredential := range dockerCredentials {
		name, err := dockerCredential.Name()
		if err != nil {
			return names, err
		}
		names = append(names, name)
	}

	return
}

func listSecrets(projectName, namespace string, config config.Config) (names []string, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	secrets, err := projectClient.Secrets(namespace)
	if err != nil {
		return
	}

	for _, secret := range secrets {
		name, err := secret.Name()
		if err != nil {
			return names, err
		}
		names = append(names, name)
	}

	return
}

func listApps(projectName, namespace string, config config.Config) (names []string, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	apps, err := projectClient.Apps()
	if err != nil {
		return
	}

	for _, app := range apps {
		name, err := app.Name()
		if err != nil {
			return names, err
		}
		names = append(names, name)
	}

	return
}

func listJobs(projectName, namespace string, config config.Config) (names []string, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	jobs, err := projectClient.Jobs(namespace)
	if err != nil {
		return
	}

	for _, job := range jobs {
		name, err := job.Name()
		if err != nil {
			return names, err
		}
		names = append(names, name)
	}

	return
}

func listCronJobs(projectName, namespace string, config config.Config) (names []string, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	cronJobs, err := projectClient.CronJobs(namespace)
	if err != nil {
		return
	}

	for _, cronJob := range cronJobs {
		name, err := cronJob.Name()
		if err != nil {
			return names, err
		}
		names = append(names, name)
	}

	return
}

func listDeployments(projectName, namespace string, config config.Config) (names []string, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	deployments, err := projectClient.Deployments(namespace)
	if err != nil {
		return
	}

	for _, deployment := range deployments {
		name, err := deployment.Name()
		if err != nil {
			return names, err
		}
		names = append(names, name)
	}

	return
}

func listDaemonSets(projectName, namespace string, config config.Config) (names []string, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	daemonSets, err := projectClient.DaemonSets(namespace)
	if err != nil {
		return
	}

	for _, daemonSet := range daemonSets {
		name, err := daemonSet.Name()
		if err != nil {
			return names, err
		}
		names = append(names, name)
	}

	return
}

func listStatefulSets(projectName, namespace string, config config.Config) (names []string, err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	statefulSets, err := projectClient.StatefulSets(namespace)
	if err != nil {
		return
	}

	for _, statefulSet := range statefulSets {
		name, err := statefulSet.Name()
		if err != nil {
			return names, err
		}
		names = append(names, name)
	}

	return
}
