package logging

import "go.uber.org/zap"

// Configure creates and configures a new `zap` logger and returns it.
func Configure() (*zap.Logger, error) {
	// TODO: Add necessary configuration to the zap instance.
	return zap.NewDevelopment()
}
