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

package descriptor

import (
	"fmt"

	"github.com/bitgrip/cattlectl/internal/pkg/rancher"
)

func ClusterResourceDescriptorConverger(clusterName string, partConvergers []Converger) Converger {
	return DescriptorConverger{
		InitCluster: DefaultInitCluster(
			clusterName,
		),
		InitProject: func(client rancher.Client) error {
			return nil
		},
		PartConvergers: partConvergers,
	}
}

func ProjectResourceDescriptorConverger(clusterName, projectName string, partConvergers []Converger) Converger {
	return DescriptorConverger{
		InitCluster: DefaultInitCluster(
			clusterName,
		),
		InitProject: DefaultInitProject(
			projectName,
		),
		PartConvergers: partConvergers,
	}
}

// Converger is a object which can converge github.com/bitgrip/cattlectl/internal/pkg/projectModel.Project
type Converger interface {
	Converge(client rancher.Client) error
}

type DescriptorConverger struct {
	InitCluster    func(client rancher.Client) error
	InitProject    func(client rancher.Client) error
	PartConvergers []Converger
}

func (converger DescriptorConverger) Converge(client rancher.Client) error {
	if err := converger.InitCluster(client); err != nil {
		return err
	}
	if err := converger.InitProject(client); err != nil {
		return err
	}
	for _, partConverger := range converger.PartConvergers {
		if err := partConverger.Converge(client); err != nil {
			return err
		}
	}
	return nil
}

type PartConverger struct {
	PartName   string
	HasPart    func(client rancher.Client) (bool, error)
	CreatePart func(client rancher.Client) error
	UpdatePart func(client rancher.Client) error
}

func (converger PartConverger) Converge(client rancher.Client) error {
	if hasPart, err := converger.HasPart(client); hasPart {
		return converger.UpdatePart(client)
	} else if err != nil {
		return fmt.Errorf("Failed to check for %s, %v", converger.PartName, err)
	}
	return converger.CreatePart(client)
}

func DefaultInitCluster(clusterName string) func(client rancher.Client) error {
	return func(client rancher.Client) error {
		return rancher.InitCluster(
			"",
			clusterName,
			client,
		)
	}
}

func DefaultInitProject(projectName string) func(client rancher.Client) error {
	return func(client rancher.Client) error {
		return rancher.InitProject(
			projectName,
			client,
		)
	}
}
