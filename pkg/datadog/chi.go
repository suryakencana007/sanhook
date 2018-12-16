/*  mux.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 14, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 14/10/18 23:51 
 */

package datadog

import (
    "fmt"
    "net/http"
    "strconv"

    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func DataDog(service string) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        fn := func(w http.ResponseWriter, r *http.Request) {
            resource := r.Method + " Datadog"
            opts := []ddtrace.StartSpanOption{
                tracer.ServiceName(service),
                tracer.ResourceName(resource),
                tracer.SpanType(ext.SpanTypeWeb),
                tracer.Tag(ext.HTTPMethod, r.Method),
                tracer.Tag(ext.HTTPURL, r.URL.Path),
            }
            if spanCtx, err := tracer.Extract(tracer.HTTPHeadersCarrier(r.Header)); err == nil {
                opts = append(opts, tracer.ChildOf(spanCtx))
            }

            span, ctx := tracer.StartSpanFromContext(r.Context(), "http.request", opts...)
            defer span.Finish()

            // pass the span through the request context
            w = newResponseWriter(w, span)

            // serve the request to the next middleware
            next.ServeHTTP(w, r.WithContext(ctx))

        }
        return http.HandlerFunc(fn)
    }
}

// responseWriter is a small wrapper around an http response writer that will
// intercept and store the status of a request.
type responseWriter struct {
    http.ResponseWriter
    span   ddtrace.Span
    status int
}

func newResponseWriter(w http.ResponseWriter, span ddtrace.Span) *responseWriter {
    return &responseWriter{w, span, 0}
}

// Write writes the data to the connection as part of an HTTP reply.
// We explicitely call WriteHeader with the 200 status code
// in order to get it reported into the span.
func (w *responseWriter) Write(b []byte) (int, error) {
    if w.status == 0 {
        w.WriteHeader(http.StatusOK)
    }
    return w.ResponseWriter.Write(b)
}

// WriteHeader sends an HTTP response header with status code.
// It also sets the status code to the span.
func (w *responseWriter) WriteHeader(status int) {
    w.ResponseWriter.WriteHeader(status)
    w.status = status
    w.span.SetTag(ext.HTTPCode, strconv.Itoa(status))
    if status >= 500 && status < 600 {
        w.span.SetTag(ext.Error, fmt.Errorf("%d: %s", status, http.StatusText(status)))
    }
}
