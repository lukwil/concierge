package function

import (
	"context"
	"fmt"
	"log"
	"strings"

	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func createVirtualSvc(name, namespace, urlPrefix string) error {
	name = strings.TrimSpace(name)
	namespace = strings.TrimSpace(namespace)

	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	ic, err := versionedclient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create istio client: %s", err)
		return err
	}

	virtualServicesClient := ic.NetworkingV1alpha3().VirtualServices(namespace)

	vs := &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: networkingv1alpha3.VirtualService{
			Hosts:    []string{"*"},
			Gateways: []string{"default/gateway"},
			Http: []*networkingv1alpha3.HTTPRoute{
				{
					Match: []*networkingv1alpha3.HTTPMatchRequest{
						{
							Uri: &networkingv1alpha3.StringMatch{
								MatchType: &networkingv1alpha3.StringMatch_Prefix{
									Prefix: fmt.Sprintf("/%s/", name),
								},
							},
						},
					},
					Rewrite: &networkingv1alpha3.HTTPRewrite{
						Uri: fmt.Sprintf("%s/", urlPrefix),
					},
					Route: []*networkingv1alpha3.HTTPRouteDestination{
						{
							Destination: &networkingv1alpha3.Destination{
								Port: &networkingv1alpha3.PortSelector{
									Number: 8888,
								},
								Host: name,
							},
						},
					},
				},
			},
		},
	}

	fmt.Println("Creating VirtualService...")
	result, err := virtualServicesClient.Create(context.TODO(), vs, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Created VirtualService %q.\n", result.GetObjectMeta().GetName())
	return nil
}
