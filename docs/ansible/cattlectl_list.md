cattleclt_list
===============

Synopsis
--------

Parameters
----------

| Parameter | Choices/<span style="color:blue">Defaults</span> | Comments |
|---|---|---|
| rancher_url | | The URL to access the Rancher server<br>Read from `config_file` if absent |
| insecure_api<br><span style="color:blue">boolean</span> | __Choices:__<br><span style="color:blue">no ←</span><br>yes | If cattlectl is going to accept insecure api<br>Read from `config_file` if absent |
| ca_certs | | Certs of a private CA if requierd<br>Read from `config_file` if absent |
| access_key | | Access key to gain access to the Rancher server<br>Read from `config_file` if absent |
| secret_key | | Secret key to authenticate the `access_key`<br>Read from `config_file` if absent |
| cluster_name | | The name of the Cluster to access via Rancher<br>Read from `config_file` if absent |
| config_file | __Default:__<br><span style="color:blue">~/.cattlectl.yaml</span>| The location of the cattlectl config file to use |
| dry_run<br><span style="color:blue">boolean</span> | __Choices:__<br><span style="color:blue">no ←</span><br>yes | If true all `write` operations are logged only |

Examples
--------
