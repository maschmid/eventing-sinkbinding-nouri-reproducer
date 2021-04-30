# eventing-sinkbinding-nouri-reproducer
Attempt at reproducing a Knative Eventing SinkBinding "NoURI" issue

tldr; on a fresh install, run

```
go test
```

(Or, run `oc delete pod -n knative-eventing --all` before running the tests, to make sure the webhook didn't attempt to watch addressables before running the test)

When the problem occurs, the `TestSinkBinding/Thing_that_causes_future_SinkBinding_with_IMC_destination_to_hang_with_'No_URI'` test will PASS, but will cause the `v1.AddressableType` watch to break, as can be seen in eventing-webhook logs as

```
W0430 08:04:36.732657       1 reflector.go:405] knative.dev/pkg/apis/duck/typed.go:67: watch of *v1.AddressableType ended with: an error on the server ("unable to decode an event from the watch stream: context canceled") has prevented the request from succeeding
E0430 08:04:37.767383       1 reflector.go:178] knative.dev/pkg/apis/duck/typed.go:67: Failed to list *v1.AddressableType: context canceled
E0430 08:04:40.454141       1 reflector.go:178] knative.dev/pkg/apis/duck/typed.go:67: Failed to list *v1.AddressableType: context canceled
...
```

The subsequent `TestSinkBinding/SinkBinding_should_not_hang_with_'No_URI'/Requirement/SinkBinding_becomes_Ready` test will fail, and the SinkBinding will never become Ready, hung with `NoURI: URI could not be extracted from destination`

Note, this seems to be only reproducible on a fresh eventing-webhook. Restart (pod delete) the eventing-webhook before running the test.
