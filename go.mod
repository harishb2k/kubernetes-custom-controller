module awesomeProject

go 1.13

require (
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9 // indirect
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	k8s.io/api v0.18.3 // indirect
	k8s.io/apimachinery v0.0.0-20200519081849-bdcc9f4ab675
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog/v2 v2.2.0
	k8s.io/utils v0.0.0-20200603063816-c1c6865ac451 // indirect
)

replace (
	golang.org/x/sys => golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a // pinned to release-branch.go1.13
	golang.org/x/tools => golang.org/x/tools v0.0.0-20190821162956-65e3620a7ae7 // pinned to release-branch.go1.13
	k8s.io/api => k8s.io/api v0.0.0-20200519082056-2543aba0e237
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20200519081849-bdcc9f4ab675
	k8s.io/client-go => k8s.io/client-go v0.0.0-20200519082352-455d6109ca5a
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20200519081644-3bc239a9bae4
)
