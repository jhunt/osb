package api

type ProvisionSpec struct {
	ServiceID string `json:"service_id"`
	PlanID    string `json:"plan_id"`

	Context          map[string]interface{} `json:"context"`
	OrganizationGUID string                 `json:"organization_guid"`
	SpaceGUID        string                 `json:"space_guid"`

	Parameters map[string]interface{} `json:"parameters"`
}

type ProvisionStatus struct {
	InstanceID string `json:"-"`
	Status     string `json:"-"`

	DashboardURL string `json:"dashboard_url"`
	Operation    string `json:"operation"`
}

func (c *Client) Provision(id string, spec ProvisionSpec) (*ProvisionStatus, error) {
	if id == "" {
		id = randomID()
	}
	res, err := c.put("/v2/service_instances/"+id, spec)
	if err != nil {
		return nil, err
	}

	var status ProvisionStatus
	status.InstanceID = id

	switch res.StatusCode {
	case 200:
		status.Status = "already existed"
		return &status, c.parse(res, &status)

	case 201:
		status.Status = "provisioned"
		return &status, c.parse(res, &status)

	case 202:
		status.Status = "provisioning"
		return &status, c.parse(res, &status)
	}

	return nil, c.err(res)
}
