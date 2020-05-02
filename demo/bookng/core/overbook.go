package core

import (
	"fmt"
	"time"
)

type Overbook struct {
	resource *Resource
	startTime, endTime time.Time
	quantity int64
}

func NewOverbook(resource *Resource, startTime, endTime time.Time, quantity int64) *Overbook {
	return &Overbook{resource: resource, startTime: startTime, endTime: endTime, quantity: quantity}
}

func (self *Overbook) Error() string {
	return fmt.Sprintf("Available quantity exceeded (%v) for resource %v between %v and %v",
		self.resource.Name, self.startTime, self.endTime, self.quantity)
}
