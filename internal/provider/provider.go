package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/provider"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("APPHOST_USERNAME", nil),
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  `The Username of the IOSXE switch. E.g.: "admin". This can also be set by environment variable "APPHOST_USERNAME".`,
			},
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("APPHOST_PASSWORD", nil),
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  `The Password of the IOSXE switch. E.g.: "somePassword". This can also be set by environment variable "APPHOST_PASSWORD".`,
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Allow insecure TLS. Default: true, means the API call is insecure.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Timeout for HTTP requests. Default value: 30.",
			},
			"ca_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPHOST_CA_FILE", nil),
				Description: "The path to CA certificate file (PEM). In case, certificate is based on legacy CN instead of ASN, set env. variable `GODEBUG=x509ignoreCN=0`. This can also be set by environment variable `APPHOST_CA_FILE`.",
			},
			"proxy_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPHOST_PROXY_URL", nil),
				Description: "Proxy Server URL with port number. This can also be set by environment variable `APPHOST_PROXY_URL`.",
			},
			"proxy_creds": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPHOST_PROXY_CREDS", nil),
				Description: "Proxy credential in format `username:password`. This can also be set by environment variable `APPHOST_PROXY_CREDS`.",
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Debug JSON Payloads in to debug folder",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ciscoapphosting_iox":  resourceIox(),
			"ciscoapphosting_app":  resourceApp(),
			"ciscoapphosting_copy": resourceCopy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ciscoapphosting_app": dataApp(),
		},
	}
	p.ConfigureContextFunc = providerConfigure(p)
	return p
}

func providerConfigure(p *schema.Provider) func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics
		if diags = providerValidateInput(d); diags.HasError() {
			return nil, diags
		}
		return &provider.ProviderClient{
			Provider: *d,
		}, diags
	}
}

func providerValidateInput(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	if v, ok := d.GetOk("username"); !ok && v == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Username is required",
			Detail:   "Username must be set for Cisco AppHosting Provider",
		})
	}
	if v, ok := d.GetOk("password"); !ok && v == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Password is required",
			Detail:   "Password must be set for Cisco AppHosting Provider",
		})
	}
	return diags
}
