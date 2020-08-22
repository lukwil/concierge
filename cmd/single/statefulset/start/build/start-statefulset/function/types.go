package function

type output struct {
	Replicas int
}

type input struct {
	NameK8s string
}

type mutation struct {
	StartStatefulSet *output
}

type startStatefulSetArgs struct {
	NameK8s input
}
