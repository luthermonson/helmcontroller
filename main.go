//go:generate go run pkg/codegen/cleanup/main.go
//go:generate /bin/rm -rf pkg/generated
//go:generate go run pkg/codegen/main.go

package main

import (
	"context"
	batchv1 "github.com/rancher/helm-controller/pkg/generated/controllers/batch"
	corev1 "github.com/rancher/helm-controller/pkg/generated/controllers/core"
	helmv1 "github.com/rancher/helm-controller/pkg/generated/controllers/helm.cattle.io"
	rbacv1 "github.com/rancher/helm-controller/pkg/generated/controllers/rbac"
	helmcontroller "github.com/rancher/helm-controller/pkg/helm"
	"github.com/rancher/wrangler/pkg/apply"
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
	app.Name = "helm-controller"
	app.Version = VERSION
	app.Usage = "helm controller, to help with helm deployments."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "kubeconfig, k",
			EnvVar: "KUBECONFIG",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "master, m",
			EnvVar: "MASTERURL",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "namespace, n",
			EnvVar: "NAMESPACE",
			Value:  "",
		},
		cli.IntFlag{
			Name:   "threads, t",
			EnvVar: "THREADINESS",
			Value:  4,
		},
	}
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		klog.Fatal(err)
	}
}

func run(c *cli.Context) error {
	masterURL := c.String("master")
	kubeconfig := c.String("kubeconfig")
	namespace := c.String("namespace")
	threadiness := c.Int("threads")

	ctx := signals.SetupSignalHandler(context.Background())

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	helms, err := helmv1.NewFactoryFromConfigWithNamespace(cfg, namespace)
	if err != nil {
		klog.Fatalf("Error building sample controllers: %s", err.Error())
	}

	batches, err := batchv1.NewFactoryFromConfigWithNamespace(cfg, namespace)
	if err != nil {
		klog.Fatalf("Error building sample controllers: %s", err.Error())
	}

	rbacs, err := rbacv1.NewFactoryFromConfigWithNamespace(cfg, namespace)
	if err != nil {
		klog.Fatalf("Error building sample controllers: %s", err.Error())
	}

	cores, err := corev1.NewFactoryFromConfigWithNamespace(cfg, namespace)
	if err != nil {
		klog.Fatalf("Error building sample controllers: %s", err.Error())
	}

	discoverClient, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building discovery client: %s", err.Error())
	}

	objectSetApply := apply.New(discoverClient, apply.NewClientFactory(cfg))

	helmcontroller.Register(ctx, objectSetApply,
		helms.Helm().V1().HelmChart(),
		batches.Batch().V1().Job(),
		rbacs.Rbac().V1().ClusterRoleBinding(),
		cores.Core().V1().ServiceAccount())

	if err := start.All(ctx, threadiness, helms, batches, rbacs, cores); err != nil {
		klog.Fatalf("Error starting: %s", err.Error())
	}

	<-ctx.Done()
	return nil
}
