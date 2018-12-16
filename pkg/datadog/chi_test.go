// +build integration

/*  chi_test.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 15, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 15/10/18 01:59 
 */

package datadog

import (
    "errors"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-chi/chi"
    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

    "github.com/stretchr/testify/assert"
)

func TestChildSpan(t *testing.T) {
    assert := assert.New(t)
    mt := mocktracer.Start()
    defer mt.Stop()

    router := chi.NewRouter()
    router.Use(DataDog("foobar"))
    router.Get("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
        _, ok := tracer.SpanFromContext(r.Context())
        assert.True(ok)
    })

    r := httptest.NewRequest("GET", "/user/123", nil)
    w := httptest.NewRecorder()

    router.ServeHTTP(w, r)
}

func TestTrace200(t *testing.T) {
    assert := assert.New(t)
    mt := mocktracer.Start()
    defer mt.Stop()

    router := chi.NewRouter()
    router.Use(DataDog("foobar"))

    router.Get("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
        _, ok := tracer.SpanFromContext(r.Context())
        assert.True(ok)

        id := chi.URLParam(r, "id")
        w.Write([]byte(id))
    })

    r := httptest.NewRequest("GET", "/user/123", nil)
    w := httptest.NewRecorder()

    // do and verify the request
    router.ServeHTTP(w, r)
    response := w.Result()
    assert.Equal(response.StatusCode, 200)

    // verify traces look good
    spans := mt.FinishedSpans()
    assert.Len(spans, 1)
    if len(spans) < 1 {
        t.Fatalf("no spans")
    }
    span := spans[0]
    assert.Equal("http.request", span.OperationName())
    assert.Equal(ext.SpanTypeWeb, span.Tag(ext.SpanType))
    assert.Equal("foobar", span.Tag(ext.ServiceName))
    assert.Contains(span.Tag(ext.ResourceName), "GET Datadog")
    assert.Equal("200", span.Tag(ext.HTTPCode))
    assert.Equal("GET", span.Tag(ext.HTTPMethod))
    // TODO(x) would be much nicer to have "/user/:id" here
    assert.Equal("/user/123", span.Tag(ext.HTTPURL))
}

func TestError(t *testing.T) {
    assert := assert.New(t)
    mt := mocktracer.Start()
    defer mt.Stop()

    // setup
    router := chi.NewRouter()
    router.Use(DataDog("foobar"))
    wantErr := errors.New("oh no")

    // a handler with an error and make the requests
    router.Get("/err", func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, wantErr.Error(), 500)
    })
    r := httptest.NewRequest("GET", "/err", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, r)
    response := w.Result()
    assert.Equal(response.StatusCode, 500)

    // verify the errors and status are correct
    spans := mt.FinishedSpans()
    assert.Len(spans, 1)
    if len(spans) < 1 {
        t.Fatalf("no spans")
    }
    span := spans[0]
    assert.Equal("http.request", span.OperationName())
    assert.Equal("foobar", span.Tag(ext.ServiceName))
    assert.Equal("500", span.Tag(ext.HTTPCode))
    assert.NotEqual(wantErr.Error(), span.Tag(ext.Error).(error).Error())
}
