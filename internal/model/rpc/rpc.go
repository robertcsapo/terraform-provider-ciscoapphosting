package rpc

type AppHosting struct {
	CiscoIOSXERPCAppHosting CiscoIOSXERPCAppHosting `json:"Cisco-IOS-XE-rpc:app-hosting"`
}
type Install struct {
	Appid   string `json:"appid,omitempty"`
	Package string `json:"package,omitempty"`
}
type Activate struct {
	Appid string `json:"appid,omitempty"`
}
type Deactivate struct {
	Appid string `json:"appid,omitempty"`
}
type Start struct {
	Appid string `json:"appid,omitempty"`
}
type Stop struct {
	Appid string `json:"appid,omitempty"`
}
type Uninstall struct {
	Appid string `json:"appid,omitempty"`
}
type CiscoIOSXERPCAppHosting struct {
	Install    *Install    `json:"install,omitempty"`
	Activate   *Activate   `json:"activate,omitempty"`
	Deactivate *Deactivate `json:"deactivate,omitempty"`
	Start      *Start      `json:"start,omitempty"`
	Stop       *Stop       `json:"stop,omitempty"`
	Uninstall  *Uninstall  `json:"uninstall,omitempty"`
}

type Default struct {
	CiscoIOSXERPCDefault CiscoIOSXERPCDefault `json:"Cisco-IOS-XE-rpc:default"`
}
type CiscoIOSXERPCDefault struct {
	Interface string `json:"interface"`
}

type Copy struct {
	CiscoIOSXERPCCopy CiscoIOSXERPCCopy `json:"Cisco-IOS-XE-rpc:copy"`
}
type CiscoIOSXERPCCopy struct {
	Source      string `json:"source-drop-node-name,omitempty"`
	Destination string `json:"destination-drop-node-name,omitempty"`
}

type Delete struct {
	CiscoIOSXERPCDelete CiscoIOSXERPCDelete `json:"Cisco-IOS-XE-rpc:delete"`
}
type CiscoIOSXERPCDelete struct {
	Filename string `json:"filename-drop-node-name,omitempty"`
}
