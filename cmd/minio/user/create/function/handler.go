package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/lukwil/concierge/cmd/common/nats"
	"github.com/lukwil/concierge/cmd/minio/common"
)

type secretMessage struct {
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

	var payload common.BucketPayload
	if err := json.Unmarshal(input, &payload); err != nil {
		log.Println(err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	name := payload.Event.Data.New.Name
	endpoint := "127.0.0.1:9000"
	if val, ok := os.LookupEnv("minio_url"); ok {
		endpoint = val
	}
	accessKey := "AKIAIOSFODNN7EXAMPLE"
	if val, ok := os.LookupEnv("minio_access_key"); ok {
		accessKey = val
	}
	secret := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	if val, ok := os.LookupEnv("minio_secret"); ok {
		secret = val
	}

	minioInstance := &common.MinioInstance{
		Endpoint:  endpoint,
		AccessKey: accessKey,
		Secret:    secret,
	}

	user, pwd, err := createUser(minioInstance, name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Cannot create MinIO user: %s", err)
		log.Println(errMsg)
		w.Write([]byte(errMsg))
		return
	}

	subjectPolicy := "nats-set-minio-policy"
	if val, ok := os.LookupEnv("topic_set_minio_policy"); ok {
		subjectPolicy = val
	}

	if err := nats.Send(subjectPolicy, input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Cannot connect/publish to message queue: %s", err)
		log.Println(errMsg)
		w.Write([]byte(errMsg))
		return
	}

	subjectSecret := "nats-create-secret"
	if val, ok := os.LookupEnv("topic_create_secret"); ok {
		subjectSecret = val
	}

	msg := secretMessage{}
	msg.Name = name
	msg.Secrets[0].Key = "access_key"
	msg.Secrets[0].Value = user
	msg.Secrets[1].Key = "secret"
	msg.Secrets[1].Value = pwd

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		http.Error(w, "Cannot marshal json", http.StatusInternalServerError)
		return
	}

	if err := nats.Send(subjectSecret, msgBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Cannot connect/publish to message queue: %s", err)
		log.Println(errMsg)
		w.Write([]byte(errMsg))
		return
	}

	w.WriteHeader(http.StatusOK)
}
