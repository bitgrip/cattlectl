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
	"github.com/sirupsen/logrus"
)

var (
	deletableProjectResouceTypes = map[string]func(string, string, config.Config) error{
		"app": deleteApp,
	}
)

// DeleteProjectResouce is deleting one project resource from project
//
// * projectName: the project to delete the resource from
// * resourceType: the type of the resource to delete
// * name: the name of the resource to delete
func DeleteProjectResouce(projectName, resouceType, name string, config config.Config) (err error) {
	deleteFunc, supportedType := deletableProjectResouceTypes[resouceType]
	if !supportedType {
		return fmt.Errorf("Not supported resouce type [%s]", resouceType)
	}
	return deleteFunc(projectName, name, config)
}

func deleteApp(projectName, name string, config config.Config) (err error) {
	_, _, projectClient, err := getProjectClient(projectName, config)
	if err != nil {
		return
	}

	app, err := projectClient.App(name)
	if err != nil {
		return
	}

	if exists, err := app.Exists(); err != nil || !exists {
		if err != nil {
			return err
		}
		logrus.
			WithField("project-name", projectName).
			WithField("resouce-name", name).
			WithField("cluster-name", config.ClusterName()).
			Info("No app skip delete")
		return nil
	}

	err = app.Delete()

	return
}
