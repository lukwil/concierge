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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func createMPIJob(p *hasura.SingleDeploymentPayload, namespace string) (string, string, error) {
	name := fmt.Sprintf("%v-%s", namespace, uuid.New())
	launcherName := fmt.Sprintf("%v-launcher", name)
	workerName := fmt.Sprintf("%v-worker", name)
	namespace = strings.TrimSpace(namespace)

	cpu := p.Event.Data.New.CPU
	ram := p.Event.Data.New.RAM
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
	cpuStrLauncher := fmt.Sprintf("%vm", cpu)
	ramStrLauncher := fmt.Sprintf("%vMi", ram)

	workerReplicas := int32(2)
	cpuStrWorker := fmt.Sprintf("%vm", cpu)
	ramStrWorker := fmt.Sprintf("%vMi", ram)

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
										Command: []string{"mpirun"},
										Args: []string{
											"-np",
											strconv.Itoa(int(workerReplicas)),
											"--allow-run-as-root",
											"-bind-to",
											"none",
											"-map-by",
											"slot",
											"-x",
											"LD_LIBRARY_PATH",
											"-x",
											"PATH",
											"-mca",
											"pml",
											"ob1",
											"-mca",
											"btl",
											"^openib",
											"python",
											"/examples/tensorflow2_mnist.py",
										},
										Resources: corev1.ResourceRequirements{
											Requests: corev1.ResourceList{
												corev1.ResourceCPU:    resource.MustParse(cpuStrLauncher),
												corev1.ResourceMemory: resource.MustParse(ramStrLauncher),
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
												corev1.ResourceCPU:    resource.MustParse(cpuStrWorker),
												corev1.ResourceMemory: resource.MustParse(ramStrWorker),
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
	return name, urlPrefix, nil
}
