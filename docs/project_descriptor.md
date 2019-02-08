Project Descriptor data model
=============================

ProjectDescriptor Structur:
---------------------------

### Toplevel ProjectDescriptor

| Field                  | Description                                                           |
|------------------------|-----------------------------------------------------------------------|
| __api_version__        | The __\<major\>.\<minor\>__ version used for this descriptor.         |
| __kind__               | The kind of descriptor in this file (`Project`)                       |
| __metadata__           | Metainformation about this descriptor e.g.: name and cluster_name     |
| __namespaces__         | List of namespaces to be part of this project                         |
| __resources__          | List of resources to be part of this project                          |
| __storage_classes__    | List of storage classes on cluster level required for this project    |
| __persistent_volumes__ | List of persistent volumes on cluster level required for this project |
| __apps__               | List of rancher apps to be deployed to this project                   |

### metadata

* In the descriptor only `name` should be set.
* All other fields are read from configuration or rancher

| Field            | Description                                                                              |
|------------------|------------------------------------------------------------------------------------------|
| __name__         | The name of the project **REQUIRED**                                                     |
| __id__           | Rancher internal ID of this project (**read from rancher**)                              |
| __rancher_url__  | The URL to reach the rancher (**placed from cattleclt configuration**)                   |
| __access_key__   | The access key to access rancher with (**placed from cattleclt configuration**)          |
| __secret_key__   | The secret key to access rancher with (**placed from cattleclt configuration**)          |
| __token_key__    | The token key to access rancher with (**placed from cattleclt configuration**)           |
| __cluster_name__ | The name of the cluster the project is part of (**placed from cattleclt configuration**) |
| __cluster_id__   | The ID of the cluster the project is part of (**read from rancher**)                     |

#### namespaces

| Field    | Description               |
|----------|---------------------------|
| __name__ | The name of the namespace |

#### resources

| Field                 | Description                          |
|-----------------------|--------------------------------------|
| __certificates__      | Array of certificates to deploy      |
| __config_maps__       | Array of config maps to deploy       |
| __docker_credential__ | Array of docker credential to deploy |
| __secrets__           | Array of secrets to deploy           |

#### certificate

| Field         | Description                                                           |
|---------------|-----------------------------------------------------------------------|
| __name__      | The name of the certificate resource                                  |
| __key__       | Private key to the certificate                                        |
| __certs__     | One string with all certs                                             |
| __namespace__ | if empty the certificate is deployed to all namespaces of the project |

#### config_map

| Field         | Description                                                           |
|---------------|-----------------------------------------------------------------------|
| __name__      | The name of the config map                                            |
| __data__      | map[string]string structure representing the config map payload       |

#### docker_credential

| Field          | Description                                                           |
|----------------|-----------------------------------------------------------------------|
| __name__       | The name of the docker credential                                     |
| __namespace__  | if empty the certificate is deployed to all namespaces of the project |
| __registries__ | Array of registry credentials for this resource                       |

#### registry

| Field        | Description                  |
|--------------|------------------------------|
| __name__     | The name of the registry     |
| __username__ | The username of the registry |
| __password__ | The password of the registry |

#### secret

| Field    | Description                                                     |
|----------|-----------------------------------------------------------------|
| __name__ | The name of the secret                                          |
| __data__ | map[string]string structure representing the config map payload |

#### storage classes

| Field                         | Description                                                           |
|-------------------------------|-----------------------------------------------------------------------|
| __name__                      | The name of the storage class                                         |
| __provisioner__               | The provisioner used by this storage class                            |
| __reclaim_policy__            | The reclaim policy of this storage class                              |
| __volume_bind_mode__          | The volume bind mode of this storage class                            |
| __parameters__                | key value map of parameters                                           |
| __mount_options__             | array of options to mount the volumes of this storage class           |
| __create_persistent_volumes__ | If persistent volumes should be precreated for this storage class     |
| __persistent_volume_groups__  | The list of persistent volume groups to create for this storage class |

#### persistent volume group

| Field             | Description                                                             |
|-------------------|-------------------------------------------------------------------------|
| __name__          | The name of the persistent volume group                                 |
| __type__          | The Type of the persistent volume group (__only "local" is supportet__) |
| __path__          | For local persistent volumes the path pattern to use                    |
| __create_script__ | For log informations the hint how to create the required directories.   |
| __access_modes__  | List of access modes (string array)                                     |
| __capacity__      | The capacity of the persistent volumes to create                        |
| __nodes__         | For local persistent volumes the ist of nodes to bound (string array).  |

#### persistent volume

| Field                  | Description                                                           |
|------------------------|-----------------------------------------------------------------------|
| __name__               | The name of the persistent volume                                     |
| __path__               | For local persistent volumes the path pattern to use                  |
| __node__               | For local persistent volumes the node to bound.                       |
| __storage_class_name__ | Name of the storage class this pv is available for.                   |
| __access_modes__       | Array of access modes for this pv.                                    |
| __capacity__           | Capacity of this pv.                                                  |
| __init_script__        | For log informations the hint how to create the required directories. |


#### apps

| Field         | Description                                           |
|---------------|-------------------------------------------------------|
| __name__      | The name of the app/deployment                        |
| __catalog__   | The catalog to find the rancher chart in              |
| __template__  | The name of the rancher chart to be used              |
| __version__   | The version of the rancher chart                      |
| __namespace__ | The namespace to deploy the app in                    |
| __answers__   | The answers to the rancher questions as key-value map |

Example:
--------
```yaml
---
api_version: v1.0
kind: Project
metadata:
  name: my-wordpress-blog
namespaces:
  - name: my-wordpress-blog-web
storage_classes:
  - name: my-wordpress-blog-local-mariadb
    provisioner: kubernetes.io/no-provisioner
    reclaim_policy: Delete
    volume_bind_mode: WaitForFirstConsumer
persistent_volumes:
  - name: my-wordpress-blog-mariadb-node-1
    type: local
    node: node-1
    path: /var/data/my-wordpress-blog-mariadb
    capacity: "3Gi"
    access_modes:
      - "ReadWriteOnce"
    create_script: ssh ${node} sudo mkdir -p ${path}
apps:
- name: editorial-namespace
  catalog: library
  chart: wordpress
  version: "2.1.10"
  namespace: my-wordpress-blog-web
  answers:
    wordpressUsername: user
    wordpressPassword: ""
    wordpressEmail: user@example.com
    mariadb.enabled: true
    mariadb.db.name: wordpress
    mariadb.db.user: wordpress
    mariadb.master.persistence.enabled: 'true'
    mariadb.master.persistence.size: 8Gi
    mariadb.master.persistence.storageClass: "my-wordpress-blog-local-mariadb"
    ingress.enabled: false
    serviceType: ClusterIP
```
