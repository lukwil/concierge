package function

import (
	"context"
	"log"
	"strings"

	"github.com/lukwil/concierge/cmd/common/clientset"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createService(name, namespace string) error {
	name = strings.TrimSpace(name)
	namespace = strings.TrimSpace(namespace)

	clientset, err := clientset.SetupInternal()
	if err != nil {
		return err
	}

	serviceClient := clientset.CoreV1().Services(namespace)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 8888,
				},
			},
			Selector: map[string]string{
				"app": name,
			},
		},
	}

	log.Println("Creating service...")
	result, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	serviceName := result.GetObjectMeta().GetName()
	log.Printf("Created service %q.\n", serviceName)

	return nil
}
