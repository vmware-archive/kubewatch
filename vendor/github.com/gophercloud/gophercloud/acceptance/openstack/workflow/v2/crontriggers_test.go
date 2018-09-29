package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCronTriggersCreateGetDelete(t *testing.T) {
	client, err := clients.NewWorkflowV2Client()
	th.AssertNoErr(t, err)

	workflow, err := CreateWorkflow(t, client)
	th.AssertNoErr(t, err)
	defer DeleteWorkflow(t, client, workflow)

	trigger, err := CreateCronTrigger(t, client, workflow)
	th.AssertNoErr(t, err)
	defer DeleteCronTrigger(t, client, trigger)

	gettrigger, err := GetCronTrigger(t, client, trigger.ID)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, trigger.ID, gettrigger.ID)

	tools.PrintResource(t, trigger)
}
