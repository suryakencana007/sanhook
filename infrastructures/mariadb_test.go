// +build integration

/*  mariadb_test.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 18, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 18/10/18 00:47 
 */

package infrastructures

import (
    "fmt"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/suryakencana007/sanhook/configs"
    "github.com/suryakencana007/sanhook/pkg/sql"
)

func setupMysql() (string, *configs.Config) {
    c := configs.New("app.config",
        "./configs", "../configs", "../../configs")
    // slaveUser, slavePassword, slaveHost, slavePort, slaveDBName, slaveCharset)
    return fmt.Sprintf(MysqlDataSourceFormat, c.MariaDb.User, c.MariaDb.Password, c.MariaDb.Host, c.MariaDb.Port, c.MariaDb.DbName, c.MariaDb.Charset), c
}

func setupFailMysql(port int, dbname string) (string, *configs.Config) {
    c := configs.New("app.config",
        "./configs", "../configs", "../../configs")
    return fmt.Sprintf(MysqlDataSourceFormat, c.MariaDb.User, c.MariaDb.Password, c.MariaDb.Host, port, dbname, c.MariaDb.Charset), c
}

func openTestMysqlConn(t *testing.T) sql.DBFactory {
    connInfo, c := setupMysql()
    db := NewMysqlSQL()
    db.OpenConnection(connInfo, c)
    return db
}

func openTestMysqlConnFail(t *testing.T) sql.DBFactory {
    connInfo, c := setupFailMysql(3307, "")
    db := NewMysqlSQL()
    db.OpenConnection(connInfo, c)
    return db
}

func TestMysqlGetDB(t *testing.T) {
    db := openTestMysqlConn(t)
    defer db.Close()
    var sqlDb *sql.DB
    sqlDb, err := db.GetDB()
    assert.IsType(t, sqlDb, sqlDb)
    if err != nil {
        panic(fmt.Errorf("failed dial database : %v", err))
    }
}

func TestMysqlGetDBFail(t *testing.T) {
    f := func() {
        db := openTestMysqlConnFail(t)
        defer db.Close()
        dbType, _ := db.GetDB()
        assert.Nil(t, dbType)
    }
    assert.Panics(t, f)
    assert.PanicsWithValue(t, "dial tcp [::1]:3307: connect: connection refused", f)
}

func TestMysqlReconnect(t *testing.T) {
    db1 := openTestMysqlConn(t)
    defer db1.Close()
    tx, err := db1.Begin()
    if err != nil {
        t.Fatal(err)
    }
    var pid1 int
    err = tx.QueryRow("SELECT CONNECTION_ID()").Scan(&pid1)
    if err != nil {
        t.Fatal(err)
    }
    db2 := openTestMysqlConn(t)
    defer db2.Close()
    _, err = db2.Exec("Select concat('KILL ',id,';') from information_schema.processlist where id = ?;", pid1)
    if err != nil {
        t.Fatal(err)
    }
    // The rollback will probably "fail" because we just killed
    // its connection above
    _ = tx.Rollback()

    const expected int = 42
    var result int
    row, _ := db1.QueryRow(fmt.Sprintf("SELECT %d", expected))
    err = row.Scan(&result)
    if err != nil {
        t.Fatal(err)
    }
    if result != expected {
        t.Errorf("got %v; expected %v", result, expected)
    }
}

func TestMysqlCommitInFailedTransaction(t *testing.T) {
    db := openTestMysqlConn(t)
    defer db.Close()

    txn, err := db.Begin()
    if err != nil {
        t.Fatal(err)
    }
    rows, err := txn.Query("SELECT error")
    if err == nil {
        rows.Close()
        t.Fatal("expected failure")
    }
    err = txn.Commit()
    if err != nil {
        t.Fatalf("expected ErrInFailedTransaction; got %#v", err)
    }
}
