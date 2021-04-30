package ksvc

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime/schema"
	duck "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/tracker"
	"knative.dev/reconciler-test/pkg/feature"
	"knative.dev/reconciler-test/pkg/k8s"
	"knative.dev/reconciler-test/pkg/manifest"
	"time"
)

func Gvr() schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: "serving.knative.dev", Version: "v1", Resource: "services"}
}

func Install(name string, image string, opts ...manifest.CfgFn) feature.StepFn {
	cfg := map[string]interface{}{
		"name": name,
	}

	cfg["image"] = image

	for _, fn := range opts {
		fn(cfg)
	}
	return func(ctx context.Context, t feature.T) {
		if _, err := manifest.InstallLocalYaml(ctx, cfg); err != nil {
			t.Fatal(err)
		}
	}
}

func AsTrackerReference(name string) *tracker.Reference {
	return &tracker.Reference{
		Kind:       "Service",
		Name:       name,
		APIVersion: "serving.knative.dev/v1",
	}
}

func AsRef(name string) *duck.KReference {
	return &duck.KReference{
		Kind:       "Service",
		APIVersion: "serving.knative.dev/v1",
		Name:       name,
	}
}

func IsReady(name string, timing ...time.Duration) feature.StepFn {
	return k8s.IsReady(Gvr(), name, timing...)
}