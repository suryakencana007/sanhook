/*  logrus.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               October 15, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 15/10/18 23:32 
 */

package log

import (
    "fmt"
    "io"
    "os"
    "strings"

    "github.com/go-stack/stack"
    "github.com/sirupsen/logrus"
)

var (
    logger = logrus.New()
)

func Init(dir, filename string, debug bool) {
    SetFormatter(&logrus.JSONFormatter{})
    if debug {
        FileHandler(dir, filename) // log file handler
    }

}

// FileHandler handles log to file
func FileHandler(dir, filename string) {
    path := strings.Join([]string{dir, filename}, "/")
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        Info("Create Dir", logrus.Fields{
            "path":  path,
            "error": err,
        })
        err = os.MkdirAll(dir, 0755)
        if err != nil {
            panic(err)
        }
    }
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
    if err != nil {
        panic(fmt.Sprintf("error opening file: %v", err))
    }
    SetOutput(file)
}

// SetOutput sets the standard logger output.
func SetOutput(out io.Writer) {
    logger.SetOutput(out)
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter logrus.Formatter) {
    logger.SetFormatter(formatter)
}

// SetLevel sets the standard logger level.
func SetLevel(level logrus.Level) {
    logger.SetLevel(level)
}

// GetLevel returns the standard logger level.
func GetLevel() logrus.Level {
    return logger.GetLevel()
}

// IsLevelEnabled checks if the log level of the standard logger is greater than the level param
func IsLevelEnabled(level logrus.Level) bool {
    return logger.IsLevelEnabled(level)
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook logrus.Hook) {
    logger.AddHook(hook)
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *logrus.Entry {
    return logger.WithField(logrus.ErrorKey, err)
}

// Debug logs information interesting for Developers
func Debug(msg interface{}, fields interface{}) {
    withFields(fields).Debug(fmt.Sprint(msg))
}

// Info logs information interesting for Support staff trying to figure out the context of a given error
func Info(msg interface{}, fields interface{}) {
    withFields(fields).Info(fmt.Sprint(msg))
}

// Error logs information of an error occurred in error handling
func Error(msg interface{}, fields interface{}) {
    withFields(fields).Error(fmt.Sprint(msg))
}

// Panic logs a message at level Panic on the standard logger.
func Panic(msg interface{}, fields interface{}) {
    withFields(fields).Panic(fmt.Sprint(msg))
}

func withFields(fields interface{}) *logrus.Entry {
    strGoStack := "%n"
    f := fields.(logrus.Fields)
    f["caller"] = stack.Caller(2)
    f[ "function"] = fmt.Sprintf(strGoStack, stack.Caller(2))
    return logger.WithFields(f)
}
