package bookng

import (
	"fmt"
	"time"
)

type Overbook struct {
	resource *Resource
	startTime, endTime time.Time
	available int64
}

func NewOverbook(resource *Resource, startTime, endTime time.Time, available int64) *Overbook {
	return &Overbook{resource: resource, startTime: startTime, endTime: endTime, available: available}
}

func (self *Overbook) Error() string {
	return fmt.Sprintf("Available quantity exceeded for resource '%v' between %v and %v: %v",
		self.resource.Name, self.startTime, self.endTime, self.available)
}
