package function

import (
	"context"
	"fmt"
	"log"

	"github.com/lukwil/concierge/cmd/minio/common"
)

func deletePolicy(mi *common.MinioInstance, bucketName string) error {
	client, err := mi.ConnectAdmin()
	if err != nil {
		return err
	}

	policyName := fmt.Sprintf("%v-policy", bucketName)

	if err := client.RemoveCannedPolicy(context.Background(), policyName); err != nil {
		log.Printf("MinIO policy for bucket %v could not be deleted: %v\n", bucketName, err)
		return err
	}

	log.Printf("Successfully deleted MinIO policy for bucket %v.\n", bucketName)
	return nil
}
