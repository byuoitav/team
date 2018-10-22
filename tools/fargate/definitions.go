package main

// FargateTaskDefinition .
type FargateTaskDefinition struct {
	AWSTemplateFormatVersion string                 `json:"AWSTemplateFormatVersion,omitempty"`
	Description              string                 `json:"Description,omitempty"`
	Parameters               map[string]interface{} `json:"Parameters,omitempty"`
	Resources                struct {
		ECSCluster struct {
			Type       string `json:"Type,omitempty"`
			Properties struct {
				ClusterName string `json:"ClusterName,omitempty"`
			} `json:"Properties,omitempty"`
		} `json:"ECSCluster,omitempty"`
		Task Task `json:"Task,omitempty"`
	} `json:"Resources,omitempty"`
}

//Task is a fargate task
type Task struct {
	Type       string `json:"Type,omitempty"`
	Properties struct {
		Family                  string                `json:"Family,omitempty"`
		CPU                     string                `json:"Cpu,omitempty"`
		Memory                  string                `json:"Memory,omitempty"`
		NetworkMode             string                `json:"NetworkMode,omitempty"`
		TaskRoleArn             string                `json:"TaskRoleArn,omitempty"`
		ExecutionRoleArn        string                `json:"ExecutionRoleArn,omitempty"`
		RequiresCompatibilities []string              `json:"RequiresCompatibilities,omitempty"`
		ContainerDefinitions    []ContainerDefinition `json:"ContainerDefinitions,omitempty"`
	} `json:"Properties,omitempty"`
}

//ContainerDefinition .
type ContainerDefinition struct {
	Name         string           `json:"Name,omitempty"`
	CPU          string           `json:"Cpu,omitempty"`
	Memory       string           `json:"Memory,omitempty"`
	Image        string           `json:"Image,omitempty"`
	PortMappings []PortMapping    `json:"PortMappings,omitempty"`
	Environment  []EnvironmentVar `json:"Environment,omitempty"`
}

//PortMapping .
type PortMapping struct {
	ContainerPort string `json:"ContainerPort,omitempty"`
	Protocol      string `json:"Protocol,omitempty"`
}

//EnvironmentVar .
type EnvironmentVar struct {
	Name  string `json:"Name,omitempty"`
	Value string `json:"Value,omitempty"`
}

//ConfigDefinition .
type ConfigDefinition struct {
	Name                 string   `json:"name,omitempty"`
	Port                 string   `json:"port,omitempty"`
	EnvironmentVariables []string `json:"env-vars,omitempty"`
}

//ConfigInfoWrapper .
type ConfigInfoWrapper struct {
	ID        string                   `json:"_id,omitempty"`
	Rev       string                   `json:"_rev,omitempty"`
	AWSStages map[string]AWSConfigInfo `json:"aws-stages,omitempty"`
	Stages    map[string]ConfigInfo    `json:"stages,omitempty"`
}

//AWSConfigInfo .
type AWSConfigInfo struct {
	TaskRoleArn      string `json:"task-role-arn,omitempty,omitempty"`
	ExecutionRoleArn string `json:"execution-role-arn,omitempty"`
	CPU              string `json:"cpu,omitempty"`
	Memory           string `json:"memory,omitempty"`
	Task             string `json:"task,omitempty"` //taskName -if any
	ConfigInfo
}

//AWSTaskWrapper .
type AWSTaskWrapper struct {
	ID        string                 `json:"_id,omitempty"`
	Rev       string                 `json:"_rev,omitempty"`
	AWSStages map[string]AWSTaskInfo `json:"aws-stages,omitempty"`
}

//AWSTaskInfo .
type AWSTaskInfo struct {
	Name             string   `json:"string,omitempty"`
	Description      string   `json:"description,omitempty"`
	Services         []string `json:"services,omitempty"`
	CPU              string   `json:"cpu,omitempty"`
	Memory           string   `json:"memory,omitempty"`
	TaskRoleArn      string   `json:"task-role-arn,omitempty"`
	ExecutionRoleArn string   `json:"execution-role-arn,omitempty"`
}

//ConfigInfo .
type ConfigInfo struct {
	Port              string            `json:"port,omitempty"`
	EnvironmentValues map[string]string `json:"environment-values,omitempty"`
}
