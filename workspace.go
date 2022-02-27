package glockify

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Workspace struct {
	HourlyRate        HourlyRate        `json:"hourlyRate"`
	ID                string            `json:"id"`
	ImageURL          string            `json:"imageUrl"`
	Memberships       []Memberships     `json:"memberships"`
	Name              string            `json:"name"`
	WorkspaceSettings WorkspaceSettings `json:"workspaceSettings"`
}

type HourlyRate struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type Memberships struct {
	HourlyRate       HourlyRate `json:"hourlyRate"`
	MembershipStatus string     `json:"membershipStatus"`
	MembershipType   string     `json:"membershipType"`
	TargetID         string     `json:"targetId"`
	UserID           string     `json:"userId"`
}

type AutomaticLock struct {
	ChangeDay       string `json:"changeDay"`
	DayOfMonth      string `json:"dayOfMonth"`
	FirstDay        string `json:"firstDay"`
	OlderThanPeriod string `json:"olderThanPeriod"`
	OlderThanValue  string `json:"olderThanValue"`
	Type            string `json:"type"`
}

type Round struct {
	Minutes string `json:"minutes"`
	Round   string `json:"round"`
}

type WorkspaceSettings struct {
	AdminOnlyPages                     []interface{} `json:"adminOnlyPages"`
	AutomaticLock                      AutomaticLock `json:"automaticLock"`
	CanSeeTimeSheet                    string        `json:"canSeeTimeSheet"`
	CanSeeTracker                      string        `json:"canSeeTracker"`
	DefaultBillableProjects            string        `json:"defaultBillableProjects"`
	ForceDescription                   string        `json:"forceDescription"`
	ForceProjects                      string        `json:"forceProjects"`
	ForceTags                          string        `json:"forceTags"`
	ForceTasks                         string        `json:"forceTasks"`
	LockTimeEntries                    time.Time     `json:"lockTimeEntries"`
	OnlyAdminsCreateProject            string        `json:"onlyAdminsCreateProject"`
	OnlyAdminsCreateTag                string        `json:"onlyAdminsCreateTag"`
	OnlyAdminsCreateTask               string        `json:"onlyAdminsCreateTask"`
	OnlyAdminsSeeAllTimeEntries        string        `json:"onlyAdminsSeeAllTimeEntries"`
	OnlyAdminsSeeBillableRates         string        `json:"onlyAdminsSeeBillableRates"`
	OnlyAdminsSeeDashboard             string        `json:"onlyAdminsSeeDashboard"`
	OnlyAdminsSeePublicProjectsEntries string        `json:"onlyAdminsSeePublicProjectsEntries"`
	ProjectFavorites                   string        `json:"projectFavorites"`
	ProjectGroupingLabel               string        `json:"projectGroupingLabel"`
	ProjectPickerSpecialFilter         string        `json:"projectPickerSpecialFilter"`
	Round                              Round         `json:"round"`
	TimeRoundingInReports              string        `json:"timeRoundingInReports"`
	TrackTimeDownToSecond              string        `json:"trackTimeDownToSecond"`
	IsProjectPublicByDefault           string        `json:"isProjectPublicByDefault"`
	FeatureSubscriptionType            string        `json:"featureSubscriptionType"`
}

const workspacesPath = "/workspace"

func (g *Glockify) Workspaces(ctx context.Context) ([]Workspace, error) {
	res, err := g.get(ctx, nil, workspacesPath)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	result := make([]Workspace, 0)
	if err := json.Unmarshal(res, &result); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return result, nil
}
