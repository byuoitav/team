package main

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
	Task           string   `json:"task,omitempty"`           //taskName -if any
	Image          string   `json:"image,omitempty"`          // docker image name
	Credentials    string   `json:"credentials,omitempty"`    // arn to secret in secrets manager
	HealthCheckCmd []string `json:"healthCheckCmd,omitempty"` // leave blank to disable health checks
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
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Services    []string `json:"services,omitempty"`
	CPU         string   `json:"cpu,omitempty"`
	Memory      string   `json:"memory,omitempty"`

	TaskRoleArn      string `json:"task-role-arn,omitempty"`
	ExecutionRoleArn string `json:"execution-role-arn,omitempty"`
	ServiceRoleArn   string `json:"service-role-arn,omitempty"`
	CertificateArn   string `json:"certificate-arn,omitempty"`

	PublicPort    int    `json:"public-port,omitempty"`
	PublicService string `json:"public-service,omitempty"`
	Public        bool

	PublicSubnets  []string          `json:"public-subnets"`
	PrivateSubnets []string          `json:"private-subnets"`
	Tags           map[string]string `json:"tags"`
	VPCID          string            `json:"vpc-id"`

	InstanceCount int `json:"instance-count"`
}

//ConfigInfo .
type ConfigInfo struct {
	Port              string            `json:"port,omitempty"`
	EnvironmentValues map[string]string `json:"environment-values,omitempty"`
}
