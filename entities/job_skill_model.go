package entities

import "github.com/google/uuid"

type JobSkill struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key;not null" json:"id"`
	JobID   uuid.UUID `gorm:"type:uuid" json:"job_id"`
	SkillID uuid.UUID `gorm:"type:uuid" json:"skill_id"`

	Job   *Job   `gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE"`
	Skill *Skill `gorm:"foreignKey:SkillID;constraint:OnDelete:CASCADE"`
	Timestamp
}
