package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type payload struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
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

	if err := createVirtualSvc(payload.Name, payload.Namespace, payload.URLPrefix); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Cannot create VirtualService: %s", err)
		log.Println(errMsg)
		w.Write([]byte(errMsg))
		return
	}

	w.WriteHeader(http.StatusOK)
}
