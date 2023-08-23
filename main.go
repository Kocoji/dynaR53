package main

import (
	"context"
	"fmt"
	"log"

	// "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// https://www.go-on-aws.com/aws-go-sdk-v2/sdkv2/process/
var client *ec2.Client

func init(){
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(),config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Fatal(err)
	}
	client = ec2.NewFromConfig(cfg)
	
}

func main() {
	instanceId := &ec2.DescribeInstancesInput{
		InstanceIds: []string{"i-0fa921b8db6b34212"},
	}
	
	result, err := client.DescribeInstances(context.TODO(),instanceId)
	if err != nil {
		panic(err)
	}
	count := len(result.Reservations)
	fmt.Println("Instances: ",count)

	// log.Print(result)
	for i, reservation := range result.Reservations {
		for k, instance := range reservation.Instances {
			fmt.Println("Instance number: ",i,"-",k	, "Id: ", instance.InstanceId)
			fmt.Printf("Instance ID: %s, State: %s\n", *instance.InstanceId, instance.State.Name)
		}
	}


	// output, err := ec2client.DescribeNetworkInterfaces(context.TODO(), &ec2.DescribeNetworkInterfacesInput{})
}
