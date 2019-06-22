Cluster Descriptor data model
=============================

ClusterDescriptor Structur:
---------------------------

### Toplevel ClusterDescriptor

| Field           | Description                                                           |
|-----------------|-----------------------------------------------------------------------|
| __api_version__ | The __\<major\>.\<minor\>__ version used for this descriptor.         |
| __kind__        | The kind of descriptor in this file (`Project`)                       |
| __metadata__    | Metainformation about this descriptor e.g.: name and cluster_name     |
| __catalogs__    | List of namespaces to be part of this project                         |

### metadata

* In the descriptor only `name` should be set.
* All other fields are read from configuration or rancher

| Field            | Description                                                                              |
|------------------|------------------------------------------------------------------------------------------|
| __name__         | The name of the cluster **Optional or placed from cattleclt configuration**              |
| __id__           | Rancher internal ID of this project (**read from rancher**)                              |
| __rancher_url__  | The URL to reach the rancher (**placed from cattleclt configuration**)                   |
| __access_key__   | The access key to access rancher with (**placed from cattleclt configuration**)          |
| __secret_key__   | The secret key to access rancher with (**placed from cattleclt configuration**)          |
| __token_key__    | The token key to access rancher with (**placed from cattleclt configuration**)           |

#### catalogs

| Field        | Description                 |
|--------------|-----------------------------|
| __name__     | The name of the catalog     |
| __url__      | The URL of the catalog      |
| __branch__   | The branch of the catalog   |
| __username__ | The username of the catalog |
| __password__ | The password of the catalog |