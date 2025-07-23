package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/hanshuaikang/gin-prometheus"
	"go.opentelemetry.io/otel/attribute"
)

const serviceName = "gin-prometheus-demo"

func main() {
	fmt.Println("initializing")

	r := gin.New()
	globalAttributes := []attribute.KeyValue{
		attribute.String("k8s.pod.name", "pod-1"),
		attribute.String("k8s.namespace.name", "test"),
		attribute.String("service.name", serviceName),
	}
	r.Use(ginprometheus.Middleware(
		// Custom attributes
		ginprometheus.WithAttributes(func(route string, request *http.Request) []attribute.KeyValue {
			attrs := []attribute.KeyValue{
				attribute.String("http.method", request.Method),
			}
			if route != "" {
				attrs = append(attrs, attribute.String("http.route", route))
			}
			return attrs
		}),
		ginprometheus.WithGlobalAttributes(globalAttributes),
		ginprometheus.WithService(serviceName, "v0.0.1"),
		ginprometheus.WithMetricPrefix("infra"),
		ginprometheus.WithPrometheusPort(4433),
		ginprometheus.WithSystemMetricDisabled(),
	))
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
