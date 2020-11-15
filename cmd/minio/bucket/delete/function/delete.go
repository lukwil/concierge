package function

import (
	"context"
	"log"

	"github.com/lukwil/concierge/cmd/minio/common"
)

func deleteBucket(mi *common.MinioInstance, name string) error {
	client, err := mi.Connect()
	if err != nil {
		return err
	}

	if err := client.RemoveBucket(context.Background(), name); err != nil {
		log.Printf("MinIO bucket %v could not be deleted: %v\n", name, err)
		return err
	}

	log.Printf("Successfully deleted MinIO bucket %v.\n", name)
	return nil
}
