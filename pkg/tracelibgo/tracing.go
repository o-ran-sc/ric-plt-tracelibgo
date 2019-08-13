/*
 * Copyright (c) 2019 AT&T Intellectual Property.
 * Copyright (c) 2018-2019 Nokia.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package tracelibgo implements a function to create a configured tracer instance
package tracelibgo

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func createTracer(name string) (opentracing.Tracer, io.Closer) {
	tracer, closer := jaeger.NewTracer(name,
		jaeger.NewConstSampler(false),
		jaeger.NewNullReporter())
	return tracer, closer
}

// CreateTracer creates a tracer entry
func CreateTracer(name string) (opentracing.Tracer, io.Closer) {
	return createTracer(name)
}
