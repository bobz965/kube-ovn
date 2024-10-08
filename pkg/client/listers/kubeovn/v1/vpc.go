/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// VpcLister helps list Vpcs.
// All objects returned here must be treated as read-only.
type VpcLister interface {
	// List lists all Vpcs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Vpc, err error)
	// Get retrieves the Vpc from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Vpc, error)
	VpcListerExpansion
}

// vpcLister implements the VpcLister interface.
type vpcLister struct {
	listers.ResourceIndexer[*v1.Vpc]
}

// NewVpcLister returns a new VpcLister.
func NewVpcLister(indexer cache.Indexer) VpcLister {
	return &vpcLister{listers.New[*v1.Vpc](indexer, v1.Resource("vpc"))}
}
