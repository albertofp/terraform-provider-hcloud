package location_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hetznercloud/terraform-provider-hcloud/internal/location"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/teste2e"
	"github.com/hetznercloud/terraform-provider-hcloud/internal/testtemplate"
)

func TestAccLocationDataSource(t *testing.T) {
	tmplMan := testtemplate.Manager{}

	lByName := &location.DData{
		LocationName: "fsn1",
	}
	lByName.SetRName("l_by_name")
	lByID := &location.DData{
		LocationID: "1",
	}
	lByID.SetRName("l_by_id")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 teste2e.PreCheck(t),
		ProtoV6ProviderFactories: teste2e.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: tmplMan.Render(t,
					"testdata/d/hcloud_location", lByName,
					"testdata/d/hcloud_location", lByID,
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(lByName.TFID(), "id", "1"),
					resource.TestCheckResourceAttr(lByName.TFID(), "name", "fsn1"),
					resource.TestCheckResourceAttr(lByName.TFID(), "description", "Falkenstein DC Park 1"),
					resource.TestCheckResourceAttr(lByName.TFID(), "network_zone", "eu-central"),

					resource.TestCheckResourceAttr(lByID.TFID(), "id", "1"),
					resource.TestCheckResourceAttr(lByID.TFID(), "name", "fsn1"),
					resource.TestCheckResourceAttr(lByID.TFID(), "description", "Falkenstein DC Park 1"),
					resource.TestCheckResourceAttr(lByID.TFID(), "network_zone", "eu-central"),
				),
			},
		},
	})
}

func TestAccLocationDataSource_UpgradePluginFramework(t *testing.T) {
	tmplMan := testtemplate.Manager{}

	lByName := &location.DData{
		LocationName: "fsn1",
	}
	lByName.SetRName("l_by_name")
	lByID := &location.DData{
		LocationID: "1",
	}
	lByID.SetRName("l_by_id")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: teste2e.PreCheck(t),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"hcloud": {
						VersionConstraint: "1.44.1",
						Source:            "hetznercloud/hcloud",
					},
				},

				Config: tmplMan.Render(t,
					"testdata/d/hcloud_location", lByName,
					"testdata/d/hcloud_location", lByID,
					"testdata/r/terraform_data_resource", lByName,
					"testdata/r/terraform_data_resource", lByID,
				),
			},
			{
				ProtoV6ProviderFactories: teste2e.ProtoV6ProviderFactories(),

				Config: tmplMan.Render(t,
					"testdata/d/hcloud_location", lByName,
					"testdata/d/hcloud_location", lByID,
					"testdata/r/terraform_data_resource", lByName,
					"testdata/r/terraform_data_resource", lByID,
				),

				PlanOnly: true,
			},
		},
	})
}

func TestAccLocationDataSourceList(t *testing.T) {
	tmplMan := testtemplate.Manager{}

	locationsDS := &location.DDataList{}
	locationsDS.SetRName("ds")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 teste2e.PreCheck(t),
		ProtoV6ProviderFactories: teste2e.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: tmplMan.Render(t,
					"testdata/d/hcloud_locations", locationsDS,
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(locationsDS.TFID(), "location_ids.0", "1"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "location_ids.1", "2"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "location_ids.2", "3"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "location_ids.3", "4"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "location_ids.4", "5"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "location_ids.5", "6"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "names.0", "fsn1"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "names.1", "nbg1"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "names.2", "hel1"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "names.3", "ash"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "names.4", "hil"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "names.5", "sin"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "descriptions.0", "Falkenstein DC Park 1"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "descriptions.1", "Nuremberg DC Park 1"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "descriptions.2", "Helsinki DC Park 1"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "descriptions.3", "Ashburn, VA"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "descriptions.4", "Hillsboro, OR"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "descriptions.5", "Singapore"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "locations.#", "6"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "locations.0.name", "fsn1"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "locations.0.network_zone", "eu-central"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "locations.1.name", "nbg1"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "locations.2.name", "hel1"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "locations.3.name", "ash"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "locations.3.network_zone", "us-east"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "locations.4.name", "hil"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "locations.4.network_zone", "us-west"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "locations.5.name", "sin"),
					resource.TestCheckResourceAttr(locationsDS.TFID(), "locations.5.network_zone", "ap-southeast"),
				),
			},
		},
	})
}

func TestAccLocationDataSourceList_UpgradePluginFramework(t *testing.T) {
	tmplMan := testtemplate.Manager{}

	locationsDS := &location.DDataList{}
	locationsDS.SetRName("ds")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: teste2e.PreCheck(t),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"hcloud": {
						VersionConstraint: "1.44.1",
						Source:            "hetznercloud/hcloud",
					},
				},

				Config: tmplMan.Render(t,
					"testdata/d/hcloud_locations", locationsDS,
					"testdata/r/terraform_data_resource", locationsDS,
				),
			},
			{
				ProtoV6ProviderFactories: teste2e.ProtoV6ProviderFactories(),

				Config: tmplMan.Render(t,
					"testdata/d/hcloud_locations", locationsDS,
					"testdata/r/terraform_data_resource", locationsDS,
				),

				PlanOnly: true,
			},
		},
	})
}
