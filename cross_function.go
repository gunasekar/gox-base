package gox

import (
	"github.com/devlibx/gox-base/metrics"
	"github.com/devlibx/gox-base/util"
	"go.uber.org/zap"
)

// Implementation of cross function
type crossFunction struct {
	logger *zap.Logger
	metrics.Scope
	TimeService
	config      StringObjectMap
	timeTracker util.TimeTracker
}

func (c *crossFunction) TimeTracker() util.TimeTracker {
	return c.timeTracker
}

func (c *crossFunction) Metric() metrics.Scope {
	return c.Scope
}

func (c *crossFunction) Logger() *zap.Logger {
	return c.logger
}

func (c *crossFunction) Config() StringObjectMap {
	return c.config
}

// Create a new cross function object
func NewCrossFunction(args ...interface{}) CrossFunction {
	obj := crossFunction{}
	for _, arg := range args {
		switch o := arg.(type) {
		case *zap.Logger:
			obj.logger = o
		case metrics.Scope:
			obj.Scope = o
		case StringObjectMap:
			obj.config = o
		case util.TimeTracker:
			obj.timeTracker = o
		}
	}

	// Set default time-service
	if obj.TimeService == nil {
		obj.TimeService = &DefaultTimeService{}
	}

	// Setup no-op logger if it is not passed
	if obj.logger == nil {
		obj.logger = zap.NewNop()
	}

	// Setup no-op metrics
	if obj.Scope == nil {
		obj.Scope = metrics.NoOpMetric()
	}

	// Set default config
	if obj.config == nil {
		obj.config = StringObjectMap{}
	}

	// Set dummy trim tracker if not provided
	if obj.timeTracker == nil {
		obj.timeTracker = util.NewNoOpTimeTracker()
	}

	return &obj
}

// A No Op cross function
func NewNoOpCrossFunction(args ...interface{}) CrossFunction {
	obj := crossFunction{TimeService: &DefaultTimeService{}}
	obj.logger = zap.NewNop()
	obj.TimeService = &DefaultTimeService{}
	obj.Scope = metrics.NoOpMetric()
	obj.config = StringObjectMap{}
	obj.timeTracker = util.NewNoOpTimeTracker()
	return &obj
}
