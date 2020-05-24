cattleclt_delete
================

Synopsis
--------

* Deletes on or more resources from a Rancher server

Parameters
----------

### Specific parameters

| Parameter | Choices/<span style="color:blue">Defaults</span> | Comments |
|---|---|---|
| project_name<br><span style="color:blue">string</span> | __Default:__<br><span style="color:blue">""</span> | The project to delete resources from<br> ignored for kind project |
| namespace<br><span style="color:blue">string</span> | __Default:__<br><span style="color:blue">""</span> | The namespace to delete resources from<br>ignored for kind namespace and project |
| kind<br><span style="color:blue">string</span> | __Default:__<br><span style="color:blue">""</span> | The kind of resource to delete |
| names<br><span style="color:blue">list</span> | __Default:__<br><span style="color:blue">[]</span> | The names of the resources to delete |

### General parameters

| Parameter | Choices/<span style="color:blue">Defaults</span> | Comments |
|---|---|---|
| rancher_url<br><span style="color:blue">string</span> | __Default:__<br><span style="color:blue">""</span> | The URL to access the Rancher server<br>Read from `config_file` if absent |
| insecure_api<br><span style="color:blue">boolean</span> | __Choices:__<br><span style="color:blue">no ←</span><br>yes | If cattlectl is going to accept insecure api<br>Read from `config_file` if absent |
| ca_certs<br><span style="color:blue">string</span> | __Default:__<br><span style="color:blue">""</span> | Certs of a private CA if requierd<br>Read from `config_file` if absent |
| access_key<br><span style="color:blue">string</span> | __Default:__<br><span style="color:blue">""</span> | Access key to gain access to the Rancher server<br>Read from `config_file` if absent |
| secret_key<br><span style="color:blue">string</span> | __Default:__<br><span style="color:blue">""</span> | Secret key to authenticate the `access_key`<br>Read from `config_file` if absent |
| cluster_name<br><span style="color:blue">string</span> | __Default:__<br><span style="color:blue">""</span> | The name of the Cluster to access via Rancher<br>Read from `config_file` if absent |
| config_file<br><span style="color:blue">string</span> | __Default:__<br><span style="color:blue">~/.cattlectl.yaml</span>| The location of the cattlectl config file to use |
| dry_run<br><span style="color:blue">boolean</span> | __Choices:__<br><span style="color:blue">no ←</span><br>yes | If true all `write` operations are logged only |

Examples
--------

```yaml
- name: delete namespaces
  cattlectl_delete:
    project_name: my-project
    kind: namespace
    names:
      - my-namespace01
      - my-namespace02
```
