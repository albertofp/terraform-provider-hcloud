package server_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/network"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/server"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/sshkey"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/teste2e"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/testsupport"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/testtemplate"
)

func TestAccServerNetworkResource_NetworkID(t *testing.T) {
	var (
		nw hcloud.Network
		s  hcloud.Server
	)

	sk := sshkey.NewRData(t, "server-network-id")
	netRes := &network.RData{
		Name:    "test-network",
		IPRange: "10.0.0.0/16",
	}
	netRes.SetRName("test-network")
	subNetRes := &network.RDataSubnet{
		Type:        "cloud",
		NetworkID:   netRes.TFID() + ".id",
		NetworkZone: "eu-central",
		IPRange:     "10.0.1.0/24",
	}
	subNetRes.SetRName("test-network-subnet")
	sRes := &server.RData{
		Name:       "s-network-test",
		Type:       teste2e.TestServerType,
		Datacenter: teste2e.TestDataCenter,
		Image:      teste2e.TestImage,
		SSHKeys:    []string{sk.TFID() + ".id"},
	}
	sRes.SetRName("s-network-test")
	sNRes := &server.RDataNetwork{
		Name:      "test-network",
		ServerID:  sRes.TFID() + ".id",
		NetworkID: netRes.TFID() + ".id",
		IP:        "10.0.1.5",
		DependsOn: []string{subNetRes.TFID()},
	}
	tmplMan := testtemplate.Manager{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 teste2e.PreCheck(t),
		ProtoV6ProviderFactories: teste2e.ProtoV6ProviderFactories(),
		CheckDestroy:             testsupport.CheckResourcesDestroyed(server.ResourceType, server.ByID(t, nil)),
		Steps: []resource.TestStep{
			{
				Config: tmplMan.Render(t,
					"testdata/r/hcloud_ssh_key", sk,
					"testdata/r/hcloud_network", netRes,
					"testdata/r/hcloud_network_subnet", subNetRes,
					"testdata/r/hcloud_server", sRes,
					"testdata/r/hcloud_server_network", sNRes,
				),
				Check: resource.ComposeTestCheckFunc(
					testsupport.CheckResourceExists(netRes.TFID(), network.ByID(t, &nw)),
					testsupport.CheckResourceExists(sRes.TFID(), server.ByID(t, &s)),
					testsupport.LiftTCF(hasServerNetwork(t, &s, &nw, "10.0.1.5")),
					resource.TestCheckResourceAttr(
						server.NetworkResourceType+".test-network", "ip", "10.0.1.5"),
				),
			},
			{
				// Try to import the newly created Server
				ResourceName:      server.NetworkResourceType + ".test-network",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(_ *terraform.State) (string, error) {
					return fmt.Sprintf("%d-%d", s.ID, nw.ID), nil
				},
			},
		},
	})
}

func TestAccServerNetworkResource_SubNetID(t *testing.T) {
	var (
		nw hcloud.Network
		s  hcloud.Server
	)

	sk := sshkey.NewRData(t, "server-network-subnetid")
	netRes := &network.RData{
		Name:    "test-network",
		IPRange: "10.0.0.0/16",
	}
	netRes.SetRName("test-network")
	subNetRes := &network.RDataSubnet{
		Type:        "cloud",
		NetworkID:   netRes.TFID() + ".id",
		NetworkZone: "eu-central",
		IPRange:     "10.0.1.0/24",
	}
	subNetRes.SetRName("test-network-subnet")
	sRes := &server.RData{
		Name:       "s-network-test",
		Type:       teste2e.TestServerType,
		Datacenter: teste2e.TestDataCenter,
		Image:      teste2e.TestImage,
		SSHKeys:    []string{sk.TFID() + ".id"},
	}
	sRes.SetRName("s-network-test")

	tmplMan := testtemplate.Manager{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 teste2e.PreCheck(t),
		ProtoV6ProviderFactories: teste2e.ProtoV6ProviderFactories(),
		CheckDestroy:             testsupport.CheckResourcesDestroyed(server.ResourceType, server.ByID(t, nil)),
		Steps: []resource.TestStep{
			{
				Config: tmplMan.Render(t,
					"testdata/r/hcloud_ssh_key", sk,
					"testdata/r/hcloud_network", netRes,
					"testdata/r/hcloud_network_subnet", subNetRes,
					"testdata/r/hcloud_server", sRes,
					"testdata/r/hcloud_server_network", &server.RDataNetwork{
						Name:     "test-network",
						ServerID: sRes.TFID() + ".id",
						SubNetID: subNetRes.TFID() + ".id",
						IP:       "10.0.1.5",
					},
				),
				Check: resource.ComposeTestCheckFunc(
					testsupport.CheckResourceExists(netRes.TFID(), network.ByID(t, &nw)),
					testsupport.CheckResourceExists(sRes.TFID(), server.ByID(t, &s)),
					testsupport.LiftTCF(hasServerNetwork(t, &s, &nw, "10.0.1.5")),
				),
			},
		},
	})
}

func hasServerNetwork(t *testing.T, s *hcloud.Server, nw *hcloud.Network, ips ...string) func() error {
	return func() error {
		var privNet *hcloud.ServerPrivateNet

		for _, n := range s.PrivateNet {
			if n.Network.ID == nw.ID {
				privNet = &n
				break
			}
		}
		if !assert.NotNil(t, privNet, "server has no private network") {
			return nil
		}
		assert.Contains(t, ips, privNet.IP.String())
		if len(ips) > 1 {
			for _, aip := range privNet.Aliases {
				assert.Contains(t, ips, aip.String())
			}
		}

		return nil
	}
}
