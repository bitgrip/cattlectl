// Copyright © 2019 Bitgrip <berlin@bitgrip.de>
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

package rancher

import (
	"github.com/imdario/mergo"
)

// MergeProject merge two projects.
// the result represents the parent + all fields of child
// which have no equivalent Name in parent
func MergeProject(child Project, parent Project) (Project, error) {
	var dst Project
	if err := mergo.Merge(&dst, parent); err != nil {
		return parent, err
	}
	dst.Namespaces = mergeNamespaces(child.Namespaces, dst.Namespaces)
	dst.Resources.Certificates = mergeCertificates(child.Resources.Certificates, parent.Resources.Certificates)
	dst.Resources.ConfigMaps = mergeConfigMaps(child.Resources.ConfigMaps, parent.Resources.ConfigMaps)
	dst.Resources.DockerCredentials = mergeDockerCredentials(child.Resources.DockerCredentials, parent.Resources.DockerCredentials)
	dst.Resources.Secrets = mergeConfigMaps(child.Resources.Secrets, parent.Resources.Secrets)
	dst.StorageClasses = mergeStorageClasses(child.StorageClasses, parent.StorageClasses)
	dst.PersistentVolumes = mergePersistentVolumes(child.PersistentVolumes, parent.PersistentVolumes)
	dst.Apps = mergeApps(child.Apps, parent.Apps)
	return dst, nil
}

func mergeNamespaces(childNamespaces, parentNamespaces []Namespace) []Namespace {
	dst := parentNamespaces
CHILD_LOOP:
	for _, childNamespace := range childNamespaces {
		for _, parentNamespace := range parentNamespaces {
			if childNamespace.Name == parentNamespace.Name {
				continue CHILD_LOOP
			}
		}
		dst = append(dst, childNamespace)
	}
	return dst
}

func mergeCertificates(childCertificates, parentCertificates []Certificate) []Certificate {
	dst := parentCertificates
CHILD_LOOP:
	for _, childCertificate := range childCertificates {
		for _, parentCertificate := range parentCertificates {
			if childCertificate.Name == parentCertificate.Name {
				continue CHILD_LOOP
			}
		}
		dst = append(dst, childCertificate)
	}
	return dst
}

func mergeConfigMaps(childConfigMaps, parentConfigMaps []ConfigMap) []ConfigMap {
	dst := parentConfigMaps
CHILD_LOOP:
	for _, childConfigMap := range childConfigMaps {
		for _, parentConfigMap := range parentConfigMaps {
			if childConfigMap.Name == parentConfigMap.Name {
				continue CHILD_LOOP
			}
		}
		dst = append(dst, childConfigMap)
	}
	return dst
}

func mergeDockerCredentials(childDockerCredentials, parentDockerCredentials []DockerCredential) []DockerCredential {
	dst := parentDockerCredentials
CHILD_LOOP:
	for _, childDockerCredential := range childDockerCredentials {
		for _, parentDockerCredential := range parentDockerCredentials {
			if childDockerCredential.Name == parentDockerCredential.Name {
				continue CHILD_LOOP
			}
		}
		dst = append(dst, childDockerCredential)
	}
	return dst
}

func mergeStorageClasses(childStorageClasses, parentStorageClasses []StorageClass) []StorageClass {
	dst := parentStorageClasses
CHILD_LOOP:
	for _, childStorageClass := range childStorageClasses {
		for _, parentStorageClass := range parentStorageClasses {
			if childStorageClass.Name == parentStorageClass.Name {
				continue CHILD_LOOP
			}
		}
		dst = append(dst, childStorageClass)
	}
	return dst
}

func mergePersistentVolumes(childPersistentVolumes, parentPersistentVolumes []PersistentVolume) []PersistentVolume {
	dst := parentPersistentVolumes
CHILD_LOOP:
	for _, childPersistentVolume := range childPersistentVolumes {
		for _, parentPersistentVolume := range parentPersistentVolumes {
			if childPersistentVolume.Name == parentPersistentVolume.Name {
				continue CHILD_LOOP
			}
		}
		dst = append(dst, childPersistentVolume)
	}
	return dst
}

func mergeApps(childApps, parentApps []App) []App {
	dst := parentApps
CHILD_LOOP:
	for _, childApp := range childApps {
		for _, parentApp := range parentApps {
			if childApp.Name == parentApp.Name {
				continue CHILD_LOOP
			}
		}
		dst = append(dst, childApp)
	}
	return dst
}
