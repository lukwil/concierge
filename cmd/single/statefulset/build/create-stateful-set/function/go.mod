module github.com/lukwil/concierge/cmd/single/statefulset/create

go 1.14

require (
	github.com/google/uuid v1.1.1
	github.com/nats-io/nats-server/v2 v2.1.7 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/shurcooL/graphql v0.0.0-20181231061246-d48a9a75455f
	github.com/lukwil/concierge/cmd/common v0.0.0-20200817145652-ebeac7c626f5
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8 // indirect
	k8s.io/utils v0.0.0-20200815180417-3bc9d57fc792 // indirect
)

// replace github.com/lukwil/concierge/cmd/common => ../../../common
