/*  config_test.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 10, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 10/10/18 00:34 
 */

package configs

import (
    "bytes"
    "testing"

    "github.com/spf13/viper"
    "github.com/stretchr/testify/assert"
)

var config = []byte(`
[App] # App configuration
name = "µs-wyonna"
port = 8080
read_timeout = 5 # seconds
write_timeout = 10 # seconds
timezone = "Asia/Jakarta"
debug = false # (true|false)
env = "development" # ('development'|'staging'|'production')
secret_key = ""

[Api] # Rest API configuration
prefix = "/api"

[Log] # Log configuration
dir = "logs"
filename = "wyonna.log"

[CB]
retry_count = 3
db_timeout = 1500
max_concurrent = 5

[Nats]
host = "0.0.0.0"
port = 4222
no_log = true
no_signal = true
max_control_line = 2048

[Postgres]
db_name = "wyonna"
host = "localhost"
port = 54322
user = "root"
password = "root"
max_life_time = 30
max_idle_connection = 0
max_open_connection = 0

[Mariadb]
db_name = "wyonna"
host = "localhost"
port = 3306
user = "root"
password = "root"
max_life_time = 30
max_idle_connection = 0
max_open_connection = 0
charset = "utf8"
`)

var configNoLog = []byte(`
[App] # App configuration
name = "µs-wyonna"
port = 8080
read_timeout = 5 # seconds
write_timeout = 10 # seconds
timezone = "Asia/Jakarta"
debug = false # (true|false)
env = "development" # ('development'|'staging'|'production')
secret_key = ""

[Api] # Rest API configuration
prefix = "/api"`)

func initConfig(t *testing.T, c []byte) {
    viper.Reset()
    viper.New()
    viper.SetConfigType("toml")
    r := bytes.NewReader(c)

    var (
        err error
        n   int64
    )
    buf := new(bytes.Buffer)
    n, err = buf.ReadFrom(r)

    assert.IsType(t, int64(345), n)
    assert.Nil(t, err)
    err = viper.ReadConfig(buf)
    assert.Nil(t, err)
}

func TestNew(t *testing.T) {
    initConfig(t, config)
    var constants Constants
    err := viper.Unmarshal(&constants)
    assert.Nil(t, err)

    c := Config{}
    c.Constants = constants
    configuration := New("app.config",
        "./configs", "../configs", "../../configs")
    assert.Equal(t, &c, configuration)

}

func TestNewNoLog(t *testing.T) {
    initConfig(t, configNoLog)
    var constants Constants
    err := viper.Unmarshal(&constants)
    assert.Nil(t, err)

    c := Config{}
    c.Constants = constants
    configuration := New("app.config",
        "./configs", "../configs", "../../configs")
    assert.NotEqual(t, &c, configuration)
}

func TestNewFailTOML(t *testing.T) {
    f := func() {
        New("app.config.test",
            "./configs", "../configs", "../../configs")
    }
    assert.Panics(t, f)
    assert.PanicsWithValue(t, "1 error(s) decoding:\n\n* '' has invalid keys: appa", f)
}

func TestNewConfigNotFound(t *testing.T) {
    f := func() {
        New("", "")
    }
    assert.Panics(t, f)
    assert.PanicsWithValue(t, "Config File \"app.config\" Not Found in \"[]\"", f)
}

func TestInitViperConfigNotFound(t *testing.T) {
    f := func() {
        initViper("app.config")
    }
    assert.Panics(t, f)
    assert.PanicsWithValue(t, "Config File \"app.config\" Not Found in \"[]\"", f)
}
