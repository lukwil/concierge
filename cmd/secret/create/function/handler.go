package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type payload struct {
	Name    string `json:"name"`
	Secrets []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"secrets"`
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

	namespaceSingle := "container"
	if val, ok := os.LookupEnv("namespace_single"); ok {
		namespaceSingle = val
	}
	namespaceDistributed := "dist"
	if val, ok := os.LookupEnv("namespace_dist"); ok {
		namespaceDistributed = val
	}

	name := payload.Name
	secrets := make(map[string][]byte)
	for _, s := range payload.Secrets {
		secrets[s.Key] = []byte(s.Value)
	}

	if err := createSecret(name, namespaceSingle, secrets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Cannot create secret %v for namespace %v: %v", name, namespaceSingle, err)
		log.Println(errMsg)
		w.Write([]byte(errMsg))
		return
	}

	if err := createSecret(name, namespaceDistributed, secrets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Cannot create secret %v for namespace %v: %v", name, namespaceDistributed, err)
		log.Println(errMsg)
		w.Write([]byte(errMsg))
		return
	}
	w.WriteHeader(http.StatusOK)
}
