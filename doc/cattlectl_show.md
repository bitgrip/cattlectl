## cattlectl show

Show the resulting project descriptor

### Synopsis

Show the resulting project descriptor after values.yaml is applyed

This is useful for debugging the template without interacting with an actual rancher.

```
cattlectl show [flags]
```

### Options

```
  -f, --file string     project file to show (default "project.yaml")
  -h, --help            help for show
      --values string   values file to show (default "values.yaml")
```

### Options inherited from parent commands

```
      --access-key string    The access key to access rancher with
      --cluster-id string    The ID of the cluster the project is part of
      --config string        config file (default is $HOME/.cattlectl.yaml)
      --log-json             if to log using json format
      --rancher-url string   The URL to reach the rancher
      --secret-key string    The secret key to access rancher with
  -v, --verbosity int        verbosity level to use
```

### SEE ALSO

* [cattlectl](cattlectl.md)	 - controll your cattle on the ranch

