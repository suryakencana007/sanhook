/*  transport.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               December 17, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2018-12-17 00:52 
 */

package nats

import (
    "fmt"
    "unsafe"

    "github.com/go-chi/chi"
    "github.com/suryakencana007/sanhook/configs"
    "github.com/suryakencana007/sanhook/internal/container"
)

// handleInternal handler for internal route
func MakeHandler(c *configs.Config) *chi.Mux {
    injector := container.Instance(c).Build()                        // inject just service needed
    fmt.Println("alamat memory injector:", unsafe.Pointer(injector)) // service needed
    router := chi.NewRouter()
    router.Get("/publish/{subject:[a-z-]+}", getPublishMessage(injector.NewMessageServiceJSON()))
    router.Get("/inbox/{id:[a-z-0-9]+}", getInboxMessage(injector.NewMessageServiceJSON()))
    router.Get("/inbox", getInboxAll(injector.NewMessageServiceJSON()))

    return router
}
