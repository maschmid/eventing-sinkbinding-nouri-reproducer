package inmemorychannel

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime/schema"
	duck "knative.dev/pkg/apis/duck/v1"
	"knative.dev/reconciler-test/pkg/feature"
	"knative.dev/reconciler-test/pkg/k8s"
	"knative.dev/reconciler-test/pkg/manifest"
	"time"
)

func Gvr() schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: "messaging.knative.dev", Version: "v1", Resource: "inmemorychannels"}
}

func Install(name string, opts ...manifest.CfgFn) feature.StepFn {
	cfg := map[string]interface{}{
		"name": name,
	}

	for _, fn := range opts {
		fn(cfg)
	}
	return func(ctx context.Context, t feature.T) {
		if _, err := manifest.InstallLocalYaml(ctx, cfg); err != nil {
			t.Fatal(err)
		}
	}
}

func AsRef(name string) *duck.KReference {
	return &duck.KReference{
		Kind:       "InMemoryChannel",
		APIVersion: "messaging.knative.dev/v1",
		Name:       name,
	}
}

func IsReady(name string, timing ...time.Duration) feature.StepFn {
	return k8s.IsReady(Gvr(), name, timing...)
}