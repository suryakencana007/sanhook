/*  routes_test.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 10, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 10/10/18 14:43 
 */

package http

import (
    "io"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/suryakencana007/sanhook/configs"
)

func TestRoutes(t *testing.T) {
    // create instance mux route
    c := configs.New("app.config",
        "./configs", "../configs", "../../configs")
    r := Routes(c)
    ts := httptest.NewServer(r)
    defer ts.Close()

    /**
     * Routes Mount
     * /v1/api/status
     */

    // :GET /health
    res, _ := testRequest(t, ts, "GET", "/v1/api/status/health", nil)
    if got, want := res.StatusCode, http.StatusOK; got != want {
        t.Fatalf("status code got: %d, want %d", got, want)
    }
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
    req, err := http.NewRequest(method, ts.URL+path, body)
    if err != nil {
        t.Fatal(err)
        return nil, ""
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatal(err)
        return nil, ""
    }

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        t.Fatal(err)
        return nil, ""
    }
    defer resp.Body.Close()

    return resp, string(respBody)
}
