package main

import (
	"encoding/json"
	"log"
	"order/domain"
	"order/infra/couchbase"
	"order/infra/kafka"
	"order/services"

	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tracerEndpoint := "localhost:4318"
	tp := initTracer(tracerEndpoint)
	defer func() { _ = tp.Shutdown(context.Background()) }()
	tracer := otel.Tracer("order-tracer")

	kafkaRepo := kafka.InitKafka()
	couchBaseRepo := couchbase.InitCouchBase()
	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Concurrency:  256 * 1024,
	})

	app.Post("/order", func(c fiber.Ctx) error {

		orderItem := new(domain.OrderItem)
		s := services.NewOrderService(kafkaRepo, couchBaseRepo)

		if err := json.Unmarshal(c.Body(), orderItem); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid JSON"})
		}
		s.SendKafkaToStoreCouchbase(orderItem, tracer)

		return c.JSON(fiber.Map{"message": "order created"})
	})
	defer kafka.CloseKafkaConnection(kafkaRepo)
	log.Fatal(app.Listen(":3000"))
}

func initTracer(endpoint string) *sdktrace.TracerProvider {
	headers := map[string]string{
		"content-type": "application/json",
	}

	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithHeaders(headers),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("order-service"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	return tp
}

func IdGenerate() string {
	return uuid.New().String()
}
