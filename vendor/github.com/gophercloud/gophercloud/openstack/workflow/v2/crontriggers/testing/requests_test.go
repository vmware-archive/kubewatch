package testing

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/workflow/v2/crontriggers"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreateCronTrigger(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/cron_triggers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusCreated)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
			{
				"created_at": "2018-09-12 15:48:18",
				"first_execution_time": "2018-09-12 17:48:00",
				"id": "0520ffd8-f7f1-4f2e-845b-55d953a1cf46",
				"name": "crontrigger",
				"next_execution_time": "2018-09-12 17:48:00",
				"pattern": "0 0 1 1 *",
				"project_id": "778c0f25df0d492a9a868ee9e2fbb513",
				"remaining_executions": 42,
				"scope": "private",
				"updated_at": null,
				"workflow_id": "604a3a1e-94e3-4066-a34a-aa56873ef236",
				"workflow_input": "{\"msg\": \"hello\"}",
				"workflow_name": "workflow_echo",
				"workflow_params": "{\"msg\": \"world\"}"
			}
		`)
	})

	firstExecution := time.Date(2018, time.September, 12, 17, 48, 0, 0, time.UTC)
	opts := &crontriggers.CreateOpts{
		WorkflowID:         "w1",
		Name:               "trigger",
		FirstExecutionTime: &firstExecution,
		WorkflowParams: map[string]interface{}{
			"msg": "world",
		},
		WorkflowInput: map[string]interface{}{
			"msg": "hello",
		},
	}

	actual, err := crontriggers.Create(fake.ServiceClient(), opts).Extract()
	if err != nil {
		t.Fatalf("Unable to create cron trigger: %v", err)
	}

	expected := &crontriggers.CronTrigger{
		ID:                  "0520ffd8-f7f1-4f2e-845b-55d953a1cf46",
		Name:                "crontrigger",
		Pattern:             "0 0 1 1 *",
		ProjectID:           "778c0f25df0d492a9a868ee9e2fbb513",
		RemainingExecutions: 42,
		Scope:               "private",
		WorkflowID:          "604a3a1e-94e3-4066-a34a-aa56873ef236",
		WorkflowName:        "workflow_echo",
		WorkflowParams: map[string]interface{}{
			"msg": "world",
		},
		WorkflowInput: map[string]interface{}{
			"msg": "hello",
		},
		CreatedAt:          time.Date(2018, time.September, 12, 15, 48, 18, 0, time.UTC),
		FirstExecutionTime: &firstExecution,
		NextExecutionTime:  &firstExecution,
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}

func TestDeleteCronTrigger(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/cron_triggers/1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})

	res := crontriggers.Delete(fake.ServiceClient(), "1")
	th.AssertNoErr(t, res.Err)
}

func TestGetCronTrigger(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/cron_triggers/1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
			{
				"created_at": "2018-09-12 15:48:18",
				"first_execution_time": "2018-09-12 17:48:00",
				"id": "0520ffd8-f7f1-4f2e-845b-55d953a1cf46",
				"name": "crontrigger",
				"next_execution_time": "2018-09-12 17:48:00",
				"pattern": "0 0 1 1 *",
				"project_id": "778c0f25df0d492a9a868ee9e2fbb513",
				"remaining_executions": 42,
				"scope": "private",
				"updated_at": null,
				"workflow_id": "604a3a1e-94e3-4066-a34a-aa56873ef236",
				"workflow_input": "{\"msg\": \"hello\"}",
				"workflow_name": "workflow_echo",
				"workflow_params": "{\"msg\": \"world\"}"
			}
		`)
	})
	actual, err := crontriggers.Get(fake.ServiceClient(), "1").Extract()
	if err != nil {
		t.Fatalf("Unable to get cron trigger: %v", err)
	}

	firstExecution := time.Date(2018, time.September, 12, 17, 48, 0, 0, time.UTC)

	expected := &crontriggers.CronTrigger{
		ID:                  "0520ffd8-f7f1-4f2e-845b-55d953a1cf46",
		Name:                "crontrigger",
		Pattern:             "0 0 1 1 *",
		ProjectID:           "778c0f25df0d492a9a868ee9e2fbb513",
		RemainingExecutions: 42,
		Scope:               "private",
		WorkflowID:          "604a3a1e-94e3-4066-a34a-aa56873ef236",
		WorkflowName:        "workflow_echo",
		WorkflowParams: map[string]interface{}{
			"msg": "world",
		},
		WorkflowInput: map[string]interface{}{
			"msg": "hello",
		},
		CreatedAt:          time.Date(2018, time.September, 12, 15, 48, 18, 0, time.UTC),
		FirstExecutionTime: &firstExecution,
		NextExecutionTime:  &firstExecution,
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}
