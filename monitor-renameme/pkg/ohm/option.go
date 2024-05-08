package ohm

import (
	"fmt"
	"strings"
)

type SensorFilters struct {
	Identifier       *string
	Name             *string
	ParentIdentifier *string
	Type             *SensorType
}

type SensorFilter func(filter *SensorFilters)

func narrowSensorsQuery(opts ...SensorFilter) string {
	if len(opts) == 0 {
		return ""
	}

	filter := SensorFilters{}

	for _, opt := range opts {
		opt(&filter)
	}

	condition := &conditionQuery{}

	if filter.Identifier != nil {
		condition.and(fmt.Sprintf("Identifier = '%s' ", *filter.Identifier))
	}
	if filter.Name != nil {
		condition.and(fmt.Sprintf("Name = '%s' ", *filter.Name))
	}
	if filter.ParentIdentifier != nil {
		condition.and(fmt.Sprintf("Parent = '%s' ", *filter.ParentIdentifier))
	}
	if filter.Type != nil {
		condition.and(fmt.Sprintf("SensorType = '%s' ", *filter.Type))
	}

	return condition.String()
}

func SensorsWithIdentifier(val string) SensorFilter {
	return func(filter *SensorFilters) {
		filter.Identifier = &val
	}
}

func SensorsWithName(val string) SensorFilter {
	return func(filter *SensorFilters) {
		filter.Name = &val
	}
}

func SensorsWithType(val SensorType) SensorFilter {
	return func(filter *SensorFilters) {
		filter.Type = &val
	}
}

func SensorsWithParentIdentifier(val string) SensorFilter {
	return func(filter *SensorFilters) {
		filter.ParentIdentifier = &val
	}
}

type HardwareFilters struct {
	Type *HardwareType
}

type HardwareFilter func(filter *HardwareFilters)

func narrowHardwareQuery(opts ...HardwareFilter) string {
	if len(opts) == 0 {
		return ""
	}

	filter := HardwareFilters{}

	for _, opt := range opts {
		opt(&filter)
	}

	condition := &conditionQuery{}

	if filter.Type != nil {
		condition.and(fmt.Sprintf("HardwareType = '%s' ", *filter.Type))
	}

	return condition.String()
}

func HardwareWithType(val HardwareType) HardwareFilter {
	return func(filter *HardwareFilters) {
		filter.Type = &val
	}
}

type conditionQuery struct {
	b         strings.Builder
	hasPrefix bool
}

func (c *conditionQuery) and(val string) {
	if c.hasPrefix {
		c.b.WriteString("AND ")
	} else {
		c.b.WriteString("WHERE ")
		c.hasPrefix = true
	}
	c.b.WriteString(val)
}

func (c *conditionQuery) String() string {
	return c.b.String()
}
