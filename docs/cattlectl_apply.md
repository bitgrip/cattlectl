## cattlectl apply

Apply a project descriptor to your rancher

### Synopsis

Apply a deployment descriptor to your rancher.

### Descriptor Format

* The descriptor file (project.yaml by default) contains one or more YAML objects

### Supported deployment descriptors

* Project - A Rancher project and its members [more...](project_descriptor.md)
* Workload - A Rancher workload and its configuration [more...](project_workload_descriptor.md)
  * Job
  * CronJob
  * Deployment
  * DaemonSet
  * StatefulSet

### Directory layout

By default this command expects in the current folder:

| File             | Description                                                 |
|------------------|-------------------------------------------------------------|
| __project.yaml__ | The project descriptor with support for go template syntax. |
| __values.yaml__  | The set of values used to render the project descriptor.    |

See the [Project Descriptor data model](project_descriptor.md)

```
cattlectl apply [flags]
```

### Options

```
  -f, --file string      project file to apply (default "project.yaml")
  -h, --help             help for apply
      --merge-answers    If answers of existing apps should be merged with the new apply answers
      --values strings   values file(s) to apply (default [values.yaml])
```

### Options inherited from parent commands

```
      --access-key string     The access key to access rancher with
      --cluster-id string     The ID of the cluster the project is part of
      --cluster-name string   The name of the cluster the project is part of
      --config string         config file (default is $HOME/.cattlectl.yaml)
      --insecure-api          If Rancher uses a self signed certificate
      --log-json              if to log using json format
      --rancher-url string    The URL to reach the rancher
      --secret-key string     The secret key to access rancher with
  -v, --verbosity int         verbosity level to use
```

### SEE ALSO

* [cattlectl](cattlectl.md)	 - controll your cattle on the ranch

