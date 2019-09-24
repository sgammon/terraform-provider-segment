package segment

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

//ListTrackingPlans returns all tracking plans for a workspace
func (c *Client) ListTrackingPlans() (TrackingPlans, error) {
	var p TrackingPlans
	data, err := c.doRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s", WorkspacesEndpoint, c.workspace, TrackingPlanEndpoint),
		nil)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(data, &p)
	if err != nil {
		return p, errors.Wrap(err, "failed to unmarshal tracking plans response")
	}

	return p, nil
}

// GetTrackingPlan returns information about a tracking plan
func (c *Client) GetTrackingPlan(planName string) (TrackingPlan, error) {
	var p TrackingPlan
	data, err := c.doRequest(http.MethodGet,
		fmt.Sprintf("%s/%s/%s/%s", WorkspacesEndpoint, c.workspace, TrackingPlanEndpoint, planName),
		nil)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(data, &p)
	if err != nil {
		return p, errors.Wrap(err, "failed to unmarshal tracking plan response")
	}
	return p, nil
}

// CreateTrackingPlan creates a new tracking plan
func (c *Client) CreateTrackingPlan(displayName string, rules Rules) (TrackingPlan, error) {
	var p TrackingPlan
	plan := trackingPlanCreateRequest{
		TrackingPlan: TrackingPlan{
			DisplayName: displayName,
			Rules:       rules,
		}}
	data, err := c.doRequest(http.MethodPost,
		fmt.Sprintf("%s/%s/%s/",
			WorkspacesEndpoint, c.workspace, TrackingPlanEndpoint),
		plan)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(data, &p)
	if err != nil {
		return p, errors.Wrap(err, "failed to unmarshall tracking plan response")
	}

	return p, nil

}

// UpdateTrackingPlan updates an existing tracking plan
func (c *Client) UpdateTrackingPlan(planName string, paths []string, updatedPlan TrackingPlan) (TrackingPlan, error) {
	var p TrackingPlan
	req := trackingPlanUpdateRequest{TrackingPlan: updatedPlan, UpdateMask: UpdateMask{Paths: paths}}
	data, err := c.doRequest(http.MethodPut, fmt.Sprintf("%s/%s/%s/%s/", WorkspacesEndpoint, c.workspace, TrackingPlanEndpoint, planName), req)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(data, &p)
	if err != nil {
		return p, errors.Wrap(err, "failed to unmarshal tracking plan response")
	}

	return p, nil
}

// CreateTrackingPlanSourceConnection connects a source to a tracking plan
func (c *Client) CreateTrackingPlanSourceConnection(planName string, srcName string) (trackingPlanSourceConnection, error) {
	var p trackingPlanSourceConnection
	req := trackingPlanSourceConnection{SourceName: srcName}
	endpoint := fmt.Sprintf("%s/%s/%s/%s/%s/", WorkspacesEndpoint, c.workspace, TrackingPlanEndpoint,
		planName, TrackingPlanSourceConnectionEndpoint)
	data, err := c.doRequest(http.MethodPost, endpoint, req)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(data, &p)
	if err != nil {
		return p, errors.Wrap(err, "failed to unmarshal tracking plan source connection response")
	}

	return p, nil
}

// ListTrackingPlanSourceConnections lists the source connections for a tracking plan
func (c *Client) ListTrackingPlanSourceConnections(planName string) (trackingPlanSourceConnections, error) {
	var p trackingPlanSourceConnections
	data, err := c.doRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s/%s/%s/", WorkspacesEndpoint, c.workspace, TrackingPlanEndpoint, planName, TrackingPlanSourceConnectionEndpoint), nil)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(data, &p)
	if err != nil {
		return p, errors.Wrap(err, "failed to unmarshal tracking plan response")
	}
	return p, nil
}

// DeleteTrackingPlan deletes a destination for a source from the workspace
func (c *Client) DeleteTrackingPlan(planName string) error {
	_, err := c.doRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s/%s/", WorkspacesEndpoint, c.workspace, TrackingPlanEndpoint, planName), nil)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTrackingPlanSourceConnection deletes a source connection for a tracking plan
func (c *Client) DeleteTrackingPlanSourceConnection(planName string, srcName string) error {
	_, err := c.doRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s/%s/%s/%s/", WorkspacesEndpoint, c.workspace, TrackingPlanEndpoint,
		planName, TrackingPlanSourceConnectionEndpoint, srcName), nil)
	if err != nil {
		return err
	}

	return nil
}
