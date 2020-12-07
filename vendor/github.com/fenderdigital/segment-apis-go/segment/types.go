package segment

import "time"

// Workspace defines the struct for the workspace object
type Workspace struct {
	Name        string     `json:"name,omitempty"`
	DisplayName string     `json:"display_name,omitempty"`
	ID          string     `json:"id,omitempty"`
	CreateTime  *time.Time `json:"create_time,omitempty"`
}

// Sources defines the struct for the sources object
type Sources struct {
	Sources []Source `json:"sources,omitempty"`
}

// Source defines the struct for the source object
type Source struct {
	Name          string        `json:"name,omitempty"`
	CatalogName   string        `json:"catalog_name,omitempty"`
	Parent        string        `json:"parent,omitempty"`
	WriteKeys     []string      `json:"write_keys,omitempty"`
	LibraryConfig LibraryConfig `json:"library_config,omitempty"`
	CreateTime    *time.Time    `json:"create_time,omitempty"`
}

// LibraryConfig contains information about a source's library
type LibraryConfig struct {
	MetricsEnabled       bool   `json:"metrics_enabled,omitempty"`
	RetryQueue           bool   `json:"retry_queue,omitempty"`
	CrossDomainIDEnabled bool   `json:"cross_domain_id_enabled,omitempty"`
	APIHost              string `json:"api_host,omitempty"`
}

// Destinations defines the struct for the destination object
type Destinations struct {
	Destinations []Destination `json:"destinations,omitempty"`
}

// Destination defines the struct for the destination object
type Destination struct {
	Name           string              `json:"name,omitempty"`
	Parent         string              `json:"parent,omitempty"`
	DisplayName    string              `json:"display_name,omitempty"`
	Enabled        bool                `json:"enabled,omitempty"`
	ConnectionMode string              `json:"connection_mode,omitempty"`
	Configs        []DestinationConfig `json:"config,omitempty"`
	CreateTime     *time.Time          `json:"create_time,omitempty"`
	UpdateTime     *time.Time          `json:"update_time,omitempty"`
}

// TrackingPlans defines the struct for the tracking plan object
type TrackingPlans struct {
	TrackingPlans []TrackingPlan `json:"tracking_plans,omitempty"`
}

// TrackingPlan defines the struct for the destination object
type TrackingPlan struct {
	Name        string     `json:"name,omitempty"`
	DisplayName string     `json:"display_name,omitempty"`
	Rules       Rules      `json:"rules,omitempty"`
	CreateTime  *time.Time `json:"create_time,omitempty"`
	UpdateTime  *time.Time `json:"update_time,omitempty"`
}

// Rules contains the information about all the rules of a tracking plan
type Rules struct {
	Global         Rule          `json:"global,omitempty"`
	Events         []Event       `json:"events,omitempty"`
	Identify       Rule          `json:"identify,omitempty"`
	Group          Rule          `json:"group,omitempty"`
	IdentifyTraits []interface{} `json:"identify_traits"`
	GroupTraits    []interface{} `json:"group_traits"`
}

// Rule contains the information about the rule definition
type Rule struct {
	Description string                 `json:"description,omitempty"`
	Enum        []interface{}          `json:"enum,omitempty"`
	Labels      map[string]interface{} `json:"labels,omitempty"`
	Pattern     interface{}            `json:"pattern,omitempty"`
	Properties  map[string]Rule        `json:"properties,omitempty"`
	Required    []string               `json:"required,omitempty"`
	Type        interface{}            `json:"type,omitempty"`
	Schema      string                 `json:"$schema,omitempty"`
}

// Event contains the rules for each tracking event
type Event struct {
	Name        string `json:"name,omitempty"`
	Version     int    `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	Rules       Rule   `json:"rules,omitempty"`
}

// DestinationConfig contains information about how a Destination is configured
type DestinationConfig struct {
	Name        string      `json:"name,omitempty"`
	DisplayName string      `json:"display_name,omitempty"`
	Value       interface{} `json:"value,omitempty"`
	Type        string      `json:"type,omitempty"`
}

// UpdateMask contains information for updating Destinations
type UpdateMask struct {
	Paths []string `json:"paths,omitempty"`
}

type sourceCreateRequest struct {
	Source Source `json:"source,omitempty"`
}

type destinationCreateRequest struct {
	Destination Destination `json:"destination,omitempty"`
}

type destinationUpdateRequest struct {
	Destination Destination `json:"destination,omitempty"`
	UpdateMask  UpdateMask  `json:"update_mask,omitempty"`
}

type trackingPlanCreateRequest struct {
	TrackingPlan TrackingPlan `json:"tracking_plan,omitempty"`
}

type trackingPlanUpdateRequest struct {
	TrackingPlan TrackingPlan `json:"tracking_plan,omitempty"`
	UpdateMask   UpdateMask   `json:"update_mask,omitempty"`
}

type trackingPlanSourceConnection struct {
	SourceName     string `json:"source_name,omitempty"`
	TrackingPlanID string `json:"tracking_plan_id,omitempty"`
}

type trackingPlanSourceConnections struct {
	Connections []trackingPlanSourceConnection `json:"connections,omitempty"`
}
