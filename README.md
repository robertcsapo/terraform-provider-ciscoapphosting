# Terraform Provider Cisco AppHosting
Tech Preview (Early field trial)

terraform-provider-ciscoapphosting is a Terraform Provider for Cisco Catalyst 9000 Switches.

## Requirements for Development

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.18
- Cisco IOS-XE 17.8+
- Catalyst 9000 Switching Series
    - Catalyst 9300

## Using the provider

Use ```terraform init``` to download the plugin from Terrafrom Registry.

Configure the provider to connect towards your Cisco Catalyst 9000 Switches
```
terraform {
  required_providers {
    ciscoapphosting = {
      source  = "robertcsapo/ciscoapphosting"
      version = "1.0.0"
    }
  }
}

provider "ciscoapphosting" {
  username = var.username
  password = var.password
  insecure = var.insecure
  timeout  = var.timeout
}

resource "ciscoapphosting_app" "app" {
  host = "127.0.0.1"
  image = "https://downloads.thousandeyes.com/enterprise-agent/thousandeyes-enterprise-agent-4.2.2.cisco.tar"
  app_gigabit_ethernet = "1/0/1"
  vlan_trunk = false
  vlan = 1
  env = {
    TEAGENT_ACCOUNT_TOKEN = "token"
  }
}
```

Examples can be found in [examples/](./examples/).