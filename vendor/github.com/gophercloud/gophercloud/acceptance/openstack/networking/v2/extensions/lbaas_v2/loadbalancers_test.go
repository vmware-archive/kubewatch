// +build acceptance networking lbaas_v2 loadbalancers

package lbaas_v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/listeners"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/monitors"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/pools"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestLoadbalancersList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	allPages, err := loadbalancers.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	allLoadbalancers, err := loadbalancers.ExtractLoadBalancers(allPages)
	th.AssertNoErr(t, err)

	for _, lb := range allLoadbalancers {
		tools.PrintResource(t, lb)
	}
}

func TestLoadbalancersCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	subnet, err := networking.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, client, subnet.ID)

	lb, err := CreateLoadBalancer(t, client, subnet.ID)
	th.AssertNoErr(t, err)
	defer DeleteLoadBalancer(t, client, lb.ID)

	newLB, err := loadbalancers.Get(client, lb.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newLB)

	// Because of the time it takes to create a loadbalancer,
	// this test will include some other resources.

	// Listener
	listener, err := CreateListener(t, client, lb)
	th.AssertNoErr(t, err)
	defer DeleteListener(t, client, lb.ID, listener.ID)

	updateListenerOpts := listeners.UpdateOpts{
		Description: "Some listener description",
	}
	_, err = listeners.Update(client, listener.ID, updateListenerOpts).Extract()
	th.AssertNoErr(t, err)

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE", loadbalancerActiveTimeoutSeconds); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newListener, err := listeners.Get(client, listener.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newListener)

	// Pool
	pool, err := CreatePool(t, client, lb)
	th.AssertNoErr(t, err)
	defer DeletePool(t, client, lb.ID, pool.ID)

	updatePoolOpts := pools.UpdateOpts{
		Description: "Some pool description",
	}
	_, err = pools.Update(client, pool.ID, updatePoolOpts).Extract()
	th.AssertNoErr(t, err)

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE", loadbalancerActiveTimeoutSeconds); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newPool, err := pools.Get(client, pool.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPool)

	// Member
	member, err := CreateMember(t, client, lb, newPool, subnet.ID, subnet.CIDR)
	th.AssertNoErr(t, err)
	defer DeleteMember(t, client, lb.ID, pool.ID, member.ID)

	newWeight := tools.RandomInt(11, 100)
	updateMemberOpts := pools.UpdateMemberOpts{
		Weight: newWeight,
	}
	_, err = pools.UpdateMember(client, pool.ID, member.ID, updateMemberOpts).Extract()
	th.AssertNoErr(t, err)

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE", loadbalancerActiveTimeoutSeconds); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newMember, err := pools.GetMember(client, pool.ID, member.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newMember)

	// Monitor
	monitor, err := CreateMonitor(t, client, lb, newPool)
	th.AssertNoErr(t, err)
	defer DeleteMonitor(t, client, lb.ID, monitor.ID)

	newDelay := tools.RandomInt(20, 30)
	updateMonitorOpts := monitors.UpdateOpts{
		Delay: newDelay,
	}
	_, err = monitors.Update(client, monitor.ID, updateMonitorOpts).Extract()
	th.AssertNoErr(t, err)

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE", loadbalancerActiveTimeoutSeconds); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newMonitor, err := monitors.Get(client, monitor.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newMonitor)
}
