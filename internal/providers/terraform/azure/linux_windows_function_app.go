package azure

import (
	"strings"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

func GetAzureRMLinuxFunctionAppRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "azurerm_linux_function_app",
		RFunc: NewAzureRMLinuxWindowsFunctionApp,
		ReferenceAttributes: []string{
			"service_plan_id",
		},
	}
}

func GetAzureRMWindowsFunctionAppRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "azurerm_windows_function_app",
		RFunc: NewAzureRMLinuxWindowsFunctionApp,
		ReferenceAttributes: []string{
			"service_plan_id",
		},
	}
}

func NewAzureRMLinuxWindowsFunctionApp(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	region := lookupRegion(d, []string{})

	var memorySize, executionTime, executions, gbSeconds *decimal.Decimal
	var skuCPU *int64
	var skuMemory *float64
	var skuName string
	var maximumWorkerCount int64 = 1

	if u != nil && u.Get("monthly_executions").Type != gjson.Null {
		executions = decimalPtr(decimal.NewFromInt(u.Get("monthly_executions").Int()))
	}
	if u != nil && u.Get("execution_duration_ms").Type != gjson.Null &&
		u.Get("memory_mb").Type != gjson.Null &&
		executions != nil {

		memorySize = decimalPtr(decimal.NewFromInt(u.Get("memory_mb").Int()))
		executionTime = decimalPtr(decimal.NewFromInt(u.Get("execution_duration_ms").Int()))
		gbSeconds = decimalPtr(calculateFunctionAppGBSeconds(*memorySize, *executionTime, *executions))
	}

	skuMapCPU := map[string]int64{
		"ep1": 1,
		"ep2": 2,
		"ep3": 4,
	}
	skuMapMemory := map[string]float64{
		"ep1": 3.5,
		"ep2": 7.0,
		"ep3": 14.0,
	}

	servicePlanID := d.References("service_plan_id")

	if len(servicePlanID) > 0 {
		skuName = strings.ToLower(servicePlanID[0].Get("sku_name").String())
		maximumWorkerCount = servicePlanID[0].Get("maximum_elastic_worker_count").Int()
	}

	if val, ok := skuMapCPU[skuName]; ok {
		skuCPU = &val
	}
	if val, ok := skuMapMemory[skuName]; ok {
		skuMemory = &val
	}

	instances := decimal.NewFromInt(maximumWorkerCount)
	if u != nil && u.Get("instances").Type != gjson.Null {
		instances = decimal.NewFromInt(u.Get("instances").Int())
	}

	costComponents := make([]*schema.CostComponent, 0)

	// TODO: add pre-warmed instances (shared by plan): instances x hours
	// premium
	if (strings.ToLower(skuName)[0] == 'e') && skuCPU != nil && skuMemory != nil {
		costComponents = append(costComponents, AppFunctionPremiumCPUCostComponent(skuName, instances, skuCPU, region))
		costComponents = append(costComponents, AppFunctionPremiumMemoryCostComponent(skuName, instances, skuMemory, region))
	} else if strings.ToLower(skuName)[0] == 'y' {
		// consumption plan: The first 400,000 GB/s of execution and 1,000,000 executions are free.
		// Memory size x execution time (ms) x executions per month
		if gbSeconds != nil {
			costComponents = append(costComponents, AppFunctionConsumptionExecutionTimeCostComponent(gbSeconds, region))
		}
		if executions != nil {
			costComponents = append(costComponents, AppFunctionConsumptionExecutionsCostComponent(executions, region))
		}
	}

	return &schema.Resource{
		Name:           d.Address,
		CostComponents: costComponents,
	}
}
