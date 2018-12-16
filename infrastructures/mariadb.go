/*  mariadb.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               September 11, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 11/09/18 21:26 
 */

package infrastructures

import (
    _ "github.com/go-sql-driver/mysql"
    "github.com/sirupsen/logrus"
    "github.com/suryakencana007/sanhook/configs"
    "github.com/suryakencana007/sanhook/pkg/log"
    "github.com/suryakencana007/sanhook/pkg/sql"
)

var MysqlDataSourceFormat = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local"

type MysqlInfrastructure struct {
    *sql.DB
}

func NewMysqlSQL() sql.DBFactory {
    return &MysqlInfrastructure{}
}

// OpenConnection gets a handle for a database
func (s *MysqlInfrastructure) OpenConnection(connString string, config *configs.Config) {

    db, err := sql.Open(sql.MYSQL, connString, config)
    if err != nil {
        panic(err)
    }
    s.DB = db
    _, err = s.GetDB()
    if err != nil {
        panic(err.Error())
    }
}

// GetDB gets database connection
func (s *MysqlInfrastructure) GetDB() (*sql.DB, error) {
    err := s.Ping()
    if err != nil {
        log.Error("Get DB", logrus.Fields{
            "error": err.Error(),
        })
        return nil, err
    }

    return s.DB, nil
}
