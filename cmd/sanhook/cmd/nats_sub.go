/*  nats_sub.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               December 17, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2018-12-17 02:03 
 */

package cmd

import (
    "fmt"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"

    "github.com/nats-io/go-nats"
    "github.com/spf13/cobra"
    "github.com/spf13/pflag"
    "github.com/suryakencana007/sanhook/configs"
    "github.com/suryakencana007/sanhook/internal/container"
    "github.com/suryakencana007/sanhook/pkg/log"
)

type subscribeCmd struct {
    stop <-chan bool

    BaseCmd       *cobra.Command
    configuration *configs.Config
    container     *container.Container
    filename      string
    subject       string
}

func NewSubscribeCmd(
    configuration *configs.Config,
) *subscribeCmd {
    return NewSubscribeCmdSignaled(configuration, nil)
}

func NewSubscribeCmdSignaled(
    configuration *configs.Config,
    stop <-chan bool,
) *subscribeCmd {
    cc := &subscribeCmd{stop: stop}
    cc.configuration = configuration
    cc.BaseCmd = &cobra.Command{
        Use:   "subscribe",
        Short: "Used to run the NATS Subscribe",
        RunE:  cc.subscribe,
    }
    fs := pflag.NewFlagSet("Root", pflag.ContinueOnError)
    fs.StringVarP(&cc.filename, "file", "f", "", "Custom configuration filename")
    fs.StringVarP(&cc.subject, "subject", "s", "", "subject for subscribe message")
    cc.BaseCmd.Flags().AddFlagSet(fs)
    return cc
}

func (h *subscribeCmd) subscribe(cmd *cobra.Command, args []string) (err error) {
    if len(h.filename) > 1 {
        h.configuration = configs.New(h.filename,
            "./configs",
            "../configs",
            "../../configs",
            "../../../configs")
    }

    // Description µ micro service
    fmt.Println(
        fmt.Sprintf(
            WelkomText,
            h.configuration.Nats.Port,
            strings.Join([]string{
                h.configuration.Log.Dir,
                h.configuration.Log.Filename}, "/"),
        ))

    return h.serve()
}

func (h *subscribeCmd) serve() error {
    errCh := make(chan error, 1)
    quit := make(chan os.Signal)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

    urls := h.configuration.Nats.Host
    // Connect Options.

    // Connect Options.
    opts := []nats.Option{nats.Name("NATS Sample Subscriber")}
    opts = setupConnOptions(opts)
    // Connect to NATS
    nc, err := nats.Connect(urls, opts...)
    if err != nil {
        log.Error(
            "Subscribe Error",
            log.Field("error", err.Error()),
        )
    }
    defer nc.Close()

    // Run server in Go routine.
    go func() {
        i := 0
        if sub, err := nc.Subscribe(
            h.subject,
            func(msg *nats.Msg) {
                i += 1
                printMsg(msg, i)
            });
            err != nil {
            sub.Unsubscribe()
            log.Error(
                "Subscribe Error",
                log.Field("error", err.Error()),
            )
            errCh <- err
        }
        if err := nc.Flush();
            err != nil {
            log.Error(
                "Subscribe Error",
                log.Field("error", err.Error()),
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

func setupConnOptions(opts []nats.Option) []nats.Option {
    totalWait := 10 * time.Minute
    reconnectDelay := time.Second

    opts = append(opts, nats.ReconnectWait(reconnectDelay))
    opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
    opts = append(opts, nats.DisconnectHandler(func(nc *nats.Conn) {
        log.Info(
            "Setup Opts",
            log.Field("message", fmt.Sprintf(
                "Disconnected: will attempt reconnects for %.0fm",
                totalWait.Minutes()),
            ),
        )
    }))
    opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
        log.Info(
            fmt.Sprintf("Reconnected [%s]", nc.ConnectedUrl()),
        )
    }))
    opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
        log.Info(
            "Exiting, no servers available",
        )
    }))
    return opts
}

func printMsg(m *nats.Msg, i int) {
    log.Info(
        "Subscribe",
        log.Field(m.Subject, fmt.Sprintf("[#%d] Received on [%s]: '%s'", i, m.Subject, string(m.Data))),
    )
}
