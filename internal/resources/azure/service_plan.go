package azure

import (
	"fmt"
	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"math"
	"strings"

	"github.com/shopspring/decimal"
)

type ServicePlan struct {
	Address            string
	SKUName            string
	OSType             string
	Region             string
	WorkerCount        int64
	MaximumWorkerCount int64
}

var ServicePlanUsageSchema = []*schema.UsageItem{}

func (r *ServicePlan) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

func (r *ServicePlan) BuildResource() *schema.Resource {
	sku := ""
	productName := ""
	var capacity int64 = int64(math.Max(float64(r.WorkerCount), float64(r.MaximumWorkerCount)))

	// free plan
	if strings.ToLower(r.SKUName) == "f1" {
		return &schema.Resource{
			Name:      r.Address,
			IsSkipped: true,
			NoPrice:   true, UsageSchema: ServicePlanUsageSchema,
		}
	}

	// P plan
	sku = strings.ToUpper(r.SKUName)
	if strings.ToLower(r.SKUName[0:1]) == "p" {
		sku = r.SKUName[:2] + " " + r.SKUName[2:]
		switch strings.ToLower(r.SKUName[2:]) {
		case "v2":
			productName = "Premium v2 Plan"
		case "v3":
			productName = "Premium v3 Plan"
		}
	} else if strings.ToLower(r.SKUName[0:1]) == "b" {
		productName = "Basic Plan"
	} else if strings.ToLower(r.SKUName[0:1]) == "s" {
		productName = "Standard Plan"
	} else if strings.ToLower(r.SKUName[0:1]) == "i" {
		if len(r.SKUName) <= 3 {
			productName = "Isolated Plan"
		} else if strings.ToLower(r.SKUName[2:]) == "v2" {
			productName = "Isolated v2 Plan"
		}
	}

	if strings.ToLower(r.OSType) != "windows" {
		productName += " - Linux"
	}

	return &schema.Resource{
		Name:           r.Address,
		CostComponents: []*schema.CostComponent{r.servicePlanCostComponent(fmt.Sprintf("Instance usage (%s)", sku), productName, sku, capacity)},
		UsageSchema:    ServicePlanUsageSchema,
	}
}

func (r *ServicePlan) servicePlanCostComponent(name, productName, skuRefactor string, capacity int64) *schema.CostComponent {
	return &schema.CostComponent{
		Name:           name,
		Unit:           "hours",
		UnitMultiplier: decimal.NewFromInt(1),
		HourlyQuantity: decimalPtr(decimal.NewFromInt(capacity)),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Azure App Service"),
			ProductFamily: strPtr("Compute"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "productName", Value: strPtr("Azure App Service " + productName)},
				{Key: "skuName", ValueRegex: strPtr(fmt.Sprintf("/%s/i", skuRefactor))},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
	}
}
