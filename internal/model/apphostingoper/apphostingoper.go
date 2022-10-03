package apphostingoper

type CiscoIOSXEAppHostingOper struct {
	CiscoIOSXEAppHostingOperApp []CiscoIOSXEAppHostingOperApp `json:"Cisco-IOS-XE-app-hosting-oper:app"`
}
type Application struct {
	Name              string `json:"name"`
	InstalledVersion  string `json:"installed-version"`
	Description       string `json:"description"`
	Type              string `json:"type"`
	Owner             string `json:"owner"`
	ActivationAllowed bool   `json:"activation-allowed"`
	Author            string `json:"author"`
}
type Signing struct {
	KeyType string `json:"key-type"`
	Method  string `json:"method"`
}
type Licensing struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
type PackageInformation struct {
	Name        string      `json:"name"`
	Path        string      `json:"path"`
	Application Application `json:"application"`
	Signing     Signing     `json:"signing"`
	Licensing   Licensing   `json:"licensing"`
	URLPath     string      `json:"url-path"`
}
type Processes struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Pid    string `json:"pid"`
	Uptime string `json:"uptime"`
	Memory string `json:"memory"`
}
type DetailedGuestStatus struct {
	Processes Processes `json:"processes"`
}
type ResourceReservation struct {
	Disk       string `json:"disk"`
	Memory     string `json:"memory"`
	CPU        string `json:"cpu"`
	Vcpu       string `json:"vcpu"`
	CPUPercent int    `json:"cpu-percent"`
}
type ResourceAdmission struct {
	State     string `json:"state"`
	DiskSpace string `json:"disk-space"`
	Memory    string `json:"memory"`
	CPU       string `json:"cpu"`
	Vcpus     string `json:"vcpus"`
}
type Details struct {
	State                 string              `json:"state"`
	PackageInformation    PackageInformation  `json:"package-information"`
	DetailedGuestStatus   DetailedGuestStatus `json:"detailed-guest-status"`
	ActivatedProfileName  string              `json:"activated-profile-name"`
	ResourceReservation   ResourceReservation `json:"resource-reservation"`
	GuestInterface        string              `json:"guest-interface"`
	ResourceAdmission     ResourceAdmission   `json:"resource-admission"`
	DockerRunOpts         string              `json:"docker-run-opts"`
	Command               string              `json:"command"`
	EntryPoint            string              `json:"entry-point"`
	HealthStatus          int                 `json:"health-status"`
	LastHealthProbeError  string              `json:"last-health-probe-error"`
	LastHealthProbeOutput string              `json:"last-health-probe-output"`
	PkgRunOpt             string              `json:"pkg-run-opt"`
	IeobcMacAddr          string              `json:"ieobc-mac-addr"`
}
type CPUUtil struct {
	RequestedApplicationUtil string `json:"requested-application-util"`
	ActualApplicationUtil    string `json:"actual-application-util"`
	CPUState                 string `json:"cpu-state"`
}
type MemoryUtil struct {
	MemoryAllocation string `json:"memory-allocation"`
	MemoryUsed       string `json:"memory-used"`
}
type Utilization struct {
	Name       string     `json:"name"`
	CPUUtil    CPUUtil    `json:"cpu-util"`
	MemoryUtil MemoryUtil `json:"memory-util"`
}
type StorageUtil struct {
	Name       string `json:"name"`
	Alias      string `json:"alias"`
	RdBytes    string `json:"rd-bytes"`
	RdRequests string `json:"rd-requests"`
	Errors     string `json:"errors"`
	WrBytes    string `json:"wr-bytes"`
	WrRequests string `json:"wr-requests"`
	Capacity   string `json:"capacity"`
	Available  string `json:"available"`
	Used       string `json:"used"`
	Usage      string `json:"usage"`
}
type StorageUtils struct {
	StorageUtil []StorageUtil `json:"storage-util"`
}
type AttachedDevice struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Alias string `json:"alias"`
}
type AttachedDevices struct {
	AttachedDevice []AttachedDevice `json:"attached-device"`
}
type CiscoIOSXEAppHostingOperApp struct {
	Name            string          `json:"name"`
	Details         Details         `json:"details"`
	Utilization     Utilization     `json:"utilization"`
	StorageUtils    StorageUtils    `json:"storage-utils"`
	AttachedDevices AttachedDevices `json:"attached-devices"`
	PkgPolicy       string          `json:"pkg-policy"`
}
