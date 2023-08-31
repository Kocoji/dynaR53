package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"os"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// https://www.go-on-aws.com/aws-go-sdk-v2/sdkv2/process/
var ec2Client *ec2.Client
var r53Client *route53.Client

func init(){
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(),config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Fatal(err)
	}
	ec2Client = ec2.NewFromConfig(cfg)
	r53Client = route53.NewFromConfig(cfg)
	
}

func main() {
	subs := strings.Split(strings.TrimSpace(os.Getenv("SUB_DOMAIN")), ",")
	ip := getEc2Ip(strings.TrimSpace(os.Getenv("INSTANCE_ID")))
	hostedZoneId := strings.TrimSpace(os.Getenv("HOSTEDZONE_ID"))
	
	setHost(hostedZoneId, ip,subs )
}


func getEc2Ip(iId string) string{
	instanceId := &ec2.DescribeInstancesInput{
		InstanceIds: []string{iId},
	}	
	result, err := ec2Client.DescribeInstances(context.TODO(),instanceId)
	if err != nil {
		panic(err)
	}
	count := len(result.Reservations)
	fmt.Println("Instances: ",count)
	publicIp := *result.Reservations[0].Instances[0].NetworkInterfaces[0].PrivateIpAddresses[0].Association.PublicIp
	fmt.Println("Public IP is: ",publicIp)
	return publicIp
}

func getHost(hostedZoneId string )string{
	input := &route53.GetHostedZoneInput{Id: aws.String(hostedZoneId),}
	rs, e := r53Client.GetHostedZone(context.TODO(),input)
	if e != nil {
		log.Fatal("Cannot get Host")
	}
	result := *rs.HostedZone.Name
	return result
}

// https://stackoverflow.com/questions/64530843/aws-route53-adding-simple-record
func setHost(hostedZoneId, ip string, subs []string){
	hostedZone := getHost(hostedZoneId)
	for _, v := range subs {
		input := &route53.ChangeResourceRecordSetsInput{
			ChangeBatch: &types.ChangeBatch{
				Changes: []types.Change{
					{
						Action: types.ChangeActionUpsert,
						ResourceRecordSet: &types.ResourceRecordSet{
							Name: aws.String(v+"."+hostedZone),
							Type: types.RRTypeA,
							ResourceRecords: []types.ResourceRecord{
								{
									Value: aws.String(ip),
								},
							},
							TTL: aws.Int64(300),
						},
					},
				},Comment: aws.String("Updated by dynaR53"),
			},
			HostedZoneId: aws.String(hostedZoneId),
		}

		set, err := r53Client.ChangeResourceRecordSets(context.TODO(), input)
		if err != nil {
			log.Panic(err)
		}
		fmt.Print(set.ChangeInfo.Status)
	}
}