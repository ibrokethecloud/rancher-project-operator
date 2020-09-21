package management

import (
	"fmt"

	"reflect"

	"encoding/json"

	managementv1alpha1 "github.com/ibrokethecloud/rancher-project-operator/pkg/api/v1alpha1"
	"github.com/rancher/norman/types"
	managementv3 "github.com/rancher/types/client/management/v3"
)

const (
	ProjectNotFound        = "ProjectNotFound"
	MultipleProjectFound   = "MultipleProjectFound"
	ProjectCreationError   = "ProjectCreationError"
	ProjectSynced          = "ProjectSynced"
	ProjectCreated         = "ProjectCreated"
	ProjectUpdated         = "ProjectUpdated"
	ProjectExists          = "ProjectExists"
	ProjectMonitoringError = "ProjectMonitoringError"
	ProjectUpdateError     = "ProjectUpdateError"
	NoUserFound            = "NoUserFound"
	MultipleUserFound      = "MultipleUserFound"
	PRTBCreated            = "ProjectRoleTemplateBindingCreated"
)

func (a *ApiWrapper) ManageProject(project *managementv1alpha1.Project) (status *managementv1alpha1.ProjectStatus, err error) {
	status = &managementv1alpha1.ProjectStatus{}

	projectID, err := a.findProject(project.Spec.DisplayName)
	if status.ID == "" && err.Error() == ProjectNotFound {
		// Project doesnt already exist so lets create it //
		inputProject := a.generateProject(*project)
		outputProject, err := a.m.Project.Create(&inputProject)
		if err != nil {
			status.Status = ProjectCreationError
			return status, err
		}

		if !reflect.ValueOf(project.Spec.MonitoringInput).IsZero() && project.Spec.EnableProjectMonitoring {
			monitoringInput := generateMonitoringInput(project)
			err := a.m.Project.ActionEditMonitoring(outputProject, &monitoringInput)
			if err != nil {
				status.Status = ProjectMonitoringError
				return status, err
			}
		}

		status.ID = outputProject.ID
		status.Status = ProjectSynced
		status.State = ProjectCreated
	} else if status.ID == projectID {
		// Need to Update the project
		currentProject, err := a.m.Project.ByID(status.ID)
		if err != nil {
			return status, err
		}
		update := map[string]interface{}{
			"name":                          project.Spec.DisplayName,
			"description":                   project.Spec.Description,
			"containerDefaultResourceLimit": project.Spec.ContainerDefaultResourceLimit,
			"namespaceDefaultResourceQuota": project.Spec.NamespaceDefaultResourceQuota,
			"resourcequota":                 project.Spec.ResourceQuota,
			"labels":                        project.Spec.Labels,
		}
		newProject, err := a.m.Project.Update(currentProject, update)
		if err != nil {
			status.Status = err.Error()
			status.State = ProjectUpdateError
			return status, err
		}

		// If there is a change to PSP
		if currentProject.PodSecurityPolicyTemplateName != project.Spec.PodSecurityPolicyTemplateName {
			pspInput := &managementv3.SetPodSecurityPolicyTemplateInput{
				PodSecurityPolicyTemplateName: project.Spec.PodSecurityPolicyTemplateName,
			}
			_, err = a.m.Project.ActionSetpodsecuritypolicytemplate(newProject, pspInput)
			if err != nil {
				status.Status = err.Error()
				return status, err
			}
		}

		// If there is a change to monitoring state
		if currentProject.EnableProjectMonitoring != project.Spec.EnableProjectMonitoring {
			if !project.Spec.EnableProjectMonitoring {
				err = a.m.Project.ActionDisableMonitoring(newProject)
				if err != nil {
					status.Status = err.Error()
					return status, err
				}
			}

			if project.Spec.EnableProjectMonitoring {
				monitoringInput := generateMonitoringInput(project)
				err = a.m.Project.ActionEnableMonitoring(newProject, &monitoringInput)
				if err != nil {
					status.Status = ProjectMonitoringError
					return status, err
				}
			}
		}

		status.Status = ProjectSynced
		status.State = ProjectUpdated
	} else {
		status.Status = err.Error()
		return status, err
	}

	return status, nil
}

func (a *ApiWrapper) RemoveProject(project *managementv1alpha1.Project) (err error) {
	// check if project status has an ID
	if project.Status.ID == "" {
		// nothing to do.. lets ignore this project //
		return nil
	}

	projectID, err := a.findProject(project.Spec.DisplayName)
	if err != nil {
		if err.Error() == ProjectNotFound {
			return nil
		}
		return err
	}

	if projectID == project.Status.ID {
		inputProject, err := a.m.Project.ByID(projectID)
		if err != nil {
			return err
		}
		err = a.m.Project.Delete(inputProject)
	}

	return err
}

func (a *ApiWrapper) findProject(project string) (projectID string, err error) {
	projects, err := a.m.Project.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"cluster_id": a.clusterID,
			"name":       project,
		},
	})

	if err != nil {
		return "", err
	}

	count := len(projects.Data)
	if count <= 0 {
		//project not found
		return "", fmt.Errorf(ProjectNotFound)
	}

	if count > 1 {
		return "", fmt.Errorf(MultipleProjectFound)
	}

	return projects.Data[0].ID, nil
}

func (a *ApiWrapper) generateProject(input managementv1alpha1.Project) (output managementv3.Project) {
	output = managementv3.Project{
		ClusterID:               a.clusterID,
		Name:                    input.Spec.DisplayName,
		EnableProjectMonitoring: input.Spec.EnableProjectMonitoring,
		Description:             input.Spec.Description,
	}

	if !reflect.ValueOf(input.Spec.ResourceQuota).IsZero() {
		tmpProjectResourceQuota, err := convertProjectResourceQuote(input.Spec.ResourceQuota)
		if err != nil {
			output.ResourceQuota = &tmpProjectResourceQuota
		}
	}

	if !reflect.ValueOf(input.Spec.NamespaceDefaultResourceQuota).IsZero() {
		tmpNamespaceDefaultResourceQuota, err := convertNamespaceQuota(input.Spec.NamespaceDefaultResourceQuota)
		if err != nil {
			output.NamespaceDefaultResourceQuota = &tmpNamespaceDefaultResourceQuota
		}
	}

	if !reflect.ValueOf(input.Spec.ContainerDefaultResourceLimit).IsZero() {
		tmpContainerDefaultResourceLimit, err := convertContainerResourceLimit(input.Spec.ContainerDefaultResourceLimit)
		if err != nil {
			output.ContainerDefaultResourceLimit = &tmpContainerDefaultResourceLimit
		}
	}

	if len(input.Labels) != 0 {
		output.Labels = input.Spec.Labels
	}

	if len(input.Spec.PodSecurityPolicyTemplateName) != 0 {
		output.PodSecurityPolicyTemplateName = input.Spec.PodSecurityPolicyTemplateName
	}

	return output
}

func convertResourceQuote(input managementv1alpha1.ResourceQuotaLimit) (output managementv3.ResourceQuotaLimit, err error) {
	outputByte, err := json.Marshal(input)
	if err != nil {
		return output, err
	}
	output = managementv3.ResourceQuotaLimit{}
	err = json.Unmarshal(outputByte, &output)

	return output, err
}

func convertNamespaceQuota(input managementv1alpha1.NamespaceResourceQuota) (output managementv3.NamespaceResourceQuota, err error) {
	outputByte, err := json.Marshal(input)
	if err != nil {
		return output, err
	}
	output = managementv3.NamespaceResourceQuota{}
	err = json.Unmarshal(outputByte, &output)

	return output, err
}

func convertContainerResourceLimit(input managementv1alpha1.ContainerResourceLimit) (output managementv3.ContainerResourceLimit, err error) {
	outputByte, err := json.Marshal(input)
	if err != nil {
		return output, err
	}
	output = managementv3.ContainerResourceLimit{}
	err = json.Unmarshal(outputByte, &output)

	return output, err
}

func convertProjectResourceQuote(input managementv1alpha1.ProjectResourceQuota) (output managementv3.ProjectResourceQuota, err error) {
	outputByte, err := json.Marshal(input)
	if err != nil {
		return output, err
	}
	output = managementv3.ProjectResourceQuota{}
	err = json.Unmarshal(outputByte, &output)

	return output, err
}

func generateMonitoringInput(project *managementv1alpha1.Project) (monitoringInput managementv3.MonitoringInput) {
	if len(project.Spec.MonitoringInput.Version) != 0 {
		monitoringInput.Version = project.Spec.MonitoringInput.Version
	}

	if !reflect.ValueOf(project.Spec.MonitoringInput.Answers).IsZero() {
		monitoringInput.Answers = project.Spec.MonitoringInput.Answers
	}

	return monitoringInput
}
