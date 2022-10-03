package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/hostname"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/iox"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/provider"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/provider/iosxe"
)

func resourceIox() *schema.Resource {
	return &schema.Resource{
		Description: "Cisco AppHosting (IOX) Service",
		Read:        resourceIoxRead,
		Create:      resourceIoxCreate,
		Update:      resourceIoxUpdate,
		Delete:      resourceIoxDelete,
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIoxRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Cisco IOX READ")
	var hostname hostname.Hostname
	var resp string
	var bytes []byte
	var err error

	c, _ := meta.(*provider.ProviderClient)
	host := d.Get("host").(string)

	c.Method = "GET"
	c.Path = "/data/Cisco-IOS-XE-native:native/iox"
	_, err = iosxe.Session(host, c)
	if err != nil {
		d.Set("enable", false)
		return err // TODO handle drift change for state
	}
	d.Set("enable", true)

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

func resourceIoxCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Cisco IOX CREATE")
	var data iox.CiscoIOSXENative
	var hostname hostname.Hostname
	var resp string
	var bytes []byte
	var err error

	c, _ := meta.(*provider.ProviderClient)
	host := d.Get("host").(string)

	c.Method = "PATCH"
	c.Path = "/data/Cisco-IOS-XE-native:native"
	data.CiscoIOSXENativeNative.Iox = null()
	if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
		c.Payload = string(b)
	}
	if c.Provider.Get("debug").(bool) {
		debugJson(fmt.Sprintf("iox-%v", host), c.Payload)
	}
	_, err = iosxe.Session(host, c)
	if err != nil {
		return err
	}
	d.Set("enable", true)

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

	// Possible to add a function to handle IOX readiness
	// "IOX may take upto 5 mins to be ready"
	// Takes longer if USB is avaiable
	// "Migrating IOX from Bootflash to Harddisk"

	return err
}

func resourceIoxUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Cisco AppHosting IOX UPDATE")
	return fmt.Errorf("not possible to update this resource. if you want disable iox then delete this resource")
}

func resourceIoxDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Cisco IOX DELETE")
	var err error

	c, _ := meta.(*provider.ProviderClient)
	host := d.Get("host").(string)

	c.Method = "DELETE"
	c.Path = "/data/Cisco-IOS-XE-native:native/iox"
	_, err = iosxe.Session(host, c)
	if err != nil {
		return err
	}
	d.SetId("")
	return err
}
