package management

import (
	"fmt"
	"strings"

	"github.com/rancher/norman/types"

	managementv1alpha1 "github.com/ibrokethecloud/rancher-project-operator/pkg/api/v1alpha1"
	managementv3 "github.com/rancher/types/client/management/v3"
)

func (a *ApiWrapper) ManageRoleTemplateBinding(prtb *managementv1alpha1.ProjectRoleTemplateBinding) (status *managementv1alpha1.ProjectRoleTemplateBindingStatus, err error) {
	prtbObject, err := a.generatePRTB(prtb)
	status = prtb.Status.DeepCopy()
	if err != nil {
		status.Status = err.Error()
		return status, err
	}

	newprtbObject, err := a.m.ProjectRoleTemplateBinding.Create(&prtbObject)
	if err != nil {
		status.Status = err.Error()
		return status, err
	}

	status.Status = PRTBCreated
	status.ID = newprtbObject.ID
	return status, nil
}

func (a *ApiWrapper) generatePRTB(prtb *managementv1alpha1.ProjectRoleTemplateBinding) (object managementv3.ProjectRoleTemplateBinding, err error) {
	projectID, err := a.findProject(prtb.Spec.ProjectName)
	if err != nil {
		return object, err
	}
	object.ProjectID = projectID

	// check for users //
	if len(prtb.Spec.UserName) != 0 && len(prtb.Spec.UserPrincipalName) == 0 && len(prtb.Spec.GroupName) == 0 && len(prtb.Spec.GroupPrincipalName) == 0 {
		// need to find the principalName
		userID, err := a.findUserPrincipalID(prtb.Spec.UserName)
		if err != nil {
			return object, err
		}
		object.UserPrincipalID = userID
	} else if len(prtb.Spec.UserName) == 0 && len(prtb.Spec.UserPrincipalName) != 0 && len(prtb.Spec.GroupName) == 0 && len(prtb.Spec.GroupPrincipalName) == 0 {
		object.UserPrincipalID = prtb.Spec.UserPrincipalName
	} else if len(prtb.Spec.UserName) == 0 && len(prtb.Spec.UserPrincipalName) == 0 && len(prtb.Spec.GroupName) != 0 && len(prtb.Spec.GroupPrincipalName) == 0 {
		// find group name
	} else if len(prtb.Spec.UserName) == 0 && len(prtb.Spec.UserPrincipalName) == 0 && len(prtb.Spec.GroupName) == 0 && len(prtb.Spec.GroupPrincipalName) != 0 {
		// find group name
	} else {
		// error condition
		err = fmt.Errorf("Only one of UserName or Group needs to be specified in a prtb")
		return object, err
	}

	object.RoleTemplateID = prtb.Spec.RoleTemplateName

	return object, nil
}

func (a *ApiWrapper) findUserPrincipalID(userName string) (userPrincipalID string, err error) {
	userList, err := a.m.User.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"username": userName,
		},
	})

	if err != nil {
		return userPrincipalID, err
	}

	if len(userList.Data) == 0 {
		return "", fmt.Errorf(NoUserFound)
	}

	if len(userList.Data) > 1 {
		return "", fmt.Errorf(MultipleUserFound)
	}

	userPrincipalID = userList.Data[0].PrincipalIDs[0]

	return userPrincipalID, err
}

func (a *ApiWrapper) DeletePRTB(prtbID string) (err error) {
	prtbObject, err := a.m.ProjectRoleTemplateBinding.ByID(prtbID)

	if err != nil {
		// In case the prtb has already been manually cleaned up we can ignore
		if strings.Contains(err.Error(), "404 Not Found") {
			return nil
		}
		return err
	}

	err = a.m.ProjectRoleTemplateBinding.Delete(prtbObject)

	return err
}
