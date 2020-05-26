// Copyright © 2018 Bitgrip <berlin@bitgrip.de>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

type WorkloadMetadata struct {
	ProjectName string `yaml:"project_name,omitempty"`
	ProjectID   string `yaml:"project_id,omitempty"`
	Namespace   string `yaml:"namespace,omitempty"`
	NamespaceID string `yaml:"namespace_id,omitempty"`
	RancherURL  string `yaml:"rancher_url,omitempty"`
	AccessKey   string `yaml:"access_key,omitempty"`
	SecretKey   string `yaml:"secret_key,omitempty"`
	TokenKey    string `yaml:"token_key,omitempty"`
	ClusterName string `yaml:"cluster_name,omitempty"`
	ClusterID   string `yaml:"cluster_id,omitempty"`
}

type baseWorkload struct {
	ActiveDeadlineSeconds         *int64                 `json:"activeDeadlineSeconds,omitempty" yaml:"activeDeadlineSeconds,omitempty"`
	Annotations                   map[string]string      `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	AutomountServiceAccountToken  *bool                  `json:"automountServiceAccountToken,omitempty" yaml:"automountServiceAccountToken,omitempty"`
	Containers                    []Container            `json:"containers,omitempty" yaml:"containers,omitempty"`
	DNSConfig                     *PodDNSConfig          `json:"dnsConfig,omitempty" yaml:"dnsConfig,omitempty"`
	DNSPolicy                     string                 `json:"dnsPolicy,omitempty" yaml:"dnsPolicy,omitempty"`
	EnableServiceLinks            *bool                  `json:"enableServiceLinks,omitempty" yaml:"enableServiceLinks,omitempty"`
	HostAliases                   []HostAlias            `json:"hostAliases,omitempty" yaml:"hostAliases,omitempty"`
	HostIPC                       bool                   `json:"hostIPC,omitempty" yaml:"hostIPC,omitempty"`
	HostNetwork                   bool                   `json:"hostNetwork,omitempty" yaml:"hostNetwork,omitempty"`
	HostPID                       bool                   `json:"hostPID,omitempty" yaml:"hostPID,omitempty"`
	Hostname                      string                 `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	ImagePullSecrets              []LocalObjectReference `json:"imagePullSecrets,omitempty" yaml:"imagePullSecrets,omitempty"`
	Labels                        map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name                          string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Overhead                      map[string]string      `json:"overhead,omitempty" yaml:"overhead,omitempty"`
	PreemptionPolicy              string                 `json:"preemptionPolicy,omitempty" yaml:"preemptionPolicy,omitempty"`
	RestartPolicy                 string                 `json:"restartPolicy,omitempty" yaml:"restartPolicy,omitempty"`
	RunAsGroup                    *int64                 `json:"runAsGroup,omitempty" yaml:"runAsGroup,omitempty"`
	RunAsNonRoot                  *bool                  `json:"runAsNonRoot,omitempty" yaml:"runAsNonRoot,omitempty"`
	Scheduling                    *Scheduling            `json:"scheduling,omitempty" yaml:"scheduling,omitempty"`
	Selector                      *LabelSelector         `json:"selector,omitempty" yaml:"selector,omitempty"`
	ServiceAccountName            string                 `json:"serviceAccountName,omitempty" yaml:"serviceAccountName,omitempty"`
	ShareProcessNamespace         *bool                  `json:"shareProcessNamespace,omitempty" yaml:"shareProcessNamespace,omitempty"`
	Subdomain                     string                 `json:"subdomain,omitempty" yaml:"subdomain,omitempty"`
	TerminationGracePeriodSeconds *int64                 `json:"terminationGracePeriodSeconds,omitempty" yaml:"terminationGracePeriodSeconds,omitempty"`
	Volumes                       []Volume               `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	WorkloadAnnotations           map[string]string      `json:"workloadAnnotations,omitempty" yaml:"workloadAnnotations,omitempty"`
	WorkloadLabels                map[string]string      `json:"workloadLabels,omitempty" yaml:"workloadLabels,omitempty"`
}

// Container are created in the context of a workload
type Container struct {
	AllowPrivilegeEscalation *bool                 `json:"allowPrivilegeEscalation,omitempty" yaml:"allowPrivilegeEscalation,omitempty"`
	Command                  []string              `json:"command,omitempty" yaml:"command,omitempty"`
	Entrypoint               []string              `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
	Environment              map[string]string     `json:"environment,omitempty" yaml:"environment,omitempty"`
	EnvironmentFrom          []EnvironmentFrom     `json:"environmentFrom,omitempty" yaml:"environmentFrom,omitempty"`
	Image                    string                `json:"image,omitempty" yaml:"image,omitempty"`
	ImagePullPolicy          string                `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
	InitContainer            bool                  `json:"initContainer,omitempty" yaml:"initContainer,omitempty"`
	LivenessProbe            *Probe                `json:"livenessProbe,omitempty" yaml:"livenessProbe,omitempty"`
	Name                     string                `json:"name,omitempty" yaml:"name,omitempty"`
	Ports                    []ContainerPort       `json:"ports,omitempty" yaml:"ports,omitempty"`
	PostStart                *Handler              `json:"postStart,omitempty" yaml:"postStart,omitempty"`
	PreStop                  *Handler              `json:"preStop,omitempty" yaml:"preStop,omitempty"`
	Privileged               *bool                 `json:"privileged,omitempty" yaml:"privileged,omitempty"`
	ProcMount                string                `json:"procMount,omitempty" yaml:"procMount,omitempty"`
	ReadOnly                 *bool                 `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	ReadinessProbe           *Probe                `json:"readinessProbe,omitempty" yaml:"readinessProbe,omitempty"`
	Resources                *ResourceRequirements `json:"resources,omitempty" yaml:"resources,omitempty"`
	RunAsGroup               *int64                `json:"runAsGroup,omitempty" yaml:"runAsGroup,omitempty"`
	RunAsNonRoot             *bool                 `json:"runAsNonRoot,omitempty" yaml:"runAsNonRoot,omitempty"`
	Stdin                    bool                  `json:"stdin,omitempty" yaml:"stdin,omitempty"`
	StdinOnce                bool                  `json:"stdinOnce,omitempty" yaml:"stdinOnce,omitempty"`
	TTY                      bool                  `json:"tty,omitempty" yaml:"tty,omitempty"`
	TerminationMessagePath   string                `json:"terminationMessagePath,omitempty" yaml:"terminationMessagePath,omitempty"`
	TerminationMessagePolicy string                `json:"terminationMessagePolicy,omitempty" yaml:"terminationMessagePolicy,omitempty"`
	VolumeDevices            []VolumeDevice        `json:"volumeDevices,omitempty" yaml:"volumeDevices,omitempty"`
	VolumeMounts             []VolumeMount         `json:"volumeMounts,omitempty" yaml:"volumeMounts,omitempty"`
	WorkingDir               string                `json:"workingDir,omitempty" yaml:"workingDir,omitempty"`
}

// ContainerPort represents a network port in a single container.
type ContainerPort struct {
	ContainerPort int64  `json:"containerPort,omitempty" yaml:"containerPort,omitempty"`
	DNSName       string `json:"dnsName,omitempty" yaml:"dnsName,omitempty"`
	HostIP        string `json:"hostIp,omitempty" yaml:"hostIp,omitempty"`
	Kind          string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	Protocol      string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	SourcePort    int64  `json:"sourcePort,omitempty" yaml:"sourcePort,omitempty"`
}

// EnvironmentFrom represents the source of a set of ConfigMaps
type EnvironmentFrom struct {
	Optional   bool   `json:"optional,omitempty" yaml:"optional,omitempty"`
	Prefix     string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	Source     string `json:"source,omitempty" yaml:"source,omitempty"`
	SourceKey  string `json:"sourceKey,omitempty" yaml:"sourceKey,omitempty"`
	SourceName string `json:"sourceName,omitempty" yaml:"sourceName,omitempty"`
	TargetKey  string `json:"targetKey,omitempty" yaml:"targetKey,omitempty"`
}

// Handler defines a specific action that should be taken
type Handler struct {
	Command     []string     `json:"command,omitempty" yaml:"command,omitempty"`
	HTTPHeaders []HTTPHeader `json:"httpHeaders,omitempty" yaml:"httpHeaders,omitempty"`
	Host        string       `json:"host,omitempty" yaml:"host,omitempty"`
	Path        string       `json:"path,omitempty" yaml:"path,omitempty"`
	Port        IntOrString  `json:"port,omitempty" yaml:"port,omitempty"`
	Scheme      string       `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	TCP         bool         `json:"tcp,omitempty" yaml:"tcp,omitempty"`
}

// HostAlias holds the mapping between IP and hostnames that will be injected as an entry in the pod's hosts file.
type HostAlias struct {
	Hostnames []string `json:"hostnames,omitempty" yaml:"hostnames,omitempty"`
	IP        string   `json:"ip,omitempty" yaml:"ip,omitempty"`
}

// HTTPHeader describes a custom header to be used in HTTP probes
type HTTPHeader struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

// A LabelSelector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
type LabelSelector struct {
	MatchExpressions []LabelSelectorRequirement `json:"matchExpressions,omitempty" yaml:"matchExpressions,omitempty"`
	MatchLabels      map[string]string          `json:"matchLabels,omitempty" yaml:"matchLabels,omitempty"`
}

// A LabelSelectorRequirement is a selector that contains values, a key, and an operator that relates the key and values.
type LabelSelectorRequirement struct {
	Key      string   `json:"key,omitempty" yaml:"key,omitempty"`
	Operator string   `json:"operator,omitempty" yaml:"operator,omitempty"`
	Values   []string `json:"values,omitempty" yaml:"values,omitempty"`
}

// LocalObjectReference contains enough information to let you locate the referenced object inside the same namespace.
type LocalObjectReference struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type NodeScheduling struct {
	NodeID     string   `json:"nodeId,omitempty" yaml:"nodeId,omitempty"`
	Preferred  []string `json:"preferred,omitempty" yaml:"preferred,omitempty"`
	RequireAll []string `json:"requireAll,omitempty" yaml:"requireAll,omitempty"`
	RequireAny []string `json:"requireAny,omitempty" yaml:"requireAny,omitempty"`
}

// PodDNSConfig defines the DNS parameters of a pod in addition to those generated from DNSPolicy.
type PodDNSConfig struct {
	Nameservers []string             `json:"nameservers,omitempty" yaml:"nameservers,omitempty"`
	Options     []PodDNSConfigOption `json:"options,omitempty" yaml:"options,omitempty"`
	Searches    []string             `json:"searches,omitempty" yaml:"searches,omitempty"`
}

// PodDNSConfigOption defines DNS resolver options of a pod.
type PodDNSConfigOption struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

// Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.
type Probe struct {
	Command             []string     `json:"command,omitempty" yaml:"command,omitempty"`
	FailureThreshold    int64        `json:"failureThreshold,omitempty" yaml:"failureThreshold,omitempty"`
	HTTPHeaders         []HTTPHeader `json:"httpHeaders,omitempty" yaml:"httpHeaders,omitempty"`
	Host                string       `json:"host,omitempty" yaml:"host,omitempty"`
	InitialDelaySeconds int64        `json:"initialDelaySeconds,omitempty" yaml:"initialDelaySeconds,omitempty"`
	Path                string       `json:"path,omitempty" yaml:"path,omitempty"`
	PeriodSeconds       int64        `json:"periodSeconds,omitempty" yaml:"periodSeconds,omitempty"`
	Port                IntOrString  `json:"port,omitempty" yaml:"port,omitempty"`
	Scheme              string       `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	SuccessThreshold    int64        `json:"successThreshold,omitempty" yaml:"successThreshold,omitempty"`
	TCP                 bool         `json:"tcp,omitempty" yaml:"tcp,omitempty"`
	TimeoutSeconds      int64        `json:"timeoutSeconds,omitempty" yaml:"timeoutSeconds,omitempty"`
}

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements struct {
	Limits   map[string]string `json:"limits,omitempty" yaml:"limits,omitempty"`
	Requests map[string]string `json:"requests,omitempty" yaml:"requests,omitempty"`
}

type Scheduling struct {
	Node              *NodeScheduling `json:"node,omitempty" yaml:"node,omitempty"`
	Priority          *int64          `json:"priority,omitempty" yaml:"priority,omitempty"`
	PriorityClassName string          `json:"priorityClassName,omitempty" yaml:"priorityClassName,omitempty"`
	Scheduler         string          `json:"scheduler,omitempty" yaml:"scheduler,omitempty"`
	Tolerate          []Toleration    `json:"tolerate,omitempty" yaml:"tolerate,omitempty"`
}

// Toleration of the pod this Toleration is attached to tolerates any taint that matches the triple <key,value,effect> using the matching operator <operator>.
type Toleration struct {
	Effect            string `json:"effect,omitempty" yaml:"effect,omitempty"`
	Key               string `json:"key,omitempty" yaml:"key,omitempty"`
	Operator          string `json:"operator,omitempty" yaml:"operator,omitempty"`
	TolerationSeconds *int64 `json:"tolerationSeconds,omitempty" yaml:"tolerationSeconds,omitempty"`
	Value             string `json:"value,omitempty" yaml:"value,omitempty"`
}

// VolumeDevice describes a mapping of a raw block device within a container.
type VolumeDevice struct {
	DevicePath string `json:"devicePath,omitempty" yaml:"devicePath,omitempty"`
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
}

// VolumeMount describes a mounting of a Volume within a container.
type VolumeMount struct {
	MountPath        string `json:"mountPath,omitempty" yaml:"mountPath,omitempty"`
	MountPropagation string `json:"mountPropagation,omitempty" yaml:"mountPropagation,omitempty"`
	Name             string `json:"name,omitempty" yaml:"name,omitempty"`
	ReadOnly         bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	SubPath          string `json:"subPath,omitempty" yaml:"subPath,omitempty"`
}

type Volume struct {
	AWSElasticBlockStore  *AWSElasticBlockStoreVolumeSource  `json:"awsElasticBlockStore,omitempty" yaml:"awsElasticBlockStore,omitempty"`
	AzureDisk             *AzureDiskVolumeSource             `json:"azureDisk,omitempty" yaml:"azureDisk,omitempty"`
	AzureFile             *AzureFileVolumeSource             `json:"azureFile,omitempty" yaml:"azureFile,omitempty"`
	CephFS                *CephFSVolumeSource                `json:"cephfs,omitempty" yaml:"cephfs,omitempty"`
	Cinder                *CinderVolumeSource                `json:"cinder,omitempty" yaml:"cinder,omitempty"`
	ConfigMap             *ConfigMapVolumeSource             `json:"configMap,omitempty" yaml:"configMap,omitempty"`
	DownwardAPI           *DownwardAPIVolumeSource           `json:"downwardAPI,omitempty" yaml:"downwardAPI,omitempty"`
	EmptyDir              *EmptyDirVolumeSource              `json:"emptyDir,omitempty" yaml:"emptyDir,omitempty"`
	FC                    *FCVolumeSource                    `json:"fc,omitempty" yaml:"fc,omitempty"`
	FlexVolume            *FlexVolumeSource                  `json:"flexVolume,omitempty" yaml:"flexVolume,omitempty"`
	Flocker               *FlockerVolumeSource               `json:"flocker,omitempty" yaml:"flocker,omitempty"`
	GCEPersistentDisk     *GCEPersistentDiskVolumeSource     `json:"gcePersistentDisk,omitempty" yaml:"gcePersistentDisk,omitempty"`
	GitRepo               *GitRepoVolumeSource               `json:"gitRepo,omitempty" yaml:"gitRepo,omitempty"`
	Glusterfs             *GlusterfsVolumeSource             `json:"glusterfs,omitempty" yaml:"glusterfs,omitempty"`
	HostPath              *HostPathVolumeSource              `json:"hostPath,omitempty" yaml:"hostPath,omitempty"`
	ISCSI                 *ISCSIVolumeSource                 `json:"iscsi,omitempty" yaml:"iscsi,omitempty"`
	NFS                   *NFSVolumeSource                   `json:"nfs,omitempty" yaml:"nfs,omitempty"`
	Name                  string                             `json:"name,omitempty" yaml:"name,omitempty"`
	PersistentVolumeClaim *PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty" yaml:"persistentVolumeClaim,omitempty"`
	PhotonPersistentDisk  *PhotonPersistentDiskVolumeSource  `json:"photonPersistentDisk,omitempty" yaml:"photonPersistentDisk,omitempty"`
	PortworxVolume        *PortworxVolumeSource              `json:"portworxVolume,omitempty" yaml:"portworxVolume,omitempty"`
	Projected             *ProjectedVolumeSource             `json:"projected,omitempty" yaml:"projected,omitempty"`
	Quobyte               *QuobyteVolumeSource               `json:"quobyte,omitempty" yaml:"quobyte,omitempty"`
	RBD                   *RBDVolumeSource                   `json:"rbd,omitempty" yaml:"rbd,omitempty"`
	ScaleIO               *ScaleIOVolumeSource               `json:"scaleIO,omitempty" yaml:"scaleIO,omitempty"`
	Secret                *SecretVolumeSource                `json:"secret,omitempty" yaml:"secret,omitempty"`
	StorageOS             *StorageOSVolumeSource             `json:"storageos,omitempty" yaml:"storageos,omitempty"`
	VsphereVolume         *VsphereVirtualDiskVolumeSource    `json:"vsphereVolume,omitempty" yaml:"vsphereVolume,omitempty"`
}

type AWSElasticBlockStoreVolumeSource struct {
	FSType    string `json:"fsType,omitempty" yaml:"fsType,omitempty"`
	Partition int64  `json:"partition,omitempty" yaml:"partition,omitempty"`
	ReadOnly  bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	VolumeID  string `json:"volumeID,omitempty" yaml:"volumeID,omitempty"`
}

type AzureDiskVolumeSource struct {
	CachingMode string `json:"cachingMode,omitempty" yaml:"cachingMode,omitempty"`
	DataDiskURI string `json:"diskURI,omitempty" yaml:"diskURI,omitempty"`
	DiskName    string `json:"diskName,omitempty" yaml:"diskName,omitempty"`
	FSType      string `json:"fsType,omitempty" yaml:"fsType,omitempty"`
	Kind        string `json:"kind,omitempty" yaml:"kind,omitempty"`
	ReadOnly    *bool  `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
}

type AzureFileVolumeSource struct {
	ReadOnly   bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	SecretName string `json:"secretName,omitempty" yaml:"secretName,omitempty"`
	ShareName  string `json:"shareName,omitempty" yaml:"shareName,omitempty"`
}

type CephFSVolumeSource struct {
	Monitors   []string              `json:"monitors,omitempty" yaml:"monitors,omitempty"`
	Path       string                `json:"path,omitempty" yaml:"path,omitempty"`
	ReadOnly   bool                  `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	SecretFile string                `json:"secretFile,omitempty" yaml:"secretFile,omitempty"`
	SecretRef  *LocalObjectReference `json:"secretRef,omitempty" yaml:"secretRef,omitempty"`
	User       string                `json:"user,omitempty" yaml:"user,omitempty"`
}

type CinderVolumeSource struct {
	FSType    string                `json:"fsType,omitempty" yaml:"fsType,omitempty"`
	ReadOnly  bool                  `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	SecretRef *LocalObjectReference `json:"secretRef,omitempty" yaml:"secretRef,omitempty"`
	VolumeID  string                `json:"volumeID,omitempty" yaml:"volumeID,omitempty"`
}

type ConfigMapProjection struct {
	Items    []KeyToPath `json:"items,omitempty" yaml:"items,omitempty"`
	Name     string      `json:"name,omitempty" yaml:"name,omitempty"`
	Optional *bool       `json:"optional,omitempty" yaml:"optional,omitempty"`
}

type ConfigMapVolumeSource struct {
	DefaultMode *int64      `json:"defaultMode,omitempty" yaml:"defaultMode,omitempty"`
	Items       []KeyToPath `json:"items,omitempty" yaml:"items,omitempty"`
	Name        string      `json:"name,omitempty" yaml:"name,omitempty"`
	Optional    *bool       `json:"optional,omitempty" yaml:"optional,omitempty"`
}

type DownwardAPIVolumeSource struct {
	DefaultMode *int64                  `json:"defaultMode,omitempty" yaml:"defaultMode,omitempty"`
	Items       []DownwardAPIVolumeFile `json:"items,omitempty" yaml:"items,omitempty"`
}

type DownwardAPIProjection struct {
	Items []DownwardAPIVolumeFile `json:"items,omitempty" yaml:"items,omitempty"`
}

type DownwardAPIVolumeFile struct {
	FieldRef         *ObjectFieldSelector   `json:"fieldRef,omitempty" yaml:"fieldRef,omitempty"`
	Mode             *int64                 `json:"mode,omitempty" yaml:"mode,omitempty"`
	Path             string                 `json:"path,omitempty" yaml:"path,omitempty"`
	ResourceFieldRef *ResourceFieldSelector `json:"resourceFieldRef,omitempty" yaml:"resourceFieldRef,omitempty"`
}

type EmptyDirVolumeSource struct {
	Medium    string `json:"medium,omitempty" yaml:"medium,omitempty"`
	SizeLimit string `json:"sizeLimit,omitempty" yaml:"sizeLimit,omitempty"`
}

type FCVolumeSource struct {
	FSType     string   `json:"fsType,omitempty" yaml:"fsType,omitempty"`
	Lun        *int64   `json:"lun,omitempty" yaml:"lun,omitempty"`
	ReadOnly   bool     `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	TargetWWNs []string `json:"targetWWNs,omitempty" yaml:"targetWWNs,omitempty"`
	WWIDs      []string `json:"wwids,omitempty" yaml:"wwids,omitempty"`
}

type FlexVolumeSource struct {
	Driver    string                `json:"driver,omitempty" yaml:"driver,omitempty"`
	FSType    string                `json:"fsType,omitempty" yaml:"fsType,omitempty"`
	Options   map[string]string     `json:"options,omitempty" yaml:"options,omitempty"`
	ReadOnly  bool                  `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	SecretRef *LocalObjectReference `json:"secretRef,omitempty" yaml:"secretRef,omitempty"`
}

type FlockerVolumeSource struct {
	DatasetName string `json:"datasetName,omitempty" yaml:"datasetName,omitempty"`
	DatasetUUID string `json:"datasetUUID,omitempty" yaml:"datasetUUID,omitempty"`
}

type GCEPersistentDiskVolumeSource struct {
	FSType    string `json:"fsType,omitempty" yaml:"fsType,omitempty"`
	PDName    string `json:"pdName,omitempty" yaml:"pdName,omitempty"`
	Partition int64  `json:"partition,omitempty" yaml:"partition,omitempty"`
	ReadOnly  bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
}

type GitRepoVolumeSource struct {
	Directory  string `json:"directory,omitempty" yaml:"directory,omitempty"`
	Repository string `json:"repository,omitempty" yaml:"repository,omitempty"`
	Revision   string `json:"revision,omitempty" yaml:"revision,omitempty"`
}

type GlusterfsVolumeSource struct {
	EndpointsName string `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	Path          string `json:"path,omitempty" yaml:"path,omitempty"`
	ReadOnly      bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
}

type HostPathVolumeSource struct {
	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
}

type ISCSIVolumeSource struct {
	DiscoveryCHAPAuth bool                  `json:"chapAuthDiscovery,omitempty" yaml:"chapAuthDiscovery,omitempty"`
	FSType            string                `json:"fsType,omitempty" yaml:"fsType,omitempty"`
	IQN               string                `json:"iqn,omitempty" yaml:"iqn,omitempty"`
	ISCSIInterface    string                `json:"iscsiInterface,omitempty" yaml:"iscsiInterface,omitempty"`
	InitiatorName     string                `json:"initiatorName,omitempty" yaml:"initiatorName,omitempty"`
	Lun               int64                 `json:"lun,omitempty" yaml:"lun,omitempty"`
	Portals           []string              `json:"portals,omitempty" yaml:"portals,omitempty"`
	ReadOnly          bool                  `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	SecretRef         *LocalObjectReference `json:"secretRef,omitempty" yaml:"secretRef,omitempty"`
	SessionCHAPAuth   bool                  `json:"chapAuthSession,omitempty" yaml:"chapAuthSession,omitempty"`
	TargetPortal      string                `json:"targetPortal,omitempty" yaml:"targetPortal,omitempty"`
}

type KeyToPath struct {
	Key  string `json:"key,omitempty" yaml:"key,omitempty"`
	Mode *int64 `json:"mode,omitempty" yaml:"mode,omitempty"`
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
}

type NFSVolumeSource struct {
	Path     string `json:"path,omitempty" yaml:"path,omitempty"`
	ReadOnly bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	Server   string `json:"server,omitempty" yaml:"server,omitempty"`
}

type ObjectFieldSelector struct {
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	FieldPath  string `json:"fieldPath,omitempty" yaml:"fieldPath,omitempty"`
}

type PersistentVolumeClaimVolumeSource struct {
	PersistentVolumeClaimID string `json:"persistentVolumeClaimId,omitempty" yaml:"persistentVolumeClaimId,omitempty"`
	ReadOnly                bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
}

type PhotonPersistentDiskVolumeSource struct {
	FSType string `json:"fsType,omitempty" yaml:"fsType,omitempty"`
	PdID   string `json:"pdID,omitempty" yaml:"pdID,omitempty"`
}

type PortworxVolumeSource struct {
	FSType   string `json:"fsType,omitempty" yaml:"fsType,omitempty"`
	ReadOnly bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	VolumeID string `json:"volumeID,omitempty" yaml:"volumeID,omitempty"`
}

type ProjectedVolumeSource struct {
	DefaultMode *int64             `json:"defaultMode,omitempty" yaml:"defaultMode,omitempty"`
	Sources     []VolumeProjection `json:"sources,omitempty" yaml:"sources,omitempty"`
}

type ResourceFieldSelector struct {
	ContainerName string `json:"containerName,omitempty" yaml:"containerName,omitempty"`
	Divisor       string `json:"divisor,omitempty" yaml:"divisor,omitempty"`
	Resource      string `json:"resource,omitempty" yaml:"resource,omitempty"`
}

type QuobyteVolumeSource struct {
	Group    string `json:"group,omitempty" yaml:"group,omitempty"`
	ReadOnly bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	Registry string `json:"registry,omitempty" yaml:"registry,omitempty"`
	User     string `json:"user,omitempty" yaml:"user,omitempty"`
	Volume   string `json:"volume,omitempty" yaml:"volume,omitempty"`
}

type RBDVolumeSource struct {
	CephMonitors []string              `json:"monitors,omitempty" yaml:"monitors,omitempty"`
	FSType       string                `json:"fsType,omitempty" yaml:"fsType,omitempty"`
	Keyring      string                `json:"keyring,omitempty" yaml:"keyring,omitempty"`
	RBDImage     string                `json:"image,omitempty" yaml:"image,omitempty"`
	RBDPool      string                `json:"pool,omitempty" yaml:"pool,omitempty"`
	RadosUser    string                `json:"user,omitempty" yaml:"user,omitempty"`
	ReadOnly     bool                  `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	SecretRef    *LocalObjectReference `json:"secretRef,omitempty" yaml:"secretRef,omitempty"`
}

type ScaleIOVolumeSource struct {
	FSType           string                `json:"fsType,omitempty" yaml:"fsType,omitempty"`
	Gateway          string                `json:"gateway,omitempty" yaml:"gateway,omitempty"`
	ProtectionDomain string                `json:"protectionDomain,omitempty" yaml:"protectionDomain,omitempty"`
	ReadOnly         bool                  `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	SSLEnabled       bool                  `json:"sslEnabled,omitempty" yaml:"sslEnabled,omitempty"`
	SecretRef        *LocalObjectReference `json:"secretRef,omitempty" yaml:"secretRef,omitempty"`
	StorageMode      string                `json:"storageMode,omitempty" yaml:"storageMode,omitempty"`
	StoragePool      string                `json:"storagePool,omitempty" yaml:"storagePool,omitempty"`
	System           string                `json:"system,omitempty" yaml:"system,omitempty"`
	VolumeName       string                `json:"volumeName,omitempty" yaml:"volumeName,omitempty"`
}

type SecretProjection struct {
	Items    []KeyToPath `json:"items,omitempty" yaml:"items,omitempty"`
	Name     string      `json:"name,omitempty" yaml:"name,omitempty"`
	Optional *bool       `json:"optional,omitempty" yaml:"optional,omitempty"`
}

type SecretVolumeSource struct {
	DefaultMode *int64      `json:"defaultMode,omitempty" yaml:"defaultMode,omitempty"`
	Items       []KeyToPath `json:"items,omitempty" yaml:"items,omitempty"`
	Optional    *bool       `json:"optional,omitempty" yaml:"optional,omitempty"`
	SecretName  string      `json:"secretName,omitempty" yaml:"secretName,omitempty"`
}

type ServiceAccountTokenProjection struct {
	Audience          string `json:"audience,omitempty" yaml:"audience,omitempty"`
	ExpirationSeconds *int64 `json:"expirationSeconds,omitempty" yaml:"expirationSeconds,omitempty"`
	Path              string `json:"path,omitempty" yaml:"path,omitempty"`
}

type StorageOSVolumeSource struct {
	FSType          string                `json:"fsType,omitempty" yaml:"fsType,omitempty"`
	ReadOnly        bool                  `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	SecretRef       *LocalObjectReference `json:"secretRef,omitempty" yaml:"secretRef,omitempty"`
	VolumeName      string                `json:"volumeName,omitempty" yaml:"volumeName,omitempty"`
	VolumeNamespace string                `json:"volumeNamespace,omitempty" yaml:"volumeNamespace,omitempty"`
}

type VolumeProjection struct {
	ConfigMap           *ConfigMapProjection           `json:"configMap,omitempty" yaml:"configMap,omitempty"`
	DownwardAPI         *DownwardAPIProjection         `json:"downwardAPI,omitempty" yaml:"downwardAPI,omitempty"`
	Secret              *SecretProjection              `json:"secret,omitempty" yaml:"secret,omitempty"`
	ServiceAccountToken *ServiceAccountTokenProjection `json:"serviceAccountToken,omitempty" yaml:"serviceAccountToken,omitempty"`
}

type VsphereVirtualDiskVolumeSource struct {
	FSType            string `json:"fsType,omitempty" yaml:"fsType,omitempty"`
	StoragePolicyID   string `json:"storagePolicyID,omitempty" yaml:"storagePolicyID,omitempty"`
	StoragePolicyName string `json:"storagePolicyName,omitempty" yaml:"storagePolicyName,omitempty"`
	VolumePath        string `json:"volumePath,omitempty" yaml:"volumePath,omitempty"`
}
