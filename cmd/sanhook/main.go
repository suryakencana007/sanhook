/*  main.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               September 30, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 30/09/18 01:32 
 */

package main

import (
    "os"
    "runtime"

    "github.com/spf13/pflag"
    "github.com/suryakencana007/sanhook/cmd/sanhook/cmd"
    "github.com/suryakencana007/sanhook/configs"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU()) // optimize runtime
    var filename string
    root := cmd.RootCmd()
    fs := pflag.NewFlagSet("Root", pflag.ContinueOnError)
    fs.StringVarP(&filename,
        "file",
        "f",
        "",
        "Custom configuration filename",
    )
    root.Flags().AddFlagSet(fs)
    configuration := configs.New(filename, cmd.ConfigPath...)
    root.AddCommand(
        cmd.NewHttpCmd(
            configuration,
        ).BaseCmd,
    )
    root.AddCommand(
        cmd.NewNatsCmd(
            configuration,
        ).BaseCmd,
    )
    root.AddCommand(
        cmd.NewSubscribeCmd(
            configuration,
        ).BaseCmd,
    )
    if err := root.Execute(); err != nil {
        panic(err.Error())
        os.Exit(1)
    }
}
