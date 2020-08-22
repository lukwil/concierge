package function

import "github.com/lukwil/concierge/cmd/single/common"

func stopStatefulSet(name, namespace string) error {
	return common.SetReplicaStatefulSet(0, name, namespace)
}
