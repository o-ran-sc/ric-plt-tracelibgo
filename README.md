# Tracing helper library

The library creates a configured tracer instance.

ToDO: configuration...

## Usage

Create a tracer instance and set it as a global tracer:

```go
import (
		"github.com/opentracing/opentracing-go"
        "gerrit.o-ran-sc.org/ric-plt/tracelibgo/pkg/tracelibgo"
        ...
)

tracer, closer := tracelibgo.CreateTracer("my-service-name")
defer closer.Close()
opentracing.SetGlobalTracer(tracer)
```

Serialize span context to a byte array that can be sent
to another component via some messaging. For example, using
the RMR library.

```go
	carrier := make(map[string]string)
	opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.TextMap,
			opentracing.TextMapCarrier(carrier))
	b, err := json.Marshal(carrier) // b is a []byte and contains serilized span context
```

Extract a span context from byte array and create a new child span from it.
The serialized span context is got, for example, from the RMR library.

```go
	var carrier map[string]string
	err = json.Unmarshal(data, &carrier) // data is []byte containing serialized span context
	if err != nil {
		...
	}
	context, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(carrier))
	if err != nil {
		...
	}
	span := opentracing.GlobalTracer().StartSpan("go test span", opentracing.ChildOf(context))
```

## Unit testing

 GO111MODULE=on go mod download
 go test ./pkg/tracelibgo

