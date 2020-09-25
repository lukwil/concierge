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

type payload struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)
		input = body
	}
	log.Println(string(input))

	var payload payload
	if err := json.Unmarshal(input, &payload); err != nil {
		log.Println(err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if err := deleteSvc(payload.Name, payload.Namespace); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Cannot delete Service: %s", err)
		log.Println(errMsg)
		w.Write([]byte(errMsg))
		return
	}

	subjectSingle := "nats-delete-statefulset"
	if val, ok := os.LookupEnv("topic_delete_statefulset"); ok {
		subjectSingle = val
	}
	subjectDistributed := "nats-delete-mpi"
	if val, ok := os.LookupEnv("topic_delete_mpi"); ok {
		subjectDistributed = val
	}

	namespaceSingle := "container"
	if val, ok := os.LookupEnv("namespace_single"); ok {
		namespaceSingle = val
	}
	namespaceDistributed := "dist"
	if val, ok := os.LookupEnv("namespace_dist"); ok {
		namespaceDistributed = val
	}

	subject := ""
	if payload.Namespace == namespaceSingle {
		subject = subjectSingle
	} else if payload.Namespace == namespaceDistributed {
		subject = subjectDistributed
	} else {
		log.Printf("namespace %s not recognized", payload.Namespace)
		http.Error(w, "invalid namespace", http.StatusBadRequest)
		return
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not marshal json", http.StatusInternalServerError)
		return
	}

	if err := nats.Send(subject, payloadBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Can not connect/publish to message queue: %s", err)
		w.Write([]byte(errMsg))
		return
	}

	w.WriteHeader(http.StatusOK)
}
