/*  main.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               November 10, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 10/11/18 11:33 
 */

package container

import (
    "sync"

    "github.com/suryakencana007/sanhook/configs"
    "github.com/suryakencana007/sanhook/internal/message"
)

type Factory interface {
    Build() *Container
    NewMessageServiceJSON() message.Service
}

type Container struct {
    *configs.Config
}

func New(config *configs.Config) Factory {
    return &Container{Config: config}
}

var once sync.Once
var instance Factory

func Instance(config *configs.Config) Factory {
    once.Do(func() {
        instance = New(config)
    })
    return instance
}

func (c *Container) Build() *Container {
    return c
}
