/*  message.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               February 08, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-02-08 03:21 
 */

package http

import (
    "fmt"
    "unsafe"

    "github.com/go-chi/chi"
    "github.com/suryakencana007/sanhook/configs"
    "github.com/suryakencana007/sanhook/internal/container"
    "github.com/suryakencana007/sanhook/internal/message"
)

// handleInternal handler for internal route
func messageMakeHandler(c *configs.Config) *chi.Mux {
    injector := container.Instance(c).Build()                        // inject just service needed
    fmt.Println("alamat memory injector:", unsafe.Pointer(injector)) // service needed
    router := chi.NewRouter()
    router.Get("/publish/{subject:[a-z-]+}", message.GetPublishMessage(c, injector.NewMessageServiceJSON()))
    router.Get("/inbox/{id:[a-z-0-9]+}", message.GetInboxMessage(injector.NewMessageServiceJSON()))
    router.Get("/inbox", message.GetInboxAll(injector.NewMessageServiceJSON()))

    return router
}
