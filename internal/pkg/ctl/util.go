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
	rancher_client "github.com/bitgrip/cattlectl/internal/pkg/rancher/client"
)

var (
	getRancherClient = doGetRancherClient
	getClusterClient = doGetClusterClient
	getProjectClient = doGetProjectClient
)

func doGetProjectClient(projectName string, config config.Config) (rancherClient rancher_client.RancherClient, clusterClient rancher_client.ClusterClient, projectClient rancher_client.ProjectClient, err error) {
	rancherClient, clusterClient, err = getClusterClient(config)
	if err != nil {
		return
	}
	projectClient, err = clusterClient.Project(projectName)
	if err != nil {
		return
	}
	if exists, err := projectClient.Exists(); err != nil || !exists {
		if err != nil {
			return nil, nil, nil, err
		}
		return nil, nil, nil, fmt.Errorf("Project not found [%s]", projectName)
	}
	return
}

func doGetClusterClient(config config.Config) (rancherClient rancher_client.RancherClient, clusterClient rancher_client.ClusterClient, err error) {
	rancherClient, err = getRancherClient(config)
	if err != nil {
		return
	}
	clusterClient, err = rancherClient.Cluster(config.ClusterName())
	if err != nil {
		return
	}
	return
}

func doGetRancherClient(config config.Config) (rancherClient rancher_client.RancherClient, err error) {
	rancherClient, err = newRancherClient(rancher_client.RancherConfig{
		RancherURL:   config.RancherURL(),
		AccessKey:    config.AccessKey(),
		SecretKey:    config.SecretKey(),
		Insecure:     config.InsecureAPI(),
		CACerts:      config.CACerts(),
		MergeAnswers: config.MergeAnswers(),
	})
	if err != nil {
		return
	}
	return
}
