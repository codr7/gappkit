package bookng

import (
	"fmt"
	"time"
)

type Overbook struct {
	resource *Resource
	startTime, endTime time.Time
}

func NewOverbook(resource *Resource, startTime, endTime time.Time) *Overbook {
	return &Overbook{resource: resource, startTime: startTime, endTime: endTime}
}

func (self *Overbook) Error() string {
	return fmt.Sprintf("Available quantity exceeded for resource '%v' between %v and %v",
		self.resource.Name, self.startTime, self.endTime)
}
