package apphosting

type CiscoIOSXEAppHosting struct {
	CiscoIOSXEAppHostingCfgAppHostingCfgData CiscoIOSXEAppHostingCfgAppHostingCfgData `json:"Cisco-IOS-XE-app-hosting-cfg:app-hosting-cfg-data"`
}
type CiscoIOSXEAppHostingApp struct {
	CiscoIOSXEAppHostingCfgApp []App `json:"Cisco-IOS-XE-app-hosting-cfg:app"`
}
type CiscoIOSXEApp struct {
	App []App `json:"Cisco-IOS-XE-app-hosting-cfg:app"`
}

type ApplicationNetworkResource struct {
	AppintfVlanMode              string `json:"appintf-vlan-mode,omitempty"`
	AppintfAccessInterfaceNumber int    `json:"appintf-access-interface-number"`
	//AppintfAccessInterfaceNumber int    `json:"appintf-access-interface-number,omitempty"` // TODO pointer if 0
	//AppintfAccessInterfaceNumber *AppintfAccessInterfaceNumber `json:"appintf-access-interface-number"`
	//AppintfAccessInterfaceNumber *AppintfAccessInterfaceNumber
}
type AppintfVlanRule struct {
	VlanID         int `json:"vlan-id,omitempty"`
	GuestInterface int `json:"guest-interface"`
}
type AppintfVlanRules struct {
	AppintfVlanRule []AppintfVlanRule `json:"appintf-vlan-rule,omitempty"`
}
type RunOpts struct {
	LineIndex   int    `json:"line-index,omitempty"`
	LineRunOpts string `json:"line-run-opts,omitempty"`
}
type RunOptss struct {
	RunOpts []RunOpts `json:"run-opts,omitempty"`
}
type App struct {
	ApplicationName            string                     `json:"application-name,omitempty"`
	ApplicationNetworkResource ApplicationNetworkResource `json:"application-network-resource,omitempty"`
	AppintfVlanRules           *AppintfVlanRules          `json:"appintf-vlan-rules,omitempty"`
	DockerResource             bool                       `json:"docker-resource,omitempty"`
	RunOptss                   RunOptss                   `json:"run-optss,omitempty"`
	PrependPkgOpts             bool                       `json:"prepend-pkg-opts,omitempty"`
}
type Apps struct {
	App []App `json:"app,omitempty"`
}

type CiscoIOSXEAppHostingCfgAppHostingCfgData struct {
	Apps Apps `json:"apps,omitempty"`
}
