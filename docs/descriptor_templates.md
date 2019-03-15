Descriptor templates
====================

* All descriptors do use [golang tempaltes](https://golang.org/pkg/text/template/) which apply before the descriptors are parsed.
* The context is build from the YAML object in the optional `values.yaml` file.
* For each simple value there is a corresponding environment variable to alter the value.

IMPORTANT: All keys in the YAML are converted to lower case.
------------------------------------------------------------

* e.g.: If you have a key in your values.yaml `someKeyWithCamelCase` this will become:
* `somekeywithcamelcase` - in the template context. So `someKeyWithCamelCase` will not be found

Access context variables in the templates
-----------------------------------------

* You can access all values from the context in the template.
* The full key path is build by joining all key parts down the yaml structure with dots.
* You can address a value from the values.yaml with `{{ .<full-key-path> }}`

```yaml
project_name: a-simple-project
storage_class:
  azure:
    provisioner: kubernetes.io/azure-disk
    volume_binding_mode: Immediate
```

This can be accessed by:

* `project_name`
* `storage_class`
* `storage_class.azure`
* `storage_class.azure.provisioner`
* `storage_class.azure. volume_binding_mode`

Build in functions
------------------

| name | description | exsample |
|---|---|---|
|read|read the content of a file| `{{ read some-content.txt }}`|
|indent|indent each line with X YAML indents |`{{ read .tls_cert_key_file | indent 3}}`|
|base64|encodes the content as base 64|`{{ read .license_file | base64}}`|
|toYaml|inserts the complete yaml branch or array of the specified value|`{{ toYaml .my_value }}`|


Corresponding environment variables
-----------------------------------

* Each key which has a simple value (string, int, bool) can be changed by a corresponding environment variable.
* The name of the environment variable is build from the full key path.
  * All dots are replaced by `_`
  * The environment variable is all upper case

```yaml
project_name: a-simple-project
storage_class:
  azure:
    provisioner: kubernetes.io/azure-disk
    volume_binding_mode: Immediate
```

This can be changed by the environment variables:

* `PROJECT_NAME`
* `STORAGE_CLASS_AZURE_PROVISIONER`
* `STORAGE_CLASS_AZURE_VOLUME_BINDING_MODE`
