package function

import (
	"context"
	"log"
	"strings"

	"github.com/lukwil/concierge/cmd/common/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func deleteSecret(name, namespace string) error {
	namespace = strings.TrimSpace(namespace)

	client, err := clientset.SetupInternal()
	if err != nil {
		return err
	}
	secretClient := client.CoreV1().Secrets(namespace)

	if err := secretClient.Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
		log.Printf("Secret %v could not be created in namespace %v: %v\n", name, namespace, err)
		return err
	}

	log.Printf("Successfully created secret %v in namespace %v.\n", name, namespace)
	return nil
}
