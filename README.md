## Rancher Project Operator

Manage Rancher projects, namespace mapping and RBAC from downstream clusters.

The operator allows users to manage their RBAC via Rancher using CRD's deployed to downstream clusters.

This allows all rbac logic to be packaged in the workload helm chart itself.

The operator currently has the following types:

1. Project: Used for creating the Rancher project:
```
apiVersion: management.cattle.io/v1alpha1
kind: Project
metadata:
  name: project-sample
spec:
  # Add fields here
  displayName: demo
```

2. ProjectNamespaceMapping: Used for assigning Namespaces to a project. 
```
apiVersion: management.cattle.io/v1alpha1
kind: ProjectNamespaceMapping
metadata:
  name: demo-ns-mapping
spec:
  projectName: demo
  namespaces:
    - demo1
    - demo2
```

*NOTE* Deleting a project will delete all associated namespaces. 
If the idea is to only remove the namespace from an existing project
then just remove it from the project namespace mapping

3. ProjectRoleTemplateBinding: Used for role binding at the project layer.
```
apiVersion: management.cattle.io/v1alpha1
kind: ProjectRoleTemplateBinding
metadata:
  name: demo-project-sample-prtb
spec:
  userName: test
  projectName: demo
  roleTemplateName: project-member
```

The operator by default looks for a secret `rancher-secret` which needs to contain the following items:
* rancher_url: Rancher endpoint url.
* rancher_api_token: Rancher global token to access the api.
* cluster_id: cluster id of downstream cluster where the operator is deployed.
* cacert: Optional ca cert to allow operator to talk to rancher endpoint url.
* insecure: Optional to allow operator to make insecure calls to the endpoint url