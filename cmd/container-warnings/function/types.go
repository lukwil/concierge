package function

type containerWarningsOutput struct {
	Timestamp string
	Reason    string
	Message   string
}

type containerWarningsInput struct {
	NameK8s string
}

type query struct {
	ContainerWarnings []*containerWarningsOutput
}

type containerWarningsArgs struct {
	NameK8s containerWarningsInput
}
