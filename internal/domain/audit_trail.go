package domain

import "time"

type AuditTrail struct {
	ID          string    `gorm:"primaryKey;column:id" json:"id"`
	ServiceName string    `gorm:"column:service_name" json:"service_name"`
	TableName   string    `gorm:"column:table_name" json:"table_name"`
	EntityID    string    `gorm:"column:entity_id" json:"entity_id"`
	Message     string    `gorm:"column:message" json:"message"`
	Data        string    `gorm:"column:data" json:"data"`
	CreatedBy   string    `gorm:"column:created_by" json:"created_by"`
	UpdatedBy   string    `gorm:"column:updated_by" json:"update_by"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type AuditTrailFilter struct {
	ID          string `json:"id"`
	EntityID    string `json:"entity_id"`
	ServiceName string `json:"service_name"`
}
