package function

import (
	"context"
	"log"
	"path/filepath"
	"strings"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func deleteVirtualSvc(name, namespace string) error {
	name = strings.TrimSpace(name)
	namespace = strings.TrimSpace(namespace)

	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	clientset, err := versionedclient.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create istio client: %s", err)
		return err
	}

	virtualServicesClient := clientset.NetworkingV1alpha3().VirtualServices(namespace)
	log.Println("Deleting VirtualService...")
	if err := virtualServicesClient.Delete(context.TODO(), name, metav1.DeleteOptions{}); err != nil {
		return err
	}
	log.Println("Deleted VirtualService.")

	return nil
}
