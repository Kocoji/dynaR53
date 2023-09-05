package main

import (
	"dynar53/utils"
	"os"
	"strings"
	"github.com/aws/aws-lambda-go/lambda"
)

func hdlr() {
	subs := strings.Split(strings.TrimSpace(os.Getenv("SUB_DOMAIN")), ",")
	ip := utils.GetEc2Ip(strings.TrimSpace(os.Getenv("INSTANCE_ID")))
	hostedZoneId := strings.TrimSpace(os.Getenv("HOSTEDZONE_ID"))

	utils.SetHost(hostedZoneId, ip, subs)
}

func main(){
	lambda.Start(hdlr)
}