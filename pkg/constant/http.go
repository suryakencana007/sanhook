/*  http.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 01, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 01/10/18 00:24 
 */

package constant

const (
    // Message for response general success
    CodeGeneralSuccess    = 1000
    MessageGeneralSuccess = "Success"

    // Message for response general error
    CodeErrorUnknown    = 2000
    MessageGeneralError = "Error message"

    // InternalServerErrorDetail constant for internal server error detail
    CodeInternalServerError    = 2003
    MessageInternalServerError = "Oops something went wrong"

    CodeInvalidAuthentication   = 2004
    MessagInvalidAuthentication = "The resource owner or authorization server denied the request"

    CodeValidationError = 2005
    CodeForbidden       = 2006
    MessageForbidden    = "Forbidden Access"

    // InvalidDataDetail constant for unprocessable error detail
    CodeInvalidData    = 2007
    MessageInvalidData = "One or more data is failing the validation rule"

    CodeUnauthorized    = 2019
    MessageUnauthorized = "not authorized to access the service"
)
