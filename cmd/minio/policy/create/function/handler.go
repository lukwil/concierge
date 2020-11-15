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

	if err := createPolicy(minioInstance, name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Cannot create MinIO policy for bucket: %s", err)
		log.Println(errMsg)
		w.Write([]byte(errMsg))
		return
	}

	subject := "nats-create-minio-user"
	if val, ok := os.LookupEnv("topic_create_minio_user"); ok {
		subject = val
	}

	if err := nats.Send(subject, input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Cannot connect/publish to message queue: %s", err)
		log.Println(errMsg)
		w.Write([]byte(errMsg))
		return
	}
	w.WriteHeader(http.StatusOK)
}
