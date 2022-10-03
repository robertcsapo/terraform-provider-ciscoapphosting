package appgigabitethernet

type CiscoIOSXEAppGigabitEthernet struct {
	CiscoIOSXENativeAppGigabitEthernet []AppGigabitEthernet `json:"Cisco-IOS-XE-native:AppGigabitEthernet"`
}
type TrunkMode struct {
}
type AccessMode struct {
}
type SwitchMode struct {
	Trunk  *TrunkMode  `json:"trunk,omitempty"` // TrunkMode?
	Access *AccessMode `json:"access,omitempty"`
}
type AccessVlan struct {
	Vlan int `json:"vlan,omitempty"`
}
type SwitchportConfigSwitchportSwitchAccess struct {
	Vlan AccessVlan `json:"vlan,omitempty"`
}
type Vlan struct {
	Vlans int `json:"vlans,omitempty"`
}
type Allowed struct {
	Vlan Vlan `json:"vlan,omitempty"`
}
type SwitchportConfigSwitchportSwitchTrunk struct {
	Allowed Allowed `json:"allowed,omitempty"`
}
type SwitchportConfigSwitchport struct {
	CiscoIOSXESwitchMode   SwitchMode                              `json:"Cisco-IOS-XE-switch:mode,omitempty"`
	CiscoIOSXESwitchTrunk  *SwitchportConfigSwitchportSwitchTrunk  `json:"Cisco-IOS-XE-switch:trunk,omitempty"`
	CiscoIOSXESwitchAccess *SwitchportConfigSwitchportSwitchAccess `json:"Cisco-IOS-XE-switch:access,omitempty"`
}
type SwitchportConfig struct {
	Switchport SwitchportConfigSwitchport `json:"switchport,omitempty"`
}
type VlanConfig struct {
	Tag bool `json:"tag,omitempty"`
}
type Native struct {
	VlanConfig VlanConfig `json:"vlan-config,omitempty"`
}
type SwitchportSwitchTrunk struct {
	Native Native `json:"native,omitempty"`
}
type Switchport struct {
	CiscoIOSXESwitchTrunk  *SwitchportSwitchTrunk                  `json:"Cisco-IOS-XE-switch:trunk,omitempty"`
	CiscoIOSXESwitchAccess *SwitchportConfigSwitchportSwitchAccess `json:"Cisco-IOS-XE-switch:access,omitempty"`
}
type AppGigabitEthernet struct {
	Name             string           `json:"name,omitempty"`
	Description      string           `json:"description,omitempty"`
	SwitchportConfig SwitchportConfig `json:"switchport-config,omitempty"`
	Switchport       Switchport       `json:"switchport,omitempty"`
}
