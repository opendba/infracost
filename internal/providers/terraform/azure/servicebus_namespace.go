package azure

import (
	"github.com/infracost/infracost/internal/resources/azure"
	"github.com/infracost/infracost/internal/schema"
)

func getServicebusNamespaceRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "azurerm_servicebus_namespace",
		RFunc: NewServicebusNamespace,
	}
}
func NewServicebusNamespace(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	r := &azure.ServicebusNamespace{
		Address:       d.Address,
		Region:        lookupRegion(d, []string{}),
		SKUName:       d.Get("sku").String(),
		Capacity:      d.Get("capacity").Int(),
		ZoneRedundant: d.Get("zone_redundant").Bool(),
	}
	r.PopulateUsage(u)
	return r.BuildResource()
}
