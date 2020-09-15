package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/shurcooL/graphql"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/lukwil/concierge/cmd/common/clientset"
	"github.com/lukwil/concierge/cmd/common/hasura"
)

const (
	StatefulSetStarting  graphql.String = "starting"
	StatefulSetStarted   graphql.String = "started"
	StatefulSetStopping  graphql.String = "stopping"
	StatefulSetStopped   graphql.String = "stopped"
	StatefulSetScheduled graphql.String = "scheduled"
)

var status struct {
	Status []struct {
		ID   graphql.Int
		Name graphql.String
	}
}

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

func main() {
	clientset, err := clientset.SetupInternal()
	if err != nil {
		log.Fatal(err)
	}

	namespace := "container"
	if val, ok := os.LookupEnv("namespace"); ok {
		namespace = val
	}

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", 8080),
		ReadTimeout:    3 * time.Second,
		WriteTimeout:   3 * time.Second,
		MaxHeaderBytes: 1 << 20, // Max header of 1MB
	}

	http.HandleFunc("/", okHandler)
	go server.ListenAndServe()

	getStatusTypes()
	watch(clientset, namespace)
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func setup() *kubernetes.Clientset {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientset
}

func getStatusTypes() {
	statusIDName = make(map[graphql.String]graphql.Int)
	client := hasura.Client()
	if err := client.Query(context.TODO(), &status, make(map[string]interface{})); err != nil {
		panic(err)
	}
	for _, s := range status.Status {
		statusIDName[s.Name] = s.ID
	}
	fmt.Println(statusIDName[StatefulSetStarting])
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

func onAdd(obj interface{}) {
	objCast := obj.(*appsv1.StatefulSet)
	name := objCast.Name
	status := getStatus(objCast)

	log.Printf("ADDED --- %v: %v\n", name, status)
	if err := updateTable(name, status); err != nil {
		log.Println("Cannot update table")
		log.Println(err)
	}
}

func onUpdate(oldObj, newObj interface{}) {
	newObjCast := newObj.(*appsv1.StatefulSet)
	name := newObjCast.Name
	status := getStatus(newObjCast)

	log.Printf("UPDATED --- %v: %v\n", name, status)
	if err := updateTable(name, status); err != nil {
		log.Println("Cannot update table")
		log.Println(err)
	}
}

func onDelete(obj interface{}) {
	objCast := obj.(*appsv1.StatefulSet)
	name := objCast.Name

	log.Printf("DELETED --- %v\n", name)
	if err := deleteObjFromTable(name); err != nil {
		log.Println("Cannot delete object")
		log.Println(err)
	}
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
			status = StatefulSetStopped
		} else {
			status = StatefulSetStopping
		}
	} else {
		if current == 1 {
			if ready == 1 {
				status = StatefulSetStarted
			} else {
				status = StatefulSetStarting
			}
		} else {
			status = StatefulSetScheduled
		}
	}
	return status
}

func updateTable(name string, status graphql.String) error {
	client := hasura.Client()

	id := statusIDName[status]

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
