/*  endpoint.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               December 17, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2018-12-17 00:51 
 */

package message

import (
    "fmt"
    "net/http"
    "strings"

    "github.com/go-chi/chi"
    "github.com/nats-io/go-nats"
    "github.com/suryakencana007/sanhook/configs"
    "github.com/suryakencana007/sanhook/pkg/entity"
    "github.com/suryakencana007/sanhook/pkg/log"
    "github.com/suryakencana007/sanhook/pkg/response"
)

// Publish Message
func GetPublishMessage(c *configs.Config, svc Service) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        messageID := make(chan string, 1)
        subject := chi.URLParam(r, "subject")
        if len(subject) < 1 {
            response.Write(w, r, response.APIErrorInvalidData)
            return
        }
        log.Debug("Get Publish Message", log.Field("subject", subject))
        params := r.URL.Query()
        msg, ok := params["message"]
        if !ok {
            msg = []string{""}
        }
        log.Debug("Param:message", log.Field("Message", msg))
        go func() {
            urls := c.Nats.Host
            // Connect Options.
            opts := []nats.Option{nats.Name("NATS Sample Publisher")}

            // Connect to NATS
            nc, err := nats.Connect(urls, opts...)
            if err != nil {
                log.Fatal("Error", log.Field("Error", err))
            }
            defer nc.Close()

            if err := nc.Publish(subject, []byte(strings.Join(msg, " ")));
                err != nil {
                log.Fatal("Error", log.Field("Error", err))
            }
            if err := nc.Flush();
                err != nil {
                log.Fatal("Error", log.Field("Error", err))
            }
            if err := nc.LastError(); err != nil {
                log.Fatal("Error", log.Field("Error", err))
            } else {
                log.Info("getPublishMessage", log.Field("message", fmt.Sprintf("Published [%s] : '%s'", subject, msg)))
            }
            store := entity.Message{}
            store.Subject = subject
            store.Content = strings.Join(msg, " ")
            if id, err := svc.StoreMessage(store); err == nil {
                log.Info("getPublishMessage", log.Field("message", fmt.Sprintf("Success Store Subject [%s] - ID: '%s'", subject, id)))
                messageID <- id
            }
        }()

        message := svc.GetMessage(<-messageID)
        resp := response.APIOK
        resp.Data = map[string]interface{}{
            "ID":      message.ID,
            "Subject": message.Subject,
            "Content": message.Content,
            "Created": message.CreatedAt,
        }
        response.Write(w, r, resp)
    }
}

func GetInboxMessage(svc Service) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")
        if len(id) < 1 {
            response.Write(w, r, response.APIErrorInvalidData)
            return
        }
        message := svc.GetMessage(id)
        resp := response.APIOK
        resp.Data = map[string]interface{}{
            "ID":      message.ID,
            "Subject": message.Subject,
            "Content": message.Content,
            "Created": message.CreatedAt,
        }
        response.Write(w, r, resp)
    }
}

func GetInboxAll(svc Service) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var messages []map[string]interface{}
        for _, message := range svc.All() {
            msg := map[string]interface{}{
                "ID":      message.ID,
                "Subject": message.Subject,
                "Content": message.Content,
                "Created": message.CreatedAt,
            }
            messages = append(messages, msg)
        }
        response.Write(w, r, messages)
    }
}
