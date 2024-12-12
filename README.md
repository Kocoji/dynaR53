# What is this?
Because Route53 does not support dynamic updates of the Public IP from a single EC2 instance, this script retrieves the EC2’s public IP and updates the A record to the hosted zone on Route53.

# Requirement
These environment variables is required for this script work:
- `INSTANCE_ID`: Instance ID
- `HOSTEDZONE_ID`: HostedZone ID
- `SUBDOMAINS`: the list of subdomain you want to update, split by commas `,` .e.g: `"abc,def,gxh"`
- `REGION_ID`: the region where EC2 instance is running.
  