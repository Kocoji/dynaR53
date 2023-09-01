package main

import (
	"dynar53/utils"
	"os"
	"strings"
)

func main() {
	subs := strings.Split(strings.TrimSpace(os.Getenv("SUB_DOMAIN")), ",")
	ip := utils.GetEc2Ip(strings.TrimSpace(os.Getenv("INSTANCE_ID")))
	hostedZoneId := strings.TrimSpace(os.Getenv("HOSTEDZONE_ID"))

	utils.SetHost(hostedZoneId, ip, subs)
}
