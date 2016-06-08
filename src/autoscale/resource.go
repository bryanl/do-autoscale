package autoscale

import (
	"time"

	"golang.org/x/net/context"
)

type ResourceHistory struct {
	Iteration int     `json:"iteration"`
	Value     float64 `json:"value"`
}

// ResourceAllocation is information about an allocated resource.
type ResourceAllocation struct {
	Name      string            `json:"name"`
	Address   string            `json:"address"`
	CreatedAt time.Time         `json:"createdAt"`
	History   []ResourceHistory `json:"history"`
}

// ResourceManager is a watched resource interface.
type ResourceManager interface {
	Count() (int, error)
	Scale(ctx context.Context, g Group, byN int, repo Repository) (bool, error)
	Allocated() ([]ResourceAllocation, error)
}
