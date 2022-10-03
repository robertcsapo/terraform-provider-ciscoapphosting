package iox

type CiscoIOSXENative struct {
	CiscoIOSXENativeNative CiscoIOSXENativeNative `json:"Cisco-IOS-XE-native:native"`
}
type CiscoIOSXENativeNative struct {
	Iox []string `json:"Cisco-IOS-XE-native:iox"`
}
