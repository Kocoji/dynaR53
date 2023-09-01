package utils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"log"
)

var ec2Client *ec2.Client
var r53Client *route53.Client

func init() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Fatal(err)
	}
	ec2Client = ec2.NewFromConfig(cfg)
	r53Client = route53.NewFromConfig(cfg)

}

func GetEc2Ip(iId string) string {
	instanceId := &ec2.DescribeInstancesInput{
		InstanceIds: []string{iId},
	}
	result, err := ec2Client.DescribeInstances(context.TODO(), instanceId)
	if err != nil {
		panic(err)
	}
	count := len(result.Reservations)
	fmt.Println("Instances: ", count)
	publicIp := *result.Reservations[0].Instances[0].NetworkInterfaces[0].PrivateIpAddresses[0].Association.PublicIp
	fmt.Println("Public IP is: ", publicIp)
	return publicIp
}

func GetHost(hostedZoneId string) string {
	input := &route53.GetHostedZoneInput{Id: aws.String(hostedZoneId)}
	rs, e := r53Client.GetHostedZone(context.TODO(), input)
	if e != nil {
		log.Fatal("Cannot get Host ", e)
	}
	result := *rs.HostedZone.Name
	return result
}

func SetHost(hostedZoneId, ip string, subs []string) {
	hostedZone := GetHost(hostedZoneId)
	for _, v := range subs {
		input := &route53.ChangeResourceRecordSetsInput{
			ChangeBatch: &types.ChangeBatch{
				Changes: []types.Change{
					{
						Action: types.ChangeActionUpsert,
						ResourceRecordSet: &types.ResourceRecordSet{
							Name: aws.String(v + "." + hostedZone),
							Type: types.RRTypeA,
							ResourceRecords: []types.ResourceRecord{
								{
									Value: aws.String(ip),
								},
							},
							TTL: aws.Int64(300),
						},
					},
				}, Comment: aws.String("Updated by dynaR53"),
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
