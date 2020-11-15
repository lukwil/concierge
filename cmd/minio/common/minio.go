package common

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio/pkg/madmin"
)

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
