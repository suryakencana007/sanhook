/*  config.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               September 30, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 30/09/18 01:43 
 */

package configs

import (
    "github.com/spf13/viper"
    "github.com/suryakencana007/sanhook/pkg/log"
)

type Constants struct {
    App struct {
        Name         string `mapstructure:"name"`
        Port         int    `mapstructure:"port"`
        ReadTimeout  int    `mapstructure:"read_timeout"`
        WriteTimeout int    `mapstructure:"write_timeout"`
        Timezone     string `mapstructure:"timezone"`
        Debug        bool   `mapstructure:"debug"`
        Env          string `mapstructure:"env"`
        SecretKey    string `mapstructure:"secret_key"`
    }
    Api struct {
        Prefix string `mapstructure:"prefix"`
    }
    Log struct {
        Dir      string `mapstructure:"dir"`
        Filename string `mapstructure:"filename"`
    }
    CB struct {
        Retry      int `mapstructure:"retry_count"`
        Timeout    int `mapstructure:"db_timeout"`
        Concurrent int `mapstructure:"max_concurrent"`
    }
    Nats struct {
        Host           string `mapstructure:"host"`
        Port           int    `mapstructure:"port"`
        NoLog          bool   `mapstructure:"no_log"`
        NoSigs         bool   `mapstructure:"no_signal"`
        MaxControlLine int    `mapstructure:"max_control_line"`
    }
    Postgres struct {
        DbName            string `mapstructure:"db_name"`
        Host              string `mapstructure:"host"`
        Port              int    `mapstructure:"port"`
        User              string `mapstructure:"user"`
        Password          string `mapstructure:"password"`
        MaxLifeTime       int    `mapstructure:"max_life_time"`
        MaxIdleConnection int    `mapstructure:"max_idle_connection"`
        MaxOpenConnection int    `mapstructure:"max_open_connection"`
    }
    MariaDb struct {
        DbName            string `mapstructure:"db_name"`
        Host              string `mapstructure:"host"`
        Port              int    `mapstructure:"port"`
        User              string `mapstructure:"user"`
        Password          string `mapstructure:"password"`
        MaxLifeTime       int    `mapstructure:"max_life_time"`
        MaxIdleConnection int    `mapstructure:"max_idle_connection"`
        MaxOpenConnection int    `mapstructure:"max_open_connection"`
        Charset           string `mapstructure:"charset"`
    }
}

type Config struct {
    Constants
}

// NewConfig is used to generate a configuration instance which will be passed around the codebase
func New(filename string, paths ...string) *Config {
    if filename == "" {
        filename = "app.config"
    }
    constants := initViper(filename, paths...)
    // init log
    // log.Init(constants.Log.Dir, constants.Log.Filename, constants.App.Debug)
    // log.LogrusInit()
    log.ZapInit()
    // init Infrastructure DB
    return &Config{
        constants,
    }
}

func initViper(filename string, paths ...string) (constants Constants) {
    vip := viper.New()
    // Search the root directory for the configuration file
    for _, path := range paths {
        vip.AddConfigPath(path)
    }
    // Configuration fileName without the .TOML or .YAML extension
    vip.SetConfigName(filename)
    if err := vip.ReadInConfig(); err != nil {
        panic(err.Error())
    }

    vip.WatchConfig() // Watch for changes to the configuration file and recompile
    if err := vip.UnmarshalExact(&constants); err != nil {
        panic(err.Error())
    }
    return constants
}
