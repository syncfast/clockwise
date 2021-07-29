package tui

import "time"

// calculateCost calculates the total cost.
func calculateCost(data *Data, averageSalary int) {
	for {
		count := data.GetCount()
		total := data.getCost()
		// FLC explained:
		// https://smallbusiness.chron.com/calculate-fully-burdened-labor-costs-33072.html
		// TODO: Make FLC configurable via the config file.
		fullyLoadedCostMultiplier := float32(1.75)
		costPer500ms := float32(count) * fullyLoadedCostMultiplier * float32(averageSalary) / 7200000 / 2 // 50 (weeks) * 40 (hours) * 60 (minutes) * 60 (seconds) / 2
		total += costPer500ms
		data.setCost(total)
		time.Sleep(time.Millisecond * 500)
	}
}
