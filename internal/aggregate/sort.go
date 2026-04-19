package aggregate

import "sort"

func sortWindows(ws []Window) {
	sort.Slice(ws, func(i, j int) bool {
		return ws[i].Start.Before(ws[j].Start)
	})
}
