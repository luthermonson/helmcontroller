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
	time "time"

	versioned "github.com/luthermonson/helmcontroller/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/luthermonson/helmcontroller/pkg/generated/informers/externalversions/internalinterfaces"
	v1 "github.com/luthermonson/helmcontroller/pkg/generated/listers/helm.cattle.io/v1"
	helmcattleiov1 "github.com/luthermonson/helmcontroller/types/apis/helm.cattle.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// HelmChartInformer provides access to a shared informer and lister for
// HelmCharts.
type HelmChartInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.HelmChartLister
}

type helmChartInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewHelmChartInformer constructs a new informer for HelmChart type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewHelmChartInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredHelmChartInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredHelmChartInformer constructs a new informer for HelmChart type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredHelmChartInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.HelmV1().HelmCharts().List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.HelmV1().HelmCharts().Watch(options)
			},
		},
		&helmcattleiov1.HelmChart{},
		resyncPeriod,
		indexers,
	)
}

func (f *helmChartInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredHelmChartInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *helmChartInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&helmcattleiov1.HelmChart{}, f.defaultInformer)
}

func (f *helmChartInformer) Lister() v1.HelmChartLister {
	return v1.NewHelmChartLister(f.Informer().GetIndexer())
}
