Project Descriptor data model
=============================

ProjectDescriptor Structur:
---------------------------

### Toplevel ProjectDescriptor

| Field               | Description                                                     |
|---------------------|-----------------------------------------------------------------|
| __api_version__     | The __\<major\>.\<minor\>__ version used for this descriptor.   |
| __kind__            | The kind of descriptor in this file                             |
| __metadata__        | Metainformation about this descriptor e.g.: name and cluster_id |
| __apps__            | List of rancher apps to be deployed to this project             |
| __namespaces__      | List of namespaces to be part of this project                   |
| __storage_classes__ | List of storage classes required for this project               |

### metadata

| Field           | Description                                                                 |
|-----------------|-----------------------------------------------------------------------------|
| __name__        | The name of the project                                                     |
| __rancher_url__ | The URL to reach the rancher (placed from cattleclt configuration)          |
| __access_key__  | The access key to access rancher with (placed from cattleclt configuration) |
| __secret_key__  | The secret key to access rancher with (placed from cattleclt configuration) |
| __cluster_id__  | The ID of the cluster the project is part of                                |

#### apps

| Field         | Description                                           |
|---------------|-------------------------------------------------------|
| __name__      | The name of the app/deployment                        |
| __catalog__   | The catalog to find the rancher chart in              |
| __template__  | The name of the rancher chart to be used              |
| __version__   | The version of the rancher chart                      |
| __namespace__ | The namespace to deploy the app in                    |
| __answers__   | The answers to the rancher questions as key-value map |

#### namespaces

| Field    | Description               |
|----------|---------------------------|
| __name__ | The name of the namespace |

#### storage classes

| Field                         | Description                                                           |
|-------------------------------|-----------------------------------------------------------------------|
| __name__                      | The name of the storage class                                         |
| __provisioner__               | The provisioner used by this storage class                            |
| __reclaim_policy__            | The reclaim policy of this storage class                              |
| __volume_bind_mode__          | The volume bind mode of this storage class                            |
| __create_persistent_volumes__ | If persistent volumes should be precreated for this storage class     |
| __persistent_volume_groups__  | The list of persistent volume groups to create for this storage class |

#### persistent volume group

| Field             | Description               |
|-------------------|---------------------------|
| __name__          | The name of the persistent volume group                                 |
| __type__          | The Type of the persistent volume group (__only "local" is supportet__) |
| __path__          | For local persistent volumes the path pattern to use                    |
| __create_script__ | For log informations the hint how to create the required directories.   |
| __access_modes__  | List of access modes (string array)                                     |
| __capacity__      | The capacity of the persistent volumes to create                        |
| __nodes__         | For local persistent volumes the ist of nodes to bound (string array).  |

Example:
--------
```yaml
---
api_version: v1.0
kind: Project
metadata:
  name: {{template "full_project_name" .}}
  rancher_url: https://ui.rancher.server
  access_key: token-12345
  secret_key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  token_key: token-12345:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  cluster_id: j-4444
namespaces:
  - name: {{template "full_project_name" .}}-web
storage_classes:
  - name: {{template "full_project_name" .}}-local-mariadb
    provisioner: kubernetes.io/no-provisioner
    reclaim_policy: Delete
    volume_bind_mode: WaitForFirstConsumer
    create_persistent_volumes: {{.use_local_volumes}}
    persistent_volume_groups:
      - name: {{template "full_project_name" .}}-mariadb
        type: local
        path: /var/data/{{template "full_project_name" .}}-mariadb
        capacity: "3Gi"
        access_modes:
          - "ReadWriteOnce"
        create_script: ssh ${node} sudo mkdir -p ${path}
        nodes:
          - node-1
          - node-2
          - node-3
apps:
- name: editorial-namespace
  catalog: library
  chart: wordpress
  version: "2.1.10"
  namespace: {{template "full_project_name" .}}-web
  answers:
    wordpressUsername: user
    wordpressPassword: ""
    wordpressEmail: user@example.com
    mariadb.enabled: true
    mariadb.db.name: wordpress
    mariadb.db.user: wordpress
    mariadb.master.persistence.enabled: 'true'
    mariadb.master.persistence.size: 8Gi
    mariadb.master.persistence.storageClass: "{{template "full_project_name" .}}-local-mariadb"
    ingress.enabled: false
    serviceType: ClusterIP
    license: {{ read .license_file | base64}}
{{/*

Create a fully qualified project name.

*/}}
{{- define "full_project_name" -}}
  {{- if eq .stage "" -}}
    {{- print  .project_name -}}
  {{- else -}}
    {{- printf "%s-%s" .project_name .stage -}}
  {{- end -}}
{{- end -}}
```
