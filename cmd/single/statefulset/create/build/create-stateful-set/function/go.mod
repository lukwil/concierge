module github.com/lukwil/concierge/cmd/single/statefulset/create/function

go 1.14

require (
	github.com/google/uuid v1.1.1
	github.com/lukwil/concierge/cmd/common v0.0.0-20200818101846-0b2b4491905d // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/shurcooL/graphql v0.0.0-20181231061246-d48a9a75455f
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
)

//replace github.com/lukwil/concierge/cmd/common => ../../../common
