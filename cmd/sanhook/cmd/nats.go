/*  nats.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               December 16, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2018-12-16 23:51 
 */

package cmd

import (
    "fmt"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "time"

    "github.com/nats-io/gnatsd/server"
    "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/spf13/pflag"
    "github.com/suryakencana007/sanhook/configs"
    "github.com/suryakencana007/sanhook/internal/container"
    "github.com/suryakencana007/sanhook/pkg/log"
)

type natsCmd struct {
    stop <-chan bool

    BaseCmd       *cobra.Command
    configuration *configs.Config
    container     *container.Container
    filename      string
}

func NewNatsCmd(
    configuration *configs.Config,
) *natsCmd {
    return NewNatsCmdSignaled(configuration, nil)
}

func NewNatsCmdSignaled(
    configuration *configs.Config,
    stop <-chan bool,
) *natsCmd {
    cc := &natsCmd{stop: stop}
    cc.configuration = configuration
    cc.BaseCmd = &cobra.Command{
        Use:   "nats",
        Short: "Used to run the NATS service",
        RunE:  cc.server,
    }
    fs := pflag.NewFlagSet("Root", pflag.ContinueOnError)
    fs.StringVarP(&cc.filename, "file", "f", "", "Custom configuration filename")
    cc.BaseCmd.Flags().AddFlagSet(fs)
    return cc
}

func (h *natsCmd) server(cmd *cobra.Command, args []string) (err error) {
    if len(h.filename) > 1 {
        h.configuration = configs.New(h.filename,
            "./configs",
            "../configs",
            "../../configs",
            "../../../configs")
    }

    var DefaultTestOptions = &server.Options{
        Host:           h.configuration.Nats.Host,
        Port:           h.configuration.Nats.Port,
        NoLog:          h.configuration.Nats.NoLog,
        NoSigs:         h.configuration.Nats.NoSigs,
        MaxControlLine: h.configuration.Nats.MaxControlLine,
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

    return h.serve(DefaultTestOptions)
}

// RunServer starts a new Go routine based server
func (h *natsCmd) serve(opts *server.Options) error {
    errCh := make(chan error, 1)
    quit := make(chan os.Signal)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    if opts == nil {
        opts = &server.Options{
            Host:           h.configuration.Nats.Host,
            Port:           h.configuration.Nats.Port,
            NoLog:          h.configuration.Nats.NoLog,
            NoSigs:         h.configuration.Nats.NoSigs,
            MaxControlLine: h.configuration.Nats.MaxControlLine,
        }
    }
    s := server.New(opts)
    if s == nil {
        panic("No NATS Server object returned.")
    }

    // Configure the logger based on the flags
    s.ConfigureLogger()

    // Run server in Go routine.
    go func() {
        // Start things up. Block here until done.
        if err := server.Run(s); err != nil {
            server.PrintAndDie(err.Error())
            errCh <- err
        }
        // Wait for accept loop(s) to be started
        if !s.ReadyForConnections(10 * time.Second) {
            panic("Unable to start NATS Server in Go Routine")
        }
        <-h.stop
    }()

    if h.stop != nil {
        select {
        case err := <-errCh:
            log.Info(
                "Server gracefully h stop stopped",
                logrus.Fields{})
            return err
        case <-h.stop:
        case <-quit:
        }
    } else {
        select {
        case err := <-errCh:
            log.Info(
                "Server gracefully stopped",
                logrus.Fields{})
            return err
        case <-quit:
        }
    }
    return nil
}
