/*  endpoint.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 01, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 01/10/18 01:11 
 */

package health

import (
    "net/http"

    "github.com/suryakencana007/sanhook/pkg/response"
)

// GetStatus gets the service health status
func getAPIStatus() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        response.Write(w, r, response.APIOK)
    }
}
