/*  response.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 01, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 01/10/18 00:21 
 */

package response

import (
    "net/http"

    "github.com/go-chi/render"
    "github.com/suryakencana007/sanhook/pkg/constant"
)

type ErrorValidation struct {
    Errors interface{} `json:"errors"`
}

// APIResponse defines attributes for api Response
type APIResponse struct {
    HTTPCode   int         `json:"-"`
    Code       int         `json:"code"`
    Message    interface{} `json:"message"`
    Data       interface{} `json:"data,omitempty"`
    Pagination interface{} `json:"pagination,omitempty"`
}

// Write writes the data to http response writer
func Write(w http.ResponseWriter, r *http.Request, response interface{}) {
    render.JSON(w, r, response)
}

// Defaults API Response with standard HTTP Status Code. The default value can
// be changed either using `ModifyMessage` or `ModifyHTTPCode`. You can call it
// directly by :
//
// response.Write(res, response.APIErrorUnknown)
// return
var (
    // A generic error message, given when an unexpected condition was encountered.
    APIErrorUnknown = APIResponse{
        HTTPCode: http.StatusInternalServerError,
        Code:     constant.CodeInternalServerError,
        Message:  constant.MessageInternalServerError,
    }

    // Standard response for successful HTTP requests.
    APIOK = APIResponse{
        HTTPCode: http.StatusOK,
        Code:     constant.CodeGeneralSuccess,
        Message:  constant.MessageGeneralSuccess,
    }

    // The request has been fulfilled, resulting in the creation of a new resource
    APICreated = APIResponse{
        HTTPCode: http.StatusCreated,
        Code:     constant.CodeGeneralSuccess,
        Message:  constant.MessageGeneralError,
    }

    // The request has been accepted for processing, but the processing has not been completed.
    APIAccepted = APIResponse{
        HTTPCode: http.StatusAccepted,
        Code:     constant.CodeGeneralSuccess,
        Message:  constant.MessageGeneralSuccess,
    }

    // The server cannot or will not process the request due to an apparent client error (e.g., malformed request syntax
    // , size too large, invalid request message framing, or deceptive request routing).
    APIErrorValidation = APIResponse{
        HTTPCode: http.StatusBadRequest,
        Code:     constant.CodeValidationError,
        Message:  constant.MessageGeneralError,
    }

    // The server cannot or will not process the request due to an apparent client error (e.g., malformed request syntax
    // , size too large, invalid request message framing, or deceptive request routing).
    APIErrorInvalidPassword = APIResponse{
        HTTPCode: http.StatusBadRequest,
        Code:     constant.CodeInvalidAuthentication,
        Message:  constant.MessageInvalidData,
    }

    APIErrorInvalidData = APIResponse{
        HTTPCode: http.StatusBadRequest,
        Code:     constant.CodeInvalidData,
        Message:  constant.MessageInvalidData,
    }

    // The request was valid, but the server is refusing action. The user might not have the necessary permissions for
    // a resource, or may need an account of some sort.
    APIErrorForbidden = APIResponse{
        HTTPCode: http.StatusForbidden,
        Code:     constant.CodeForbidden,
        Message:  constant.MessageForbidden,
    }

    // The request was valid, but the server is refusing action. The user might not have the necessary permissions for
    // a resource, or may need an account of some sort.
    APIErrorUnauthorized = APIResponse{
        HTTPCode: http.StatusUnauthorized,
        Code:     constant.CodeUnauthorized,
        Message:  constant.MessageUnauthorized,
    }
)
