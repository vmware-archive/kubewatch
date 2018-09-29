package workflows

import (
	"io"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// CreateOptsBuilder allows extension to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToWorkflowCreateParams() (io.Reader, string, error)
}

// CreateOpts specifies parameters used to create a cron trigger.
type CreateOpts struct {
	// Scope is the scope of the workflow.
	// Allowed values are "private" and "public".
	Scope string `q:"scope"`

	// Namespace will define the namespace of the workflow.
	Namespace string `q:"namespace"`

	// Definition is the workflow definition written in Mistral Workflow Language v2.
	Definition io.Reader
}

// ToWorkflowCreateParams constructs a request query string from CreateOpts.
func (opts CreateOpts) ToWorkflowCreateParams() (io.Reader, string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return opts.Definition, q.String(), err
}

// Create requests the creation of a new execution.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	url := createURL(client)
	var b io.Reader
	if opts != nil {
		tmpB, query, err := opts.ToWorkflowCreateParams()
		if err != nil {
			r.Err = err
			return
		}
		url += query
		b = tmpB
	}

	_, r.Err = client.Post(url, nil, &r.Body, &gophercloud.RequestOpts{
		RawBody: b,
		MoreHeaders: map[string]string{
			"Content-Type": "text/plain",
			"Accept":       "", // Drop default JSON Accept header
		},
	})

	return
}

// Delete deletes the specified execution.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

// Get retrieves details of a single execution.
// Use Extract to convert its result into an Workflow.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// ListOptsBuilder allows extension to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToWorkflowListQuery() (string, error)
}

// ListOpts filters the result returned by the List() function.
type ListOpts struct {
	// Name allows to filter by workflow name.
	Name string `q:"name"`
	// Namespace allows to filter by workflow namespace.
	Namespace string `q:"namespace"`
	// Definition allows to filter by workflow definition.
	Definition string `q:"definition"`
	// Scope filters by the workflow's scope.
	// Values can be "private" or "public".
	Scope string `q:"scope"`
	// SortDir allows to select sort direction.
	// It can be "asc" or "desc" (default).
	SortDir string `q:"sort_dir"`
	// SortKey allows to sort by one of the cron trigger attributes.
	SortKey string `q:"sort_key"`
	// Marker and Limit control paging.
	// Marker instructs List where to start listing from.
	Marker string `q:"marker"`
	// Limit instructs List to refrain from sending excessively large lists of
	// cron triggers.
	Limit int `q:"limit"`
}

// ToWorkflowListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToWorkflowListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List performs a call to list cron triggers.
// You may provide options to filter the results.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToWorkflowListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return WorkflowPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
