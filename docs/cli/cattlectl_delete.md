## cattlectl delete

Deletes an rancher resouce

### Synopsis

Deletes an rancher resouce.

### Supported resource types:

* namespace
* certificate - NOT YET IMPLEMENTED
* config-map - NOT YET IMPLEMENTED
* docker-credential - NOT YET IMPLEMENTED
* secret - NOT YET IMPLEMENTED
* app
* job
* cron-job - NOT YET IMPLEMENTED
* deployment - NOT YET IMPLEMENTED
* daemon-set - NOT YET IMPLEMENTED
* stateful-set - NOT YET IMPLEMENTED

```
cattlectl delete KIND NAME [flags]
```

### Options

```
  -h, --help                  help for delete
      --namespace string      The namespace of the project to delete resouces from
      --project-name string   The name of the project to delete resouces from
```

### Options inherited from parent commands

```
      --access-key string     The access key to access rancher with
      --cluster-id string     The ID of the cluster the project is part of
      --cluster-name string   The name of the cluster the project is part of
      --config string         config file (default is $HOME/.cattlectl.yaml)
      --dry-run               if do dry-run
      --insecure-api          If Rancher uses a self signed certificate
      --log-json              if to log using json format
      --rancher-url string    The URL to reach the rancher
      --secret-key string     The secret key to access rancher with
  -v, --verbosity int         verbosity level to use
```

### SEE ALSO

* [cattlectl](cattlectl.md)	 - controll your cattle on the ranch

