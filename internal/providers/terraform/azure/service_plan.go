package azure

import (
	"github.com/infracost/infracost/internal/resources/azure"
	"github.com/infracost/infracost/internal/schema"
)

func getServicePlanRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "azurerm_service_plan",
		RFunc: NewServicePlan,
	}
}
func NewServicePlan(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	r := &azure.ServicePlan{
		Address:            d.Address,
		Region:             lookupRegion(d, []string{}),
		SKUName:            d.Get("sku_name").String(),
		OSType:             d.Get("os_type").String(),
		WorkerCount:        d.Get("worker_count").Int(),
		MaximumWorkerCount: d.Get("maximum_elastic_worker_count").Int(),
	}
	if r.WorkerCount == 0 {
		r.WorkerCount = 1
	}
	if r.MaximumWorkerCount == 0 {
		r.MaximumWorkerCount = 1
	}
	r.PopulateUsage(u)
	return r.BuildResource()
}
