/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package echo providers health check filter.
// RPCService need a Echo method in consumer, if you want to use Filter
// eg: Echo func(ctx context.Context, arg any, rsp *Xxx) error
package echo

import (
	"context"
	"sync"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol/base"
	"dubbo.apache.org/dubbo-go/v3/protocol/result"
)

var (
	once sync.Once
	echo *echoFilter
)

func init() {
	extension.SetFilter(constant.EchoFilterKey, newEchoFilter)
}

type echoFilter struct{}

func newEchoFilter() filter.Filter {
	if echo == nil {
		once.Do(func() {
			echo = &echoFilter{}
		})
	}
	return echo
}

// Invoke response to the callers with its first argument.
func (f *echoFilter) Invoke(ctx context.Context, invoker base.Invoker, invocation base.Invocation) result.Result {
	if invocation.MethodName() == constant.Echo && len(invocation.Arguments()) == 1 {
		return &result.RPCResult{
			Rest:  invocation.Arguments()[0],
			Attrs: invocation.Attachments(),
		}
	}

	return invoker.Invoke(ctx, invocation)
}

// OnResponse dummy process, returns the result directly
func (f *echoFilter) OnResponse(_ context.Context, result result.Result, _ base.Invoker,
	_ base.Invocation) result.Result {

	return result
}
