## cattlectl

controll your cattle on the ranch

### Synopsis

cattlectl is a CLI controller to your rancher instance.

It allowes you to describe your full project as code and lett it apply to your
cluster.

cattlectl handles the input in a idempotend way so that it dosn't change your
deployement, if you run cattlectl twice.

### Options

```
      --access-key string     The access key to access rancher with
      --cluster-id string     The ID of the cluster the project is part of
      --cluster-name string   The name of the cluster the project is part of
      --config string         config file (default is $HOME/.cattlectl.yaml)
  -h, --help                  help for cattlectl
      --insecure-api          If Rancher uses a self signed certificate
      --log-json              if to log using json format
      --rancher-url string    The URL to reach the rancher
      --secret-key string     The secret key to access rancher with
  -v, --verbosity int         verbosity level to use
```

### SEE ALSO

* [cattlectl apply](cattlectl_apply.md)	 - Apply a project descriptor to your rancher
* [cattlectl completion](cattlectl_completion.md)	 - Generates bash completion scripts
* [cattlectl delete](cattlectl_delete.md)	 - Deletes an rancher resouce
* [cattlectl gen-doc](cattlectl_gen-doc.md)	 - genrates the markdown documentation
* [cattlectl list](cattlectl_list.md)	 - Lists an rancher resouce
* [cattlectl show](cattlectl_show.md)	 - Show the resulting project descriptor
* [cattlectl version](cattlectl_version.md)	 - version of cattlectl

