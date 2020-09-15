package function

import (
	"context"
	"fmt"
	"log"

	"github.com/lukwil/concierge/cmd/common/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func listWarnings(namespace, name string) ([]containerWarningsOutput, error) {
	clientset, err := clientset.SetupInternal()
	if err != nil {
		return nil, err
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%v", name),
	})
	if err != nil {
		log.Printf("Pods cannot be listed: %v", err)
		return nil, err
	}
	warnings := []containerWarningsOutput{}
	for _, pod := range pods.Items {
		name := pod.Name
		w, err := warningsPerPod(clientset, namespace, name)
		if err != nil {
			return nil, err
		}
		warnings = append(warnings, w...)

	}
	return warnings, nil
}

func warningsPerPod(clientset *kubernetes.Clientset, namespace, podName string) ([]containerWarningsOutput, error) {
	events, err := clientset.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%v,involvedObject.kind=Pod", podName),
	})
	if err != nil {
		log.Printf("Failed to list Events: %v", err)
		return nil, err
	}
	warnings := []containerWarningsOutput{}
	for _, event := range events.Items {
		if event.Reason == "Failed" {
			msg := event.Message
			timestamp := event.LastTimestamp

			warning := containerWarningsOutput{
				Timestamp: timestamp.String(),
				Message:   msg,
			}
			warnings = append(warnings, warning)

			fmt.Printf("%v -- %v -- %v\n\n", timestamp.String(), event.InvolvedObject.Name, msg)
		}
	}
	return warnings, nil
}
