package dynamic

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// SetupInternal returns a new dynamic Client for internal cluster use.
func SetupInternal() (dynamic.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
