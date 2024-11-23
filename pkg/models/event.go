package models

type Event struct {
	Id            string `json:"id"`
	Version       int64  `json:"version" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Data          []byte `json:"data" binding:"required"`
	AggregateId   string `json:"aggregateId"`
	AggregateType string `json:"aggregateType" binding:"required"`
}
