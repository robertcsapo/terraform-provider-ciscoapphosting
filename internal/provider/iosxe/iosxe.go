package iosxe

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CiscoDevNet/iosxe-go-client/client"
	"github.com/CiscoDevNet/iosxe-go-client/container"
	"github.com/CiscoDevNet/iosxe-go-client/models"
	"github.com/robertcsapo/terraform-provider-ciscoapphosting/internal/model/provider"
)

type HttpClient struct {
	Client   *client.V2
	Provider *provider.ProviderClient
	Host     string
	Method   string
	Path     string
	Payload  string
}

func Session(host string, p *provider.ProviderClient) (string, error) {
	var err error
	var body string
	s := &HttpClient{
		Provider: p,
		Host:     fmt.Sprintf("%v", host),
		Method:   p.Method,
		Path:     p.Path,
		Payload:  p.Payload,
	}
	body, err = s.methods(s.Method)
	if err != nil {
		log.Println("[DEBUG] ERROR Session: ", err)
		return body, err
	}
	return body, nil
}

func NewClient(host string, d *provider.ProviderClient) (*client.V2, error) {
	var err error
	iosxeV2Client, err := client.NewV2(
		host,
		d.Provider.Get("username").(string),
		d.Provider.Get("password").(string),
		d.Provider.Get("timeout").(int),
		d.Provider.Get("insecure").(bool),
		d.Provider.Get("proxy_url").(string),
		d.Provider.Get("proxy_creds").(string),
		d.Provider.Get("ca_file").(string),
	)

	if err != nil {
		return nil, err
	}
	return iosxeV2Client, err

}

func (s *HttpClient) methods(method string) (string, error) {
	var err error
	var body string
	var resp *http.Response
	var container *container.Container
	var httpRetry int
	var httpSleep time.Duration

	httpRetry = 20
	httpSleep = 10 * time.Second

	iosxeGM := &models.GenericModel{
		JSONPayload: s.Payload,
	}

	c, _ := NewClient(
		fmt.Sprintf("https://%v", s.Host),
		s.Provider,
	)

	switch method {
	case "GET":
		log.Printf("[DEBUG] IOS-XE GET %v on: %v \n", s.Path, s.Host)
		for i := 0; ; i++ {
			resp, container, err = c.Get(s.Path, nil)
			if resp != nil {
				if resp.StatusCode == 409 {
					log.Println(s.httpErrorMsg(resp.StatusCode))
				} else {
					break
				}
			}
			if i >= (httpRetry - 1) {
				break
			}
			log.Printf("[DEBUG] IOS-XE Retry: (%v/%v) waiting: %v\n", i, httpRetry, httpSleep)
			if err != nil {
				log.Println("[DEBUG] ERROR GET: ", err)
				break
			}
			time.Sleep(httpSleep)
		}
		if err != nil {
			return body, err
		}
		body = container.String()
	case "POST":
		log.Printf("[DEBUG] IOS-XE POST %v on: %v \n", s.Path, s.Host)
		for i := 0; ; i++ {
			resp, container, err = c.Create(s.Path, iosxeGM)
			if resp != nil {
				if resp.StatusCode == 409 {
					log.Println(s.httpErrorMsg(resp.StatusCode))
				} else {
					break
				}
			}
			if i >= (httpRetry - 1) {
				break
			}
			log.Printf("[DEBUG] IOS-XE Retry: (%v/%v) waiting: %v\n", i, httpRetry, httpSleep)
			time.Sleep(httpSleep)
		}
		if err != nil {
			log.Println("[DEBUG] ERROR PATCH: ", err, resp.Status)
			return body, err
		}
		body = container.String()
	case "PATCH":
		log.Printf("[DEBUG] IOS-XE PATCH %v on: %v \n", s.Path, s.Host)
		for i := 0; ; i++ {
			resp, container, err = c.Patch(s.Path, iosxeGM)
			if resp != nil {
				if resp.StatusCode == 409 {
					log.Println(s.httpErrorMsg(resp.StatusCode))
				} else {
					break
				}
			}
			if i >= (httpRetry - 1) {
				break
			}
			log.Printf("[DEBUG] IOS-XE Retry: (%v/%v) waiting: %v\n", i, httpRetry, httpSleep)
			time.Sleep(httpSleep)
		}
		if err != nil {
			log.Println("[DEBUG] ERROR PATCH: ", err, resp.Status)
			return body, err
		}
		body = container.String()
	case "UPDATE":
		log.Printf("[DEBUG] IOS-XE UPDATE %v on: %v \n", s.Path, s.Host)
		for i := 0; ; i++ {
			resp, err = c.Update(s.Path, iosxeGM)
			if resp != nil {
				if resp.StatusCode == 409 {
					log.Println(s.httpErrorMsg(resp.StatusCode))
				} else {
					break
				}
			}
			if i >= (httpRetry - 1) {
				break
			}
			log.Printf("IOS-XE Retry: (%v/%v) waiting: %v\n", i, httpRetry, httpSleep)
			time.Sleep(httpSleep)
		}
		if err != nil {
			log.Println("[DEBUG] ERROR UPDATE: ", err, resp.Status)
			return body, err
		}
	case "DELETE":
		log.Printf("[DEBUG] IOS-XE DELETE %v on: %v \n", s.Path, s.Host)
		for i := 0; ; i++ {
			resp, err = c.Delete(s.Path)
			if resp != nil {
				if resp.StatusCode == 409 {
					log.Println(s.httpErrorMsg(resp.StatusCode))
				} else {
					break
				}
			}
			if i >= (httpRetry - 1) {
				break
			}
			log.Printf("[DEBUG] IOS-XE Retry: (%v/%v) waiting: %v\n", i, httpRetry, httpSleep)
			time.Sleep(httpSleep)
		}
		if err != nil {
			log.Println("[DEBUG] ERROR DELETE: ", err, resp.Status)
			return body, err
		}
	}
	return body, nil
}

func (*HttpClient) httpErrorMsg(status int) string {
	var msg string
	if status == 409 {
		msg = "[DEBUG] IOS-XE configuration database is unavailable"
	}
	return msg
}
