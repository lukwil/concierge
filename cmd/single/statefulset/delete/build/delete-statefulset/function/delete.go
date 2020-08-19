package function

import (
	"context"
	"log"
	"strings"

	"github.com/lukwil/concierge/cmd/common/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func deleteStatefulSet(name, namespace string) error {
	name = strings.TrimSpace(name)
	namespace = strings.TrimSpace(namespace)

	clientset, err := clientset.SetupInternal()
	if err != nil {
		return err
	}

	ssClient := clientset.AppsV1().StatefulSets(namespace)
	log.Println("Deleting StatefulSet...")
	if err := ssClient.Delete(context.TODO(), name, metav1.DeleteOptions{}); err != nil {
		return err
	}
	log.Println("Deleted StatefulSet.")

	return nil
}
