package api

import (
	"fmt"
)

type DeprovisionSpec struct {
	InstanceID string
	ServiceID string
	PlanID string
}

type DeprovisionStatus struct {
	Status string `json:"-"`

	Operation string `json:"operation"`
}

func (c *Client) Deprovision(spec DeprovisionSpec) (*DeprovisionStatus, error) {
	if spec.InstanceID == "" {
		return nil, fmt.Errorf("instance ID is required for deprovisioning")
	}

	if spec.ServiceID == "" {
		spec.ServiceID = "oops-unknown-service-id"
	}
	if spec.PlanID == "" {
		spec.PlanID = "oops-unknown-plan-id"
	}

	res, err := c.del("/v2/service_instances/" + spec.InstanceID+"?service_id="+spec.ServiceID+"&plan_id="+spec.PlanID)
	if err != nil {
		return nil, err
	}

	var status DeprovisionStatus

	switch res.StatusCode {
	case 410:
		fallthrough
	case 200:
		status.Status = "deprovisioned"
		return &status, nil

	case 202:
		status.Status = "deprovisioning"
		return &status, c.parse(res, &status)
	}

	return nil, c.err(res)
}
