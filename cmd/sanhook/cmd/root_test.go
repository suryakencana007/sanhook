/*  root_test.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 19, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 19/10/18 13:36 
 */

package cmd

import (
    "bytes"
    "testing"

    "github.com/spf13/cobra"
    "github.com/stretchr/testify/assert"
)

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
    buf := new(bytes.Buffer)
    root.SetOutput(buf)
    root.SetArgs(args)

    c, err = root.ExecuteC()

    return c, buf.String(), err
}


func TestRootCmd(t *testing.T) {
    rootCmd := RootCmd()

    _, output, err := executeCommandC(rootCmd)
    assert.Contains(t, output, "sanhook is Micro Service for Inventory Management")
    assert.Nil(t, err)
}
