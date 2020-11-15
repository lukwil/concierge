package function

import (
	"context"
	"fmt"
	"log"

	"github.com/lukwil/concierge/cmd/minio/common"
)

func setPolicy(mi *common.MinioInstance, bucketName string) error {
	client, err := mi.ConnectAdmin()
	if err != nil {
		return err
	}

	policyName := fmt.Sprintf("%v-policy", bucketName)
	user := fmt.Sprintf("%v-user", bucketName)

	if err := client.SetPolicy(context.Background(), policyName, user, false); err != nil {
		log.Printf("MinIO policy %v could not be set for user %v: %v\n", policyName, user, err)
		return err
	}

	log.Printf("Successfully set MinIO policy %v for user %v.\n", policyName, user)
	return nil
}
