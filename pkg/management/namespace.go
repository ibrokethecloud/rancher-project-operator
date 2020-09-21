package management

import (
	"fmt"

	managementv1alpha1 "github.com/ibrokethecloud/rancher-project-operator/pkg/api/v1alpha1"

	"github.com/rancher/norman/types"
	clusterClient "github.com/rancher/types/client/cluster/v3"
)

const (
	NamespaceNotFound = "NamespaceNotFound"
)

func (a *ApiWrapper) moveNamespace(namespace string, projectID string) (status string, err error) {

	nsObject, err := a.findNamespace(namespace)
	if err != nil {
		return err.Error(), err
	}

	if nsObject.ProjectID != projectID {
		err = a.c.Namespace.ActionMove(&nsObject, &clusterClient.NamespaceMove{
			ProjectID: projectID,
		})

		if err != nil {
			return err.Error(), err
		}
	}

	return ProjectSynced, nil
}

func (a *ApiWrapper) findNamespace(namespace string) (nsObject clusterClient.Namespace, err error) {
	nsCollection, err := a.c.Namespace.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": namespace,
		},
	})
	if err != nil {
		return nsObject, err
	}

	if len(nsCollection.Data) == 0 {
		return nsObject, fmt.Errorf(NamespaceNotFound)
	}

	// ns is going to be unique in the cluster //
	nsObject = nsCollection.Data[0]
	return nsObject, nil
}

// Reconcile the mapping request//
func (a *ApiWrapper) ManageProjectNameSpaceMapping(pnsMapping *managementv1alpha1.ProjectNamespaceMapping) (status *managementv1alpha1.ProjectNamespaceMappingStatus,
	err error) {
	status = pnsMapping.Status.DeepCopy()
	nsObjectList := []clusterClient.Namespace{}

	for _, ns := range pnsMapping.Spec.Namespaces {
		nsObject, err := a.findNamespace(ns)
		if err != nil {
			status.Status = err.Error()
			return status, err
		}
		nsObjectList = append(nsObjectList, nsObject)
	}

	projectID, err := a.findProject(pnsMapping.Spec.ProjectName)
	if err != nil {
		status.Status = err.Error()
		return status, err
	}
	projectNSObjectList, err := a.c.Namespace.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"projectId": projectID,
		},
	})

	if err != nil {
		status.Status = err.Error()
		return status, err
	}

	if len(projectNSObjectList.Data) == 0 {
		// No existing namespaces mapped to project
		// just move ns in the list to the project
		for _, ns := range pnsMapping.Spec.Namespaces {
			status.Status, err = a.moveNamespace(ns, projectID)
			if err != nil {
				return status, err
			}
		}
	} else {
		// Need to reconcile current project NS mapping
		// with crd spec. Needed to ensure updates are
		// managed correctly.
		curNSList := []string{}
		for _, nsObject := range projectNSObjectList.Data {
			curNSList = append(curNSList, nsObject.Name)
		}

		moveList := []string{}
		for _, v := range pnsMapping.Spec.Namespaces {
			if !stringExists(curNSList, v) {
				moveList = append(moveList, v)
			}
		}

		removeList := []string{}
		for _, v := range curNSList {
			if !stringExists(pnsMapping.Spec.Namespaces, v) {
				removeList = append(removeList, v)
			}
		}

		if len(moveList) != 0 {
			// Need to move to project
			for _, ns := range moveList {
				status.Status, err = a.moveNamespace(ns, projectID)
				if err != nil {
					return status, err
				}
			}
		}

		if len(removeList) != 0 {
			// Need to move the ns out of the project
			for _, ns := range removeList {
				status.Status, err = a.moveNamespace(ns, "")
				if err != nil {
					return status, err
				}
			}
		}
	}
	return status, nil
}

func stringExists(list []string, item string) (ok bool) {
	for _, v := range list {
		if v == item {
			ok = true
			return ok
		}
	}
	return false
}

func (a *ApiWrapper) RemoveProjectNamespaceMapping(pnsMapping *managementv1alpha1.ProjectNamespaceMapping) (err error) {
	if len(pnsMapping.Spec.Namespaces) != 0 {
		for _, ns := range pnsMapping.Spec.Namespaces {
			_, err = a.moveNamespace(ns, "")
		}
		if err != nil {
			return err
		}
	}

	return nil
}
