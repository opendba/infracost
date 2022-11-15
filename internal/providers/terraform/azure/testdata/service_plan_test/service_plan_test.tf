provider "azurerm" {
  skip_provider_registration = true
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "exampleRG1"
  location = "eastus"
}

resource "azurerm_service_plan" "standard_s1" {
  name                = "api-appserviceplan-pro"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  os_type             = "Windows"

  sku_name     = "S1"
  worker_count = 1
}

resource "azurerm_service_plan" "standard_s2" {
  name                = "api-appserviceplan-pro"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  os_type             = "Windows"

  sku_name     = "S1"
  worker_count = 5
}

resource "azurerm_service_plan" "premium_v2" {
  name                = "api-appserviceplan-pro"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  os_type             = "Linux"

  sku_name     = "P1v2"
  worker_count = 10
}

resource "azurerm_service_plan" "basic" {
  name                = "api-appserviceplan-pro"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  os_type             = "Linux"

  sku_name     = "B2"
  worker_count = 1
}

resource "azurerm_service_plan" "free" {
  name                = "api-appserviceplan-pro"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  os_type             = "Linux"

  sku_name     = "F1"
  worker_count = 10
}

resource "azurerm_service_plan" "default_capacity" {
  name                = "api-appserviceplan-pro"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  os_type             = "Linux"

  sku_name = "B2"
}
