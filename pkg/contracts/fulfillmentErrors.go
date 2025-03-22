package contracts

import (
	"fmt"
)

type OutOfRangeTimeWindowError struct {
	startC string
	endC   string
	startO string
	endO   string
}

func (e *OutOfRangeTimeWindowError) Error() string {
	return fmt.Sprintf("violated time window - contracted: %s - %s does not contain obliged: %s - %s", e.startC, e.endC, e.startO, e.endO)
}

type RequestedCapacityExceededError struct {
	contracted float64
	requested  float64
}

func (e *RequestedCapacityExceededError) Error() string {
	return fmt.Sprintf("exceeded capacity - contracted: %.2f < requsted: %2.f", e.contracted, e.requested)
}

type InsufficientCapacityError struct {
	dispatched float64
	needed     float64
}

func (e *InsufficientCapacityError) Error() string {
	return fmt.Sprintf("insufficient capacity - dispatched: %.2f < needed: %2.f", e.dispatched, e.needed)
}

type IncorrectDirectionError struct {
	dispatched string
	needed     string
}

func (e *IncorrectDirectionError) Error() string {
	return fmt.Sprintf("incorrect direction - dispatched %s != obliged: %s", e.dispatched, e.needed)
}

type ViolatedTimeWindowError struct {
	startD string
	endD   string
	startO string
	endO   string
}

func (e *ViolatedTimeWindowError) Error() string {
	return fmt.Sprintf("violated time window - dispatch: %s - %s != obliged: %s - %s", e.startD, e.endD, e.startO, e.endO)
}
