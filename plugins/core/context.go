// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package core

var (
	GetGLS = func() interface{} { return nil }
	SetGLS = func(interface{}) {}
)

type TracingContext struct {
	ActiveSpan Span
	Runtime    *RuntimeContext
}

type RuntimeContext struct {
	data map[string]interface{}
}

func NewTracingContext() *TracingContext {
	return &TracingContext{
		Runtime: &RuntimeContext{
			data: make(map[string]interface{}),
		},
	}
}

func (r *RuntimeContext) clone() *RuntimeContext {
	newData := make(map[string]interface{})
	for k, v := range r.data {
		newData[k] = v
	}
	return &RuntimeContext{
		data: newData,
	}
}

func GetTracingContext() *TracingContext {
	gls := GetGLS()
	if gls == nil {
		return nil
	}
	return gls.(*TracingContext)
}

func (r *RuntimeContext) Get(key string) interface{} {
	return r.data[key]
}

func (r *RuntimeContext) Set(key string, value interface{}) {
	r.data[key] = value
}

func TaskTracingContextSnapshot(val interface{}) interface{} {
	if val == nil {
		return nil
	}
	context, ok := val.(*TracingContext)
	if !ok {
		return nil
	}
	snapshot := newSnapshotSpan(context.ActiveSpan)
	return &TracingContext{
		ActiveSpan: snapshot,
		Runtime:    context.Runtime.clone(),
	}
}