/*  main_test.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 18, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 18/10/18 02:16 
 */

package main

import (
    "os"
    "strings"
    "testing"
)

func TestAppMain(t *testing.T) {
    // don't launch etcd server when invoked via go test
    if strings.HasSuffix(os.Args[0], "main") {
        return
    }
    main()
}
