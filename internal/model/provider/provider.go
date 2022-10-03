package provider

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type ProviderClient struct {
	Provider schema.ResourceData
	Method   string
	Path     string
	Payload  string
}
