Descriptor templates
====================

* All descriptors do use [golang tempaltes](https://golang.org/pkg/text/template/) which apply before the descriptors are parsed.
* The context is build from the YAML object in the optional `values.yaml` file.
* For each simple value there is a corresponding environment variable to alter the value.

IMPORTANT: All keys in the YAML are converted to lower case.
------------------------------------------------------------

Corresponding environment variables
-----------------------------------

Access context variables in the templates
-----------------------------------------

Build in functions
------------------
