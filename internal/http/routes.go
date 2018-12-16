/*  routes.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               September 11, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 11/09/18 21:24 
 */

package http

import (
    "fmt"

    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "github.com/go-chi/render"
    "github.com/suryakencana007/sanhook/configs"
    "github.com/suryakencana007/sanhook/internal/health"
    "github.com/suryakencana007/sanhook/internal/nats"
    "github.com/suryakencana007/sanhook/pkg/datadog"
)

// Main Router
func Routes(configuration *configs.Config) *chi.Mux {
    router := chi.NewMux()
    router.Use(
        render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
        middleware.Logger,                             // Log API request calls
        middleware.DefaultCompress,                    // Compress results, mostly gzipping assets and json
        middleware.RedirectSlashes,                    // Redirect slashes to no slash URL versions
        middleware.Recoverer,                          // Recover from panics without crashing server
        datadog.DataDog("chi.mux"),                    // Datadog Tracer
    )

    // Rest API ver. 1
    router.Route("/v1", func(r chi.Router) {
        r.Mount(
            fmt.Sprintf(`%s/%s`, configuration.Api.Prefix, "status"),
            health.MakeHandler(configuration),
        )
        r.Mount(
            fmt.Sprintf(`%s/%s`, configuration.Api.Prefix, "message"),
            nats.MakeHandler(configuration),
        )
    })

    return router
}
