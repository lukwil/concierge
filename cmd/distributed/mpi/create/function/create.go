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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var distributedDeploymentPayload struct {
	DistributedDeploymentByPk struct {
		ID graphql.Int `graphql:"id"`
	} `graphql:"update_distributed_deployment_by_pk(pk_columns: $pkColumns, _set: $set)"`
}

var distributedEnvironmentVariablesPayload struct {
	DistributedEnvironmentVariables []struct {
		Name  graphql.String `graphql:"name"`
		Value graphql.String `graphql:"value"`
	} `graphql:"distributed_environment_variables(where: $where)"`
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

// Non-idiomatic Go naming, but needed by graphql library
type distributed_environment_variables_bool_exp struct {
	DistributedDeploymentID struct {
		EQ int `json:"_eq"`
	} `json:"distributed_deployment_id"`
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
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"slotsPerWorker": &slotsPerWorker,
				"cleanPodPolicy": &cleanPodPolicy,
				"mpiReplicaSpecs": map[string]interface{}{
					"Launcher": map[string]interface{}{
						"replicas": &launcherReplicas,
						"template": map[string]interface{}{
							"metadata": map[string]interface{}{
								"labels": map[string]interface{}{
									"app": launcherName,
								},
								"annotations": map[string]interface{}{
									"sidecar.istio.io/inject": "true",
								},
							},
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":            launcherName,
										"image":           image,
										"imagePullPolicy": "IfNotPresent",
										"resources": map[string]interface{}{
											"requests": map[string]interface{}{
												"cpu":    launcherCPUStr,
												"memory": launcherRAMStr,
											},
										},
										"ports": []interface{}{
											map[string]interface{}{
												"name":          "http",
												"protocol":      "TCP",
												"containerPort": int64(8888),
											},
										},
										"env": []interface{}{},
									},
								},
							},
						},
					},
					"Worker": map[string]interface{}{
						"replicas": &workerReplicas,
						"template": map[string]interface{}{
							"metadata": map[string]interface{}{
								"labels": map[string]interface{}{
									"app": workerName,
								},
							},
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":            workerName,
										"image":           image,
										"imagePullPolicy": "IfNotPresent",
										"resources": map[string]interface{}{
											"requests": map[string]interface{}{
												"cpu":    workerCPUStr,
												"memory": workerRAMStr,
											},
											"limits": map[string]interface{}{
												"nvidia.com/gpu": strconv.Itoa(workerGPU),
											},
										},
										"env": []interface{}{},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	vars, err := getEnvVariables(id)
	launcherVars := vars
	launcherVars = append(launcherVars, map[string]interface{}{"name": "URL_PREFIX", "value": urlPrefix})
	launcherVars = append(launcherVars, map[string]interface{}{"name": "IS_LAUNCHER", "value": "1"})
	workerVars := vars
	workerVars = append(workerVars, map[string]interface{}{"name": "IS_LAUNCHER", "value": "0"})

	launcherContainers, found, err := unstructured.NestedSlice(mpi.Object, "spec", "mpiReplicaSpecs", "Launcher", "template", "spec", "containers")
	if err != nil || !found || launcherContainers == nil {
		log.Printf("deployment launcher containers not found or error in spec: %v", err)
		return "", "", err
	}
	if err := unstructured.SetNestedSlice(launcherContainers[0].(map[string]interface{}), launcherVars, "env"); err != nil {
		log.Printf("environment variables could not be set for launcher: %v", err)
		return "", "", err
	}
	if err := unstructured.SetNestedField(mpi.Object, launcherContainers, "spec", "mpiReplicaSpecs", "Launcher", "template", "spec", "containers"); err != nil {
		log.Printf("deployment launcher containers not be set or error in spec: %v", err)
		return "", "", err
	}

	workerContainers, found, err := unstructured.NestedSlice(mpi.Object, "spec", "mpiReplicaSpecs", "Worker", "template", "spec", "containers")
	if err != nil || !found || workerContainers == nil {
		log.Printf("deployment worker containers not found or error in spec: %v", err)
		return "", "", err
	}
	if err := unstructured.SetNestedSlice(workerContainers[0].(map[string]interface{}), workerVars, "env"); err != nil {
		log.Printf("environment variables could not be set for workers: %v", err)
		return "", "", err
	}
	if err := unstructured.SetNestedField(mpi.Object, workerContainers, "spec", "mpiReplicaSpecs", "Worker", "template", "spec", "containers"); err != nil {
		log.Printf("deployment worker containers not be set or error in spec: %v", err)
		return "", "", err
	}

	log.Println(mpi)

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

func getEnvVariables(id int) ([]interface{}, error) {
	client := hasura.Client()

	eq := distributed_environment_variables_bool_exp{}
	eq.DistributedDeploymentID.EQ = id

	variables := map[string]interface{}{
		"where": eq,
	}
	if err := client.Query(context.TODO(), &distributedEnvironmentVariablesPayload, variables); err != nil {
		log.Println(err)
		return nil, err
	}

	var vars []interface{}
	for _, env := range distributedEnvironmentVariablesPayload.DistributedEnvironmentVariables {
		vars = append(vars, map[string]interface{}{"name": string(env.Name), "value": string(env.Value)})
	}
	return vars, nil
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
