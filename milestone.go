package gl

import (
	"time"
)

type MilestoneStateEvent string

const (
	MilestoneActivate = MilestoneStateEvent("activate")
	MilestoneClose    = MilestoneStateEvent("close")
)

type Milestone struct {
	Id          int       `json:"id,omitempty"`
	Iid         int       `json:"iid,omitempty"`
	ProjectId   int       `json:"project_id,omitempty"`
	Description string    `json:"description,omitempty"`
	DueDate     *JsonDate `json:"due_date,omitempty"`
	State       string    `json:"state,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	Title       string    `json:"title,omitempty"`
}

const (
	milestones_url = "/projects/:id/milestones"
	milestone_url  = "/projects/:id/milestones/:milestone_id"
)

type Milestones []Milestone
