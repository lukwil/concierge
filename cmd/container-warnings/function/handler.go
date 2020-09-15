package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type actionPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            containerWarningsArgs  `json:"input"`
}

type graphQLError struct {
	Message string `json:"message"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		input = body
	}
	log.Println(string(input))

	var actionPayload actionPayload
	if err := json.Unmarshal(input, &actionPayload); err != nil {
		log.Println(err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	namespace := "container"
	if val, ok := os.LookupEnv("namespace"); ok {
		namespace = val
	}

	res, err := listWarnings(namespace, actionPayload.Input.NameK8s.NameK8s)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorObject := graphQLError{
			Message: err.Error(),
		}
		errorBody, _ := json.Marshal(errorObject)
		w.Write(errorBody)

		errMsg := fmt.Sprintf("Cannot retreive container warnings: %s", err)
		log.Println(errMsg)
		return
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		http.Error(w, "Cannot marshal json", http.StatusInternalServerError)
		return
	}
	w.Write(resBytes)
	w.WriteHeader(http.StatusOK)
}
