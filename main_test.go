package eventing_sinkbinding_nouri_reproducer

import (
	"flag"
	"knative.dev/pkg/injection"
	"knative.dev/reconciler-test/pkg/environment"
	"os"
	"testing"
)

var global environment.GlobalEnvironment

func init() {
	// environment.InitFlags registers state and level filter flags.
	environment.InitFlags(flag.CommandLine)
}

func TestMain(m *testing.M) {
	// We get a chance to parse flags to include the framework flags for the
	// framework as well as any additional flags included in the integration.
	flag.Parse()

	// EnableInjectionOrDie will enable client injection, this is used by the
	// testing framework for namespace management, and could be leveraged by
	// features to pull Kubernetes clients or the test environment out of the
	// context passed in the features.
	ctx, startInformers := injection.EnableInjectionOrDie(nil, nil) //nolint
	startInformers()

	// global is used to make instances of Environments, NewGlobalEnvironment
	// is passing and saving the client injection enabled context for use later.
	global = environment.NewGlobalEnvironment(ctx)

	// Run the tests.
	os.Exit(m.Run())
}
