cattleclt_apply
===============

Synopsis
--------

* Applies a descriptor to a Rancher server
* Uses a set of value files and a dict with highest precedence to execute the descriptor template

Parameters
----------

### Specific parameters

| Parameter | Choices/<span style="color:blue">Defaults</span> | Comments |
|---|---|---|
|file<br><span style="color:blue">string</span>| __Default:__<br><span style="color:blue">""</span> | The file containing the descriptor template to apply|
|value_files<br><span style="color:blue">list</span>| __Default:__<br><span style="color:blue">[]</span> | The set of value files to use when executing the descriptor template|
|values<br><span style="color:blue">dict</span>| __Default:__<br><span style="color:blue">{}</span> | Dict of values with highest precedence when executing the descriptor template|
|working_directory<br><span style="color:blue">string</span>| __Default:__<br><span style="color:blue">""</span> |If set all relative files are relative to `working_directory`<br>Relative to the playbook directory otherwais |

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
  - name: apply project descriptor
    cattlectl_apply:
      file: "descriptors/project.yaml"
      value_files:
        - values/general.yaml
        - values/{{cluster_name}}/values.yaml
        - values/{{ cluster_name }}/{{ project_name }}.yaml
      values:
        my_extra_value: "extra-value"
```
