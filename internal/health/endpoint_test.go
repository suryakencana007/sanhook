/*  endpoint_test.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 08, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 08/10/18 16:45
 */

package health

import (
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/suryakencana007/sanhook/internal/container"
    pkgHttp "github.com/suryakencana007/sanhook/pkg/http"
)

// Unit Test Health
func TestHealth(t *testing.T) {
    res, _ := testHandler(t,
        getAPIStatus(container.New().NewItemServiceJSON()),
        http.MethodGet, "/health", nil)
    if got, want := res.StatusCode, http.StatusOK; got != want {
        t.Fatalf("status code got: %d, want %d", got, want)
    }
}

func testHandler(
    t *testing.T,
    h http.Handler,
    method, path string,
    body io.Reader,
) (*http.Response, string) {
    bodyReader, err := pkgHttp.BodyToJson(body)
    if err != nil {
        panic(err)
    }
    r, _ := http.NewRequest(method, path, bodyReader)
    w := httptest.NewRecorder()
    h.ServeHTTP(w, r)
    return w.Result(), w.Body.String()
}
