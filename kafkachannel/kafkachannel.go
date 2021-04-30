package kafkachannel

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
	return schema.GroupVersionResource{Group: "messaging.knative.dev", Version: "v1beta1", Resource: "kafkachannels"}
}


func Install(name string, numPartitions, replicationFactor int, opts ...manifest.CfgFn) feature.StepFn {
	cfg := map[string]interface{}{
		"name": name,
	}

	cfg["numPartitions"] = numPartitions
	cfg["replicationFactor"] = replicationFactor

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
		Kind:       "KafkaChannel",
		APIVersion: "messaging.knative.dev/v1beta1",
		Name:       name,
	}
}

func IsReady(name string, timing ...time.Duration) feature.StepFn {
	return k8s.IsReady(Gvr(), name, timing...)
}