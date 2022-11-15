package azure

import (
	"fmt"
	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

type ServicebusNamespace struct {
	Address       string
	Region        string
	SKUName       string
	Capacity      int64
	ZoneRedundant bool
}

var ServicebusNamespaceUsageSchema = []*schema.UsageItem{}

func (r *ServicebusNamespace) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *ServicebusNamespace) BuildResource() *schema.Resource {
	//sku := ""
	//productName := ""
	//extra := ""
	/*
		var capacity int64 = int64(math.Max(float64(r.WorkerCount), float64(r.MaximumWorkerCount)))

		// free plan
		if strings.ToLower(r.SKUName) == "f1" {
			return &schema.Resource{
				Name:      r.Address,
				IsSkipped: true,
				NoPrice:   true, UsageSchema: ServicebusNamespaceUsageSchema,
			}
		}
	*/
	sku := "Premium"
	productName := ""
	capacity := r.Capacity
	if r.Capacity < 1 {
		capacity = 1
	}

	return &schema.Resource{
		Name:           r.Address,
		CostComponents: []*schema.CostComponent{r.servicebusNamespaceCostComponent(fmt.Sprintf("Instance usage (%s)", sku), productName, sku, capacity)},
		UsageSchema:    ServicebusNamespaceUsageSchema,
	}
}

func (r *ServicebusNamespace) servicebusNamespaceCostComponent(name, productName, skuRefactor string, capacity int64) *schema.CostComponent {
	return &schema.CostComponent{
		Name:           name,
		Unit:           "hours",
		UnitMultiplier: decimal.NewFromInt(1),
		HourlyQuantity: decimalPtr(decimal.NewFromInt(capacity)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Service Bus"),
			ProductFamily: strPtr("Integration"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "productName", Value: strPtr("Service Bus")},
				{Key: "skuName", ValueRegex: strPtr(fmt.Sprintf("/%s/i", skuRefactor))},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
	}
}
