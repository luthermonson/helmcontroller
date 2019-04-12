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
	"github.com/luthermonson/helmcontroller/pkg/generated/clientset/versioned/scheme"
	v1 "github.com/luthermonson/helmcontroller/types/apis/helm.cattle.io/v1"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type HelmV1Interface interface {
	RESTClient() rest.Interface
	HelmChartsGetter
}

// HelmV1Client is used to interact with features provided by the helm.cattle.io group.
type HelmV1Client struct {
	restClient rest.Interface
}

func (c *HelmV1Client) HelmCharts() HelmChartInterface {
	return newHelmCharts(c)
}

// NewForConfig creates a new HelmV1Client for the given config.
func NewForConfig(c *rest.Config) (*HelmV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &HelmV1Client{client}, nil
}

// NewForConfigOrDie creates a new HelmV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *HelmV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new HelmV1Client for the given RESTClient.
func New(c rest.Interface) *HelmV1Client {
	return &HelmV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *HelmV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
