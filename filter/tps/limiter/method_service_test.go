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

package limiter

import (
	"net/url"
	"testing"
)

import (
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/filter/tps/strategy"
	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"
)

func TestMethodServiceTpsLimiterImplIsAllowableOnlyServiceLevel(t *testing.T) {
	methodName := "hello"
	invoc := invocation.NewRPCInvocation(methodName, []any{"OK"}, make(map[string]any))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invokeUrl := common.NewURLWithOptions(
		common.WithParams(url.Values{}),
		common.WithParamsValue(constant.InterfaceKey, methodName),
		common.WithParamsValue(constant.TPSLimitRateKey, "20"),
		common.WithParamsValue(constant.TPSLimitIntervalKey, "60000"))

	mockStrategyImpl := strategy.NewMockTpsLimitStrategy(ctrl)
	mockStrategyImpl.EXPECT().IsAllowable().Return(true).Times(1)

	extension.SetTpsLimitStrategy(constant.DefaultKey, &mockStrategyCreator{
		rate:     20,
		interval: 60000,
		t:        t,
		strategy: mockStrategyImpl,
	})

	limiter := GetMethodServiceTpsLimiter()
	result := limiter.IsAllowable(invokeUrl, invoc)
	assert.True(t, result)
}

func TestMethodServiceTpsLimiterImplIsAllowableNoConfig(t *testing.T) {
	methodName := "hello1"
	invoc := invocation.NewRPCInvocation(methodName, []any{"OK"}, make(map[string]any))
	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()

	invokeUrl := common.NewURLWithOptions(
		common.WithParams(url.Values{}),
		common.WithParamsValue(constant.InterfaceKey, methodName),
		common.WithParamsValue(constant.TPSLimitRateKey, ""))

	limiter := GetMethodServiceTpsLimiter()
	result := limiter.IsAllowable(invokeUrl, invoc)
	assert.True(t, result)
}

func TestMethodServiceTpsLimiterImplIsAllowableMethodLevelOverride(t *testing.T) {
	methodName := "hello2"
	methodConfigPrefix := "methods." + methodName + "."
	invoc := invocation.NewRPCInvocation(methodName, []any{"OK"}, make(map[string]any))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invokeUrl := common.NewURLWithOptions(
		common.WithParams(url.Values{}),
		common.WithParamsValue(constant.InterfaceKey, methodName),
		common.WithParamsValue(constant.TPSLimitRateKey, "20"),
		common.WithParamsValue(constant.TPSLimitIntervalKey, "3000"),
		common.WithParamsValue(constant.TPSLimitStrategyKey, "invalid"),
		common.WithParamsValue(methodConfigPrefix+constant.TPSLimitRateKey, "40"),
		common.WithParamsValue(methodConfigPrefix+constant.TPSLimitIntervalKey, "7000"),
		common.WithParamsValue(methodConfigPrefix+constant.TPSLimitStrategyKey, "default"),
	)

	mockStrategyImpl := strategy.NewMockTpsLimitStrategy(ctrl)
	mockStrategyImpl.EXPECT().IsAllowable().Return(true).Times(1)

	extension.SetTpsLimitStrategy(constant.DefaultKey, &mockStrategyCreator{
		rate:     40,
		interval: 7000,
		t:        t,
		strategy: mockStrategyImpl,
	})

	limiter := GetMethodServiceTpsLimiter()
	result := limiter.IsAllowable(invokeUrl, invoc)
	assert.True(t, result)
}

func TestMethodServiceTpsLimiterImplIsAllowableBothMethodAndService(t *testing.T) {
	methodName := "hello3"
	methodConfigPrefix := "methods." + methodName + "."
	invoc := invocation.NewRPCInvocation(methodName, []any{"OK"}, make(map[string]any))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invokeUrl := common.NewURLWithOptions(
		common.WithParams(url.Values{}),
		common.WithParamsValue(constant.InterfaceKey, methodName),
		common.WithParamsValue(constant.TPSLimitRateKey, "20"),
		common.WithParamsValue(constant.TPSLimitIntervalKey, "3000"),
		common.WithParamsValue(methodConfigPrefix+constant.TPSLimitRateKey, "40"),
	)

	mockStrategyImpl := strategy.NewMockTpsLimitStrategy(ctrl)
	mockStrategyImpl.EXPECT().IsAllowable().Return(true).Times(1)

	extension.SetTpsLimitStrategy(constant.DefaultKey, &mockStrategyCreator{
		rate:     40,
		interval: 3000,
		t:        t,
		strategy: mockStrategyImpl,
	})

	limiter := GetMethodServiceTpsLimiter()
	result := limiter.IsAllowable(invokeUrl, invoc)
	assert.True(t, result)
}

type mockStrategyCreator struct {
	rate     int
	interval int
	t        *testing.T
	strategy filter.TpsLimitStrategy
}

func (creator *mockStrategyCreator) Create(rate int, interval int) filter.TpsLimitStrategy {
	assert.Equal(creator.t, creator.rate, rate)
	assert.Equal(creator.t, creator.interval, interval)
	return creator.strategy
}
