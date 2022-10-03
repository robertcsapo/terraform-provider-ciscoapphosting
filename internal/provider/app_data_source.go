package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/appgigabitethernet"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/apphosting"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/hostname"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/provider"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/provider/iosxe"
)

func dataApp() *schema.Resource {
	return &schema.Resource{
		Read:        dataAppRead,
		Description: "Cisco AppHosting (IOX) App",
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vlan_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"guest_interface": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"docker": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"env": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataAppRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Cisco AppHosting Data READ")
	var resp string
	var bytes []byte
	var data apphosting.CiscoIOSXEAppHostingApp
	var eth appgigabitethernet.CiscoIOSXEAppGigabitEthernet
	var hostname hostname.Hostname
	var err error

	c, _ := meta.(*provider.ProviderClient)
	host := d.Get("host").(string)

	c.Method = "GET"
	c.Path = fmt.Sprintf("/data/Cisco-IOS-XE-app-hosting-cfg:app-hosting-cfg-data/apps/app=%v", d.Get("name").(string))
	resp, err = iosxe.Session(host, c)
	if err != nil {
		return err
	}
	bytes = []byte(resp)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}

	d.Set("name", data.CiscoIOSXEAppHostingCfgApp[0].ApplicationName)
	d.Set("vlan_mode", data.CiscoIOSXEAppHostingCfgApp[0].ApplicationNetworkResource.AppintfVlanMode)
	if data.CiscoIOSXEAppHostingCfgApp[0].ApplicationNetworkResource.AppintfVlanMode == "appintf-trunk" {
		d.Set("vlan", data.CiscoIOSXEAppHostingCfgApp[0].AppintfVlanRules.AppintfVlanRule[0].VlanID)
		d.Set("guest_interface", data.CiscoIOSXEAppHostingCfgApp[0].AppintfVlanRules.AppintfVlanRule[0].GuestInterface)
	} else {
		d.Set("guest_interface", data.CiscoIOSXEAppHostingCfgApp[0].ApplicationNetworkResource.AppintfAccessInterfaceNumber)
		appInt := strings.ReplaceAll(fmt.Sprintf("1/0/%v", data.CiscoIOSXEAppHostingCfgApp[0].ApplicationNetworkResource.AppintfAccessInterfaceNumber+1), "/", "%2F")
		c.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/interface/AppGigabitEthernet=%v", appInt)
		resp, err := iosxe.Session(host, c)
		if err != nil {
			return err
		}

		bytes := []byte(resp)
		err = json.Unmarshal(bytes, &eth)
		if err != nil {
			return err
		}
		d.Set("vlan", eth.CiscoIOSXENativeAppGigabitEthernet[0].SwitchportConfig.Switchport.CiscoIOSXESwitchAccess.Vlan.Vlan)
	}
	if data.CiscoIOSXEAppHostingCfgApp[0].DockerResource {
		d.Set("docker", data.CiscoIOSXEAppHostingCfgApp[0].DockerResource)
	} else {
		d.Set("docker", false)
	}
	options := make(map[string]string)
	for _, v := range data.CiscoIOSXEAppHostingCfgApp[0].RunOptss.RunOpts {
		env := strings.Split(v.LineRunOpts[3:], "=")
		options[env[0]] = env[1]
	}
	d.Set("env", options)

	c.Method = "GET"
	c.Path = "/data/Cisco-IOS-XE-native:native/hostname"
	resp, err = iosxe.Session(host, c)
	if err != nil {
		return err
	}
	bytes = []byte(resp)
	err = json.Unmarshal(bytes, &hostname)
	if err != nil {
		return err
	}
	d.Set("hostname", hostname.CiscoIOSXENativeHostname)

	d.SetId(d.Get("host").(string))
	return err
}
