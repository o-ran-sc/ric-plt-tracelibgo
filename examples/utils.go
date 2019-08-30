package main

/*
#include <time.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <rmr/rmr.h>

void write_bytes_array(unsigned char *dst, void *data, int len) {
    memcpy((void *)dst, (void *)data, len);
}

#cgo CFLAGS: -I../
#cgo LDFLAGS: -lrmr_nng -lnng
*/
import "C"

import (
	"time"
	"fmt"
	"unsafe"
	"github.com/opentracing/opentracing-go"
	"encoding/json"
)

func initRmr(protoport string) unsafe.Pointer {
	protport := C.CString(protoport)
	context := C.rmr_init(protport, C.int(1000), C.int(0))
	if context == nil {
		panic("Cannot init rmr")
	}
	for {
		if ready := int(C.rmr_ready(context)); ready == 1 {
			break
		}
		time.Sleep(10*time.Second)
		fmt.Println("Wait for rmr")
	}
	return context
}

func sendMsg(rmrcontext unsafe.Pointer, mtype int, traceContext []byte) {
	var buf *C.rmr_mbuf_t

	if len(traceContext) > 0 {
		buf = C.rmr_tralloc_msg(rmrcontext, C.int(10), C.int(len(traceContext)), (*C.uchar)(unsafe.Pointer(&traceContext[0])))
	} else {
		buf = C.rmr_alloc_msg(rmrcontext, C.int(10))
	}
	buf.mtype = C.int(mtype)
	buf.sub_id = C.int(0)
	buf.len = C.int(10)
	txBuffer := C.rmr_send_msg(rmrcontext, buf)
	if txBuffer.state != C.RMR_OK {
		fmt.Println("RMR send error")
	}
}

func recvMsg(rmrcontext unsafe.Pointer) *C.rmr_mbuf_t {
	return  C.rmr_rcv_msg(rmrcontext, nil)
}

func getTraceContext(rmrbuf *C.rmr_mbuf_t) []byte{
	var traceLen C.int
	tracecont := C.rmr_trace_ref(rmrbuf, &traceLen)
	if tracecont != nil && int(traceLen) > 0 {
		return C.GoBytes(tracecont, traceLen)
	}
	return []byte{}
}

func newSpan(spanName string, traceContext []byte) opentracing.Span {
	if len(traceContext) > 0 {
		var carrier map[string]string
		err := json.Unmarshal(traceContext, &carrier)
		if err != nil {
			fmt.Println("Json unmarshal error: ", err.Error())
		} else {
			cont, err := opentracing.GlobalTracer().Extract(opentracing.TextMap,
					opentracing.TextMapCarrier(carrier))
			if err == nil {
				return opentracing.StartSpan(spanName, opentracing.ChildOf(cont))
			} else {
				fmt.Println("span extract error: ", err.Error())
			}
		}
	}
	return opentracing.StartSpan(spanName)
}

func startSpan(spanName string, traceContext []byte) (opentracing.Span, []byte) {
	span := newSpan(spanName, traceContext)
	carrier := make(map[string]string)
	opentracing.GlobalTracer().Inject(span.Context(),
				opentracing.TextMap,
			opentracing.TextMapCarrier(carrier))
	fmt.Println(carrier)
	newContext, _ := json.Marshal(carrier)
	return span, newContext
}

func freeBuffer(buf *C.rmr_mbuf_t) {
	C.rmr_free_msg(buf)
}