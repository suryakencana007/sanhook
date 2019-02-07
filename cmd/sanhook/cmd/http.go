/*  http.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               September 30, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 30/09/18 01:37 
 */

package cmd

import (
    "fmt"
    "net"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "syscall"
    "time"

    "github.com/go-chi/chi"
    "github.com/olekukonko/tablewriter"
    "github.com/spf13/pflag"
    "github.com/suryakencana007/sanhook/configs"
    "github.com/suryakencana007/sanhook/internal/container"
    internalHttp "github.com/suryakencana007/sanhook/internal/http"

    "github.com/spf13/cobra"
    "github.com/suryakencana007/sanhook/pkg/log"
)

type httpCmd struct {
    stop <-chan bool

    BaseCmd       *cobra.Command
    configuration *configs.Config
    container     *container.Container
    filename      string
}

func NewHttpCmd(
    configuration *configs.Config,
) *httpCmd {
    return NewHttpCmdSignaled(configuration, nil)
}

func NewHttpCmdSignaled(
    configuration *configs.Config,
    stop <-chan bool,
) *httpCmd {
    cc := &httpCmd{stop: stop}
    cc.configuration = configuration
    cc.BaseCmd = &cobra.Command{
        Use:   "http",
        Short: "Used to run the http service",
        RunE:  cc.server,
    }
    fs := pflag.NewFlagSet("Root", pflag.ContinueOnError)
    fs.StringVarP(&cc.filename, "file", "f", "", "Custom configuration filename")
    cc.BaseCmd.Flags().AddFlagSet(fs)
    return cc
}

func (h *httpCmd) server(cmd *cobra.Command, args []string) (err error) {
    if len(h.filename) > 1 {
        h.configuration = configs.New(h.filename,
            "./configs",
            "../configs",
            "../../configs",
            "../../../configs")
    }

    router := internalHttp.Routes(
        h.configuration,
    )
    // Description µ micro service
    fmt.Println(
        fmt.Sprintf(
            WelkomText,
            h.configuration.App.Port,
            strings.Join([]string{
                h.configuration.Log.Dir,
                h.configuration.Log.Filename}, "/"),
        ))
    tableRoute(router) // Prettier Route Pattern

    return h.serve(router)
}

func (h *httpCmd) serve(router *chi.Mux) error {
    errCh := make(chan error, 1)
    quit := make(chan os.Signal)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    addr := net.JoinHostPort("",
        strconv.Itoa(h.configuration.App.Port))
    s := StartWebServer(
        addr,
        h.configuration.App.ReadTimeout,
        h.configuration.App.WriteTimeout,
        router,
    )

    go func() {
        if err := s.ListenAndServe(); err != nil {
            log.Info(
                "Server gracefully ListenAndServe",
            )
            errCh <- err
        }
        <-h.stop
    }()

    if h.stop != nil {
        select {
        case err := <-errCh:
            log.Info(
                "Server gracefully h stop stopped",
            )
            return err
        case <-h.stop:
        case <-quit:
        }
    } else {
        select {
        case err := <-errCh:
            log.Info(
                "Server gracefully stopped",
            )
            return err
        case <-quit:
        }
    }
    return nil
}

func tableRoute(router *chi.Mux) {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Handler", "Url", "Method"})

    walkFunc := func(method string, route string, handler http.Handler, middleware ...func(http.Handler) http.Handler) error {
        table.Append([]string{"", route, method}) // Append Walk all routes
        return nil
    }
    _ = chi.Walk(router, walkFunc)
    table.Render() // Send output
}

// StartWebServer starts a web server
func StartWebServer(addr string, readTimeout, writeTimeout int, handler http.Handler) *http.Server {
    return &http.Server{
        Addr:         addr,
        Handler:      handler,
        ReadTimeout:  time.Duration(readTimeout) * time.Second,
        WriteTimeout: time.Duration(writeTimeout) * time.Second,
    }
}
