// Copyright 2018 The rethinkdb-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Copyright 2018 The vault-operator Authors
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

package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultBaseImage = "jmckind/rethinkdb"
	defaultVersion = "latest"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RethinkDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []RethinkDB `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RethinkDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              RethinkDBSpec   `json:"spec"`
	Status            RethinkDBStatus `json:"status,omitempty"`
}

type RethinkDBSpec struct {
	// Number of nodes to deploy for a RethinkDB deployment.
	// Default: 1.
	Size int32 `json:"nodes,omitempty"`

	// Base image to use for a RethinkDB deployment.
	BaseImage string `json:"baseImage"`

	// Version of RethinkDB to be deployed.
	Version string `json:"version"`

	// Flag to indicate whether or not a Service will be created for the deployment.
	WithService bool `json:"withService"`

	// Name of ConfigMap to use or create.
	ConfigMapName string `json:"configMapName"`

	// Name of Secret to use or create.
	SecretName string `json:"secretName"`

	// Pod defines the policy for pods owned by rethinkdb operator.
	// This field cannot be updated once the CR is created.
	Pod *PodPolicy `json:"pod,omitempty"`
}

type RethinkDBStatus struct {
	// StatefulSet is the name of the rethinkdb StatefulSet
	StatefulSet string `json:"nodes"`

	// Pods are the names of the rethinkdb pods
	Pods []string `json:"nodes"`
}

// PodPolicy defines the policy for pods owned by rethinkdb operator.
type PodPolicy struct {
	// Resources is the resource requirements for the containers.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`

	// PersistentVolumeClaimSpec is the spec to describe PVC for the rethinkdb container
	// This field is optional. If no PVC spec, rethinkdb container will use emptyDir as volume
	PersistentVolumeClaimSpec *v1.PersistentVolumeClaimSpec `json:"persistentVolumeClaimSpec,omitempty"`
}

// SetDefaults sets the default vaules for the cuberite spec and returns true if the spec was changed
func (r *RethinkDB) SetDefaults() bool {
	changed := false
	rs := &r.Spec
	if rs.Size == 0 {
		rs.Size = 1
		changed = true
	}
	if len(rs.BaseImage) == 0 {
		rs.BaseImage = defaultBaseImage
		changed = true
	}
	if len(rs.Version) == 0 {
		rs.Version = defaultVersion
		changed = true
	}
	if len(rs.ConfigMapName) == 0 {
		rs.ConfigMapName = r.Name
		changed = true
	}
	if len(rs.SecretName) == 0 {
		rs.SecretName = r.Name
		changed = true
	}
	return changed
}

func (r *RethinkDB) IsPVEnabled() bool {
	if podPolicy := r.Spec.Pod; podPolicy != nil {
		return podPolicy.PersistentVolumeClaimSpec != nil
	}
	return false
}
