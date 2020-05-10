## cattlectl list

Lists an rancher resouce

### Synopsis

Lists an rancher resouce

```
cattlectl list KIND [flags]
```

### Options

```
  -h, --help                  help for list
      --namespace string      The namespace of the project to list resouces from
      --pattern string        Match pattern to filter resouce names
      --project-name string   The name of the project to list resouces from
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

