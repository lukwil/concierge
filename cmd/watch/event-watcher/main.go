package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/shurcooL/graphql"

	"github.com/lukwil/concierge/cmd/common/hasura"
	"github.com/lukwil/concierge/cmd/watch/event-watcher/distributed"
	"github.com/lukwil/concierge/cmd/watch/event-watcher/single"
)

var status struct {
	Status []struct {
		ID   graphql.Int
		Name graphql.String
	}
}

var statusIDName map[graphql.String]graphql.Int

func main() {
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", 8080),
		ReadTimeout:    3 * time.Second,
		WriteTimeout:   3 * time.Second,
		MaxHeaderBytes: 1 << 20, // Max header of 1MB
	}

	getStatusTypes()

	go single.Watch(statusIDName)
	go distributed.Watch(statusIDName)

	http.HandleFunc("/", okHandler)
	server.ListenAndServe()
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func getStatusTypes() {
	statusIDName = make(map[graphql.String]graphql.Int)
	client := hasura.Client()
	if err := client.Query(context.TODO(), &status, make(map[string]interface{})); err != nil {
		panic(err)
	}
	for _, s := range status.Status {
		statusIDName[s.Name] = s.ID
	}
}
