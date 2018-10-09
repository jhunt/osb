package api

import (
	"fmt"
)

type UnbindSpec struct {
	InstanceID string `json:"-"`
	BindingID  string `json:"-"`

	Context    map[string]interface{} `json:"context"`
	ServiceID  string                 `json:"service_id"`
	PlanID     string                 `json:"plan_id"`
	Parameters map[string]interface{} `json:"parameters"`
}

type UnbindStatus struct {
	InstanceID string `json:"-"`
	BindingID  string `json:"-"`
	Status     string `json:"-"`

	Operation string `json:"operation"`
}

func (c *Client) Unbind(spec UnbindSpec) (*UnbindStatus, error) {
	if spec.InstanceID == "" {
		return nil, fmt.Errorf("instance ID is required for unbinding")
	}

	if spec.BindingID == "" {
		return nil, fmt.Errorf("binding ID is required for unbinding")
	}

	res, err := c.del("/v2/service_instances/" + spec.InstanceID + "/service_bindings/" + spec.BindingID)
	if err != nil {
		return nil, err
	}

	var status UnbindStatus
	status.InstanceID = spec.InstanceID
	status.BindingID = spec.BindingID

	switch res.StatusCode {
	case 200:
		status.Status = "already unbound"
		return &status, c.parse(res, &status)

	case 201:
		status.Status = "unbound"
		return &status, c.parse(res, &status)

	case 202:
		status.Status = "unbinding"
		return &status, c.parse(res, &status)
	}

	return nil, c.err(res)
}
