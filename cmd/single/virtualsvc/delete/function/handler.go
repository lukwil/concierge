package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/lukwil/concierge/cmd/common/nats"
)

type message struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type tableSpec struct {
	Event struct {
		Data struct {
			Old struct {
				NameK8S string `json:"name_k8s"`
			} `json:"old"`
		} `json:"data"`
	} `json:"event"`
	Table struct {
		Schema string `json:"schema"`
		Name   string `json:"name"`
	} `json:"table"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)
		input = body
	}
	log.Println(string(input))

	var payload tableSpec
	if err := json.Unmarshal(input, &payload); err != nil {
		log.Println(err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	namespaceSingle := "container"
	if val, ok := os.LookupEnv("namespace_single"); ok {
		namespaceSingle = val
	}
	namespaceDistributed := "dist"
	if val, ok := os.LookupEnv("namespace_dist"); ok {
		namespaceDistributed = val
	}

	namespace := ""
	if payload.Table.Name == "single_deployment" {
		namespace = namespaceSingle
	} else if payload.Table.Name == "distributed_deployment" {
		namespace = namespaceDistributed
	} else {
		log.Printf("table name %s not recognized", payload.Table.Name)
		http.Error(w, "invalid table name", http.StatusBadRequest)
		return
	}

	if err := deleteVirtualSvc(payload.Event.Data.Old.NameK8S, namespace); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Cannot delete VirtualService: %s", err)
		log.Println(errMsg)
		w.Write([]byte(errMsg))
		return
	}

	subject := "nats-delete-svc"
	if val, ok := os.LookupEnv("topic_delete_svc"); ok {
		subject = val
	}
	msg := message{
		Name:      payload.Event.Data.Old.NameK8S,
		Namespace: namespace,
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		http.Error(w, "Cannot marshal json", http.StatusInternalServerError)
		return
	}

	if err := nats.Send(subject, msgBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Cannot connect/publish to message queue: %s", err)
		w.Write([]byte(errMsg))
		return
	}

	w.WriteHeader(http.StatusOK)
}
