/*
 * Copyright 2022 CloudWeGo Authors
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

package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	hertzlogrus "github.com/hertz-contrib/obs-opentelemetry/logging/logrus"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"sync"

	"github.com/spf13/viper"
)

var wg sync.WaitGroup
var cli, _ = client.NewClient() // fixme: client 全局唯一

func ping1(c context.Context, ctx *app.RequestContext) {
	// 自定义span
	ctx2, span := otel.Tracer("github.com/hertz-contrib/obs-opentelemetry").
		Start(c, "multiple-ping1") // 进程内使用 go 的 context.Context 传递span

	_, b, err := cli.Get(ctx2, nil, "http://0.0.0.0:8081/ping?foo=bar")
	if err != nil {
		hlog.CtxErrorf(ctx2, err.Error())
	}

	span.SetAttributes(attribute.String("msg", string(b)))

	hlog.CtxInfof(ctx2, "multiple-ping1 %s", string(b))
	span.End()

	ctx.JSON(consts.StatusOK, utils.H{"ping1": "pong1"}) //http间使用 http header 传递 span
}

func hertz1() {
	tracer, cfg := hertztracing.NewServerTracer() // server span 必须用到的代码

	hlog.Debugf("tracer:%s, cfg:%s ", tracer, cfg)

	h := server.Default(tracer, server.WithHostPorts(":8080")) // server span 必须用到的代码

	h.Use(hertztracing.ServerMiddleware(cfg)) // server span 必须用到的代码
	h.GET("/ping", ping1)
	h.Spin()
}

func ping2(c context.Context, ctx *app.RequestContext) {

	ctx2, span := otel.Tracer("github.com/hertz-contrib/obs-opentelemetry").
		Start(c, "multiple-ping2")

	_, b, err := cli.Get(ctx2, nil, "http://0.0.0.0:8888/ping?foo=bar")
	if err != nil {
		hlog.CtxErrorf(ctx2, err.Error())
	}

	span.SetAttributes(attribute.String("msg", string(b)))

	hlog.CtxInfof(ctx2, "multiple-ping2 %s", string(b))
	span.End()

	ctx.JSON(consts.StatusOK, utils.H{"ping2": "pong2"})
}

func hertz2() {
	tracer, cfg := hertztracing.NewServerTracer()
	h := server.Default(tracer, server.WithHostPorts(":8081")) //

	h.Use(hertztracing.ServerMiddleware(cfg))
	h.GET("/ping", ping2)
	h.Spin()
}

func main() {
	cli.Use(hertztracing.ClientMiddleware()) //client 端要使用中间件

	hlog.SetLogger(hertzlogrus.NewLogger())
	hlog.SetLevel(hlog.LevelDebug)

	err := viper.BindEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if err != nil {
		hlog.Errorf("BindEnv:", err)
		return
	}
	endpoint := viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT")
	hlog.Debugf("OTEL_EXPORTER_OTLP_ENDPOIN", endpoint)

	serviceName := "demo-multiple-server"

	// provider 进程中唯一,每个进程中都有
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(serviceName),
		// Support setting ExportEndpoint via environment variables: OTEL_EXPORTER_OTLP_ENDPOINT
		provider.WithExportEndpoint(endpoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	wg.Add(2)
	go func() {
		defer wg.Done()
		hertz1()
	}()
	go func() {
		defer wg.Done()
		hertz2()
	}()
	wg.Wait()
}
