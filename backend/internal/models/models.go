package models

import "time"

// EmergencyPlan 应急预案主体
type EmergencyPlan struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:200;not null" json:"name"`
	PlanType  string    `gorm:"column:plan_type;size:40;not null" json:"plan_type"`
	Scenario  string    `gorm:"type:text" json:"scenario"`
	CreatedAt time.Time `json:"created_at"`
}

func (EmergencyPlan) TableName() string { return "emergency_plans" }

// PlanVersion 预案版本与审批状态
type PlanVersion struct {
	ID              uint64     `gorm:"primaryKey" json:"id"`
	EmergencyPlanID uint64     `gorm:"column:emergency_plan_id;not null;index" json:"emergency_plan_id"`
	VersionNo       string     `gorm:"column:version_no;size:32;not null" json:"version_no"`
	RevisionRecord  string     `gorm:"column:revision_record;type:text" json:"revision_record"`
	ContentMd       string     `gorm:"column:content_md;type:text" json:"content_md"`
	PublishedDate   *time.Time `gorm:"column:published_date;type:date" json:"published_date"`
	Preparer        string     `gorm:"size:120" json:"preparer"`
	PreparedDate    *time.Time `gorm:"column:prepared_date;type:date" json:"prepared_date"`
	Approver        string     `gorm:"size:120" json:"approver"`
	ApprovedDate    *time.Time `gorm:"column:approved_date;type:date" json:"approved_date"`
	ApprovalStatus  string     `gorm:"column:approval_status;size:20;not null;default:draft" json:"approval_status"`
	IsCurrent       bool       `gorm:"column:is_current;not null;default:false" json:"is_current"`
	CreatedAt       time.Time  `json:"created_at"`
	Plan            *EmergencyPlan `gorm:"foreignKey:EmergencyPlanID" json:"plan,omitempty"`
}

func (PlanVersion) TableName() string { return "plan_versions" }

// Drill 演练计划与完成情况
type Drill struct {
	ID                   uint64     `gorm:"primaryKey" json:"id"`
	EmergencyPlanID      uint64     `gorm:"column:emergency_plan_id;not null;index" json:"emergency_plan_id"`
	DrillKind            string     `gorm:"column:drill_kind;size:32;not null" json:"drill_kind"`
	ScheduledDate        time.Time  `gorm:"column:scheduled_date;type:date;not null" json:"scheduled_date"`
	Location             string     `gorm:"size:300" json:"location"`
	OrgDept              string     `gorm:"column:org_dept;size:200" json:"org_dept"`
	ParticipantScope     string     `gorm:"column:participant_scope;size:500" json:"participant_scope"`
	Objectives           string     `gorm:"type:text" json:"objectives"`
	NotifyDepts          string     `gorm:"column:notify_depts;type:text" json:"notify_depts"`
	DrillStatus          string     `gorm:"column:drill_status;size:24;not null;default:planned" json:"drill_status"`
	ActualParticipants   *int       `gorm:"column:actual_participants" json:"actual_participants"`
	DurationMinutes      *int       `gorm:"column:duration_minutes" json:"duration_minutes"`
	ProcessDescription   string     `gorm:"column:process_description;type:text" json:"process_description"`
	ProblemList          string     `gorm:"column:problem_list;type:text" json:"problem_list"`
	Evaluation           string     `gorm:"size:20" json:"evaluation"`
	PhotoPaths           string     `gorm:"column:photo_paths;type:text" json:"photo_paths"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	Plan                 *EmergencyPlan `gorm:"foreignKey:EmergencyPlanID" json:"plan,omitempty"`
}

func (Drill) TableName() string { return "drills" }

// DrillIssue 演练发现问题
type DrillIssue struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	DrillID     uint64    `gorm:"column:drill_id;not null;index" json:"drill_id"`
	Description string    `gorm:"type:text;not null" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Drill       *Drill    `gorm:"foreignKey:DrillID" json:"drill,omitempty"`
}

func (DrillIssue) TableName() string { return "drill_issues" }

// Rectification 整改闭环
type Rectification struct {
	ID                uint64     `gorm:"primaryKey" json:"id"`
	DrillIssueID      uint64     `gorm:"column:drill_issue_id;not null;index" json:"drill_issue_id"`
	ResponsiblePerson string     `gorm:"column:responsible_person;size:120" json:"responsible_person"`
	CorrectiveMeasure string     `gorm:"column:corrective_measure;type:text" json:"corrective_measure"`
	DueDate           *time.Time `gorm:"column:due_date;type:date" json:"due_date"`
	Status            string     `gorm:"size:20;not null;default:pending" json:"status"`
	CompletedAt       *time.Time `gorm:"column:completed_at;type:date" json:"completed_at"`
	CreatedAt         time.Time  `json:"created_at"`
	Issue             *DrillIssue `json:"issue,omitempty" gorm:"foreignKey:DrillIssueID"`
}

func (Rectification) TableName() string { return "rectifications" }
