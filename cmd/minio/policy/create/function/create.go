package function

import (
	"context"
	"fmt"
	"log"

	"github.com/lukwil/concierge/cmd/minio/common"
	"github.com/minio/minio/pkg/bucket/policy"
	"github.com/minio/minio/pkg/bucket/policy/condition"
	iampolicy "github.com/minio/minio/pkg/iam/policy"
)

func createPolicy(mi *common.MinioInstance, bucketName string) error {
	client, err := mi.ConnectAdmin()
	if err != nil {
		return err
	}

	policyName := fmt.Sprintf("%v-policy", bucketName)

	p := iampolicy.Policy{
		Version: iampolicy.DefaultVersion,
		Statements: []iampolicy.Statement{
			iampolicy.NewStatement(
				policy.Allow,
				iampolicy.NewActionSet(iampolicy.GetObjectAction),
				iampolicy.NewResourceSet(iampolicy.NewResource(fmt.Sprintf("%v/*", bucketName), "")),
				condition.NewFunctions(),
			)},
	}

	if err := client.AddCannedPolicy(context.Background(), policyName, &p); err != nil {
		log.Printf("MinIO policy for bucket %v could not be created: %v\n", bucketName, err)
		return err
	}

	log.Printf("Successfully created MinIO policy for bucket %v.\n", bucketName)
	return nil
}
