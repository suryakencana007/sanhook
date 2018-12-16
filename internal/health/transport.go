/*  transport.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 01, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 01/10/18 01:11 
 */

package health

import (
    "github.com/go-chi/chi"
    "github.com/suryakencana007/sanhook/configs"
)

// handleInternal handler for internal route
func MakeHandler(c *configs.Config) *chi.Mux {
    router := chi.NewRouter()
    router.Get("/health", getAPIStatus())
    return router
}
