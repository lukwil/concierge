package function

import (
	"context"
	"log"
	"strings"

	"github.com/lukwil/concierge/cmd/common/dynamic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func deleteMPIJob(name, namespace string) error {
	name = strings.TrimSpace(name)
	namespace = strings.TrimSpace(namespace)

	client, err := dynamic.SetupInternal()
	if err != nil {
		return err
	}

	mpiResource := schema.GroupVersionResource{Group: "kubeflow.org", Version: "v1alpha2", Resource: "mpijobs"}

	log.Println("Deleting MPIJob...")
	if err := client.Resource(mpiResource).Namespace(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{}); err != nil {
		log.Printf("Cannot delete MPIJob: %v", err)
		return err
	}
	log.Printf("Deleted MPIJob %q.\n", name)

	return nil
}
