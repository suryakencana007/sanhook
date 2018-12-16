/*  service.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               December 17, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2018-12-17 02:32 
 */

package message

import (
    "github.com/suryakencana007/sanhook/pkg/entity"
)

type Service interface {
    All() []entity.Message
    GetMessage(id string) entity.Message
    InsertSample() (string, error)
    StoreMessage(msg entity.Message) (string, error)
}

type service struct {
    msg entity.RepositoryMessage
}

func NewService(mainRepo entity.RepositoryMessage) Service {
    return &service{msg: mainRepo}
}

func (s *service) All() []entity.Message {
    messages, err := s.msg.FindAll()
    if err != nil {
        return nil
    }
    return messages
}

func (s *service) GetMessage(id string) entity.Message {
    msg, err := s.msg.Find(id)
    if err != nil {
        return msg
    }
    return msg
}

func (s *service) InsertSample() (string, error) {
    msg := entity.Message{}
    msg.Subject = "Item-B002"
    msg.Content = "Super Item Message"
    id, err := s.msg.Store(msg)
    if err != nil {
        return "", err
    }
    return id, nil
}

func (s *service) StoreMessage(msg entity.Message) (string, error) {
    id, err := s.msg.Store(msg)
    if err != nil {
        return "", err
    }
    return id, nil
}

