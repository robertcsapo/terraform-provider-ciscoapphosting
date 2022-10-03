package rpc

import "time"

type Filesystem struct {
	CiscoIOSXEPlatformSoftwareOperQFilesystem []CiscoIOSXEPlatformSoftwareOperQFilesystem `json:"Cisco-IOS-XE-platform-software-oper:q-filesystem"`
}
type Thresholds struct {
	WarningThresholdPercent  int `json:"warning-threshold-percent"`
	CriticalThresholdPercent int `json:"critical-threshold-percent"`
}
type PartitionContent struct {
	FullPath     string    `json:"full-path"`
	Size         string    `json:"size"`
	Type         string    `json:"type"`
	ModifiedTime time.Time `json:"modified-time"`
}
type Partitions struct {
	Name             string             `json:"name"`
	TotalSize        string             `json:"total-size"`
	UsedSize         string             `json:"used-size"`
	UsedPercent      int                `json:"used-percent"`
	DiskStatus       string             `json:"disk-status"`
	Thresholds       Thresholds         `json:"thresholds"`
	IsPrimary        bool               `json:"is-primary"`
	IsWritable       bool               `json:"is-writable"`
	PartitionContent []PartitionContent `json:"partition-content"`
}
type ImageFiles struct {
	FullPath string `json:"full-path"`
	FileSize string `json:"file-size"`
	Sha1Sum  string `json:"sha1sum"`
}
type CiscoIOSXEPlatformSoftwareOperQFilesystem struct {
	Fru        string       `json:"fru"`
	Slot       int          `json:"slot"`
	Bay        int          `json:"bay"`
	Chassis    int          `json:"chassis"`
	Partitions []Partitions `json:"partitions"`
	ImageFiles []ImageFiles `json:"image-files"`
}
