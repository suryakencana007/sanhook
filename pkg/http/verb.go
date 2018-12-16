/*  verb.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 08, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 08/10/18 16:52 
 */

package http

import (
    "bytes"
    "encoding/json"
    "io"
    "strings"
)

/**
* Convert Body to Json
* :body: interface
* :return: io.Reader
*/
func BodyToJson(body interface{}) (io.Reader, error) {
    if body == nil {
        return nil, nil
    }

    switch v := body.(type) {
    case string:
        // return as is
        return strings.NewReader(v), nil
    default:
        b, err := json.Marshal(v)
        if err != nil {
            return nil, err
        }

        return bytes.NewReader(b), nil
    }
}
