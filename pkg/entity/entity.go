/*  entity.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               November 06, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 06/11/18 01:16 
 */

package entity

import (
    "errors"
    "time"

    "github.com/satori/go.uuid"
)

func StringID() string {
    return uuid.NewV4().String()
}

type Event struct {
    CreatedAt time.Time `json:"created_at" rql:"filter"`
    CreatedBy string     `json:"created_by" rql:"filter"`
    UpdatedAt time.Time `json:"updated_at" rql:"filter"`
    UpdatedBy string     `json:"updated_by" rql:"filter"`
    DeletedAt time.Time `json:"deleted_at"`
}

var (
    ErrUnknown   = errors.New("unknown Error")
    ErrDuplicate = errors.New("it already exists")
)
