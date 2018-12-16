/*  repository_json.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               December 17, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2018-12-17 02:37 
 */

package message

import (
    "encoding/json"
    "errors"
    "time"

    "github.com/nanobox-io/golang-scribble"
    "github.com/suryakencana007/sanhook/pkg/entity"
)

const location = "./storage/json/" // location defines where the files are stored

// JSON is the data storage layered using JSON file
type JSON struct {
    db *scribble.Driver
}

func NewJSON() (entity.RepositoryMessage, error) {
    db, err := scribble.New(location, nil)
    if err != nil {
        return nil, err
    }

    return &JSON{db: db}, nil
}

func (s *JSON) Find(id string) (entity.Message, error) {
    var msg entity.Message

    if err := s.db.Read(msg.TableName(), id, &msg); err != nil {
        return msg, entity.ErrUnknown
    }
    return msg, nil
}

func (s *JSON) Update(msg entity.Message) error {
    return nil
}

func (s *JSON) Store(msg entity.Message) (string, error) {
    var resource = entity.StringID()

    i, err := s.Find(msg.ID)
    if err != nil && len(i.ID) > 0 {

        return "", entity.ErrDuplicate
    }

    msg.ID = resource
    msg.CreatedAt = time.Now()
    msg.UpdatedAt = time.Now()

    if err := s.db.Write(msg.TableName(), resource, msg); err != nil {
        return "", err
    }
    return msg.ID, nil
}

func (s *JSON) FindAll() ([]entity.Message, error) {
    var list []entity.Message
    var msg entity.Message

    records, err := s.db.ReadAll(msg.TableName())
    if err != nil {
        return nil, errors.New("error while fetching items from the JSON file storage: %v" + err.Error())
    }

    for _, b := range records {
        if err := json.Unmarshal([]byte(b), &msg); err != nil {
            return list, nil
        }
        list = append(list, msg)
    }

    return list, nil
}
