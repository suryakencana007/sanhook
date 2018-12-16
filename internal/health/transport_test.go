/*  transport_test.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 09, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 09/10/18 13:55 
 */

package health

import (
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-chi/chi"
    "github.com/stretchr/testify/assert"
    "github.com/suryakencana007/sanhook/configs"
)

type attributes struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func TestTransport(t *testing.T) {
    // create instance mux route
    c := configs.New("app.config",
        "./configs", "../configs", "../../configs")
    r := chi.NewRouter()
    // API ver. 1
    r.Route("/v1", func(r chi.Router) {
        r.Mount(fmt.Sprintf(`%s/%s`, c.Api.Prefix, "status"), MakeHandler(c))
    })

    ts := httptest.NewServer(r)
    defer ts.Close()

    // check that we didn't break correct routes
    // Route: /v1/api/status/health
    res, b := testRequest(t, ts, "GET", "/v1/api/status/health", nil)
    if got, want := res.StatusCode, http.StatusOK; got != want {
        t.Fatalf("status code got: %d, want %d", got, want)
    }
    resp := &attributes{}
    err := json.Unmarshal([]byte(b), resp)
    if err != nil {
        t.Fatal(err)
    }
    assert.Equal(t, &attributes{
        Code:    1000,
        Message: "Success",
    }, resp)
    // :GET /v1/api/status/health
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
