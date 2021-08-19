package tui

import (
	"github.com/spf13/viper"
)

// Static configuration variables, enumerated at the start of the application
type Config struct {
	manualMode     bool
	averageSalary  int
	currencySymbol string
}

// Receive configuration values from Viper
func (c *Config) SetFromViperConfig(manualMode bool) {
	c.manualMode = manualMode
	c.averageSalary = viper.GetViper().GetInt("averageSalary")
	c.currencySymbol = viper.GetViper().GetString("currencySymbol")
}

// Set manual mode
func (c *Config) SetManualMode(value bool) {
	c.manualMode = value
}

// Get manual mode
func (c *Config) GetManualMode() bool {
	return c.manualMode
}

// Set average salary
func (c *Config) SetAverageSalary(value int) {
	c.averageSalary = value
}

// Get average salary
func (c *Config) GetAverageSalary() int {
	return c.averageSalary
}

// Set currency symbol
func (c *Config) SetCurrencySymbol(value string) {
	c.currencySymbol = value
}

// Get currency symbol
func (c *Config) GetCurrencySymbol() string {
	return c.currencySymbol
}
