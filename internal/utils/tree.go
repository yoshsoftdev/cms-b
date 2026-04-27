package utils

import "cms/internal/model"

type NavigationTree struct {
	ID       string           `json:"id"`
	Label    string           `json:"label"`
	URL      string           `json:"url"`
	Order    int              `json:"order"`
	Type     string           `json:"type"`
	Target   string           `json:"target"`
	IsActive bool             `json:"is_active"`
	Location string           `json:"location"`
	Children []NavigationTree `json:"children"`
}

func BuildNavigationTree(items []model.Navigation) []NavigationTree {
	childrenMap := map[string][]model.Navigation{}
	roots := make([]model.Navigation, 0)

	for _, item := range items {
		if item.ParentID == nil {
			roots = append(roots, item)
			continue
		}
		childrenMap[item.ParentID.String()] = append(childrenMap[item.ParentID.String()], item)
	}

	var build func(node model.Navigation) NavigationTree
	build = func(node model.Navigation) NavigationTree {
		childNodes := childrenMap[node.ID.String()]
		children := make([]NavigationTree, 0, len(childNodes))
		for _, child := range childNodes {
			children = append(children, build(child))
		}
		return NavigationTree{
			ID:       node.ID.String(),
			Label:    node.Label,
			URL:      node.URL,
			Order:    node.Order,
			Type:     node.Type,
			Target:   node.Target,
			IsActive: node.IsActive,
			Location: node.Location,
			Children: children,
		}
	}

	out := make([]NavigationTree, 0, len(roots))
	for _, root := range roots {
		out = append(out, build(root))
	}
	return out
}
