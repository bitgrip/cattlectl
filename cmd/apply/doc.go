// Copyright © 2018 Bitgrip <berlin@bitgrip.de>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apply

var applyLongDescription = `Apply a project descriptor to your rancher.

### Directory layout

By default this command expects in the current folder:

| File             | Description                                                 |
|------------------|-------------------------------------------------------------|
| __project.yaml__ | The project descriptor with support for go template syntax. |
| __values.yaml__  | The set of values used to render the project descriptor.    |

See the [Project Descriptor data model](project_descriptor.md)`
