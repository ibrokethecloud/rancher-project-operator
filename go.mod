module github.com/ibrokethecloud/rancher-project-operator

go 1.13

replace (
	k8s.io/client-go => k8s.io/client-go v0.18.0
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.6.2
)

require (
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/brancz/gojsontoyaml v0.0.0-20190425155809-e8bd32d46b3d // indirect
	github.com/go-bindata/go-bindata v3.1.2+incompatible // indirect
	github.com/go-logr/logr v0.1.0
	github.com/json-iterator/go v1.1.10
	github.com/jsonnet-bundler/jsonnet-bundler v0.2.0 // indirect
	github.com/knative/pkg v0.0.0-20190817231834-12ee58e32cc8 // indirect
	github.com/mitchellh/hashstructure v0.0.0-20170609045927-2bca23e0e452 // indirect
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/openshift/prom-label-proxy v0.1.1-0.20191016113035-b8153a7f39f1 // indirect
	github.com/rancher/norman v0.0.0-20200609224801-7afd2e9bf37f
	github.com/rancher/types v0.0.0-20200609171948-b18f4c194419
	github.com/rancher/wrangler-api v0.6.1-0.20200515193802-dcf70881b087 // indirect
	github.com/terraform-providers/terraform-provider-rancher2 v1.9.0
	github.com/thanos-io/thanos v0.10.1 // indirect
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.5.0
	sigs.k8s.io/controller-tools v0.2.4 // indirect
)
