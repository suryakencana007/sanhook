/*  message.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               December 17, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2018-12-17 02:29 
 */

package entity

type RepositoryMessage interface {
    Find(id string) (Message, error)
    FindAll() ([]Message, error)
    Update(msg Message) error
    Store(msg Message) (string, error)
}

type Message struct {
    ID      string `json:"id"`
    Subject string `json:"subject"`
    Content string `json:"content"`

    Event
}

// TableName of the role
func (Message) TableName() string {
    return "tb_message"
}
