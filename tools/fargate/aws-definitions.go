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
	Name             string           `json:"Name,omitempty"`
	CPU              string           `json:"Cpu,omitempty"`
	Memory           string           `json:"Memory,omitempty"`
	Image            string           `json:"Image,omitempty"`
	HealthCheck      *HealthCheck     `json:"HealthCheck,omitempty"`
	LogConfiguration LogConfiguration `json:"LogConfiguration"`
	PortMappings     []PortMapping    `json:"PortMappings,omitempty"`
	Environment      []EnvironmentVar `json:"Environment,omitempty"`

	RepositoryCredentials *Credentials `json:"RepositoryCredentials,omitempty"`
}

type Credentials struct {
	CredentialsParameters string `json:"CredentialsParameter,omitempty"`
}

//HealthCheck .
type HealthCheck struct {
	Command     []string `json:"Command,omitempty"`
	Interval    int      `json:"Interval,omitempty"`
	Retries     int      `json:"Retries,omitempty"`
	StartPeriod int      `json:"StartPeriod,omitempty"`
	Timeout     int      `json:"Timeout,omitempty"`
}

//LogConfiguration .
type LogConfiguration struct {
	LogDriver string            `json:"LogDriver"`
	Options   map[string]string `json:"Options"`
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

//Tag .
type Tag struct {
	Key   string `json:"Key,omitempty"`
	Value string `json:"Value,omitempty"`
}

//SecurityGroup .
type SecurityGroup struct {
	CidrIP      string `json:"CidrIp"`
	Description string `json:"Description"`
	FromPort    int    `json:"FromPort"`
	IPProtocol  string `json:"IpProtocol"`
	ToPort      int    `json:"ToPort"`
}

//Cluster  .
type Cluster struct {
	Properties struct {
		ClusterName string `json:"ClusterName"`
	} `json:"Properties"`
	Type string `json:"Type"`
}

//NewServiceCloudformationStack .
type NewServiceCloudformationStack struct {
	AWSTemplateFormatVersion string `json:"AWSTemplateFormatVersion"`
	Resources                struct {
		Cluster  *Cluster `json:"Cluster,omitempty"`
		Listener struct {
			Properties struct {
				Certificates []struct {
					CertificateArn string `json:"CertificateArn,omitempty"`
				} `json:"Certificates,omitempty"`
				DefaultActions []struct {
					TargetGroupArn struct {
						Ref string `json:"Ref"`
					} `json:"TargetGroupArn"`
					Type string `json:"Type"`
				} `json:"DefaultActions"`
				LoadBalancerArn struct {
					Ref string `json:"Ref"`
				} `json:"LoadBalancerArn"`
				Port     int    `json:"Port"`
				Protocol string `json:"Protocol"`
			} `json:"Properties"`
			Type string `json:"Type"`
		} `json:"Listener"`
		LoadBalancer struct {
			Properties struct {
				Name           string `json:"Name"`
				Scheme         string `json:"Scheme"`
				SecurityGroups []struct {
					Ref string `json:"Ref"`
				} `json:"SecurityGroups"`
				Subnets []string `json:"Subnets"`
				Tags    []Tag    `json:"Tags"`
			} `json:"Properties"`
			Type string `json:"Type"`
		} `json:"LoadBalancer"`
		LoadBalancerSecurityGroup struct {
			Properties struct {
				GroupDescription     string          `json:"GroupDescription"`
				GroupName            string          `json:"GroupName"`
				SecurityGroupEgress  []SecurityGroup `json:"SecurityGroupEgress"`
				SecurityGroupIngress []SecurityGroup `json:"SecurityGroupIngress"`
				Tags                 []Tag           `json:"Tags"`
				VpcID                string          `json:"VpcId"`
			} `json:"Properties"`
			Type string `json:"Type"`
		} `json:"LoadBalancerSecurityGroup"`
		SecurityGroup struct {
			Properties struct {
				GroupDescription     string          `json:"GroupDescription"`
				GroupName            string          `json:"GroupName"`
				SecurityGroupEgress  []SecurityGroup `json:"SecurityGroupEgress"`
				SecurityGroupIngress []SecurityGroup `json:"SecurityGroupIngress"`
				Tags                 []Tag           `json:"Tags"`
				VpcID                string          `json:"VpcId"`
			} `json:"Properties"`
			Type string `json:"Type"`
		} `json:"SecurityGroup"`
		Service struct {
			DependsOn  []string `json:"DependsOn"`
			Properties struct {
				Cluster                 string `json:"Cluster"`
				DeploymentConfiguration struct {
					MaximumPercent        int `json:"MaximumPercent"`
					MinimumHealthyPercent int `json:"MinimumHealthyPercent"`
				} `json:"DeploymentConfiguration"`
				DesiredCount                  int    `json:"DesiredCount"`
				HealthCheckGracePeriodSeconds int    `json:"HealthCheckGracePeriodSeconds"`
				LaunchType                    string `json:"LaunchType"`
				LoadBalancers                 []struct {
					ContainerName  string `json:"ContainerName"`
					ContainerPort  int    `json:"ContainerPort"`
					TargetGroupArn struct {
						Ref string `json:"Ref"`
					} `json:"TargetGroupArn"`
				} `json:"LoadBalancers"`
				NetworkConfiguration struct {
					AwsvpcConfiguration struct {
						AssignPublicIP string `json:"AssignPublicIp"`
						SecurityGroups []struct {
							Ref string `json:"Ref"`
						} `json:"SecurityGroups"`
						Subnets []string `json:"Subnets"`
					} `json:"AwsvpcConfiguration"`
				} `json:"NetworkConfiguration"`
				Role           string `json:"Role"`
				ServiceName    string `json:"ServiceName"`
				TaskDefinition struct {
					Ref string `json:"Ref"`
				} `json:"TaskDefinition"`
			} `json:"Properties"`
			Type string `json:"Type"`
		} `json:"Service"`
		LogGroup struct {
			Type       string `json:"Type"`
			Properties struct {
				LogGroupName    string `json:"LogGroupName"`
				RetentionInDays int    `json:"RetentionInDays"`
			} `json:"Properties"`
		} `json:"LogGroup"`
		TargetGroup struct {
			Properties struct {
				HealthCheckIntervalSeconds int    `json:"HealthCheckIntervalSeconds"`
				HealthCheckPath            string `json:"HealthCheckPath"`
				HealthCheckPort            int    `json:"HealthCheckPort"`
				HealthCheckProtocol        string `json:"HealthCheckProtocol"`
				HealthCheckTimeoutSeconds  int    `json:"HealthCheckTimeoutSeconds"`
				HealthyThresholdCount      int    `json:"HealthyThresholdCount"`
				TargetGroupAttributes      []Tag  `json:"TargetGroupAttributes"`
				Matcher                    struct {
					HTTPCode string `json:"HttpCode"`
				} `json:"Matcher"`
				Name                    string `json:"Name"`
				Port                    int    `json:"Port"`
				Protocol                string `json:"Protocol"`
				Tags                    []Tag  `json:"Tags"`
				TargetType              string `json:"TargetType"`
				UnhealthyThresholdCount int    `json:"UnhealthyThresholdCount"`
				VpcID                   string `json:"VpcId"`
			} `json:"Properties"`
			Type string `json:"Type"`
		} `json:"TargetGroup"`
		TaskDef Task `json:"TaskDef"`
	} `json:"Resources"`
}
