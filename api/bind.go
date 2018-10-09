package api

import (
	"fmt"
)

type BindSpec struct {
	InstanceID string `json:"-"`
	BindingID  string `json:"-"`

	Context    map[string]interface{} `json:"context"`
	ServiceID  string                 `json:"service_id"`
	PlanID     string                 `json:"plan_id"`
	Parameters map[string]interface{} `json:"parameters"`
}

type BindStatus struct {
	InstanceID string `json:"-"`
	BindingID  string `json:"-"`
	Status     string `json:"-"`

	Operation       string                 `json:"operation"`
	Credentials     map[string]interface{} `json:"credentials"`
	SyslogDrainURL  string                 `json:"syslog_drain_url"`
	RouteServiceURL string                 `json:"route_service_url"`
	VolumeMounts    []struct {
		Driver       string `json:"driver"`
		ContainerDir string `json:"container_dir"`
		Mode         string `json:"mode"`
		DeviceType   string `json:"device_type"`
		Device       struct {
			VolumeID    string                 `json:"volume_id"`
			MountConfig map[string]interface{} `json:"mount_config"`
		} `json:"device"`
	} `json:"volume_mounts"`
}

func (c *Client) Bind(spec BindSpec) (*BindStatus, error) {
	if spec.InstanceID == "" {
		return nil, fmt.Errorf("instance ID is required for binding")
	}

	if spec.BindingID == "" {
		spec.BindingID = randomID()
	}

	res, err := c.put("/v2/service_instances/"+spec.InstanceID+"/service_bindings/"+spec.BindingID, spec)
	if err != nil {
		return nil, err
	}

	var status BindStatus
	status.InstanceID = spec.InstanceID
	status.BindingID = spec.BindingID

	switch res.StatusCode {
	case 200:
		status.Status = "already bound"
		return &status, c.parse(res, &status)

	case 201:
		status.Status = "bound"
		return &status, c.parse(res, &status)

	case 202:
		status.Status = "binding"
		return &status, c.parse(res, &status)
	}

	return nil, c.err(res)
}
