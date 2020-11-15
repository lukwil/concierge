package common

import (
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio/pkg/madmin"
)

// BucketPayload is the payload refering to the minio_bucket table.
// Hasura/Postgres sends it if an event trigger is fired.
type BucketPayload struct {
	Event struct {
		SessionVariables struct {
			XHasuraRole string `json:"x-hasura-role"`
		} `json:"session_variables"`
		Op   string `json:"op"`
		Data struct {
			Old struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"old"`
			New struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"new"`
		} `json:"data"`
	} `json:"event"`
	CreatedAt    time.Time `json:"created_at"`
	ID           string    `json:"id"`
	DeliveryInfo struct {
		MaxRetries   int `json:"max_retries"`
		CurrentRetry int `json:"current_retry"`
	} `json:"delivery_info"`
	Trigger struct {
		Name string `json:"name"`
	} `json:"trigger"`
	Table struct {
		Schema string `json:"schema"`
		Name   string `json:"name"`
	} `json:"table"`
}

// MinioInstance represents a Minio Server object.
type MinioInstance struct {
	Endpoint  string
	AccessKey string
	Secret    string
}

// Connect connects to a given Minio instance.
func (mi *MinioInstance) Connect() (*minio.Client, error) {
	mc, err := minio.New(mi.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(mi.AccessKey, mi.Secret, ""),
		Secure: false,
	})
	if err != nil {
		log.Printf("Cannot connect to MinIO server: %v\n", err)
		return nil, err
	}
	return mc, nil
}

// ConnectAdmin connects to a given Minio instance in admin mode.
func (mi *MinioInstance) ConnectAdmin() (*madmin.AdminClient, error) {
	mac, err := madmin.New(mi.Endpoint, mi.AccessKey, mi.Secret, false)
	if err != nil {
		log.Printf("Cannot connect to MinIO server: %v\n", err)
		return nil, err
	}
	return mac, nil
}
