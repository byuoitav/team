package main

var newservicetemplate = `{
    "AWSTemplateFormatVersion": "2010-09-09",
    "Resources": {
        "Cluster": {
            "ClusterName": "",
            "Type": "AWS::ECS::Cluster"
        },
        "Listener": {
            "Properties": {
                "Certificates": [
                    {
                        "CertificateArn": ""
                    }
                ],
                "DefaultActions": [
                    {
                        "TargetGroupArn": {
                            "Ref": "TargetGroup"
                        },
                        "Type": "forward"
                    }
                ],
                "LoadBalancerArn": {
                    "Ref": "LoadBalancer"
                },
                "Port": 1234,
                "Protocol": "HTTP"
            },
            "Type": "AWS::ElasticLoadBalancingV2::Listener"
        },
        "LoadBalancer": {
            "Properties": {
                "Name": "",
                "Scheme": "",
                "SecurityGroups": [
                    {
                        "Ref": "LoadBalancerSecurityGroup"
                    }
                ],
                "Subnets": [
                    "",
                    ""
                ],
                "Tags": [ ] },
            "Type": "AWS::ElasticLoadBalancingV2::LoadBalancer"
        },
        "LoadBalancerSecurityGroup": {
            "Properties": {
                "GroupDescription": "",
                "GroupName": "",
                "SecurityGroupEgress": [
                    {
                        "CidrIp": "0.0.0.0/0",
                        "Description": "All Devices",
                        "FromPort": 0,
                        "IpProtocol": "tcp",
                        "ToPort": 65535
                    }
                ],
                "SecurityGroupIngress": [
                    {
                        "CidrIp": "",
                        "Description": "",
                        "FromPort": 123,
                        "IpProtocol": "tcp",
                        "ToPort": 456
                    }
                ],
                "Tags": [
                ],
                "VpcId": ""
            },
            "Type": "AWS::EC2::SecurityGroup"
        },
        "SecurityGroup": {
            "Properties": {
                "GroupDescription": "",
                "GroupName": "",
                "SecurityGroupEgress": [
                    {
                        "CidrIp": "0.0.0.0/0",
                        "Description": "All Devices",
                        "FromPort": 0,
                        "IpProtocol": "tcp",
                        "ToPort": 65535
                    }
                ],
                "SecurityGroupIngress": [
                    {
                        "CidrIp": "",
                        "Description": "",
                        "FromPort": 123,
                        "IpProtocol": "tcp",
                        "ToPort": 456
                    }
                ],
                "Tags": [
                ],
                "VpcId": ""
            },
            "Type": "AWS::EC2::SecurityGroup"
        },
        "Service": {
            "DependsOn": [
                "Listener"
            ],
            "Properties": {
                "Cluster": "",
                "DeploymentConfiguration": {
                    "MaximumPercent": 200,
                    "MinimumHealthyPercent": 50
                },
                "DesiredCount": 1234,
                "HealthCheckGracePeriodSeconds": 180,
                "LaunchType": "FARGATE",
                "LoadBalancers": [
                    {
                        "ContainerName": "",
                        "ContainerPort": 1234,
                        "TargetGroupArn": {
                            "Ref": "TargetGroup"
                        }
                    }
                ],
                "NetworkConfiguration": {
                    "AwsvpcConfiguration": {
                        "AssignPublicIp": "DISABLED",
                        "SecurityGroups": [
                            {
                                "Ref": "SecurityGroup"
                            }
                        ],
                        "Subnets": [
                            "",
                            ""
                        ]
                    }
                },
                "Role": "",
                "ServiceName": "",
                "TaskDefinition": {
                    "Ref": "TaskDef"
                }
            },
            "Type": "AWS::ECS::Service"
        },
        "TargetGroup": {
            "Properties": {
                "HealthCheckIntervalSeconds": 10,
                "HealthCheckPath": "/status",
                "HealthCheckPort": 1234,
                "HealthCheckProtocol": "HTTP",
                "HealthCheckTimeoutSeconds": 3,
                "HealthyThresholdCount": 2,
                "Matcher": {
                    "HttpCode": "200-399"
                },
                "Name": "",
                "Port": 1234,
                "Protocol": "HTTP",
                "Tags": [
                ],
                "TargetType": "ip",
                "UnhealthyThresholdCount": 5,
                "VpcId": ""
            },
            "Type": "AWS::ElasticLoadBalancingV2::TargetGroup"
        },
		"LogGroup": {
            "Type": "AWS::Logs::LogGroup",
            "Properties": {
                "LogGroupName": "",
                "RetentionInDays": 3
			}
        },
        "TaskDef": {
            "Properties": {
                "ContainerDefinitions": [
                    {
                        "Environment": [
                            {
                                "Name": "",
                                "Value": ""
                            }
                        ],
                        "HealthCheck": {
                            "Command": [
                                ""
                            ],
                            "Interval": 1234,
                            "Retries": 1234,
                            "StartPeriod": 1234,
                            "Timeout": 1234
                        },
                        "Hostname": "",
                        "Image": "",
                        "Name": "",
                        "PortMappings": [
                            {
                                "ContainerPort": "7100",
                                "Protocol": "tcp"
                            }
                        ]
                    }
                ],
                "Cpu": "512",
                "ExecutionRoleArn": "",
                "Family": "",
                "Memory": "1024",
                "NetworkMode": "awsvpc",
                "RequiresCompatibilities": [
                    "FARGATE",
                    "EC2"
                ],
                "TaskRoleArn": ""
            },
            "Type": "AWS::ECS::TaskDefinition"
        }
    }

}`
