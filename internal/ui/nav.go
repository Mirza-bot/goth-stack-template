package ui

import (
	"net/http"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/a-h/templ"
)

type NavItem struct {
	ID        string
	Label     string
	Href      string
	OnSidebar bool
	OnMenubar bool
	Order     int
	Children  []NavItem
	Icon      templ.Component
}

type NavConfig struct {
	ID        string
	ParentID  string
	Label     string
	Href      string
	OnSidebar bool
	OnMenubar bool
	Order     int
	Hidden    func(*http.Request) bool
	Icon      templ.Component
}

type PageProvider func(req *http.Request) (templ.Component, error)

type Route struct {
	Method      string
	Path        string
	Page        PageProvider
	Middlewares []func(http.Handler) http.Handler
	NavConfig   *NavConfig
}

type RouteGroup struct {
	Title  string
	Routes []Route
	Extras []NavItem
}

func (routeGroup RouteGroup) GroupRoutes(request *http.Request) []NavItem {
	// index nodes by ID for building/merging
	nodeByID := map[string]*NavItem{}

	// seed extras
	for i := range routeGroup.Extras {
		extra := routeGroup.Extras[i] // copy
		if extra.ID == "" {
			extra.ID = "extra:" + extra.Label
		}
		nodeCopy := extra
		nodeByID[extra.ID] = &nodeCopy
	}

	var rootPointers []*NavItem

	addRoot := func(nodePointer *NavItem) {
		if !containsPointer(rootPointers, nodePointer) {
			rootPointers = append(rootPointers, nodePointer)
		}
	}

	ensureNode := func(id string) *NavItem {
		if id == "" {
			return nil
		}
		if existing, ok := nodeByID[id]; ok {
			return existing
		}
		newNode := &NavItem{ID: id}
		nodeByID[id] = newNode
		return newNode
	}

	// turn routes into (maybe) nav nodes
	for _, route := range routeGroup.Routes {
		if route.NavConfig == nil {
			continue
		}
		if route.NavConfig.Hidden != nil && route.NavConfig.Hidden(request) {
			continue // hidden by predicate
		}

		nodeID := firstNonEmpty(route.NavConfig.ID, route.Path)

		// derive href: prefer explicit Href; fall back to static Pattern (no params)
		href := route.NavConfig.Href
		if href == "" && !strings.Contains(route.Path, "{") {
			href = route.Path
		}

		// derive label
		label := route.NavConfig.Label
		if label == "" {
			label = titleFromPath(firstNonEmpty(href, route.Path))
		}

		nodePointer := ensureNode(nodeID)
		nodePointer.Label = label
		nodePointer.Href = href
		nodePointer.OnSidebar = route.NavConfig.OnSidebar
		nodePointer.OnMenubar = route.NavConfig.OnMenubar
		nodePointer.Order = route.NavConfig.Order
		nodePointer.Icon = route.NavConfig.Icon

		if route.NavConfig.ParentID != "" {
			parentPointer := ensureNode(route.NavConfig.ParentID)
			parentPointer.Children = append(parentPointer.Children, *nodePointer)
		} else {
			addRoot(nodePointer)
		}
	}

	// also include any extras that didn’t get children/roots yet
	for _, nodePointer := range nodeByID {
		if !isChildOfAny(nodePointer, nodeByID) && !containsPointer(rootPointers, nodePointer) {
			addRoot(nodePointer)
		}
	}

	// sort roots and children by Order then Label
	sortNavPointerSlice(rootPointers)
	for _, nodePointer := range rootPointers {
		sortNavValueSlice(nodePointer.Children)
	}

	// materialize []NavItem roots
	navigationItems := make([]NavItem, 0, len(rootPointers))
	for _, nodePointer := range rootPointers {
		navigationItems = append(navigationItems, *nodePointer)
	}
	return navigationItems
}

func titleFromPath(p string) string {
	base := path.Base("/" + strings.TrimSpace(p))
	base = strings.Trim(base, "/")
	if base == "" || base == "." {
		return "Home"
	}
	base = strings.ReplaceAll(base, "-", " ")
	return strings.Title(base)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func isChildOfAny(nodePointer *NavItem, nodeByID map[string]*NavItem) bool {
	for _, parentPointer := range nodeByID {
		for _, child := range parentPointer.Children {
			if child.ID == nodePointer.ID {
				return true
			}
		}
	}
	return false
}

func containsPointer(list []*NavItem, target *NavItem) bool {
	return slices.Contains(list, target)
}

func sortNavPointerSlice(list []*NavItem) {
	sort.SliceStable(list, func(i, j int) bool {
		if list[i].Order != list[j].Order {
			return list[i].Order < list[j].Order
		}
		return strings.ToLower(list[i].Label) < strings.ToLower(list[j].Label)
	})
}

func sortNavValueSlice(list []NavItem) {
	sort.SliceStable(list, func(i, j int) bool {
		if list[i].Order != list[j].Order {
			return list[i].Order < list[j].Order
		}
		return strings.ToLower(list[i].Label) < strings.ToLower(list[j].Label)
	})
}
