package main

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/byuoitav/common/log"
)

//SubmitTaskConfig .
func SubmitTaskConfig(f FargateTaskDefinition) error {

	StackName := "auto-cluster-" + f.Resources.Task.Properties.Family
	request := cloudformation.UpdateStackInput{
		StackName: &StackName,
	}

	log.L.Infof("%v", request)
	return nil
}
