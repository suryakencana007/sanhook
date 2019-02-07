/*  http_test.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 15, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 15/10/18 03:38 
 */

package cmd

import (
    "bytes"
    "io"
    "os"
    "sync"
    "testing"

    "github.com/go-chi/chi"
    "github.com/spf13/cobra"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/suryakencana007/sanhook/configs"
    internalHttp "github.com/suryakencana007/sanhook/internal/http"
)

var table = `+---------+-----------------------------------------------+--------+
| HANDLER |                      URL                      | METHOD |
+---------+-----------------------------------------------+--------+
|         | /v1/*/api/message/*/inbox                     | GET    |
|         | /v1/*/api/message/*/inbox/{id:[a-z-0-9]+}     | GET    |
|         | /v1/*/api/message/*/publish/{subject:[a-z-]+} | GET    |
|         | /v1/*/api/status/*/health                     | GET    |
+---------+-----------------------------------------------+--------+
`

var (
    wg sync.WaitGroup
)

func TestTableRoute(t *testing.T) {
    configuration := configs.New("app.config",
        "./configs", "../configs", "../../../configs")
    router := internalHttp.Routes(configuration)

    oldStdout := os.Stdout
    r, w, _ := os.Pipe()

    os.Stdout = w
    outC := make(chan string)
    tableRoute(router)
    go func() {
        var buf bytes.Buffer
        if _, err := io.Copy(&buf, r); err == nil {
            outC <- buf.String()
        }
    }()
    // Reset the output again
    w.Close()
    os.Stdout = oldStdout
    out := <-outC
    assert.Equal(t, table, out)
}

func TestNewHttp(t *testing.T) {
    assert := require.New(t)
    stop := make(chan bool)
    wg.Add(1)
    configuration := configs.New("app.config",
        "./configs", "../configs", "../../../configs")

    http := NewHttpCmd(configuration)
    http.stop = stop

    go func() {
        defer wg.Done()
        err := http.BaseCmd.Execute()
        assert.NoError(err)
    }()

    stop <- true
    wg.Wait()
}

func TestHttp(t *testing.T) {
    assert := require.New(t)
    stop := make(chan bool)
    wg.Add(1)
    configuration := configs.New("app.config",
        "./configs", "../configs", "../../../configs")

    cmd := NewHttpCmdSignaled(configuration, stop).BaseCmd
    go func() {
        defer wg.Done()
        _, err := cmd.ExecuteC()
        assert.NoError(err)
    }()

    stop <- true
    wg.Wait()
}

func TestHttpFail(t *testing.T) {
    var (
        err           error
        stop          = make(chan bool)
        configuration *configs.Config
    )

    wg.Add(1)
    f := func() {
        configuration = configs.New("app.config",
            "./configs", "../configs", "../../configs")
    }
    assert.Panics(t, f)
    cmd := NewHttpCmdSignaled(configuration, stop).BaseCmd

    go func() {
        defer wg.Done()
        fn := func() {
            _, err = cmd.ExecuteC()
        }
        assert.Panics(t, fn)
        <-stop
    }()

    assert.NoError(t, err)
    stop <- true
    wg.Wait()
}

func TestNewHttpCmdWithFilename(t *testing.T) {
    stop := make(chan bool)
    configuration := configs.New("app.config", ConfigPath...)

    wg.Add(1)
    os.Args = []string{"main", "http", "-f", "app.config"}

    cmd := NewHttpCmdSignaled(configuration, stop).BaseCmd
    go func() {
        defer wg.Done()

        _, err := cmd.ExecuteC()
        assert.NoError(t, err)
    }()

    stop <- true
    wg.Wait()
    os.Args = []string{""}
}

func TestListenAndServe(t *testing.T) {
    var err error
    stop := make(chan bool)

    configuration := configs.New("app.config", ConfigPath...)

    cc := &httpCmd{stop: stop}
    cc.configuration = configuration
    cc.BaseCmd = &cobra.Command{
        Use:   "http",
        Short: "Used to run the http service",
        RunE: func(cmd *cobra.Command, args []string) (err error) {
            mux := chi.NewMux()
            return cc.serve(mux)
        },
    }

    wg.Add(1)
    go func() {
        defer wg.Done()
        err = cc.BaseCmd.Execute()
    }()
    assert.NoError(t, err)
    stop <- true
    wg.Wait()
}
//
// func TestListenAndServeInUse(t *testing.T) {
//     var err error
//     sign := make(chan os.Signal, 1)
//     stop := make(chan bool, 1)
//
//     signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
//
//     configuration := configs.New("app.config", ConfigPath...)
//
//     cc := &httpCmd{stop: stop}
//     cc.configuration = configuration
//     cc.BaseCmd = &cobra.Command{
//         Use:   "http",
//         Short: "Used to run the http service",
//         RunE: func(cmd *cobra.Command, args []string) (err error) {
//             mux := chi.NewMux()
//             return cc.serve(mux)
//         },
//     }
//
//     wg.Add(1)
//     go func() {
//         defer wg.Done()
//         err = cc.BaseCmd.Execute()
//         assert.NoError(t, err)
//     }()
//
//     go func() {
//         sig := <-sign
//         fmt.Println()
//         fmt.Println(sig)
//         err = cc.BaseCmd.Execute()
//         t.Log("defer BaseCmd 1", err)
//     }()
//
//     // The program will wait here until it gets the
//     // expected signal (as indicated by the goroutine
//     // above sending a value on `done`) and then exit.
//     fmt.Println("awaiting signal")
//     if err != nil {
//     //     fmt.Println(err.Error())
//     //     signal.Stop(sign)
//         stop <- true
//     }
//     fmt.Println("exiting")
//     <-stop
//     wg.Wait()
//
// }
