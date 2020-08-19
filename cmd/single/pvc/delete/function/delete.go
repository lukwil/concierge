package function

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/lukwil/concierge/cmd/common/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func deletePVC(name, namespace string) error {
	name = strings.TrimSpace(name)
	namespace = strings.TrimSpace(namespace)

	clientset, err := clientset.SetupInternal()
	if err != nil {
		return err
	}

	pvcClient := clientset.CoreV1().PersistentVolumeClaims(namespace)
	pvcList, err := pvcClient.List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%v", name),
	})
	if err != nil {
		return err
	}
	for _, item := range pvcList.Items {
		pvcName := item.Name
		log.Println("Deleting PVC...")
		if err := pvcClient.Delete(context.TODO(), pvcName, metav1.DeleteOptions{}); err != nil {
			return err
		}
		log.Println("Deleted PVC.")
	}

	return nil
}
