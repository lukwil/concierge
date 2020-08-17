package create

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/lukwil/concierge/cmd/common/clientset"

	"github.com/google/uuid"
	"github.com/shurcooL/graphql"
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
				ID             int    `json:"id"`
				Name           string `json:"name"`
				NameK8S        string `json:"name_k8s"`
				ContainerImage string `json:"container_image"`
				CPU            int    `json:"cpu"`
				RAM            int    `json:"ram"`
				GPU            int    `json:"gpu"`
				URLPrefix      string `json:"url_prefix"`
				StatusID       int    `json:"status_id"`
				VolumeID       int    `json:"volume_id"`
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

var volumePayload struct {
	VolumeByPk struct {
		MountPath graphql.String
		Size      graphql.Int
	} `graphql:"volume_by_pk(id: $id)"`
}

type volume struct {
	Size      int
	MountPath string
}

func (p *payload) createStatefulSet(namespace string) (string, error) {
	namespace = strings.TrimSpace(namespace)
	image := strings.TrimSpace(p.Event.Data.New.ContainerImage)
	urlPrefix := p.Event.Data.New.URLPrefix
	cpu := p.Event.Data.New.CPU
	ram := p.Event.Data.New.RAM
	volumeID := p.Event.Data.New.VolumeID

	clientset, err := clientset.SetupInternal()
	if err != nil {
		return "", err
	}

	replicas := int32(1)
	name := fmt.Sprintf("%v-%s", namespace, uuid.New())
	cpuStr := fmt.Sprintf("%vm", cpu)
	ramStr := fmt.Sprintf("%vMi", ram)
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
									Name:  "URL_PREFIX",
									Value: urlPrefix,
								},
							},
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
	if volumeID != 0 {
		vol, err := getVolumeDetails(volumeID)
		if err != nil {
			return "", err
		}

		storageStr := fmt.Sprintf("%vMi", vol.Size)

		container := &statefulSet.Spec.Template.Spec.Containers[0]
		container.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      name,
				MountPath: vol.MountPath,
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
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(storageStr),
						},
					},
				},
			},
		}
	}

	log.Println("Creating stateful set...")
	result, err := statefulSetClient.Create(context.TODO(), statefulSet, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	statefulSetName := result.GetObjectMeta().GetName()
	log.Printf("Created stateful set %q.\n", statefulSetName)

	return statefulSetName, nil
}

func getVolumeDetails(volumeID int) (*volume, error) {
	graphqlURL := "http://localhost:8080/hasura/v1/graphql"
	if val, ok := os.LookupEnv("graphql_url"); ok {
		graphqlURL = val
	}
	client := graphql.NewClient(graphqlURL, nil)

	variables := map[string]interface{}{
		"id": graphql.ID(volumeID),
	}

	if err := client.Query(context.TODO(), &volumePayload, variables); err != nil {
		return &volume{}, err
	}

	vol := &volume{
		Size:      int(volumePayload.VolumeByPk.Size),
		MountPath: string(volumePayload.VolumeByPk.MountPath),
	}
	return vol, nil
}
