package tui

import "time"

// calculateCost calculates the total cost.
func calculateCost(data *Data) {
	for {
		count := data.GetCount()
		total := data.getCost()
		// FLC explained:
		// https://smallbusiness.chron.com/calculate-fully-burdened-labor-costs-33072.html
		// TODO: Make FLC configurable via the config file.
		fullyLoadedCostMultiplier := float32(1.75)
		cps := float32(count) * fullyLoadedCostMultiplier * float32(150000) / 7200000 // 50 (weeks) * 40 (hours) * 60 (minutes) * 60 (seconds)
		total += cps
		data.setCost(total)
		time.Sleep(refreshInterval)
	}
}
