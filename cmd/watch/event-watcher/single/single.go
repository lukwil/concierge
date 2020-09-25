package single

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/lukwil/concierge/cmd/common/clientset"
	"github.com/lukwil/concierge/cmd/common/hasura"
	"github.com/shurcooL/graphql"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const (
	starting  graphql.String = "starting"
	started   graphql.String = "started"
	stopping  graphql.String = "stopping"
	stopped   graphql.String = "stopped"
	scheduled graphql.String = "scheduled"
)

var statusIDName map[graphql.String]graphql.Int

var singleDeploymentPayload struct {
	SingleDeploymentByPk struct {
		AffectedRows graphql.Int `graphql:"affected_rows"`
	} `graphql:"update_single_deployment(where: $where, _set: $set)"`
}

var singleDeploymentDeletePayload struct {
	SingleDeploymentDelete struct {
		AffectedRows graphql.Int `graphql:"affected_rows"`
	} `graphql:"delete_single_deployment(where: $where)"`
}

// Non-idiomatic Go naming, but needed by graphql library
type single_deployment_set_input struct {
	StatusID int `json:"status_id"`
}

// Non-idiomatic Go naming, but needed by graphql library
type single_deployment_bool_exp struct {
	NameK8s struct {
		EQ string `json:"_eq"`
	} `json:"name_k8s"`
}

func onAdd(obj interface{}) {
	objCast, ok := obj.(*appsv1.StatefulSet)
	if !ok {
		return
	}
	name := objCast.Name
	status := getStatus(objCast)

	log.Printf("ADDED --- %v: %v\n", name, status)
	if err := updateTable(name, status); err != nil {
		log.Println("Cannot update table")
		log.Println(err)
	}
}

func onUpdate(oldObj, newObj interface{}) {
	newObjCast, ok := newObj.(*appsv1.StatefulSet)
	if !ok {
		return
	}
	name := newObjCast.Name
	status := getStatus(newObjCast)

	log.Printf("UPDATED --- %v: %v\n", name, status)
	if err := updateTable(name, status); err != nil {
		log.Println("Cannot update table")
		log.Println(err)
	}
}

func onDelete(obj interface{}) {
	objCast, ok := obj.(*appsv1.StatefulSet)
	if !ok {
		return
	}
	name := objCast.Name

	log.Printf("DELETED --- %v\n", name)
	if err := deleteObjFromTable(name); err != nil {
		log.Println("Cannot delete object")
		log.Println(err)
	}
}

func Watch(sIDName map[graphql.String]graphql.Int) {
	statusIDName = sIDName

	namespace := "container"
	if val, ok := os.LookupEnv("namespace_single"); ok {
		namespace = val
	}

	clientset, err := clientset.SetupInternal()
	if err != nil {
		log.Fatalf("Cannot create clientset: %v", err)
	}
	watch(clientset, namespace)
}

func watch(clientset *kubernetes.Clientset, namespace string) {
	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset, 0, informers.WithNamespace(namespace))
	informer := factory.Apps().V1().StatefulSets().Informer()

	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: onDelete,
	})
	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

func getStatus(statefulSet *appsv1.StatefulSet) graphql.String {
	desired := *statefulSet.Spec.Replicas
	current := statefulSet.Status.Replicas
	ready := statefulSet.Status.ReadyReplicas
	name := statefulSet.Name

	log.Printf("%v: %v desired - %v current - %v ready\n", name, desired, current, ready)

	var status graphql.String
	if desired == 0 {
		if current == 0 {
			status = stopped
		} else {
			status = stopping
		}
	} else {
		if current == 1 {
			if ready == 1 {
				status = started
			} else {
				status = starting
			}
		} else {
			status = scheduled
		}
	}
	return status
}

func updateTable(name string, status graphql.String) error {
	client := hasura.Client()

	id, found := statusIDName[status]
	if !found {
		return fmt.Errorf("Status %v not found", status)
	}

	set := single_deployment_set_input{
		StatusID: int(id),
	}
	eq := single_deployment_bool_exp{}
	eq.NameK8s.EQ = name

	variables := map[string]interface{}{
		"where": eq,
		"set":   set,
	}
	if err := client.Mutate(context.TODO(), &singleDeploymentPayload, variables); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func deleteObjFromTable(name string) error {
	client := hasura.Client()

	eq := single_deployment_bool_exp{}
	eq.NameK8s.EQ = name

	variables := map[string]interface{}{
		"where": eq,
	}
	if err := client.Mutate(context.TODO(), &singleDeploymentDeletePayload, variables); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
