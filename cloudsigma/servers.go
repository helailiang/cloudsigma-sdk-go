package cloudsigma

import (
	"context"
	"fmt"
	"net/http"
)

const serversBasePath = "servers"

//server status
const (
	ServerStop ServerStatus = "stopped"
	ServerRunning ServerStatus = "running"
)

//IP Type
const (
	IPTypeStatic = "static"
	IPTypeDHCP  = "dhcp"
	IPTypeManual = "manual"
)

type ServerStatus string
// ServersService handles communication with the servers related methods of
// the CloudSigma API.
//
// CloudSigma API docs: https://cloudsigma-docs.readthedocs.io/en/latest/servers.html
type ServersService service

// Server represents a CloudSigma server.
type Server struct {
	//AllocationPool     string `json:"allocation_pool"`
	AutoStart          bool                   `json:"auto_start,omitempty"`
	Context            bool                   `json:"context,omitempty"`
	CPU                int                    `json:"cpu,omitempty"`
	//CPUModel           string		`json:"cpu_model"`
	CPUType            string                 `json:"cpu_type,omitempty"`
	CPUsInsteadOfCores bool                   `json:"cpus_instead_of_cores,omitempty"`
	Drives             []ServerDrive          `json:"drives,omitempty"`
	EnableNuma         bool                   `json:"enable_numa,omitempty"`
	EnclavePageCaches  []EnclavePageCache     `json:"epcs,omitempty"`
	Hypervisor         string                 `json:"hypervisor,omitempty"`
	Memory             int                    `json:"mem,omitempty"`
	Meta               map[string]interface{} `json:"meta,omitempty"`
	Name               string                 `json:"name,omitempty"`
	NICs               []ServerNIC            `json:"nics"`
	Owner              *ResourceLink          `json:"owner,omitempty"`
	PublicKeys         []Keypair              `json:"pubkeys,omitempty"`
	ResourceURI        string                 `json:"resource_uri,omitempty"`
	Runtime            *ServerRuntime         `json:"runtime,omitempty"`
	SMP                int                    `json:"smp,omitempty"`
	Status             string                 `json:"status,omitempty"`
	Tags               []Tag                  `json:"tags,omitempty"`
	UUID               string                 `json:"uuid,omitempty"`
	VNCPassword        string                 `json:"vnc_password,omitempty"`
}

// ServerDrive represents a CloudSigma drive attached to a server.
type ServerDrive struct {
	BootOrder  int    `json:"boot_order,omitempty"`
	DevChannel string `json:"dev_channel,omitempty"`
	Device     string `json:"device,omitempty"`
	Drive      *Drive `json:"drive,omitempty"`
}

// EnclavePageCache represents a protected memory region for enclaves in a server.
type EnclavePageCache struct {
	Size int `json:"size,omitempty"`
}

// ServerNIC represents a CloudSigma network interface card attached to a server.
type ServerNIC struct {
	BootOrder        int                    `json:"boot_order,omitempty"`
	FirewallPolicy   *FirewallPolicy        `json:"firewall_policy,omitempty"`
	IP4Configuration *ServerIPConfiguration `json:"ip_v4_conf,omitempty"`
	IP6Configuration *ServerIPConfiguration `json:"ip_v6_conf,omitempty"`
	MACAddress       string                 `json:"mac,omitempty"`
	Model            string                 `json:"model,omitempty"`
	VLAN             *VLAN                  `json:"vlan,omitempty"`
}

// ServerRuntime represents a CloudSigma server runtime information.
type ServerRuntime struct {
	RuntimeNICs []ServerRuntimeNIC `json:"nics,omitempty"`
}

// ServerRuntimeNIC represents a CloudSigma server's network interface card runtime.
type ServerRuntimeNIC struct {
	InterfaceType string          `json:"interface_type,omitempty"`
	IPv4          ServerRuntimeIP `json:"ip_v4,omitempty"`
	IPv6          ServerRuntimeIP `json:"ip_v6,omitempty"`
}

// ServerRuntimeIP represents a CloudSigma server's runtime IP configuration.
type ServerRuntimeIP struct {
	ResourceURI string `json:"resource_uri,omitempty"`
	UUID        string `json:"uuid,omitempty"`
}

// ServerIPConfiguration represents a CloudSigma public IP configuration.
type ServerIPConfiguration struct {
	Type      string `json:"conf,omitempty"`
	IPAddress *IP    `json:"ip,omitempty"`
}

// ServerAction represents a CloudSigma server action.
type ServerAction struct {
	Action string `json:"action,omitempty"`
	Result string `json:"result,omitempty"`
	UUID   string `json:"uuid,omitempty"`
	Job    string `json:"job,omitempty"`
	Tags   []string `json:"tags,omitempty"`
}

// ServerCreateRequest represents a request to create a server.
type ServerCreateRequest struct {
	Servers []Server `json:"objects"`
}

type ServerDetails struct {
	Servers []Server `json:"objects"`
}

// ServerUpdateRequest represents a request to update a server.
type  ServerUpdateRequest struct {
	*Server
}

type serversRootSimple struct {
	Servers []ServerSimple `json:"objects"`
}

type ServerSimple struct {
	//AllocationPool     string `json:"allocation_pool"`
	AutoStart          bool                   `json:"auto_start,omitempty"`
	Context            bool                   `json:"context,omitempty"`
	CPU                int                    `json:"cpu,omitempty"`
	//CPUModel           string		`json:"cpu_model"`
	CPUType            string                 `json:"cpu_type,omitempty"`
	CPUsInsteadOfCores bool                   `json:"cpus_instead_of_cores,omitempty"`
	Drives             []ServerDrive          `json:"drives,omitempty"`
	EnableNuma         bool                   `json:"enable_numa,omitempty"`
	EnclavePageCaches  []EnclavePageCache     `json:"epcs,omitempty"`
	Memory             int                    `json:"mem,omitempty"`
	Name               string                 `json:"name,omitempty"`
	NICs               []ServerNIC            `json:"nics,omitempty"`
	ResourceURI        string                 `json:"resource_uri,omitempty"`
	SMP                int                    `json:"smp,omitempty"`
	Status             string                 `json:"status,omitempty"`
	UUID               string                 `json:"uuid,omitempty"`
	VNCPassword        string                 `json:"vnc_password,omitempty"`
}

type serversRoot struct {
	Meta    *Meta    `json:"meta,omitempty"`
	Servers []Server `json:"objects"`
}

type ServerCloneAndStartRequest struct {
	Name string `json:"name,omitempty"`
	AutoStart bool `json:"auto_start,omitempty"`
	Count uint `json:"count"`
}

type ServerCloneRequest struct {
	Name string `json:"name,omitempty"`
	RandomVncPassword bool `json:"random_vnc_password,omitempty"`
}


// List provides a detailed list of servers to which the authenticated user
// has access.
//
// CloudSigma API docs: https://cloudsigma-docs.readthedocs.io/en/latest/servers.html#detailed-listing
func (s *ServersService) List(ctx context.Context) ([]Server, *Response, error) {
	path := fmt.Sprintf("%v/detail/", serversBasePath)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(serversRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Servers, resp, nil
}

// Get provides detailed information for server identified by uuid.
//
// CloudSigma API docs: https://cloudsigma-docs.readthedocs.io/en/latest/servers.html#server-runtime-and-server-details
func (s *ServersService) Get(ctx context.Context, uuid string) (*Server, *Response, error) {
	if uuid == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/", serversBasePath, uuid)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	server := new(Server)
	resp, err := s.client.Do(ctx, req, server)
	if err != nil {
		return nil, resp, err
	}

	return server, resp, nil
}
func (s *ServersService) GetDetail(ctx context.Context) (*serversRootSimple, *Response, error) {
	//if uuid == "" {
	//	return nil, nil, ErrEmptyArgument
	//}

	path := fmt.Sprintf("%v/detail/?limit=0&order_by=uuid", serversBasePath)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	servers := new(serversRootSimple)
	resp, err := s.client.Do(ctx, req, servers)
	if err != nil {
		return nil, resp, err
	}

	return servers, resp, nil
}
// Create makes a new virtual server with given payload.
//
// CloudSigma API docs: https://cloudsigma-docs.readthedocs.io/en/latest/servers.html#creating
func (s *ServersService) Create(ctx context.Context, createRequest *ServerCreateRequest) ([]Server, *Response, error) {
	if createRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/", serversBasePath)

	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(ServerCreateRequest)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Servers, resp, nil
}

// Update edits a server identified by uuid. Used also for attaching NIC’s
// and drives to servers. Note that if a server is running, only name, meta,
// and tags fields can be changed, and all other changes to the definition
// of a running server will be ignored.
//
// CloudSigma API docs: https://cloudsigma-docs.readthedocs.io/en/latest/servers.html#editing
func (s *ServersService) Update(ctx context.Context, uuid string, updateRequest *ServerUpdateRequest) (*Server, *Response, error) {
	if uuid == "" {
		return nil, nil, ErrEmptyArgument
	}
	if updateRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v/", serversBasePath, uuid)

	// by update UUID must be empty
	updateRequest.UUID = ""

	req, err := s.client.NewRequest(http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	server := new(Server)
	resp, err := s.client.Do(ctx, req, server)
	if err != nil {
		return nil, resp, err
	}

	return server, resp, nil
}

// Delete removes a single server identified by uuid.
//
// CloudSigma API docs: https://cloudsigma-docs.readthedocs.io/en/latest/servers.html#deleting
func (s *ServersService) Delete(ctx context.Context, uuid string) (*Response, error) {
	if uuid == "" {
		return nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/", serversBasePath, uuid)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Start sends 'start' action and starts a server with specific uuid.
//
// CloudSigma API docs: https://cloudsigma-docs.readthedocs.io/en/latest/servers.html#start
func (s *ServersService) Start(ctx context.Context, uuid string) (*ServerAction, *Response, error) {
	return s.doAction(ctx, uuid, "start", nil)
}

// Stop sends 'stop' action and stops a server with specific uuid.
//
// CloudSigma API docs: https://cloudsigma-docs.readthedocs.io/en/latest/servers.html#stop
func (s *ServersService) Stop(ctx context.Context, uuid string) (*ServerAction, *Response, error) {
	return s.doAction(ctx, uuid, "stop",nil)
}

// Shutdown sends an ACPI shutdowns to a server with specific UUID for a minute.
//
// CloudSigma API docs: https://cloudsigma-docs.readthedocs.io/en/latest/servers.html#acpi-shutdown
func (s *ServersService) Shutdown(ctx context.Context, uuid string) (*ServerAction, *Response, error) {
	return s.doAction(ctx, uuid, "shutdown", nil)
}

func (s *ServersService) doAction(ctx context.Context, uuid, action string, body interface{}) (*ServerAction, *Response, error) {
	if uuid == "" || action == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/action/?do=%v", serversBasePath, uuid, action)

	req, err := s.client.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, nil, err
	}

	serverAction := new(ServerAction)
	resp, err := s.client.Do(ctx, req, serverAction)
	if err != nil {
		return nil, resp, err
	}

	return serverAction, resp, nil
}

func (s *ServersService) Clone(ctx context.Context, uuid string, request ServerCloneRequest) (*Server, *Response, error) {
	if uuid == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/action/?do=clone", serversBasePath, uuid)

	req, err := s.client.NewRequest(http.MethodPost, path, request)
	if err != nil {
		return nil, nil, err
	}

	serverAction := new(Server)
	resp, err := s.client.Do(ctx, req, serverAction)
	if err != nil {
		return nil, resp, err
	}

	return serverAction, resp, nil
	//return s.doAction(ctx, uuid, "clone", request)
}


func (s *ServersService) CloneAndStart(ctx context.Context, uuid string, request ServerCloneAndStartRequest) (*ServerAction, *Response, error) {
	return s.doAction(ctx, uuid, "bulk_clone_start", request)
}

//func (s *ServersService) Stop(ctx context.Context, uuid string) (*ServerAction, *Response, error) {
//	return s.doAction(ctx, uuid, "stop", nil)
//}
