package helm

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	helmcontroller "github.com/rancher/helmcontroller/pkg/generated/controllers/helm.cattle.io/v1"
	batchcontroller "github.com/rancher/helmcontroller/pkg/generated/controllers/batch/v1"
	helmv1 "github.com/rancher/helmcontroller/types/apis/helm.cattle.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/rancher/wrangler/pkg/relatedresource"
	"github.com/rancher/wrangler/pkg/objectset"
	"github.com/rancher/wrangler/pkg/apply"
	rbac "k8s.io/api/rbac/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	trueVal = true
)

type Controller struct {
	helms 		helmcontroller.HelmChartController
	jobs 		batchcontroller.JobController
	apply		apply.Apply
}

const (
	namespace = "kube-system"
	image     = "rancher/klipper-helm:v0.1.5"
	label     = "helm.k3s.cattle.io/chart"
	name      = "helm.cattle.io"
)

func Register(ctx context.Context, apply apply.Apply, helms helmcontroller.HelmChartController, jobs batchcontroller.JobController) {
	relatedresource.Watch(ctx, "helm-pod-watch",
		func(namespace, name string, obj runtime.Object) ([]relatedresource.Key, error) {
			if job, ok := obj.(*batch.Job); ok {
				name := job.Labels[label]
				if name != "" {
					return []relatedresource.Key{
						{
							Name:      name,
							Namespace: namespace,
						},
					}, nil
				}
			}
			return nil, nil
		},
		helms,
		jobs)

	controller := &Controller{
		helms: helms,
		jobs:  jobs,
		apply: apply,
	}

	helms.OnChange(ctx, name, controller.OnHelmChanged)
	helms.OnRemove(ctx, name, controller.OnHelmRemove)
}

func (c *Controller) OnHelmChanged(key string, chart *helmv1.HelmChart) (*helmv1.HelmChart, error) {
	if chart.Namespace != namespace || chart.Spec.Chart == "" {
		return chart, nil
	}

	objs := objectset.NewObjectSet()
	job, configMap := job(chart)
	objs.Add(serviceAccount(chart))
	objs.Add(roleBinding(chart))
	objs.Add(job)
	if configMap != nil {
		objs.Add(configMap)
	}

	if err := c.apply.WithOwner(chart).Apply(objs); err != nil {
		return chart, err
	}

	chartCopy := chart.DeepCopy()
	chartCopy.Status.JobName = job.Name
	return c.helms.Update(chartCopy)
}

func (c *Controller) OnHelmRemove (key string, chart *helmv1.HelmChart) (*helmv1.HelmChart, error) {
	if chart.Namespace != namespace || chart.Spec.Chart == "" {
		return chart, nil
	}

	job, _ := job(chart)
	job, err := c.jobs.Cache().Get(chart.Namespace, job.Name)
	if errors.IsNotFound(err) {
		_, err := c.OnHelmChanged(key, chart)
		if err != nil {
			return chart, err
		}
	} else if err != nil {
		return nil, err
	}

	if job.Status.Succeeded <= 0 {
		return nil, fmt.Errorf("waiting for delete of helm chart %s", chart.Name)
	}

	return chart, c.apply.WithOwner(chart).Apply(objectset.NewObjectSet())
}

func job(chart *helmv1.HelmChart) (*batch.Job, *core.ConfigMap) {
	oneThousand := int32(1000)
	valuesHash := sha256.Sum256([]byte(chart.Spec.ValuesContent))

	action := "install"
	if chart.DeletionTimestamp != nil {
		action = "delete"
	}
	job := &batch.Job{
		TypeMeta: meta.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      fmt.Sprintf("helm-%s-%s", action, chart.Name),
			Namespace: chart.Namespace,
			Labels: map[string]string{
				label: chart.Name,
			},
		},
		Spec: batch.JobSpec{
			BackoffLimit: &oneThousand,
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{
						label: chart.Name,
					},
				},
				Spec: core.PodSpec{
					RestartPolicy: core.RestartPolicyOnFailure,
					Containers: []core.Container{
						{
							Name:            "helm",
							Image:           image,
							ImagePullPolicy: core.PullIfNotPresent,
							Args:            args(chart),
							Env: []core.EnvVar{
								{
									Name:  "NAME",
									Value: chart.Name,
								},
								{
									Name:  "VERSION",
									Value: chart.Spec.Version,
								},
								{
									Name:  "REPO",
									Value: chart.Spec.Repo,
								},
								{
									Name:  "VALUES_HASH",
									Value: hex.EncodeToString(valuesHash[:]),
								},
							},
						},
					},
					ServiceAccountName: fmt.Sprintf("helm-%s", chart.Name),
				},
			},
		},
	}

	configMap := configMap(chart)
	if configMap == nil {
		return job, nil
	}

	job.Spec.Template.Spec.Volumes = []core.Volume{
		{
			Name: "values",
			VolumeSource: core.VolumeSource{
				ConfigMap: &core.ConfigMapVolumeSource{
					LocalObjectReference: core.LocalObjectReference{
						Name: configMap.Name,
					},
				},
			},
		},
	}

	job.Spec.Template.Spec.Containers[0].VolumeMounts = []core.VolumeMount{
		{
			MountPath: "/config",
			Name:      "values",
		},
	}

	return job, configMap
}

func configMap(chart *helmv1.HelmChart) *core.ConfigMap {
	if chart.Spec.ValuesContent == "" {
		return nil
	}

	return &core.ConfigMap{
		TypeMeta: meta.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      fmt.Sprintf("chart-values-%s", chart.Name),
			Namespace: chart.Namespace,
		},
		Data: map[string]string{
			"values.yaml": chart.Spec.ValuesContent,
		},
	}
}

func roleBinding(chart *helmv1.HelmChart) *rbac.ClusterRoleBinding {
	return &rbac.ClusterRoleBinding{
		TypeMeta: meta.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRoleBinding",
		},
		ObjectMeta: meta.ObjectMeta{
			Name: fmt.Sprintf("helm-%s-%s", chart.Namespace, chart.Name),
		},
		RoleRef: rbac.RoleRef{
			Kind:     "ClusterRole",
			APIGroup: "rbac.authorization.k8s.io",
			Name:     "cluster-admin",
		},
		Subjects: []rbac.Subject{
			{
				Name:      fmt.Sprintf("helm-%s", chart.Name),
				Kind:      "ServiceAccount",
				Namespace: chart.Namespace,
			},
		},
	}

}

func serviceAccount(chart *helmv1.HelmChart) *core.ServiceAccount {
	return &core.ServiceAccount{
		TypeMeta: meta.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      fmt.Sprintf("helm-%s", chart.Name),
			Namespace: chart.Namespace,
		},
		AutomountServiceAccountToken: &trueVal,
	}
}

func args(chart *helmv1.HelmChart) []string {
	if chart.DeletionTimestamp != nil {
		return []string{
			"delete",
			"--purge", chart.Name,
		}
	}

	spec := chart.Spec
	args := []string{
		"install",
		"--name", chart.Name,
		spec.Chart,
	}
	if spec.TargetNamespace != "" {
		args = append(args, "--namespace", spec.TargetNamespace)
	}
	if spec.Repo != "" {
		args = append(args, "--repo", spec.Repo)
	}
	if spec.Version != "" {
		args = append(args, "--version", spec.Version)
	}

	for _, k := range keys(spec.Set) {
		val := spec.Set[k]
		if val.StrVal != "" {
			args = append(args, "--set-string", fmt.Sprintf("%s=%s", k, val.StrVal))
		} else {
			args = append(args, "--set", fmt.Sprintf("%s=%d", k, val.IntVal))
		}
	}

	return args
}

func keys(val map[string]intstr.IntOrString) []string {
	var keys []string
	for k := range val {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
