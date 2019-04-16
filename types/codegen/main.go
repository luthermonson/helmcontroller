package main

import (
	"github.com/rancher/helmcontroller/types/apis/helm.cattle.io/v1"
	"github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
	batch "k8s.io/api/batch/v1"
)

func main() {
	controllergen.Run(args.Options{
		OutputPackage: "github.com/rancher/helmcontroller/pkg/generated",
		Boilerplate:   "hack/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"helm.cattle.io": {
				Types: []interface{}{
					v1.HelmChart{},
				},
				GenerateTypes: true,
			},
			"batch": {
				Types: []interface{}{
					batch.Job{},
				},
				InformersPackage: "k8s.io/client-go/informers",
				ClientSetPackage: "k8s.io/client-go/kubernetes",
				ListersPackage:   "k8s.io/client-go/listers",
			},
		},
	})
}