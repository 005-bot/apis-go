// Package domain defines the wire-contract DTOs shared between 005-bot
// services: the outage payload published by monitor-go to the Redis outages
// channel and consumed by tg-bot-go.
package domain

import (
	"strings"
	"time"
)

// ResourceType identifies the utility service affected by an outage.
type ResourceType string

const (
	ResourceTypeColdWater   ResourceType = "Холодное водоснабжение"
	ResourceTypeHotWater    ResourceType = "Горячее водоснабжение"
	ResourceTypeElectricity ResourceType = "Электроснабжение"
	ResourceTypeGas         ResourceType = "Газоснабжение"
	ResourceTypeHeating     ResourceType = "Теплоснабжение"
)

// DetectResourceType matches a free-form resource string against the known
// resource types using a case-insensitive substring match and returns nil
// when none matches.
func DetectResourceType(resource string) *ResourceType {
	for _, rt := range []ResourceType{
		ResourceTypeColdWater,
		ResourceTypeHotWater,
		ResourceTypeElectricity,
		ResourceTypeGas,
		ResourceTypeHeating,
	} {
		if strings.Contains(strings.ToLower(resource), strings.ToLower(string(rt))) {
			return &rt
		}
	}
	return nil
}

// OrganizationInfo describes the utility organization serving the area.
type OrganizationInfo struct {
	ResourceType *ResourceType `json:"resource_type"`
	Resource     string        `json:"resource"`
	Organization string        `json:"organization"`
	Phones       []string      `json:"phones"`
}

// Street is a street name with an optional list of affected buildings.
type Street struct {
	Name      string   `json:"name"`
	Buildings []string `json:"buildings,omitempty"`
}

// String renders the street name, appending the buildings when present.
func (s Street) String() string {
	if len(s.Buildings) == 0 {
		return s.Name
	}
	return s.Name + " " + strings.Join(s.Buildings, ", ")
}

// Reason describes why the outage occurred.
type Reason struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// WaterDelivery is a scheduled water delivery associated with an outage.
type WaterDelivery struct {
	Street    string `json:"street"`
	Buildings string `json:"buildings"`
	TimeStart string `json:"time_start"`
	TimeEnd   string `json:"time_end"`
}

// OutageDetails holds the affected streets, the reason, scheduled water
// deliveries and additional comments.
type OutageDetails struct {
	Streets         []Street        `json:"streets"`
	Reason          *Reason         `json:"reason,omitempty"`
	WaterDeliveries []WaterDelivery `json:"water_deliveries,omitempty"`
	Comments        string          `json:"comments,omitempty"`
}

// Address joins the affected street addresses with newlines.
func (d OutageDetails) Address() string {
	parts := make([]string, len(d.Streets))
	for i, s := range d.Streets {
		parts[i] = s.String()
	}
	return strings.Join(parts, "\n")
}

// Outage is the wire contract published to the outages channel. The
// [time.Time] values serialize as RFC3339Nano via encoding/json.
type Outage struct {
	Area             string           `json:"area"`
	OrganizationInfo OrganizationInfo `json:"organization_info"`
	Details          OutageDetails    `json:"details"`
	Period           []time.Time      `json:"period"`
}
