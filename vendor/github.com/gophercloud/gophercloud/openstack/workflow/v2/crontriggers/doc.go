/*
Package crontriggers provides interaction with the cron triggers API in the OpenStack Mistral service.

Cron trigger is an object that allows to run Mistral workflows according to a time pattern (Unix crontab patterns format).
Once a trigger is created it will run a specified workflow according to its properties: pattern, first_execution_time and remaining_executions.

Example to list cron triggers

	listOpts := crontriggers.ListOpts{
		WorkflowID: "w1",
	}
	allPages, err := crontriggers.List(mistralClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}
	allCrontriggers, err := crontriggers.ExtractCronTriggers(allPages)
	if err != nil {
		panic(err)
	}
	for _, ex := range allCrontriggers {
		fmt.Printf("%+v\n", ex)
	}

Example to create a cron trigger

	createOpts := &crontriggers.CreateOpts{
		WorkflowID:     "w1",
		WorkflowParams: map[string]interface{}{
			"msg": "hello",
		},
		WorkflowInput: map[string]interface{}{
			"msg": "world",
		},
		Name:           "trigger",
	}
	crontrigger, err := crontriggers.Create(mistralClient, opts).Extract()
	if err != nil {
		panic(err)
	}
*/
package crontriggers
