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

locals {
  te_agent_token = "thousandeyes_agent_token"
}

resource "ciscoapphosting_app" "app" {
  host                 = var.hosts.0
  image                = "https://downloads.thousandeyes.com/enterprise-agent/thousandeyes-enterprise-agent-4.2.2.cisco.tar"
  app_gigabit_ethernet = "1/0/1"
  vlan_trunk           = false
  vlan                 = 1
  env = {
    TEAGENT_ACCOUNT_TOKEN = local.te_agent_token
  }
}