package services

import (
	"procurement-system/internal/models"
)

type NavigationService interface {
	GetMenuForRole(role string) []models.NavigationItem
}

type navigationService struct{}

func NewNavigationService() NavigationService {
	return &navigationService{}
}

func (s *navigationService) GetMenuForRole(role string) []models.NavigationItem {
	switch role {
	case "Admin":
		return getAdminMenu()
	case "Procurement Officer":
		return getProcurementOfficerMenu()
	case "Approver":
		return getApproverMenu()
	case "Vendor":
		return getVendorMenu()
	default: // Employee
		return getEmployeeMenu()
	}
}

func getAdminMenu() []models.NavigationItem {
	return []models.NavigationItem{
		{Title: "Dashboard", Path: "/dashboard", Icon: "dashboard"},
		{Title: "Procurement", Path: "/procurement", Icon: "shopping_cart", SubItems: []models.NavigationSubItem{
			{Title: "Requisitions", Path: "/procurement/requisitions"},
			{Title: "Purchase Orders", Path: "/procurement/purchase-orders"},
			{Title: "Approvals", Path: "/procurement/approvals"},
		}},
		{Title: "Vendors", Path: "/vendors", Icon: "store"},
		{Title: "Reports", Path: "/reports", Icon: "assessment"},
		{Title: "Administration", Path: "/admin", Icon: "settings", SubItems: []models.NavigationSubItem{
			{Title: "User Management", Path: "/admin/users"},
			{Title: "System Settings", Path: "/admin/settings"},
		}},
	}
}

func getProcurementOfficerMenu() []models.NavigationItem {
	// Omitting for brevity in this example
	return getAdminMenu() // For now, let's give them the same menu as admin
}

func getApproverMenu() []models.NavigationItem {
	// Omitting for brevity
	return []models.NavigationItem{
		{Title: "Dashboard", Path: "/dashboard", Icon: "dashboard"},
		{Title: "Approvals", Path: "/approvals", Icon: "check_circle"},
	}
}

func getVendorMenu() []models.NavigationItem {
	// Omitting for brevity
	return []models.NavigationItem{
		{Title: "Purchase Orders", Path: "/purchase-orders", Icon: "list_alt"},
		{Title: "Invoices", Path: "/invoices", Icon: "receipt"},
	}
}

func getEmployeeMenu() []models.NavigationItem {
	return []models.NavigationItem{
		{Title: "Dashboard", Path: "/dashboard", Icon: "dashboard"},
		{Title: "My Requisitions", Path: "/requisitions/my", Icon: "description"},
	}
}
