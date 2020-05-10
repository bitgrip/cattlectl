cattlectl
=========

[![Build Status](https://travis-ci.org/bitgrip/cattlectl.svg?branch=master)](https://travis-ci.org/bitgrip/cattlectl)

[![Docker Pulls](https://img.shields.io/docker/stars/bitgrip/cattlectl.svg)](https://store.docker.com/community/images/bitgrip/cattlectl)
[![Docker Pulls](https://img.shields.io/docker/pulls/bitgrip/cattlectl.svg)](https://store.docker.com/community/images/bitgrip/cattlectl)
[![](https://images.microbadger.com/badges/image/bitgrip/cattlectl.svg)](https://microbadger.com/images/bitgrip/cattlectl "Get your own image badge on microbadger.com")
[![](https://images.microbadger.com/badges/version/bitgrip/cattlectl.svg)](https://microbadger.com/images/bitgrip/cattlectl "Get your own version badge on microbadger.com")

[![Go Report Card](https://goreportcard.com/badge/github.com/bitgrip/cattlectl)](https://goreportcard.com/report/github.com/bitgrip/cattlectl)

Cattlectl is a tool for managing [Rancher 2](https://rancher.io) projects

Use cattlectl to:

* Apply project descriptors to a rancher managed kubernetes cluster
* Use one configuration as code to install to multiple stages
* Automate deployments to rancher managed kubernetes clusters from your CI server.

Install
-------

* Binary download of cattlectl can be found [on the Release page.](https://github.com/bitgrip/cattlectl/releases)
* Unpack the `cattlectl` binary and add it to your PATH and you are good to go!

Usage as docker image
---------------------

* You need to mount your descriptor to the directory `/data` in your container.
* `cattlectl` is the ENTRYPOINT so that you can use the cattlectl commands directly.

```bash
docker run --rm \
-v $(pwd):/data \
bitgrip/cattlectl apply
```

Build from source
-----------------

### cattlectl

```bash
go install \
-ldflags "-X github.com/bitgrip/cattlectl/internal/pkg/ctl.Version=$(git describe --tags) -s -w" \
-a -tags netgo -installsuffix netgo -mod=vendor
```

### Ansible modules

```
go build -mod=vendor -o ~/.ansible/plugins/modules/ ./ansible/...
```

Docs
----

* Get started with the [usage documentation](https://github.com/bitgrip/cattlectl/blob/master/docs/index.md)
* Read the [command documentation](https://github.com/bitgrip/cattlectl/blob/master/docs/cattlectl.md)

License
-------

Copyright © 2018 - 2019 [bitgrip GmbH](https://www.bitgrip.de/)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
