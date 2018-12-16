/*  endpoint.go
*
* @Author:             Nanang Suryadi <nanang.suryadi@kubuskotak.com>
* @Date:               December 17, 2018
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2018-12-17 00:51 
 */

package nats

import (
    "log"
    "net/http"
    "strings"

    "github.com/go-chi/chi"
    "github.com/nats-io/go-nats"
    "github.com/suryakencana007/sanhook/internal/message"
    "github.com/suryakencana007/sanhook/pkg/entity"
    "github.com/suryakencana007/sanhook/pkg/response"
)

// Publish Message
func getPublishMessage(svc message.Service) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        messageID := make(chan string, 1)
        subject := chi.URLParam(r, "subject")
        if len(subject) < 1 {
            response.Write(w, r, response.APIErrorInvalidData)
            return
        }
        log.Println(subject)
        params := r.URL.Query()
        msg, ok := params["message"]
        if !ok {
            msg = []string{""}
        }
        log.Println(msg)
        go func() {
            urls := nats.DefaultURL
            // Connect Options.
            opts := []nats.Option{nats.Name("NATS Sample Publisher")}

            // Connect to NATS
            nc, err := nats.Connect(urls, opts...)
            if err != nil {
                log.Fatal(err)
            }
            defer nc.Close()

            if err := nc.Publish(subject, []byte(strings.Join(msg, " ")));
                err != nil {
                log.Fatal(err)
            }
            if err := nc.Flush();
                err != nil {
                log.Fatal(err)
            }
            if err := nc.LastError(); err != nil {
                log.Fatal(err)
            } else {
                log.Printf("Published [%s] : '%s'\n", subject, msg)
            }
            store := entity.Message{}
            store.Subject = subject
            store.Content = strings.Join(msg, " ")
            if id, err := svc.StoreMessage(store); err == nil {
                log.Printf("Success Store Subject [%s] - ID: '%s'\n", subject, id)
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

func getInboxMessage(svc message.Service) http.HandlerFunc {
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

func getInboxAll(svc message.Service) http.HandlerFunc {
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
