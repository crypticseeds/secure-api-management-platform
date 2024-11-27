package middleware

import (
	"apisecurityplatform/pkg/observability"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := observability.GetTracer()
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())

		ctx, span := tracer.Start(c.Request.Context(), spanName)
		defer span.End()

		// Add common attributes
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.user_agent", c.Request.UserAgent()),
			attribute.String("http.remote_addr", c.ClientIP()),
		)

		// Log trace ID for debugging
		traceID := span.SpanContext().TraceID().String()
		fmt.Printf("Trace ID: %s\n", traceID)

		// Store span in context
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Add response attributes
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
			attribute.Int("http.response_size", c.Writer.Size()),
		)

		// Record error if any
		if len(c.Errors) > 0 {
			span.RecordError(c.Errors.Last().Err)
			span.SetStatus(codes.Error, c.Errors.Last().Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}
	}
}
