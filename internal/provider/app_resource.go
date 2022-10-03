package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/appgigabitethernet"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/apphosting"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/apphostingoper"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/hostname"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/provider"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/rpc"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/provider/iosxe"
)

func resourceApp() *schema.Resource {
	return &schema.Resource{
		Description: "Cisco AppHosting (IOX) App",
		Read:        resourceAppRead,
		Create:      resourceAppCreate,
		Update:      resourceAppUpdate,
		Delete:      resourceAppDelete,
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "app",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(string)
					if strings.Contains(v, "-") {
						errs = append(errs, fmt.Errorf("can't contain this char %v", v))
					}
					if strings.Contains(v, " ") {
						errs = append(errs, fmt.Errorf("can't contain space"))
					}
					return
				},
			},
			"image": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app_gigabit_ethernet": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "1/0/1",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"vlan_trunk": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"vlan": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 4094),
			},
			"guest_interface": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 1),
			},
			"docker": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"env": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"activate": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Cisco AppHosting Data READ")
	var data apphosting.CiscoIOSXEApp
	var oper apphostingoper.CiscoIOSXEAppHostingOper
	var eth appgigabitethernet.CiscoIOSXEAppGigabitEthernet
	var hostname hostname.Hostname
	var err error

	c, _ := meta.(*provider.ProviderClient)
	host := d.Get("host").(string)

	c.Method = "GET"
	c.Path = fmt.Sprintf("/data/Cisco-IOS-XE-app-hosting-cfg:app-hosting-cfg-data/apps/app=%v", d.Get("name").(string))
	resp, err := iosxe.Session(host, c)
	if err != nil {
		return err
	}

	bytes := []byte(resp)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}

	d.Set("name", data.App[0].ApplicationName)
	if data.App[0].ApplicationNetworkResource.AppintfVlanMode == "appintf-trunk" {
		d.Set("vlan_trunk", true)
		d.Set("vlan", data.App[0].AppintfVlanRules.AppintfVlanRule[0].VlanID)
		d.Set("guest_interface", data.App[0].AppintfVlanRules.AppintfVlanRule[0].GuestInterface)
	} else {
		d.Set("vlan_trunk", false)
		d.Set("guest_interface", data.App[0].ApplicationNetworkResource.AppintfAccessInterfaceNumber)
		appInt := strings.ReplaceAll(d.Get("app_gigabit_ethernet").(string), "/", "%2F")
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

	if data.App[0].DockerResource {
		d.Set("docker", data.App[0].DockerResource)
	} else {
		d.Set("docker", false)
	}

	options := make(map[string]string)
	for _, v := range data.App[0].RunOptss.RunOpts {
		env := strings.Split(v.LineRunOpts[3:], "=")
		options[env[0]] = env[1]
	}
	d.Set("env", options)

	c.Method = "GET"
	c.Path = fmt.Sprintf("/data/Cisco-IOS-XE-app-hosting-oper:app-hosting-oper-data/app=%v", d.Get("name").(string))
	resp, err = iosxe.Session(host, c)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(resp), &oper)
	if err != nil {
		return err
	}
	d.Set("state", strings.ToLower(oper.CiscoIOSXEAppHostingOperApp[0].Details.State))

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

func resourceAppCreate(d *schema.ResourceData, meta interface{}) error {
	var err error
	var resp string
	var bytes []byte
	var oper apphostingoper.CiscoIOSXEAppHostingOper
	var hostname hostname.Hostname
	var task string

	c, _ := meta.(*provider.ProviderClient)
	host := d.Get("host").(string)
	sleep := 10 * time.Second

	c.Method = "GET"
	c.Path = "/data/Cisco-IOS-XE-native:native/iox"
	_, err = iosxe.Session(host, c)
	if err != nil {
		return fmt.Errorf("apphosting not enabled (iox) - %v", err)
	}

	// Installing App
	c.Method = "POST"
	c.Path = "/data/Cisco-IOS-XE-rpc:app-hosting/"
	task = "install"
	apprpc := resourceAppRpc(d, task)
	if b, err := json.MarshalIndent(apprpc, "", "\t"); err == nil {
		c.Payload = string(b)
	}
	if c.Provider.Get("debug").(bool) {
		debugJson(fmt.Sprintf("apprpc-%v-%v", task, host), c.Payload)
	}
	resp, err = iosxe.Session(host, c)
	if err != nil {
		return err
	}
	if strings.Contains(strings.ToLower(resp), "error") {
		return fmt.Errorf(resp)
	}
	for i := 0; i < 30; {
		c.Method = "GET"
		c.Path = fmt.Sprintf("/data/Cisco-IOS-XE-app-hosting-oper:app-hosting-oper-data/app=%v", d.Get("name").(string))
		resp, err = iosxe.Session(host, c)
		if c.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("app-hosting-oper-%v-%v", task, host), resp)
		}
		if err != nil {
			if !strings.Contains(err.Error(), "not-found: /data/Cisco-IOS-XE-app-hosting-oper:app-hosting-oper-data/app") {
				return err
			}
		} else {
			err = json.Unmarshal([]byte(resp), &oper)
			if err != nil {
				return err
			}
			if strings.ToLower(oper.CiscoIOSXEAppHostingOperApp[0].Details.State) == "deployed" {
				break
			}
			if strings.ToLower(oper.CiscoIOSXEAppHostingOperApp[0].Details.State) == "stopped" {
				break
			}
		}

		i += 1
		log.Printf("[DEBUG] Waiting for app status: %v seconds\n", sleep)
		time.Sleep(sleep)
	}

	// Configure App
	c.Method = "PATCH"
	c.Path = "/data/Cisco-IOS-XE-app-hosting-cfg:app-hosting-cfg-data/"
	apphosting := resourceAppHostingYang(d)
	if b, err := json.MarshalIndent(apphosting, "", "\t"); err == nil {
		c.Payload = string(b)
	}
	if c.Provider.Get("debug").(bool) {
		debugJson(fmt.Sprintf("apphosting-%v", host), c.Payload)
	}
	_, err = iosxe.Session(host, c)
	if err != nil {
		return err
	}

	c.Method = "PATCH"
	c.Path = "/data/Cisco-IOS-XE-native:native/interface/AppGigabitEthernet"
	appethernet := resourceAppEthYang(d)
	if b, err := json.MarshalIndent(appethernet, "", "\t"); err == nil {
		c.Payload = string(b)
	}
	if c.Provider.Get("debug").(bool) {
		debugJson(fmt.Sprintf("appeth-%v", host), c.Payload)
	}
	_, err = iosxe.Session(host, c)
	if err != nil {
		return err
	}

	// Activating App
	if d.Get("activate").(bool) {
		for i := 0; i < 30; {
			c.Method = "GET"
			c.Path = fmt.Sprintf("/data/Cisco-IOS-XE-app-hosting-oper:app-hosting-oper-data/app=%v", d.Get("name").(string))
			resp, err = iosxe.Session(host, c)
			if c.Provider.Get("debug").(bool) {
				debugJson(fmt.Sprintf("app-hosting-oper-%v-%v", task, host), resp)
			}
			if err != nil {
				if !strings.Contains(err.Error(), "not-found: /data/Cisco-IOS-XE-app-hosting-oper:app-hosting-oper-data/app") {
					return err
				}
			}
			err = json.Unmarshal([]byte(resp), &oper)
			if err != nil {
				return err
			}

			switch strings.ToLower(oper.CiscoIOSXEAppHostingOperApp[0].Details.State) {
			case "deployed":
				task = "activate"
			case "deactivate":
				task = "activate"
			default:
				task = "start"
			}

			c.Method = "POST"
			c.Path = "/data/Cisco-IOS-XE-rpc:app-hosting/"
			apprpc := resourceAppRpc(d, task)
			if b, err := json.MarshalIndent(apprpc, "", "\t"); err == nil {
				c.Payload = string(b)
			}
			if c.Provider.Get("debug").(bool) {
				debugJson(fmt.Sprintf("apprpc-%v-%v", task, host), c.Payload)
			}
			resp, err = iosxe.Session(host, c)
			log.Println("[DEBUG] RPC response: ", resp)
			if err != nil {
				return err
			}
			if strings.Contains(strings.ToLower(resp), "error") {
				return fmt.Errorf(resp)
			}

			if task == "start" {
				break
			}

			log.Printf("[DEBUG] Waiting for app status: %v seconds\n", sleep)
			time.Sleep(sleep)
		}
	}

	c.Method = "GET"
	c.Path = fmt.Sprintf("/data/Cisco-IOS-XE-app-hosting-oper:app-hosting-oper-data/app=%v", d.Get("name").(string))
	resp, err = iosxe.Session(host, c)
	if err != nil {
		return err
	}
	if c.Provider.Get("debug").(bool) {
		debugJson(fmt.Sprintf("app-hosting-oper-%v-%v", task, host), resp)
	}
	err = json.Unmarshal([]byte(resp), &oper)
	if err != nil {
		return err
	}
	d.Set("state", strings.ToLower(oper.CiscoIOSXEAppHostingOperApp[0].Details.State))

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

func resourceAppUpdate(d *schema.ResourceData, meta interface{}) error {
	var resp string
	var task string
	var oper apphostingoper.CiscoIOSXEAppHostingOper
	var err error

	if d.HasChange("name") {
		oldState, _ := d.GetChange("name")
		d.Set("name", oldState)
		return fmt.Errorf("not allowed to change app name")
	}
	if d.HasChange("image") {
		oldState, _ := d.GetChange("image")
		d.Set("image", oldState)
		return fmt.Errorf("not allowed to change app image location")
	}

	c, _ := meta.(*provider.ProviderClient)
	host := d.Get("host").(string)
	sleep := 10 * time.Second

	// Configure App
	c.Method = "DELETE"
	c.Path = fmt.Sprintf("/data/Cisco-IOS-XE-app-hosting-cfg:app-hosting-cfg-data/apps/app=%v/appintf-vlan-rules", d.Get("name").(string))
	_, err = iosxe.Session(host, c)
	if err != nil {
		return err
	}

	time.Sleep(30 * time.Second)
	c.Method = "PATCH"
	c.Path = "/data/Cisco-IOS-XE-app-hosting-cfg:app-hosting-cfg-data/"
	apphosting := resourceAppHostingYang(d)
	if b, err := json.MarshalIndent(apphosting, "", "\t"); err == nil {
		c.Payload = string(b)
	}
	if c.Provider.Get("debug").(bool) {
		debugJson(fmt.Sprintf("apphosting-%v", host), c.Payload)
	}
	_, err = iosxe.Session(host, c)
	if err != nil {
		return err
	}

	c.Method = "PATCH"
	c.Path = "/data/Cisco-IOS-XE-native:native/interface/AppGigabitEthernet"
	appethernet := resourceAppEthYang(d)
	if b, err := json.MarshalIndent(appethernet, "", "\t"); err == nil {
		c.Payload = string(b)
	}
	if c.Provider.Get("debug").(bool) {
		debugJson(fmt.Sprintf("appeth-%v", host), c.Payload)
	}
	_, err = iosxe.Session(host, c)
	if err != nil {
		return err
	}

	// Installing App
	for i := 0; i < 30; {
		c.Method = "GET"
		c.Path = fmt.Sprintf("/data/Cisco-IOS-XE-app-hosting-oper:app-hosting-oper-data/app=%v", d.Get("name").(string))
		resp, err = iosxe.Session(host, c)
		if c.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("app-hosting-oper-%v-%v", task, host), resp)
		}
		if err != nil {
			if !strings.Contains(err.Error(), "not-found: /data/Cisco-IOS-XE-app-hosting-oper:app-hosting-oper-data/app") {
				return err
			}
		}
		err = json.Unmarshal([]byte(resp), &oper)
		if err != nil {
			return err
		}

		switch strings.ToLower(oper.CiscoIOSXEAppHostingOperApp[0].Details.State) {
		case "running":
			task = "stop"
		case "stopped":
			task = "deactivate"
		case "activated":
			task = "start"
		default:
			task = "activate"
		}

		c.Method = "POST"
		c.Path = "/data/Cisco-IOS-XE-rpc:app-hosting/"
		apprpc := resourceAppRpc(d, task)
		if b, err := json.MarshalIndent(apprpc, "", "\t"); err == nil {
			c.Payload = string(b)
		}
		if c.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("apprpc-%v-%v", task, host), c.Payload)
		}
		resp, err = iosxe.Session(host, c)
		log.Println("[DEBUG] RPC response: ", resp)
		if err != nil {
			return err
		}

		if task == "deactivate" {
			break
		}

		log.Printf("[DEBUG] Waiting for app status: %v seconds\n", sleep)
		time.Sleep(sleep)
	}

	if d.Get("activate").(bool) {
		for i := 0; i < 30; {
			c.Method = "GET"
			c.Path = fmt.Sprintf("/data/Cisco-IOS-XE-app-hosting-oper:app-hosting-oper-data/app=%v", d.Get("name").(string))
			resp, err = iosxe.Session(host, c)
			if c.Provider.Get("debug").(bool) {
				debugJson(fmt.Sprintf("app-hosting-oper-%v-%v", task, host), resp)
			}
			if err != nil {
				if !strings.Contains(err.Error(), "not-found: /data/Cisco-IOS-XE-app-hosting-oper:app-hosting-oper-data/app") {
					return err
				}
			}
			err = json.Unmarshal([]byte(resp), &oper)
			if err != nil {
				return err
			}

			switch strings.ToLower(oper.CiscoIOSXEAppHostingOperApp[0].Details.State) {
			case "deployed":
				task = "activate"
			case "activated":
				task = "start"
			}

			c.Method = "POST"
			c.Path = "/data/Cisco-IOS-XE-rpc:app-hosting/"
			apprpc := resourceAppRpc(d, task)
			if b, err := json.MarshalIndent(apprpc, "", "\t"); err == nil {
				c.Payload = string(b)
			}
			if c.Provider.Get("debug").(bool) {
				debugJson(fmt.Sprintf("apprpc-%v-%v", task, host), c.Payload)
			}
			resp, err = iosxe.Session(host, c)
			log.Println("[DEBUG] RPC response: ", resp)
			if err != nil {
				return err
			}
			if strings.Contains(strings.ToLower(resp), "error") {
				return fmt.Errorf(resp)
			}

			if task == "start" {
				break
			}
			log.Printf("[DEBUG] Waiting for app status: %v seconds\n", sleep)
			time.Sleep(sleep)
		}
	}

	d.SetId(d.Get("host").(string))
	return err
}

func resourceAppDelete(d *schema.ResourceData, meta interface{}) error {
	var resp string
	var err error
	var oper apphostingoper.CiscoIOSXEAppHostingOper
	var task string

	c, _ := meta.(*provider.ProviderClient)
	host := d.Get("host").(string)
	sleep := 10 * time.Second

	// Deactivate app
	for i := 0; i < 30; {
		c.Method = "GET"
		c.Path = fmt.Sprintf("/data/Cisco-IOS-XE-app-hosting-oper:app-hosting-oper-data/app=%v", d.Get("name").(string))
		resp, err = iosxe.Session(host, c)
		if err != nil {
			if !strings.Contains(err.Error(), "not-found: /data/Cisco-IOS-XE-app-hosting-oper:app-hosting-oper-data/app") {
				break
			}
		} else {
			err = json.Unmarshal([]byte(resp), &oper)
			if err != nil {
				return err
			}
			if strings.ToLower(oper.CiscoIOSXEAppHostingOperApp[0].Details.State) == "deactivate" {
				log.Println("[DEBUG] App is now deactivated")
				break
			}
		}
		switch strings.ToLower(oper.CiscoIOSXEAppHostingOperApp[0].Details.State) {
		case "running":
			task = "stop"
		case "stopped":
			task = "deactivate"
		default:
			task = "uninstall"
		}

		c.Method = "POST"
		c.Path = "/data/Cisco-IOS-XE-rpc:app-hosting/"
		apprpc := resourceAppRpc(d, task)
		if b, err := json.MarshalIndent(apprpc, "", "\t"); err == nil {
			c.Payload = string(b)
		}
		if c.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("apprpc-%v-%v", task, host), c.Payload)
		}
		resp, err = iosxe.Session(host, c)
		log.Println("[DEBUG] RPC response: ", resp)
		if err != nil {
			return err
		}
		if strings.Contains(strings.ToLower(resp), "error") {
			return fmt.Errorf(resp)
		}
		if task == "uninstall" {
			break
		}

		log.Printf("[DEBUG] Waiting for app status: %v seconds\n", sleep)
		time.Sleep(sleep)
	}

	c.Method = "DELETE"
	c.Path = fmt.Sprintf("/data/Cisco-IOS-XE-app-hosting-cfg:app-hosting-cfg-data/apps/app=%v", d.Get("name").(string))
	_, err = iosxe.Session(host, c)
	if err != nil {
		return err
	}

	c.Method = "POST"
	c.Path = "/data/Cisco-IOS-XE-rpc:default/"
	defaultrpc := resourceDefaultRpc(d)
	if b, err := json.MarshalIndent(defaultrpc, "", "\t"); err == nil {
		c.Payload = string(b)
	}
	if c.Provider.Get("debug").(bool) {
		debugJson(fmt.Sprintf("apprpc-%v-%v", "default", host), c.Payload)
	}
	resp, err = iosxe.Session(host, c)
	log.Println("[DEBUG] RPC response: ", resp)
	if err != nil {
		return err
	}

	d.SetId("")
	return err
}

func resourceAppHostingYang(d *schema.ResourceData) *apphosting.CiscoIOSXEAppHosting {
	data := &apphosting.CiscoIOSXEAppHosting{}
	app := &apphosting.App{}
	app.ApplicationName = d.Get("name").(string)
	if !d.Get("vlan_trunk").(bool) {
		app.ApplicationNetworkResource.AppintfVlanMode = "appintf-access"
		app.ApplicationNetworkResource.AppintfAccessInterfaceNumber = d.Get("guest_interface").(int)
	} else {
		app.ApplicationNetworkResource.AppintfVlanMode = "appintf-trunk"
		vlan := &apphosting.AppintfVlanRule{
			VlanID:         d.Get("vlan").(int),
			GuestInterface: d.Get("guest_interface").(int),
		}
		vlanRule := &apphosting.AppintfVlanRules{}
		vlanRule.AppintfVlanRule = append(vlanRule.AppintfVlanRule, *vlan)
		app.AppintfVlanRules = vlanRule
	}
	app.DockerResource = d.Get("docker").(bool)
	if _, ok := d.GetOk("env"); ok {
		app.PrependPkgOpts = true
		i := 1
		for k, v := range d.Get("env").(map[string]interface{}) {
			opt := &apphosting.RunOpts{
				LineIndex:   i,
				LineRunOpts: fmt.Sprintf("-e %v=%v", k, v),
			}
			app.RunOptss.RunOpts = append(app.RunOptss.RunOpts, *opt)
			i += 1
		}
	}

	data.CiscoIOSXEAppHostingCfgAppHostingCfgData.Apps.App = append(data.CiscoIOSXEAppHostingCfgAppHostingCfgData.Apps.App, *app)
	return data
}

func resourceAppEthYang(d *schema.ResourceData) *appgigabitethernet.CiscoIOSXEAppGigabitEthernet {
	data := &appgigabitethernet.CiscoIOSXEAppGigabitEthernet{}
	eth := &appgigabitethernet.AppGigabitEthernet{}
	eth.Name = d.Get("app_gigabit_ethernet").(string)
	eth.Description = "Managed by Terraform"
	if d.Get("vlan_trunk").(bool) {
		mode := &appgigabitethernet.TrunkMode{}
		eth.SwitchportConfig.Switchport.CiscoIOSXESwitchMode.Trunk = mode
		VlanCfg := &appgigabitethernet.SwitchportConfigSwitchportSwitchTrunk{}
		VlanCfg.Allowed.Vlan.Vlans = d.Get("vlan").(int)
		eth.SwitchportConfig.Switchport.CiscoIOSXESwitchTrunk = VlanCfg
		vlanTag := &appgigabitethernet.SwitchportSwitchTrunk{}
		vlanTag.Native.VlanConfig.Tag = true
		eth.Switchport.CiscoIOSXESwitchTrunk = vlanTag
		eth.Switchport.CiscoIOSXESwitchTrunk.Native.VlanConfig.Tag = true
	} else {
		mode := &appgigabitethernet.AccessMode{}
		eth.SwitchportConfig.Switchport.CiscoIOSXESwitchMode.Access = mode
		vlanCfg := &appgigabitethernet.SwitchportConfigSwitchportSwitchAccess{}
		vlanCfg.Vlan.Vlan = d.Get("vlan").(int)
		eth.SwitchportConfig.Switchport.CiscoIOSXESwitchAccess = vlanCfg
		vlan := &appgigabitethernet.SwitchportConfigSwitchportSwitchAccess{}
		vlan.Vlan.Vlan = d.Get("vlan").(int)
		eth.Switchport.CiscoIOSXESwitchAccess = vlan
		eth.Switchport.CiscoIOSXESwitchAccess.Vlan.Vlan = d.Get("vlan").(int)
	}

	data.CiscoIOSXENativeAppGigabitEthernet = append(data.CiscoIOSXENativeAppGigabitEthernet, *eth)
	return data
}

func resourceAppRpc(d *schema.ResourceData, task string) *rpc.AppHosting {
	data := &rpc.AppHosting{}
	app := d.Get("name").(string)
	switch task {
	case "install":
		data.CiscoIOSXERPCAppHosting.Install = &rpc.Install{Appid: app, Package: d.Get("image").(string)}
	case "activate":
		data.CiscoIOSXERPCAppHosting.Activate = &rpc.Activate{Appid: app}
	case "deactivate":
		data.CiscoIOSXERPCAppHosting.Deactivate = &rpc.Deactivate{Appid: app}
	case "start":
		data.CiscoIOSXERPCAppHosting.Start = &rpc.Start{Appid: app}
	case "stop":
		data.CiscoIOSXERPCAppHosting.Stop = &rpc.Stop{Appid: app}
	case "uninstall":
		data.CiscoIOSXERPCAppHosting.Uninstall = &rpc.Uninstall{Appid: app}
	case "default":
		log.Panicf("unknown rpc: %v", task)
	}
	return data
}

func resourceDefaultRpc(d *schema.ResourceData) *rpc.Default {
	data := &rpc.Default{}
	data.CiscoIOSXERPCDefault.Interface = fmt.Sprintf("AppGigabitEthernet%v", d.Get("app_gigabit_ethernet").(string))
	return data
}
