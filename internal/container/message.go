/*  message.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               February 08, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-02-08 03:20 
 */

package container

import (
    "github.com/suryakencana007/sanhook/internal/message"
)

func (c *Container) NewMessageServiceJSON() message.Service {
    repoMsg, err := message.NewJSON()
    if err != nil {
        panic(err)
    }
    return message.NewService(repoMsg)
}
