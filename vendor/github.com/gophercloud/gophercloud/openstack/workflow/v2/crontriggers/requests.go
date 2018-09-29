package crontriggers

import (
	"time"

	"github.com/gophercloud/gophercloud"
)

// CreateOptsBuilder allows extension to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToCronTriggerCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies parameters used to create a cron trigger.
type CreateOpts struct {
	// Name is the cron trigger name.
	Name string `json:"name"`

	// Pattern is a Unix crontab patterns format to execute the workflow.
	Pattern string `json:"pattern"`

	// RemainingExecutions sets the number of executions for the trigger.
	RemainingExecutions int `json:"remaining_executions,omitempty"`

	// WorkflowID is the unique id of the workflow.
	WorkflowID string `json:"workflow_id,omitempty" or:"WorkflowName"`

	// WorkflowName is the name of the workflow.
	// It is recommended to refer to workflow by the WorkflowID parameter instead of WorkflowName.
	WorkflowName string `json:"workflow_name,omitempty" or:"WorkflowID"`

	// WorkflowParams defines workflow type specific parameters.
	WorkflowParams map[string]interface{} `json:"workflow_params,omitempty"`

	// WorkflowInput defines workflow input values.
	WorkflowInput map[string]interface{} `json:"workflow_input,omitempty"`

	// FirstExecutionTime defines the first execution time of the trigger.
	FirstExecutionTime *time.Time `json:"-"`
}

// ToCronTriggerCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToCronTriggerCreateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.FirstExecutionTime != nil {
		b["first_execution_time"] = opts.FirstExecutionTime.Format("2006-01-02 15:04")
	}

	return b, nil
}

// Create requests the creation of a new cron trigger.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCronTriggerCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, nil)

	return
}

// Delete deletes the specified cron trigger.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

// Get retrieves details of a single cron trigger.
// Use Extract to convert its result into an CronTrigger.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}
