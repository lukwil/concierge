package function

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/lukwil/concierge/cmd/common/clientset"
	"github.com/lukwil/concierge/cmd/common/hasura"

	"github.com/google/uuid"
	"github.com/shurcooL/graphql"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var volumePayload struct {
	VolumeByPk struct {
		MountPath graphql.String `graphql:"mount_path"`
		Size      graphql.Int
	} `graphql:"volume_by_pk(id: $id)"`
}

type volume struct {
	Size      int
	MountPath string
}

var singleDeploymentPayload struct {
	SingleDeploymentByPk struct {
		ID graphql.Int `graphql:"id"`
	} `graphql:"update_single_deployment_by_pk(pk_columns: $pkColumns, _set: $set)"`
}

// Non-idiomatic Go naming, but needed by graphql library
type single_deployment_pk_columns_input struct {
	ID int `json:"id"`
}

// Non-idiomatic Go naming, but needed by graphql library
type single_deployment_set_input struct {
	URLPrefix string `json:"url_prefix"`
	NameK8s   string `json:"name_k8s"`
}

func createStatefulSet(p *hasura.SingleDeploymentPayload, namespace string) (statefulSetName string, urlPrefixName string, err error) {
	id := p.Event.Data.New.ID
	namespace = strings.TrimSpace(namespace)
	name := fmt.Sprintf("%v-%s", namespace, uuid.New())
	cpu := p.Event.Data.New.CPU
	ram := p.Event.Data.New.RAM
	volumeID := p.Event.Data.New.VolumeID
	image := strings.TrimSpace(p.Event.Data.New.ContainerImage)

	urlPrefix := p.Event.Data.New.URLPrefix
	// If the user does not want to set his own URLPrefix but wants to use the name given to the container as URLPrefix
	// He cannot know the container name in advance (because of UUID), thus this workaround
	if urlPrefix == "name_k8s" {
		urlPrefix = fmt.Sprintf("/%v", name)
	}

	clientset, err := clientset.SetupInternal()
	if err != nil {
		return "", "", err
	}

	replicas := int32(1)
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

	// create a volume if required by the user
	if volumeID != 0 {
		vol, err := getVolumeDetails(volumeID)
		if err != nil {
			errMsg := fmt.Sprintf("Cannot set volume in database: %s", err)
			log.Println(errMsg)
			return "", "", err
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

	log.Println("Creating StatefulSet...")
	if _, err := statefulSetClient.Create(context.TODO(), statefulSet, metav1.CreateOptions{}); err != nil {
		return "", "", err
	}
	log.Printf("Created StatefulSet %q.\n", name)

	if err := updateTable(id, name, urlPrefix); err != nil {
		errMsg := fmt.Sprintf("Cannot update table (name_k8s, url_prefix) in database: %s", err)
		log.Println(errMsg)
		return "", "", err
	}

	return name, urlPrefix, nil
}

func getVolumeDetails(volumeID int) (*volume, error) {
	client := hasura.Client()

	variables := map[string]interface{}{
		"id": graphql.Int(volumeID),
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

func updateTable(id int, name, urlPrefix string) error {
	client := hasura.Client()

	pkColumns := single_deployment_pk_columns_input{
		ID: id,
	}
	set := single_deployment_set_input{
		URLPrefix: urlPrefix,
		NameK8s:   name,
	}
	variables := map[string]interface{}{
		"pkColumns": pkColumns,
		"set":       set,
	}
	// variables := map[string]interface{}{
	// 	"pkColumns": pkColumns,
	// 	"set": map[string]interface{}{
	// 		"name_k8s":   name,
	// 		"url_prefix": urlPrefix,
	// 	},
	// }

	if err := client.Mutate(context.TODO(), &singleDeploymentPayload, variables); err != nil {
		return err
	}
	return nil
}
