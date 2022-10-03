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

resource "ciscoapphosting_app" "thousandeyes" {
  for_each             = toset(var.hosts)
  host                 = each.key
  name                 = "thousandeyes"
  image                = "https://downloads.thousandeyes.com/enterprise-agent/thousandeyes-enterprise-agent-4.2.2.cisco.tar"
  app_gigabit_ethernet = "1/0/1"
  vlan_trunk           = true
  vlan                 = 102
  env = {
    TEAGENT_ACCOUNT_TOKEN = local.te_agent_token
  }
  activate = true
}