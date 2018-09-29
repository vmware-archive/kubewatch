// +build acceptance blockstorage

package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestVolumes(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volume1, err := CreateVolume(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, client, volume1)

	volume2, err := CreateVolume(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, client, volume2)

	listOpts := volumes.ListOpts{
		Limit: 1,
	}

	err = volumes.List(client, listOpts).EachPage(func(page pagination.Page) (bool, error) {
		actual, err := volumes.ExtractVolumes(page)
		th.AssertNoErr(t, err)
		th.AssertEquals(t, 1, len(actual))

		var found bool
		for _, v := range actual {
			if v.ID == volume1.ID || v.ID == volume2.ID {
				found = true
			}
		}

		th.AssertEquals(t, found, true)

		return true, nil
	})

	th.AssertNoErr(t, err)
}

func TestVolumesMultiAttach(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volumeName := tools.RandomString("ACPTTEST", 16)

	volOpts := volumes.CreateOpts{
		Size:        1,
		Name:        volumeName,
		Description: "Testing creation of multiattach enabled volume",
		Multiattach: true,
	}

	vol, err := volumes.Create(client, volOpts).Extract()
	th.AssertNoErr(t, err)

	err = volumes.WaitForStatus(client, vol.ID, "available", 60)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, vol.Multiattach, true)

	err = volumes.Delete(client, vol.ID).ExtractErr()
	th.AssertNoErr(t, err)
}
