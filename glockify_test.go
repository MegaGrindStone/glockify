package glockify

import (
	"fmt"
	"log"
	"net/http"
)

var (
	wantWorkspace = Workspace{
		ID:   "5g3g57bt0cb2548e22e6l9cd",
		Name: "Clockify workspace",
	}
	dummyAPIKey = "dummy"
)

func workspaces(w http.ResponseWriter, r *http.Request) {
	if !checkAuthHeader(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, `
		[
		  {
			"hourlyRate": {
			  "amount": 0,
			  "currency": "USD"
			},
			"id": "%s",
			"imageUrl": "https://s3.eu-central-1.amazonaws.com/test/workspaceImg.png",
			"memberships": [
			  {
				"hourlyRate": {
				  "amount": 15,
				  "currency": "USD"
				},
				"membershipStatus": "ACTIVE",
				"membershipType": "WORKSPACE",
				"targetId": "5g3g57bt0cb2548e22e6l9cd",
				"userId": "6h1e49bf0cb8790e43d6c9ab"
			  }
			],
			"name": "%s",
			"workspaceSettings": {
			  "adminOnlyPages": [],
			  "automaticLock": {
				"changeDay": "WEDNESDAY",
				"dayOfMonth": "13",
				"firstDay": "MONDAY",
				"olderThanPeriod": "DAYS",
				"olderThanValue": "14",
				"type": "OLDER_THAN"
			  },
			  "canSeeTimeSheet": false,
			  "canSeeTracker": true,
			  "defaultBillableProjects": true,
			  "forceDescription": true,
			  "forceProjects": true,
			  "forceTags": false,
			  "forceTasks": false,
			  "lockTimeEntries": "2019-01-31T23:00:00Z",
			  "onlyAdminsCreateProject": false,
			  "onlyAdminsCreateTag": false,
			  "onlyAdminsCreateTask": true,
			  "onlyAdminsSeeAllTimeEntries": true,
			  "onlyAdminsSeeBillableRates": false,
			  "onlyAdminsSeeDashboard": false,
			  "onlyAdminsSeePublicProjectsEntries": false,
			  "projectFavorites": true,
			  "projectGroupingLabel": "client",
			  "projectPickerSpecialFilter": false,
			  "round": {
				"minutes": "15",
				"round": "Round up to"
			  },
			  "timeRoundingInReports": true,
			  "trackTimeDownToSecond": false,
			  "isProjectPublicByDefault": false,
			  "canSeeTracker": false,
			  "featureSubscriptionType": "ENTERPRISE_YEAR"
			}
		  }
		]`, wantWorkspace.ID, wantWorkspace.Name)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func checkAuthHeader(r *http.Request) bool {
	api := r.Header.Get("X-Api-Key")
	return api == dummyAPIKey
}
