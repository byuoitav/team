package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/byuoitav/common/log"
)

const template = ` 
{
"AWSTemplateFormatVersion": "2010-09-09",
    "Resources": {
        "ECSCluster": {
            "Type": "AWS::ECS::Cluster"
        },
		"Service": {
            "Type": "AWS::ECS::Service" },
        "Task": {
            "Type": "AWS::ECS::TaskDefinition",
            "Properties": {
                "NetworkMode": "awsvpc",
                "RequiresCompatibilities": [
                    "FARGATE",
                    "EC2"
                ]
            }
        }
    }
}`

func buildTaskDefinitionConfig(wrap ConfigInfoWrapper, def ConfigDefinition, dbName, branch string) (FargateTaskDefinition, string, error) {
	var toReturn FargateTaskDefinition
	if _, ok := wrap.AWSStages[branch]; !ok {
		return toReturn, "", fmt.Errorf("No branch information for %v defined", branch)
	}

	err := json.Unmarshal([]byte(template), &toReturn)
	if err != nil {
		log.L.Fatalf("Invalid template, binary is corrupt")
	}

	//now we just need to fill in the information that is unique to this service
	stageInfo := wrap.AWSStages[branch]

	var taskStage AWSTaskInfo
	var ok bool

	//check to see if it's part of a task
	taskwrap, err := GetTaskInfoFromDB(stageInfo.Task)
	if err != nil {
		return toReturn, "", fmt.Errorf("Couldn't get task definition %v from database. %v", wrap.AWSStages[branch].Task, err.Error())
	}
	taskStage, ok = taskwrap.AWSStages[branch]
	if !ok {
		return toReturn, "", fmt.Errorf("No branch information for %v defined in the task definition", branch)
	}

	//we need to build it with this info

	toReturn.Description = taskStage.Description
	toReturn.Resources.ECSCluster.Properties.ClusterName = taskStage.Name
	toReturn.Resources.Task.Properties.Family = fmt.Sprintf("%v--%v", taskStage.Name, branch)
	toReturn.Resources.Task.Properties.CPU = taskStage.CPU
	toReturn.Resources.Task.Properties.Memory = taskStage.Memory
	toReturn.Resources.Task.Properties.TaskRoleArn = taskStage.TaskRoleArn
	toReturn.Resources.Task.Properties.ExecutionRoleArn = taskStage.ExecutionRoleArn

	//we need to be able to handle multiple containers in this definition - we pull them from the DB - for this one we pull the list of variables from the ConfigDefinition and throw an error if they don't match those in the Wrapper? The definition env variables becomes a documentation tool that is enforced...
	var defs []ContainerDefinition
	var env []EnvironmentVar

	log.L.Infof("%v", def.EnvironmentVariables)

	//do enforcement
	for k, v := range stageInfo.EnvironmentValues {
		if !contains(def.EnvironmentVariables, k) {
			return toReturn, "", fmt.Errorf("Environment variable %v in database not in Config Definition", k)
		}
		env = append(env, EnvironmentVar{Name: k, Value: v})
	}

	//build our container definition
	cDef := ContainerDefinition{
		Name:  def.Name,
		Image: fmt.Sprintf("byuoitav/%v:%v", def.Name, branch),
		PortMappings: []PortMapping{PortMapping{
			ContainerPort: stageInfo.Port,
			Protocol:      "tcp",
		}},
		Environment: env,
		HealthCheck: HealthCheck{
			Command:     []string{"CMD-SHELL", fmt.Sprintf("/usr/bin/wget -q -O- http://localhost:%v/status", strings.Trim(def.Port, ":"))},
			Interval:    5,
			Timeout:     2,
			StartPeriod: 60,
			Retries:     10,
		},
		LogConfiguration: LogConfiguration{
			LogDriver: "awslogs",
			Options: map[string]string{
				"awslogs-group":         fmt.Sprintf("/ecs/%v--%v", taskStage.Name, branch),
				"awslogs-region":        "us-west-2",
				"awslogs-stream-prefix": "ecs",
			},
		},
	}
	defs = append(defs, cDef)

	//if it's part of a task we need to add all the other containers running with this one
	for i := range taskStage.Services {
		//make sure it's not us
		if taskStage.Services[i] == def.Name {
			continue
		}

		cDef, err = buildContainerDefinition(taskStage.Name, taskStage.Services[i], branch, dbName)
		if err != nil {
			return toReturn, "", fmt.Errorf("Couldn't create task: %v", err.Error())
		}
		defs = append(defs, cDef)
	}

	toReturn.Resources.Task.Properties.ContainerDefinitions = defs
	return toReturn, taskStage.Name + "--" + branch, nil
}

func buildContainerDefinition(taskname, name, branch, dbName string) (ContainerDefinition, error) {
	var toReturn ContainerDefinition
	wrap, err := GetInfoFromDB(dbName, name)
	if err != nil {
		return toReturn, err
	}
	wrapStage, ok := wrap.AWSStages[branch]
	if !ok {
		return toReturn, fmt.Errorf("No definition for branch %v on service %v", branch, name)
	}

	toReturn = ContainerDefinition{
		Name:  name,
		Image: fmt.Sprintf("byuoitav/%v:%v", name, branch),
		PortMappings: []PortMapping{PortMapping{
			ContainerPort: wrapStage.Port,
			Protocol:      "tcp",
		}},
		HealthCheck: HealthCheck{
			Command:     []string{"CMD-SHELL", fmt.Sprintf("/usr/bin/wget -q -O- http://localhost:%v/status", strings.Trim(wrapStage.Port, ":"))},
			Interval:    5,
			Timeout:     2,
			StartPeriod: 60,
			Retries:     10,
		},
		LogConfiguration: LogConfiguration{
			LogDriver: "awslogs",
			Options: map[string]string{
				"awslogs-group":         fmt.Sprintf("/ecs/%v--%v", taskname, branch),
				"awslogs-region":        "us-west-2",
				"awslogs-stream-prefix": "ecs",
			},
		},
	}
	for k, v := range wrapStage.EnvironmentValues {
		toReturn.Environment = append(toReturn.Environment, EnvironmentVar{Name: k, Value: v})
	}

	return toReturn, nil
}

func contains(list []string, a string) bool {
	for i := range list {
		if list[i] == a {
			return true
		}
	}
	return false
}
