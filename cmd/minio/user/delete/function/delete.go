package function

import (
	"context"
	"fmt"
	"log"

	"github.com/lukwil/concierge/cmd/minio/common"
)

func deleteUser(mi *common.MinioInstance, bucketName string) error {
	client, err := mi.ConnectAdmin()
	if err != nil {
		return err
	}

	user := fmt.Sprintf("%v-user", bucketName)

	if err := client.RemoveUser(context.Background(), user); err != nil {
		log.Printf("User %v could not be deleted from MinIO: %v\n", user, err)
		return err
	}

	log.Printf("Successfully deleted MinIO user %v.\n", user)
	return nil
}
