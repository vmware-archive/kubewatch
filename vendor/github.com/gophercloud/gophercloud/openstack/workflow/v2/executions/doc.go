/*
Package executions provides interaction with the execution API in the OpenStack Mistral service.

An execution is a particular execution of a specific workflow. Each execution contains all information about workflow itself, about execution process, state, input and output data.

Example to list executions

	listOpts := executions.ListOpts{
		WorkflowID: "w1",
	}

	allPages, err := executions.List(mistralClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allExecutions, err := executions.ExtractExecutions(allPages)
	if err != nil {
		panic(err)
	}

	for _, ex := range allExecutions {
		fmt.Printf("%+v\n", ex)
	}

Example to create an execution

	createOpts := &executions.CreateOpts{
		WorkflowID:  "6656c143-a009-4bcb-9814-cc100a20bbfa",
		Input: map[string]interface{}{
			"msg": "Hello",
		},
		Description: "this is a description",
	}

	execution, err := executions.Create(mistralClient, opts).Extract()
	if err != nil {
		panic(err)
	}
*/
package executions
