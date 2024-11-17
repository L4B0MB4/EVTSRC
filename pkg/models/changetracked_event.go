package models

type ChangeTrackedEvent struct {
	Event
	IsNew bool `json:"-"`
}
