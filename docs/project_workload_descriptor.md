Project Workload Descriptor data model
======================================

__NEEDS IMPROVEMENT__

The workload descriptors represent the needed informations to communicate with the
Rancher API. The low level items are similar to the [Kubernetes API](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/).<br>
This documentation will not repeat the [Kubernetes API Documentation](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/) but will focus on the high level differences.

WorkloadDescriptor Structure:
-----------------------------

### Top level WorkloadDescriptor

| Field           | Description                                                                                            |
|-----------------|--------------------------------------------------------------------------------------------------------|
| __api_version__ | The __\<major\>.\<minor\>__ version used for this descriptor.                                          |
| __kind__        | The kind of descriptor in this file (on of `CronJob`, `DaemonSet`, `Deployment`, `Job`, `StatefulSet`) |
| __metadata__    | Metainformation about this descriptor e.g.: name and cluster_name                                      |
| __spec__        | Based on the __kind__ the corresponding workload spec                                                  |

### metadata

| Field            | Description                                                                              |
|------------------|------------------------------------------------------------------------------------------|
| __namespace__    | The namespace in the project **REQUIRED**                                                |
| __project_name__ | The name of the project **REQUIRED**                                                     |
| __cluster_name__ | The name of the cluster the project is part of (**placed from cattleclt configuration**) |
| __cluster_id__   | The ID of the cluster the project is part of (**read from rancher**)                     |
| __rancher_url__  | The URL to reach the rancher (**placed from cattleclt configuration**)                   |
| __access_key__   | The access key to access rancher with (**placed from cattleclt configuration**)          |
| __secret_key__   | The secret key to access rancher with (**placed from cattleclt configuration**)          |
| __token_key__    | The token key to access rancher with (**placed from cattleclt configuration**)           |

Workload Spec Common Members
----------------------------

### Common Members

All Workload Specs have this set of members.

| Field            | Description                                                                              |
|------------------|------------------------------------------------------------------------------------------|
| __activeDeadlineSeconds__ | |
| __annotations__ | |
| __automountServiceAccountToken__ | |
| __containers__ | |
| __dnsConfig__ | |
| __dnsPolicy__ | |
| __hostAliases__ | |
| __hostIPC__ | |
| __hostNetwork__ | |
| __hostPID__ | |
| __hostname__ | |
| __imagePullSecrets__ | |
| __labels__ | |
| __name__ | |
| __priority__ | |
| __priorityClassName__ | |
| __restartPolicy__ | |
| __runAsGroup__ | |
| __runAsNonRoot__ | |
| __schedulerName__ | |
| __scheduling__ | |
| __selector__ | |
| __serviceAccountName__ | |
| __shareProcessNamespace__ | |
| __subdomain__ | |
| __terminationGracePeriodSeconds__ | |
| __volumes__ | |
| __workloadAnnotations__ | |
| __workloadLabels__ | |

CronJob Spec
------------

### CronJob

| Field            | Description                                                                              |
|------------------|------------------------------------------------------------------------------------------|
| all common headers | |
| __cronJobConfig__ | |
| __TTLSecondsAfterFinished__ | |

### CronJobConfig

| Field            | Description                                                                              |
|------------------|------------------------------------------------------------------------------------------|
| __concurrencyPolicy__ | |
| __failedJobsHistoryLimit__ | |
| __jobAnnotations__ | |
| __jobConfig__ | |
| __jobLabels__ | |
| __schedule__ | |
| __startingDeadlineSeconds__ | |
| __successfulJobsHistoryLimit__ | |
| __suspend__ | |

DaemonSet Spec
------------

### DaemonSet

| Field            | Description                                                                              |
|------------------|------------------------------------------------------------------------------------------|
| all common headers | |
| __daemonSetConfig__ | |

### DaemonSetConfig

| Field            | Description                                                                              |
|------------------|------------------------------------------------------------------------------------------|
| __maxUnavailable__ | |
| __minReadySeconds__ | |
| __revisionHistoryLimit__ | |
| __strategy__ | |

Deployment Spec
------------

### Deployment

| Field            | Description                                                                              |
|------------------|------------------------------------------------------------------------------------------|
| all common headers | |
| __deploymentConfig__ | |
| __scale__ | |

### DeploymentConfig

| Field            | Description                                                                              |
|------------------|------------------------------------------------------------------------------------------|
| __maxSurge__ | |
| __maxUnavailable__ | |
| __minReadySeconds__ | |
| __progressDeadlineSeconds__ | |
| __revisionHistoryLimit__ | |
| __strategy__ | |

Job Spec
------------

### Job

| Field            | Description                                                                              |
|------------------|------------------------------------------------------------------------------------------|
| all common headers | |
| __jobConfig__ | |
| __TTLSecondsAfterFinished__ | |

### JobConfig

| Field            | Description                                                                              |
|------------------|------------------------------------------------------------------------------------------|
| __activeDeadlineSeconds__ | |
| __backoffLimit__ | |
| __completions__ | |
| __manualSelector__ | |
| __parallelism__ | |

StatefulSet Spec
------------

### StatefulSet

| Field            | Description                                                                              |
|------------------|------------------------------------------------------------------------------------------|
| all common headers | |
| __statefulSetConfig__ | |
| __scale__ | |

### StatefulSetConfig

| Field            | Description                                                                              |
|------------------|------------------------------------------------------------------------------------------|
| __partition__ | |
| __podManagementPolicy__ | |
| __revisionHistoryLimit__ | |
| __serviceName__ | |
| __strategy__ | |
| __volumeClaimTemplates__ | |

### PersistentVolumeClaim

| Field            | Description                                                                              |
|------------------|------------------------------------------------------------------------------------------|
| __accessModes__ | |
| __annotations__ | |
| __labels__ | |
| __name__ | |
| __resources__ | |
| __selector__ | |
| __storageClass__ | |





Example:
--------

### CronJob example from [automated-tasks-with-cron-jobs](https://kubernetes.io/docs/tasks/job/automated-tasks-with-cron-jobs/)

```yaml
---
api_version: v1.0
kind: CronJob
metadata:
  project_name: example-project
  namespace: example-namespace
spec:
  name: hello
  cronJobConfig:
    schedule: */1 * * * *
  containers:
  - name: hello
    image: busybox
    command:
    - /bin/sh
    - -c
    - date; echo Hello from the Kubernetes cluster
  restartPolicy: OnFailure
```
