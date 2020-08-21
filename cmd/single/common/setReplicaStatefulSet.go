package common

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/lukwil/concierge/cmd/common/clientset"
)

// SetReplicaStatefulSet sets the replica count of a given StatefulSet.
func SetReplicaStatefulSet(replicas int32, name, namespace string) {
	clientset, err := clientset.SetupInternal()
	if err != nil {
		panic(err)
	}

	//createPVC(clientset, name, namespace)
	statefulSetClient := clientset.AppsV1().StatefulSets(namespace)
	result, err := statefulSetClient.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", err))
	}
	result.Spec.Replicas = &replicas
	if _, err := statefulSetClient.Update(context.TODO(), result, metav1.UpdateOptions{}); err != nil {
		panic(err)
	}
}
