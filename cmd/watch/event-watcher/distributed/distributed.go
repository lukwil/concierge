package distributed

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/lukwil/concierge/cmd/common/dynamic"
	"github.com/lukwil/concierge/cmd/common/hasura"
	"github.com/shurcooL/graphql"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
	dynamicK8s "k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

const (
	created   graphql.String = "created"
	running   graphql.String = "running"
	succeeded graphql.String = "succeeded"
)

var statusIDName map[graphql.String]graphql.Int

var distributedDeploymentPayload struct {
	DistributedDeploymentByPk struct {
		AffectedRows graphql.Int `graphql:"affected_rows"`
	} `graphql:"update_distributed_deployment(where: $where, _set: $set)"`
}

var distributedDeploymentDeletePayload struct {
	DistributedDeploymentDelete struct {
		AffectedRows graphql.Int `graphql:"affected_rows"`
	} `graphql:"delete_distributed_deployment(where: $where)"`
}

// Non-idiomatic Go naming, but needed by graphql library
type distributed_deployment_set_input struct {
	StatusID int `json:"status_id"`
}

// Non-idiomatic Go naming, but needed by graphql library
type distributed_deployment_bool_exp struct {
	NameK8s struct {
		EQ string `json:"_eq"`
	} `json:"name_k8s"`
}

func onAdd(obj interface{}) {
	log.Println("ADD dist")
	objCast, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return
	}

	name, found, err := unstructured.NestedString(objCast.Object, "metadata", "name")
	if !found || err != nil {
		return
	}

	status, ok := getStatus(objCast)
	if !ok {
		return
	}

	log.Printf("ADDED --- %v: %v\n", name, status)
	if err := updateTable(name, status); err != nil {
		log.Println("Cannot update table")
		log.Println(err)
	}
}

func onUpdate(oldObj, newObj interface{}) {
	log.Println("UPDATE dist")
	newObjCast, ok := newObj.(*unstructured.Unstructured)
	if !ok {
		return
	}

	name, found, err := unstructured.NestedString(newObjCast.Object, "metadata", "name")
	if !found || err != nil {
		return
	}

	status, ok := getStatus(newObjCast)
	if !ok {
		return
	}

	log.Printf("UPDATED --- %v: %v\n", name, status)
	if err := updateTable(name, status); err != nil {
		log.Println("Cannot update table")
		log.Println(err)
	}
}

func onDelete(obj interface{}) {
	log.Println("DELETE dist")
	objCast, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return
	}

	name, found, err := unstructured.NestedString(objCast.Object, "metadata", "name")
	if !found || err != nil {
		return
	}

	log.Printf("DELETED --- %v\n", name)
	if err := deleteObjFromTable(name); err != nil {
		log.Println("Cannot delete object")
		log.Println(err)
	}
}

func Watch(sIDName map[graphql.String]graphql.Int) {
	log.Println("Watch main")
	statusIDName = sIDName

	namespace := "dist"
	if val, ok := os.LookupEnv("namespace_dist"); ok {
		namespace = val
	}

	client, err := dynamic.SetupInternal()
	if err != nil {
		log.Fatalf("Cannot create dynamic client: %v", err)
	}
	watch(client, namespace)
}

func watch(client dynamicK8s.Interface, namespace string) {
	log.Println("WATCH")
	log.Println(namespace)
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, 0, namespace, nil)
	mpiResource := schema.GroupVersionResource{Group: "kubeflow.org", Version: "v1alpha2", Resource: "mpijobs"}
	informer := factory.ForResource(mpiResource).Informer()

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

func getStatus(obj *unstructured.Unstructured) (graphql.String, bool) {
	conditions, found, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if !found || err != nil {
		return "", false
	}
	last := conditions[len(conditions)-1]

	lastMap, ok := last.(map[string]interface{})
	if !ok {
		return "", false
	}

	switch status := lastMap["type"]; status {
	case "Created":
		return created, true
	case "Running":
		return running, true
	case "Succeeded":
		return succeeded, true
	}
	return "", false
}

func updateTable(name string, status graphql.String) error {
	client := hasura.Client()

	id, found := statusIDName[status]
	if !found {
		return fmt.Errorf("Status %v not found", status)
	}

	set := distributed_deployment_set_input{
		StatusID: int(id),
	}
	eq := distributed_deployment_bool_exp{}
	eq.NameK8s.EQ = name

	variables := map[string]interface{}{
		"where": eq,
		"set":   set,
	}
	if err := client.Mutate(context.TODO(), &distributedDeploymentPayload, variables); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func deleteObjFromTable(name string) error {
	client := hasura.Client()

	eq := distributed_deployment_bool_exp{}
	eq.NameK8s.EQ = name

	variables := map[string]interface{}{
		"where": eq,
	}
	if err := client.Mutate(context.TODO(), &distributedDeploymentDeletePayload, variables); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
