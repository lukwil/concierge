package function

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/lukwil/concierge/cmd/common/dynamic"
	"github.com/lukwil/concierge/cmd/common/hasura"
	"github.com/shurcooL/graphql"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var distributedDeploymentPayload struct {
	DistributedDeploymentByPk struct {
		ID graphql.Int `graphql:"id"`
	} `graphql:"update_distributed_deployment_by_pk(pk_columns: $pkColumns, _set: $set)"`
}

// Non-idiomatic Go naming, but needed by graphql library
type distributed_deployment_pk_columns_input struct {
	ID int `json:"id"`
}

// Non-idiomatic Go naming, but needed by graphql library
type distributed_deployment_set_input struct {
	URLPrefix string `json:"url_prefix"`
	NameK8s   string `json:"name_k8s"`
}

func createMPIJob(p *hasura.DistributedDeploymentPayload, namespace string) (string, string, error) {
	id := p.Event.Data.New.ID
	name := fmt.Sprintf("%v-%s", namespace, uuid.New())
	launcherName := fmt.Sprintf("%v-launcher", name)
	workerName := fmt.Sprintf("%v-worker", name)
	namespace = strings.TrimSpace(namespace)

	launcherCPU := p.Event.Data.New.LauncherCPU
	launcherRAM := p.Event.Data.New.LauncherRAM
	workerCPU := p.Event.Data.New.WorkerCPU
	workerRAM := p.Event.Data.New.WorkerRAM
	workerGPU := p.Event.Data.New.WorkerGPU
	workerCount := p.Event.Data.New.WorkerCount
	image := strings.TrimSpace(p.Event.Data.New.ContainerImage)

	urlPrefix := p.Event.Data.New.URLPrefix
	// If the user does not want to set his own URLPrefix but wants to use the name given to the container as URLPrefix
	// He cannot know the container name in advance (because of UUID), thus this workaround
	if urlPrefix == "name_k8s" {
		urlPrefix = fmt.Sprintf("/%v", name)
	}

	client, err := dynamic.SetupInternal()
	if err != nil {
		return "", "", err
	}

	mpiResource := schema.GroupVersionResource{Group: "kubeflow.org", Version: "v1alpha2", Resource: "mpijobs"}

	slotsPerWorker := int32(1)
	cleanPodPolicy := "Running"
	launcherReplicas := int32(1)
	launcherCPUStr := fmt.Sprintf("%vm", launcherCPU)
	launcherRAMStr := fmt.Sprintf("%vMi", launcherRAM)

	workerReplicas := int32(workerCount)
	workerCPUStr := fmt.Sprintf("%vm", workerCPU)
	workerRAMStr := fmt.Sprintf("%vMi", workerRAM)

	mpi := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "kubeflow.org/v1alpha2",
			"kind":       "MPIJob",
			"metadata": metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			"spec": map[string]interface{}{
				"slotsPerWorker": &slotsPerWorker,
				"cleanPodPolicy": &cleanPodPolicy,
				"mpiReplicaSpecs": map[string]interface{}{
					"Launcher": map[string]interface{}{
						"replicas": &launcherReplicas,
						"template": corev1.PodTemplateSpec{
							ObjectMeta: metav1.ObjectMeta{
								Labels: map[string]string{
									"app": launcherName,
								},
							},
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Name:            launcherName,
										Image:           image,
										ImagePullPolicy: corev1.PullIfNotPresent,
										Env: []corev1.EnvVar{
											{
												Name:  "URL_PREFIX",
												Value: urlPrefix,
											},
										},
										Resources: corev1.ResourceRequirements{
											Requests: corev1.ResourceList{
												corev1.ResourceCPU:    resource.MustParse(launcherCPUStr),
												corev1.ResourceMemory: resource.MustParse(launcherRAMStr),
											},
										},
										Ports: []corev1.ContainerPort{
											{
												Name:          "http",
												Protocol:      corev1.ProtocolTCP,
												ContainerPort: 8888,
											},
										},
									},
								},
							},
						},
					},
					"Worker": map[string]interface{}{
						"replicas": &workerReplicas,
						"template": corev1.PodTemplateSpec{
							ObjectMeta: metav1.ObjectMeta{
								Labels: map[string]string{
									"app": workerName,
								},
							},
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Name:            workerName,
										Image:           image,
										ImagePullPolicy: corev1.PullIfNotPresent,
										Resources: corev1.ResourceRequirements{
											Requests: corev1.ResourceList{
												corev1.ResourceCPU:    resource.MustParse(workerCPUStr),
												corev1.ResourceMemory: resource.MustParse(workerRAMStr),
											},
											Limits: corev1.ResourceList{
												"nvidia.com/gpu": resource.MustParse(strconv.Itoa(workerGPU)),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	log.Println("Creating MPIJob...")
	if _, err := client.Resource(mpiResource).Namespace(namespace).Create(context.TODO(), mpi, metav1.CreateOptions{}); err != nil {
		log.Printf("Cannot create MPIJob: %v", err)
		return "", "", err
	}
	log.Printf("Created MPIJob %q.\n", name)

	if err := updateTable(id, name, urlPrefix); err != nil {
		errMsg := fmt.Sprintf("Cannot update table distributed_deployment (name_k8s, url_prefix) in database: %s", err)
		log.Println(errMsg)
		return "", "", err
	}

	return name, urlPrefix, nil
}

func updateTable(id int, name, urlPrefix string) error {
	client := hasura.Client()

	pkColumns := distributed_deployment_pk_columns_input{
		ID: id,
	}
	set := distributed_deployment_set_input{
		URLPrefix: urlPrefix,
		NameK8s:   name,
	}
	variables := map[string]interface{}{
		"pkColumns": pkColumns,
		"set":       set,
	}

	if err := client.Mutate(context.TODO(), &distributedDeploymentPayload, variables); err != nil {
		return err
	}
	return nil
}
