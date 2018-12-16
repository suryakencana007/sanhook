/*  root.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 19, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 19/10/18 01:48 
 */

package cmd

import (
    "github.com/spf13/cobra"
)

var WelkomText = `
========================================================================================  
                           888                        888      
                           888                        888      
                           888                        888      
.d8888b   8888b.  88888b.  88888b.   .d88b.   .d88b.  888  888 
88K          "88b 888 "88b 888 "88b d88""88b d88""88b 888 .88P 
"Y8888b. .d888888 888  888 888  888 888  888 888  888 888888K  
     X88 888  888 888  888 888  888 Y88..88P Y88..88P 888 "88b 
 88888P' "Y888888 888  888 888  888  "Y88P"   "Y88P"  888  888
========================================================================================
- port    : %d
- log     : %s
-----------------------------------------------------------------------------------------`

var ConfigPath = []string{
    "./configs",
    "../configs",
    "../../configs",
    "../../../configs"}

func RootCmd() *cobra.Command {
    root := &cobra.Command{
        Use:   "sanhook",
        Short: "sanhook - Stock::Inventory µService",
        Long:  "sanhook is Micro Service for Inventory Management",
        Args:  cobra.MinimumNArgs(1),
    }
    return root
}
