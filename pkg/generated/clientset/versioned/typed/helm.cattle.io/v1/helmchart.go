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

// Code generated by main. DO NOT EDIT.

package v1

import (
	"time"

	scheme "github.com/luthermonson/helmcontroller/pkg/generated/clientset/versioned/scheme"
	v1 "github.com/luthermonson/helmcontroller/types/apis/helm.cattle.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// HelmChartsGetter has a method to return a HelmChartInterface.
// A group's client should implement this interface.
type HelmChartsGetter interface {
	HelmCharts() HelmChartInterface
}

// HelmChartInterface has methods to work with HelmChart resources.
type HelmChartInterface interface {
	Create(*v1.HelmChart) (*v1.HelmChart, error)
	Update(*v1.HelmChart) (*v1.HelmChart, error)
	UpdateStatus(*v1.HelmChart) (*v1.HelmChart, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.HelmChart, error)
	List(opts metav1.ListOptions) (*v1.HelmChartList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.HelmChart, err error)
	HelmChartExpansion
}

// helmCharts implements HelmChartInterface
type helmCharts struct {
	client rest.Interface
}

// newHelmCharts returns a HelmCharts
func newHelmCharts(c *HelmV1Client) *helmCharts {
	return &helmCharts{
		client: c.RESTClient(),
	}
}

// Get takes name of the helmChart, and returns the corresponding helmChart object, and an error if there is any.
func (c *helmCharts) Get(name string, options metav1.GetOptions) (result *v1.HelmChart, err error) {
	result = &v1.HelmChart{}
	err = c.client.Get().
		Resource("helmcharts").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of HelmCharts that match those selectors.
func (c *helmCharts) List(opts metav1.ListOptions) (result *v1.HelmChartList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.HelmChartList{}
	err = c.client.Get().
		Resource("helmcharts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested helmCharts.
func (c *helmCharts) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("helmcharts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a helmChart and creates it.  Returns the server's representation of the helmChart, and an error, if there is any.
func (c *helmCharts) Create(helmChart *v1.HelmChart) (result *v1.HelmChart, err error) {
	result = &v1.HelmChart{}
	err = c.client.Post().
		Resource("helmcharts").
		Body(helmChart).
		Do().
		Into(result)
	return
}

// Update takes the representation of a helmChart and updates it. Returns the server's representation of the helmChart, and an error, if there is any.
func (c *helmCharts) Update(helmChart *v1.HelmChart) (result *v1.HelmChart, err error) {
	result = &v1.HelmChart{}
	err = c.client.Put().
		Resource("helmcharts").
		Name(helmChart.Name).
		Body(helmChart).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *helmCharts) UpdateStatus(helmChart *v1.HelmChart) (result *v1.HelmChart, err error) {
	result = &v1.HelmChart{}
	err = c.client.Put().
		Resource("helmcharts").
		Name(helmChart.Name).
		SubResource("status").
		Body(helmChart).
		Do().
		Into(result)
	return
}

// Delete takes name of the helmChart and deletes it. Returns an error if one occurs.
func (c *helmCharts) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("helmcharts").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *helmCharts) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("helmcharts").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched helmChart.
func (c *helmCharts) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.HelmChart, err error) {
	result = &v1.HelmChart{}
	err = c.client.Patch(pt).
		Resource("helmcharts").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
