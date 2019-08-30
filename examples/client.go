package main

import (
	"tracelibgo/pkg/tracelibgo"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"fmt"
	"os"
	"strconv"
	"unsafe"
)

func spanner(rmrcontext unsafe.Pointer, parentcontext []byte, counter int) {
	span, spancontext := startSpan("child :" + strconv.Itoa(counter), parentcontext)
	span.LogFields(log.Int("Counter value", counter))
	sendMsg(rmrcontext, 1, spancontext)
	buf := recvMsg(rmrcontext)
	span.Finish()
	spancontext = getTraceContext(buf)
	freeBuffer(buf)
	if (counter > 0) {
		spanner(rmrcontext, spancontext, counter-1)
	}
}


func main() {
	fmt.Println("Creating tracer")
	tracer, closer := tracelibgo.CreateTracer("Client process")
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	rmrcontext := initRmr("tcp:5000")
	var count int
	if len(os.Args) > 1 {
		count, _ = strconv.Atoi(os.Args[1])
	}
	if count == 0 {
		count = 1
	}
	parent, parentcontext := startSpan("Span start", []byte{})
	defer parent.Finish()
	spanner(rmrcontext, parentcontext, count)
}
