package function

import (
	"context"
	"log"

	"github.com/lukwil/concierge/cmd/minio/common"
	"github.com/minio/minio-go/v7"
)

func createBucket(mi *common.MinioInstance, name string) error {
	client, err := mi.Connect()
	if err != nil {
		return err
	}

	if err := client.MakeBucket(context.Background(), name, minio.MakeBucketOptions{}); err != nil {
		log.Printf("MinIO bucket %v could not be created: %v\n", name, err)
		return err
	}

	log.Printf("Successfully created MinIO bucket %v.\n", name)
	return nil
}
