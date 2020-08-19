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
	URLPrefix string `json:"url_prefix"`
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
	namespace := "container"
	if val, ok := os.LookupEnv("namespace"); ok {
		namespace = val
	}
	if err := createSvc(payload.Name, namespace); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Cannot create Service: %s", err)
		log.Println(errMsg)
		w.Write([]byte(errMsg))
		return
	}

	subject := "nats-create-virtualsvc"
	if val, ok := os.LookupEnv("topic_create_virtualsvc"); ok {
		subject = val
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		http.Error(w, "Cannot marshal json", http.StatusInternalServerError)
		return
	}

	if err := nats.Send(subject, payloadBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Cannot connect/publish to message queue: %s", err)
		w.Write([]byte(errMsg))
		return
	}

	w.WriteHeader(http.StatusOK)
}
