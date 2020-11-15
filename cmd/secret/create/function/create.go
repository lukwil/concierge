package function

import (
	"context"
	"log"
	"strings"

	"github.com/lukwil/concierge/cmd/common/clientset"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createSecret(name, namespace string, secretData map[string][]byte) error {
	namespace = strings.TrimSpace(namespace)

	client, err := clientset.SetupInternal()
	if err != nil {
		return err
	}

	secretClient := client.CoreV1().Secrets(namespace)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: corev1.SecretTypeOpaque,
		Data: secretData,
	}

	if _, err := secretClient.Create(context.Background(), secret, metav1.CreateOptions{}); err != nil {
		log.Printf("Secret %v could not be created in namespace %v: %v\n", name, namespace, err)
		return err
	}

	log.Printf("Successfully created secret %v in namespace %v.\n", name, namespace)
	return nil
}
