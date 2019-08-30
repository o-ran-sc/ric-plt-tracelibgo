package main

import (
	"tracelibgo/pkg/tracelibgo"
	"github.com/opentracing/opentracing-go"
	"fmt"
	//"encoding/json"
)

func main() {
	fmt.Println("Creating tracer")
	tracer, closer := tracelibgo.CreateTracer("Server process")
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	rmrcontext := initRmr("tcp:4000")
	for {
		fmt.Println("Serving")
		rxBuffer := recvMsg(rmrcontext)
		tc := getTraceContext(rxBuffer)
		span, newTraceContext := startSpan("server span", tc)
		sendMsg(rmrcontext, 2, newTraceContext)
		freeBuffer(rxBuffer)
		span.Finish()
	}
}

