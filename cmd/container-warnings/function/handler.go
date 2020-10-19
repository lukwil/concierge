package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type actionPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            containerWarningsArgs  `json:"input"`
}

type graphQLError struct {
	Message string `json:"message"`
}

type selector string

const distributedSelector selector = "mpi_job_name"
const singleSelector selector = "app"

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
	name := actionPayload.Input.NameK8s.NameK8s
	var namespace string
	var selectorType selector

	singlePrefix := "container"
	if val, ok := os.LookupEnv("single_prefix"); ok {
		singlePrefix = val
	}
	distributedPrefix := "dist"
	if val, ok := os.LookupEnv("dist_prefix"); ok {
		distributedPrefix = val
	}

	if strings.HasPrefix(name, singlePrefix) {
		namespace = singlePrefix
		selectorType = singleSelector
	} else if strings.HasPrefix(name, distributedPrefix) {
		namespace = distributedPrefix
		selectorType = distributedSelector
	} else {
		http.Error(w, "invalid name (prefix did not match)", http.StatusBadRequest)
		return
	}

	res, err := listWarnings(namespace, name, selectorType)
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
}
