package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Role      string         `json:"role" gorm:"default:'parent'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Child struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Age       int            `json:"age"`
	ParentID  uint           `json:"parent_id"`
	Parent    User           `json:"-" gorm:"foreignKey:ParentID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Progress struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	ChildID     uint       `json:"child_id"`
	LessonID    uint       `json:"lesson_id"`
	Completed   bool       `json:"completed" gorm:"default:false"`
	Score       int        `json:"score" gorm:"default:0"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
