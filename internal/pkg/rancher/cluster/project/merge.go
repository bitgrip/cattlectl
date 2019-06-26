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

package project

import (
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/project/model"
)

// MergeProject merge two projects.
// the result represents the parent + all fields of child
// which have no equivalent Name in parent
func MergeProject(child projectModel.Project, parent *projectModel.Project) error {
	parent.Namespaces = mergeNamespaces(child.Namespaces, parent.Namespaces)
	parent.Resources.Certificates = mergeCertificates(child.Resources.Certificates, parent.Resources.Certificates)
	parent.Resources.ConfigMaps = mergeConfigMaps(child.Resources.ConfigMaps, parent.Resources.ConfigMaps)
	parent.Resources.DockerCredentials = mergeDockerCredentials(child.Resources.DockerCredentials, parent.Resources.DockerCredentials)
	parent.Resources.Secrets = mergeConfigMaps(child.Resources.Secrets, parent.Resources.Secrets)
	parent.StorageClasses = mergeStorageClasses(child.StorageClasses, parent.StorageClasses)
	parent.PersistentVolumes = mergePersistentVolumes(child.PersistentVolumes, parent.PersistentVolumes)
	parent.Apps = mergeApps(child.Apps, parent.Apps)
	return nil
}

func mergeNamespaces(childNamespaces, parentNamespaces []projectModel.Namespace) []projectModel.Namespace {
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

func mergeCertificates(childCertificates, parentCertificates []projectModel.Certificate) []projectModel.Certificate {
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

func mergeConfigMaps(childConfigMaps, parentConfigMaps []projectModel.ConfigMap) []projectModel.ConfigMap {
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

func mergeDockerCredentials(childDockerCredentials, parentDockerCredentials []projectModel.DockerCredential) []projectModel.DockerCredential {
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

func mergeStorageClasses(childStorageClasses, parentStorageClasses []projectModel.StorageClass) []projectModel.StorageClass {
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

func mergePersistentVolumes(childPersistentVolumes, parentPersistentVolumes []projectModel.PersistentVolume) []projectModel.PersistentVolume {
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

func mergeApps(childApps, parentApps []projectModel.App) []projectModel.App {
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
