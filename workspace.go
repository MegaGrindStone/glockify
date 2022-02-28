package glockify

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// WorkspaceNode manipulating Workspace resource.
type WorkspaceNode struct {
	endpoint Endpoint
	apiKey   string
}

// Workspace wraps Clockify's workspace resource.
// See: https://clockify.me/developers-api#tag-Workspace
type Workspace struct {
	ID                string            `json:"id,omitempty"`
	Name              string            `json:"name,omitempty"`
	HourlyRate        HourlyRate        `json:"hourlyRate,omitempty"`
	ImageURL          string            `json:"imageUrl,omitempty"`
	Memberships       []Memberships     `json:"memberships,omitempty"`
	WorkspaceSettings WorkspaceSettings `json:"workspaceSettings,omitempty"`

	Client ClientNode `json:"-"`
}

// HourlyRate see: https://clockify.me/developers-api#tag-Workspace
type HourlyRate struct {
	Amount   int    `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

// Memberships see: https://clockify.me/developers-api#tag-Workspace
type Memberships struct {
	HourlyRate       HourlyRate `json:"hourlyRate,omitempty"`
	MembershipStatus string     `json:"membershipStatus,omitempty"`
	MembershipType   string     `json:"membershipType,omitempty"`
	TargetID         string     `json:"targetId,omitempty"`
	UserID           string     `json:"userId,omitempty"`
}

// AutomaticLock see: https://clockify.me/developers-api#tag-Workspace
type AutomaticLock struct {
	ChangeDay       string `json:"changeDay,omitempty"`
	DayOfMonth      string `json:"dayOfMonth,omitempty"`
	FirstDay        string `json:"firstDay,omitempty"`
	OlderThanPeriod string `json:"olderThanPeriod,omitempty"`
	OlderThanValue  string `json:"olderThanValue,omitempty"`
	Type            string `json:"type,omitempty"`
}

// Round see: https://clockify.me/developers-api#tag-Workspace
type Round struct {
	Minutes string `json:"minutes,omitempty"`
	Round   string `json:"round,omitempty"`
}

// WorkspaceSettings see: https://clockify.me/developers-api#tag-Workspace
type WorkspaceSettings struct {
	AdminOnlyPages                     []interface{} `json:"adminOnlyPages,omitempty"`
	AutomaticLock                      AutomaticLock `json:"automaticLock,omitempty"`
	CanSeeTimeSheet                    bool          `json:"canSeeTimeSheet,omitempty"`
	CanSeeTracker                      bool          `json:"canSeeTracker,omitempty"`
	DefaultBillableProjects            bool          `json:"defaultBillableProjects,omitempty"`
	ForceDescription                   bool          `json:"forceDescription,omitempty"`
	ForceProjects                      bool          `json:"forceProjects,omitempty"`
	ForceTags                          bool          `json:"forceTags,omitempty"`
	ForceTasks                         bool          `json:"forceTasks,omitempty"`
	LockTimeEntries                    time.Time     `json:"lockTimeEntries,omitempty"`
	OnlyAdminsCreateProject            bool          `json:"onlyAdminsCreateProject,omitempty"`
	OnlyAdminsCreateTag                bool          `json:"onlyAdminsCreateTag,omitempty"`
	OnlyAdminsCreateTask               bool          `json:"onlyAdminsCreateTask,omitempty"`
	OnlyAdminsSeeAllTimeEntries        bool          `json:"onlyAdminsSeeAllTimeEntries,omitempty"`
	OnlyAdminsSeeBillableRates         bool          `json:"onlyAdminsSeeBillableRates,omitempty"`
	OnlyAdminsSeeDashboard             bool          `json:"onlyAdminsSeeDashboard,omitempty"`
	OnlyAdminsSeePublicProjectsEntries bool          `json:"onlyAdminsSeePublicProjectsEntries,omitempty"`
	ProjectFavorites                   bool          `json:"projectFavorites,omitempty"`
	ProjectGroupingLabel               string        `json:"projectGroupingLabel,omitempty"`
	ProjectPickerSpecialFilter         bool          `json:"projectPickerSpecialFilter,omitempty"`
	Round                              Round         `json:"round,omitempty"`
	TimeRoundingInReports              bool          `json:"timeRoundingInReports,omitempty"`
	TrackTimeDownToSecond              bool          `json:"trackTimeDownToSecond,omitempty"`
	IsProjectPublicByDefault           bool          `json:"isProjectPublicByDefault,omitempty"`
	FeatureSubscriptionType            string        `json:"featureSubscriptionType,omitempty"`
}

// All get all Workspace resource.
func (w *WorkspaceNode) All(ctx context.Context) ([]Workspace, error) {
	endpoint := fmt.Sprintf("%s/workspaces", w.endpoint.Base)
	res, err := get(ctx, w.apiKey, nil, endpoint)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	result := make([]Workspace, 0)
	if err := json.Unmarshal(res, &result); err != nil {
		if jErr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("unmarshal field %v of type %v", jErr.Field, jErr.Type)
		}
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	for i := 0; i < len(result); i++ {
		result[i].Client = ClientNode{
			workspaceID:  result[i].ID,
			baseEndpoint: w.endpoint.Base,
			apiKey:       w.apiKey,
		}
	}
	return result, nil
}
