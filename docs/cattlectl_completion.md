## cattlectl completion

Generates bash completion scripts

### Synopsis



To configure your bash shell to load completions for each session add to your bashrc

* ~/.bashrc or ~/.profile

        . <(cattlectl completion)

* On Mac (with bash completion installed from brew)

        cattlectl completion > $(brew --prefix)/etc/bash_completion.d/cattlectl

* To load completion run

        . <(cattlectl completion)

This will only temporaly activate completion in the current session.


```
cattlectl completion [flags]
```

### Options

```
  -h, --help   help for completion
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

