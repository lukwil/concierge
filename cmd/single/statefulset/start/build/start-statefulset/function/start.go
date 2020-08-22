package function

import "github.com/lukwil/concierge/cmd/single/common"

func startStatefulSet(name, namespace string) error {
	return common.SetReplicaStatefulSet(1, name, namespace)
}
