/*  dir.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               February 04, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-02-04 23:37 
 */

package helper

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/suryakencana007/sanhook/pkg/log"
)

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func AbsPathify(inPath string) string {
    if strings.HasPrefix(inPath, "$") {
        end := strings.Index(inPath, string(os.PathSeparator))
        inPath = os.Getenv(inPath[1:end]) + inPath[end:]
    }

    if filepath.IsAbs(inPath) {
        return filepath.Clean(inPath)
    }

    p, err := filepath.Abs(inPath)
    if err == nil {
        return filepath.Clean(p)
    }
    return ""
}

func FindFiles(filename string, debug bool, paths ...string) {
    if debug {
        var dir []string
        // dw, _ := os.Getwd()
        for _, path := range paths {
            if path != "" {
                absin := AbsPathify(path)
                if !stringInSlice(absin, dir) {
                    FileHandler(absin, filename, true) // log file handler
                    dir = append(dir, absin)
                }
            }
        }
    }

}

func FileHandler(dir, filename string, read bool) *os.File {
    path := strings.Join([]string{dir, filename}, "/")
    if _, err := os.Stat(dir); os.IsNotExist(err) && !read {
        log.Info("Create Dir",
            log.Field("path", path),
            log.Field("error", err),
        )
        err = os.MkdirAll(dir, 0755)
        if err != nil {
            log.Error("File Handler Message:",
                log.Field("Error",
                    fmt.Sprintf("error creating file: %v", err),
                ))
        }
    }
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
    if err != nil {
        log.Error("File Handler Message:",
            log.Field("Error",
                fmt.Sprintf("error opening file: %v", err),
            ))
        return nil
    }
    return file
}
