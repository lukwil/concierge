package create

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type payload struct {
	Event struct {
		SessionVariables struct {
			XHasuraRole string `json:"x-hasura-role"`
		} `json:"session_variables"`
		Op   string `json:"op"`
		Data struct {
			Old interface{} `json:"old"`
			New struct {
				VolumeID int         `json:"volume_id"`
				Name     string      `json:"name"`
				Gpu      int         `json:"gpu"`
				ID       int         `json:"id"`
				NameK8S  interface{} `json:"name_k8s"`
				RAM      int         `json:"ram"`
				StatusID interface{} `json:"status_id"`
				CPU      int         `json:"cpu"`
			} `json:"new"`
		} `json:"data"`
	} `json:"event"`
	CreatedAt    time.Time `json:"created_at"`
	ID           string    `json:"id"`
	DeliveryInfo struct {
		MaxRetries   int `json:"max_retries"`
		CurrentRetry int `json:"current_retry"`
	} `json:"delivery_info"`
	Trigger struct {
		Name string `json:"name"`
	} `json:"trigger"`
	Table struct {
		Schema string `json:"schema"`
		Name   string `json:"name"`
	} `json:"table"`
}

func (p *payload) createStatefulSet(namespace string) string {
	namespace = strings.TrimSpace(namespace)
	image := strings.TrimSpace(p.Input.Arg1.ContainerImage)
	cpu := p.Input.Arg1.CPU
	ram := p.Input.Arg1.RAM
	//gpu := p.Input.Arg1.GPU
	volumeSize := p.Input.Arg1.VolumeSize
	volumeMountPath := p.Input.Arg1.VolumeMountPath
	fmt.Println("VOLUMEMOUNTPATH")
	fmt.Println(volumeMountPath)

	clientset := clientset.setupInternal()

	replicas := int32(1)
	name := fmt.Sprintf("%v-%s", namespace, uuid.New())
	cpuStr := fmt.Sprintf("%vm", cpu)
	ramStr := fmt.Sprintf("%vMi", ram)
	//storageClassName := "local-path"
	storageStr := fmt.Sprintf("%vMi", volumeSize)
	fsGroup := int64(100)

	statefulSetClient := clientset.AppsV1().StatefulSets(namespace)

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            name,
							Image:           image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Env: []corev1.EnvVar{
								{
									Name:  "NB_PREFIX",
									Value: "/", //fmt.Sprintf("/%s", name),
								},
							},
							Command: []string{"sh", "-c", "jupyter notebook --notebook-dir=/home/jovyan --ip=0.0.0.0 --no-browser --allow-root --port=8888 --NotebookApp.token='' --NotebookApp.password='' --NotebookApp.allow_origin='*' --NotebookApp.base_url=${NB_PREFIX}"},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(cpuStr),
									corev1.ResourceMemory: resource.MustParse(ramStr),
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
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup: &fsGroup,
					},
				},
			},
		},
	}

	// create a volume if required by user
	if volumeSize != 0 {
		container := &statefulSet.Spec.Template.Spec.Containers[0]
		container.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      name,
				MountPath: volumeMountPath,
			},
		}
		statefulSet.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{
						corev1.ReadWriteOnce,
					},
					//StorageClassName: &storageClassName,
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(storageStr),
						},
					},
				},
			},
		}
	}
	fmt.Println(statefulSet)

	fmt.Println("Creating stateful set...")
	result, err := statefulSetClient.Create(context.TODO(), statefulSet, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	statefulSetName := result.GetObjectMeta().GetName()
	fmt.Printf("Created stateful set %q.\n", statefulSetName)

	return statefulSetName
}
