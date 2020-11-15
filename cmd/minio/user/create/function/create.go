package function

import (
	"context"
	"fmt"
	"log"

	"github.com/lukwil/concierge/cmd/minio/common"
	"github.com/sethvargo/go-password/password"
)

func createUser(mi *common.MinioInstance, bucketName string) (string, string, error) {
	client, err := mi.ConnectAdmin()
	if err != nil {
		return "", "", err
	}

	user := fmt.Sprintf("%v-user", bucketName)
	pwd, err := password.Generate(30, 10, 10, false, false)
	if err != nil {
		log.Printf("Password could not be generated: %v\n", err)
		return "", "", err
	}

	if err := client.AddUser(context.Background(), user, pwd); err != nil {
		log.Printf("User %v could not be added to MinIO: %v\n", user, err)
		return "", "", err
	}

	log.Printf("Successfully created MinIO user %v.\n", user)
	return user, pwd, nil
}
