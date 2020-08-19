package function

import (
	"context"
	"log"
	"strings"

	"github.com/lukwil/concierge/cmd/common/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func deleteSvc(name, namespace string) error {
	name = strings.TrimSpace(name)
	namespace = strings.TrimSpace(namespace)

	clientset, err := clientset.SetupInternal()
	if err != nil {
		return err
	}

	serviceClient := clientset.CoreV1().Services(namespace)
	log.Println("Deleting Service...")
	if err := serviceClient.Delete(context.TODO(), name, metav1.DeleteOptions{}); err != nil {
		return err
	}
	log.Println("Deleted Service.")

	return nil
}
