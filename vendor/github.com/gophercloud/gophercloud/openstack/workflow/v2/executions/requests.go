package executions

import "github.com/gophercloud/gophercloud"

// CreateOptsBuilder allows extension to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToExecutionCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies parameters used to create an execution.
type CreateOpts struct {
	// ID is the unique ID of the execution.
	ID string `json:"id,omitempty"`

	// SourceExecutionID can be set to create an execution based on another existing execution.
	SourceExecutionID string `json:"source_execution_id,omitempty"`

	// WorkflowID is the unique id of the workflow.
	WorkflowID string `json:"workflow_id,omitempty" or:"WorkflowName"`

	// WorkflowName is the name identifier of the workflow.
	WorkflowName string `json:"workflow_name,omitempty" or:"WorkflowID"`

	// WorkflowNamespace is the namespace of the workflow.
	WorkflowNamespace string `json:"workflow_namespace,omitempty"`

	// Input is a JSON structure containing workflow input values, serialized as string.
	Input map[string]interface{} `json:"input,omitempty"`

	// Params define workflow type specific parameters.
	Params map[string]interface{} `json:"params,omitempty"`

	// Description is the description of the workflow execution.
	Description string `json:"description,omitempty"`
}

// ToExecutionCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToExecutionCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create requests the creation of a new execution.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToExecutionCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, nil)

	return
}

// Get retrieves details of a single execution.
// Use ExtractExecution to convert its result into an Execution.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// Delete deletes the specified execution.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}
