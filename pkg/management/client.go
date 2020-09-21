package management

import (
	"fmt"

	clusterClient "github.com/rancher/types/client/cluster/v3"
	managementClient "github.com/rancher/types/client/management/v3"
	"github.com/terraform-providers/terraform-provider-rancher2/rancher2"
	corev1 "k8s.io/api/core/v1"
)

type ApiWrapper struct {
	m         *managementClient.Client
	c         *clusterClient.Client
	clusterID string
}

// NewClient takes the k8s secret and generates a managementClient for operator
// to use for its operations
func NewClient(secret corev1.Secret) (a *ApiWrapper, err error) {
	config := rancher2.Config{}
	rancher_api_token, ok := secret.Data["rancher_api_token"]
	if !ok {
		return nil, fmt.Errorf("rancher_api_token key not found in secret")
	}

	config.TokenKey = string(rancher_api_token)

	rancher_url, ok := secret.Data["rancher_url"]
	if !ok {
		return nil, fmt.Errorf("rancher_url key not found in secret")
	}

	config.URL = string(rancher_url)

	cluster_id, ok := secret.Data["cluster_id"]
	if !ok {
		return nil, fmt.Errorf("cluster_id key not found in secret")
	}

	config.ClusterID = string(cluster_id)
	caCert, ok := secret.Data["cacert"]
	if ok {
		config.CACerts = string(caCert)
	}

	insecure, ok := secret.Data["insecure"]
	if ok {
		if string(insecure) == "true" {
			config.Insecure = true
		}
	}

	m, err := config.ManagementClient()
	if err != nil {
		return nil, err
	}

	c, err := config.ClusterClient(string(cluster_id))
	if err != nil {
		return nil, err
	}
	a = &ApiWrapper{m: m,
		c:         c,
		clusterID: string(cluster_id)}
	return a, nil
}
