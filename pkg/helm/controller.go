package helmcontroller

import (
	"context"
	"github.com/luthermonson/helmcontroller/pkg/generated/controllers/helm.cattle.io/v1"
	helmv1 "github.com/luthermonson/helmcontroller/types/apis/helm.cattle.io/v1"
	//typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	//"k8s.io/client-go/tools/record"
	"k8s.io/klog"
)

type Controller struct {
	helms 		v1.HelmChartController
}

const ControllerName = "helm.cattle.io"

// NewController returns a new sample controller
func Register(
	ctx context.Context,
	helms v1.HelmChartController) {

	controller := &Controller{
		helms:       helms,
	}

	helms.OnChange(ctx, ControllerName, controller.OnHelmChanged)
	helms.OnRemove(ctx, ControllerName, controller.OnHelmRemove)
}

func (c *Controller) OnHelmChanged(key string, helm *helmv1.HelmChart) (*helmv1.HelmChart, error) {
	klog.Infof("%s: On Helm Change", key)
	return helm, nil
}

func (c *Controller) OnHelmRemove (key string, helm *helmv1.HelmChart) (*helmv1.HelmChart, error) {
	klog.Infof("%s: On Helm Remove", key)
	return helm, nil
}
