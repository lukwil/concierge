module gitlab.com/masterarbeit-dl-cluster/concierge/cmd/single/statefulset/create

go 1.14

require (
	github.com/google/uuid v1.1.1
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/nats-io/nats-server/v2 v2.1.7 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/shurcooL/graphql v0.0.0-20181231061246-d48a9a75455f
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	k8s.io/utils v0.0.0-20200815180417-3bc9d57fc792 // indirect
)

//replace gitlab.com/masterarbeit-dl-cluster/concierge/cmd/single/statefulset/create/common => ./common
