package models

type NavigationItem struct {
	Title    string              `json:"title"`
	Path     string              `json:"path"`
	Icon     string              `json:"icon"`
	SubItems []NavigationSubItem `json:"subItems,omitempty"`
}

type NavigationSubItem struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

type BreadcrumbItem struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}
