//go:generate go run types/codegen/cleanup/main.go
//go:generate rm -rf ./pkg/generated/*
//go:generate go run types/codegen/main.go

package main

import (
	"context"
	"flag"
	batchv1 "github.com/rancher/helmcontroller/pkg/generated/controllers/batch"
	helmv1 "github.com/rancher/helmcontroller/pkg/generated/controllers/helm.cattle.io"
	helmcontroller "github.com/rancher/helmcontroller/pkg/helm"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/crd"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/rancher/wrangler/pkg/start"
	"github.com/urfave/cli"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"os"
)

var (
	VERSION = "v0.0.0-dev"
)

func main() {
	app := cli.NewApp()
	app.Name = "helmcontroller"
	app.Version = VERSION
	app.Usage = "helmcontroller needs help!"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "kubeconfig, k",
			EnvVar: "KUBECONFIG",
			Value:  "$HOME/.kube/config",
		},
		cli.StringFlag{
			Name:   "master, m",
			EnvVar: "MASTERURL",
			Value:  "",
		},
	}
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		klog.Fatal(err)
	}
}

func run(c *cli.Context) error {

	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klogFlags.Set("logtostderr", "true")
	klog.InitFlags(klogFlags)

	masterURL := c.String("master")
	kubeconfig := c.String("kubeconfig")
	ctx := signals.SetupSignalHandler(context.Background())

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	helmController, err := helmv1.NewFactoryFromConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building sample controllers: %s", err.Error())
	}

	batchController, err := batchv1.NewFactoryFromConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building sample controllers: %s", err.Error())
	}

	crds, err := crd.NewFactoryFromClient(cfg)
	if err != nil {
		klog.Fatalf("Error creating CRD: %s", err.Error())
	}

	crd := crd.NamespacedType("HelmChart.helm.cattle.io/v1")
	_, err = crds.CreateCRDs(ctx, crd)

	if err != nil {
		klog.Fatalf("Error building crd: %s", err.Error())
	}

	discoverClient, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building discovery client: %s", err.Error())
	}

	objectSetApply := apply.New(discoverClient, apply.NewClientFactory(cfg), helmController.Helm().V1().HelmChart(), batchController.Batch().V1().Job())

	helmcontroller.Register(ctx, objectSetApply, helmController.Helm().V1().HelmChart(), batchController.Batch().V1().Job())

	if err := start.All(ctx, 4, helmController); err != nil {
		klog.Fatalf("Error starting: %s", err.Error())
	}

	<-ctx.Done()
	return nil
}
