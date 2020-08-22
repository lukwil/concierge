package function

type output struct {
	Replicas int
}

type input struct {
	NameK8s string
}

type mutation struct {
	StopStatefulSet *output
}

type stopStatefulSetArgs struct {
	NameK8s input
}
