/*  message.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               December 17, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2018-12-17 02:31 
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
