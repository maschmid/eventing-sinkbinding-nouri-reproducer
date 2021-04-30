package eventing_sinkbinding_nouri_reproducer

import (
	"github.com/maschmid/eventing-sinkbinding-nouri-reproducer/inmemorychannel"
	"github.com/maschmid/eventing-sinkbinding-nouri-reproducer/kafkachannel"
	"github.com/maschmid/eventing-sinkbinding-nouri-reproducer/ksvc"
	"knative.dev/eventing/test/rekt/resources/sinkbinding"
	"knative.dev/eventing/test/rekt/resources/subscription"
	duck "knative.dev/pkg/apis/duck/v1"
	"knative.dev/reconciler-test/pkg/environment"
	"knative.dev/reconciler-test/pkg/feature"
	"testing"
)

// The images are not really important, as long as they deploy with the provided envs,
// no events are actually sent during the scenario, so we just use any image that
// becomes ready only with K_SINK provided, so we verify that the sinkbound ksvc has K_SINK injected properly
const senderImage = "quay.io/maschmid/sender:1.0.0"
const eventDisplayImage = "quay.io/openshift-knative/knative-eventing-sources-event-display:v0.13.2"

func ThingThatCausesFutureImcSinkBindingsToHangWithNoUri() *feature.Feature {

	f := feature.NewFeatureNamed("Thing that causes future SinkBinding with IMC destination to hang with 'No URI'")

	eventDisplay := "event-display"

	inputChannel := "input"
	outputChannel := "output"

	sinkBinding := "binding"

	sender := "sender"

	f.Setup("Create Event display ksvc", ksvc.Install(eventDisplay, eventDisplayImage))

	f.Setup("Create KafkaChannel", kafkachannel.Install(inputChannel, 1, 1))
	f.Setup("Create KafkaChannel", kafkachannel.Install(outputChannel, 1, 1))

	f.Setup("KafkaChannel becomes Ready", kafkachannel.IsReady(inputChannel))
	f.Setup("KafkaChannel becomes Ready", kafkachannel.IsReady(outputChannel))

	f.Setup("Create KafkaChannel SinkBinding", sinkbinding.Install(sinkBinding,
		&duck.Destination{Ref: kafkachannel.AsRef(outputChannel)},
		ksvc.AsTrackerReference(sender)))

	f.Setup("Create subscription", subscription.Install(inputChannel, subscription.WithChannel(kafkachannel.AsRef(inputChannel)), subscription.WithSubscriber(ksvc.AsRef(sender), "")))
	f.Setup("Create subscription", subscription.Install(outputChannel, subscription.WithChannel(kafkachannel.AsRef(outputChannel)), subscription.WithSubscriber(ksvc.AsRef(eventDisplay), "")))

	f.Setup("Create Sender Ksvc", ksvc.Install(sender, senderImage))

	f.Setup("SinkBinding becomes Ready", sinkbinding.IsReady(sinkBinding))

	f.Setup("Sender Ksvc becomes Ready", ksvc.IsReady(sender))

	f.Setup("Delete KafkaChannel-based SinkBinding", f.DeleteResources)

	// Take two, the same thing, this time with InMemoryChannels
	f.Setup("Create Event display ksvc", ksvc.Install(eventDisplay, eventDisplayImage))

	f.Setup("Create InMemoryChannel", inmemorychannel.Install(inputChannel))
	f.Setup("Create InMemoryChannel", inmemorychannel.Install(outputChannel))

	f.Setup("InMemoryChannel becomes Ready", inmemorychannel.IsReady(inputChannel))
	f.Setup("InMemoryChannel becomes Ready", inmemorychannel.IsReady(outputChannel))

	f.Setup("Create InMemoryChannel SinkBinding", sinkbinding.Install(sinkBinding,
		&duck.Destination{Ref: inmemorychannel.AsRef(outputChannel)},
		ksvc.AsTrackerReference(sender)))

	f.Setup("Create subscription", subscription.Install(inputChannel, subscription.WithChannel(inmemorychannel.AsRef(inputChannel)), subscription.WithSubscriber(ksvc.AsRef(sender), "")))
	f.Setup("Create subscription", subscription.Install(outputChannel, subscription.WithChannel(inmemorychannel.AsRef(outputChannel)), subscription.WithSubscriber(ksvc.AsRef(eventDisplay), "")))

	f.Setup("Create Sender Ksvc", ksvc.Install(sender, senderImage))

	f.Requirement("SinkBinding becomes Ready", sinkbinding.IsReady(sinkBinding))

	f.Requirement("Sender Ksvc becomes Ready", ksvc.IsReady(sender))

	return f
}

func SinkBindingDoesNotHangWithNoUri() *feature.Feature {
	f := feature.NewFeatureNamed("SinkBinding should not hang with 'No URI'")

	channel := feature.MakeRandomK8sName("channel")
	sender := feature.MakeRandomK8sName("sender")
	sinkBinding := feature.MakeRandomK8sName("sinkbinding")

	f.Setup("Create InMemoryChannel", inmemorychannel.Install(channel))

	f.Setup("Create SinkBinding", sinkbinding.Install(sinkBinding,
		&duck.Destination{Ref: inmemorychannel.AsRef(channel)},
		ksvc.AsTrackerReference(sender)))

	f.Setup("Create Sender Ksvc", ksvc.Install(sender, senderImage))

	f.Requirement("SinkBinding becomes Ready", sinkbinding.IsReady(sinkBinding))
	f.Requirement("Sender Ksvc becomes Ready", ksvc.IsReady(sender))

	return f
}

func TestSinkBinding(t *testing.T) {
	ctx, env := global.Environment(environment.Managed(t) /*, optional environment options */)

	// The first test passes, but breaks eventing-webhook, notice the errors like:
	// reflector.go:405] knative.dev/pkg/apis/duck/typed.go:67: watch of *v1.AddressableType ended with: an error on the server ("unable to decode an event from the watch stream: context canceled") has prevented the request from succeeding
	env.Test(ctx, t, ThingThatCausesFutureImcSinkBindingsToHangWithNoUri())

	// The following test normally passes, but not when run after ThingThatCausesFutureImcSinkBindingsToHangWithNoUri
	env.Test(ctx, t, SinkBindingDoesNotHangWithNoUri())
}
