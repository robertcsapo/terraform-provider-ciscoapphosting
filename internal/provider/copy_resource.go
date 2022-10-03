package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/provider"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/rpc"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/provider/iosxe"
)

func resourceCopy() *schema.Resource {
	return &schema.Resource{
		Description: "Cisco Filesystem Copy",
		Read:        resourceCopyRead,
		Create:      resourceCopyCreate,
		Update:      resourceCopyUpdate,
		Delete:      resourceCopyDelete,
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"src": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dst": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "flash:",
			},
			"partition": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filename": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCopyRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Cisco COPY READ")
	var data rpc.Filesystem
	var resp string
	var err error

	c, _ := meta.(*provider.ProviderClient)
	host := d.Get("host").(string)

	c.Method = "GET"
	c.Path = "/data/Cisco-IOS-XE-platform-software-oper:cisco-platform-software/q-filesystem"
	resp, err = iosxe.Session(host, c)
	if err != nil {
		return err
	}

	bytes := []byte(resp)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}

	resourceCopySplit(d)

	fileExists := false

	for _, slot := range data.CiscoIOSXEPlatformSoftwareOperQFilesystem {
		for _, p := range slot.Partitions {
			if strings.Contains(p.Name, d.Get("partition").(string)) {
				for _, content := range p.PartitionContent {
					if !strings.Contains(content.FullPath, "user/apps") {
						if strings.Contains(content.FullPath, d.Get("filename").(string)) {
							log.Println("File", content)
							fileExists = true
						}
					}
				}
			}
		}
	}
	if !fileExists {
		d.SetId("")
	}
	return err
}
func resourceCopyUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Cisco COPY UPDATE")
	return errors.New("not supported to update this resource")
}

func resourceCopyCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Cisco COPY CREATE")
	var resp string
	var err error

	c, _ := meta.(*provider.ProviderClient)
	host := d.Get("host").(string)

	c.Method = "POST"
	c.Path = "/data/Cisco-IOS-XE-rpc:copy/"
	copyrpc := resourceCopyRpc(d)
	if b, err := json.MarshalIndent(copyrpc, "", "\t"); err == nil {
		c.Payload = string(b)
	}
	if c.Provider.Get("debug").(bool) {
		debugJson(fmt.Sprintf("apprpc-%v-%v", "copy", host), c.Payload)
	}
	resp, err = iosxe.Session(host, c)
	log.Println("[DEBUG] RPC response: ", resp)
	if err != nil {
		return err
	}
	if strings.Contains(strings.ToLower(resp), "error") {
		return fmt.Errorf(resp)
	}

	resourceCopySplit(d)
	d.SetId(d.Get("host").(string))
	return err
}

func resourceCopyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Cisco COPY DELETE")
	var resp string
	var err error

	c, _ := meta.(*provider.ProviderClient)
	host := d.Get("host").(string)

	c.Method = "POST"
	c.Path = "/data/Cisco-IOS-XE-rpc:delete/"
	deleterpc := resourceCopyRpcDelete(d)
	if b, err := json.MarshalIndent(deleterpc, "", "\t"); err == nil {
		c.Payload = string(b)
	}
	if c.Provider.Get("debug").(bool) {
		debugJson(fmt.Sprintf("apprpc-%v-%v", "copy", host), c.Payload)
	}
	resp, err = iosxe.Session(host, c)
	log.Println("[DEBUG] RPC response: ", resp)
	if err != nil {
		return err
	}
	if strings.Contains(strings.ToLower(resp), "error") {
		return fmt.Errorf(resp)
	}

	d.SetId("")
	return err
}

func resourceCopyRpc(d *schema.ResourceData) *rpc.Copy {
	data := &rpc.Copy{}
	data.CiscoIOSXERPCCopy.Source = d.Get("src").(string)
	data.CiscoIOSXERPCCopy.Destination = d.Get("dst").(string)
	return data
}

func resourceCopySplit(d *schema.ResourceData) {
	partition := strings.Split(d.Get("dst").(string), ":")
	d.Set("partition", partition[0])

	filename := strings.Split(d.Get("src").(string), "/")
	d.Set("filename", filename[len(filename)-1])
}

func resourceCopyRpcDelete(d *schema.ResourceData) *rpc.Delete {
	data := &rpc.Delete{}
	data.CiscoIOSXERPCDelete.Filename = fmt.Sprintf("%v:%v", d.Get("partition").(string), d.Get("filename").(string))
	return data
}
