package gl

import (
	"net/url"
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

func (g *Client) Milestones(pid string, pg *Page) (Milestones, *Pagination, error) {
	var p Milestones
	u := expandUrl(milestones_url, map[string]interface{}{":id": pid})
	pager, e := g.get(u, nil, pg, &p)
	if e != nil {
		return nil, nil, e
	}
	return p, pager, nil
}

func (g *Client) AllMilestones(pid string) (Milestones, error) {
	var r Milestones
	err := fetchAll(func(pg *Page) (interface{}, *Pagination, error) {
		return g.Milestones(pid, pg)
	}, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (g *Client) Milestone(pid string, mid int) (*Milestone, error) {
	u := expandUrl(milestone_url, map[string]interface{}{":id": pid, ":milestone_id": mid})
	var m Milestone
	_, e := g.get(u, nil, nil, &m)
	return &m, e
}

func (g *Client) CreateMilestone(pid string, title string, description *string, duedate *time.Time) (*Milestone, error) {
	u := expandUrl(milestones_url, map[string]interface{}{":id": pid})
	vals := make(url.Values)
	vals.Set("title", title)
	addString(vals, "description", description)
	if duedate != nil {
		dd := duedate.Format(dateLayout)
		vals.Set("due_date", dd)
	}
	var m Milestone
	e := g.post(u, vals, &m)
	return &m, e
}

func (g *Client) UpdateMilestone(pid string, mid int, title, description *string, duedate *time.Time, state *MilestoneStateEvent) (*Milestone, error) {
	u := expandUrl(milestone_url, map[string]interface{}{":id": pid, ":milestone_id": mid})
	vals := make(url.Values)
	addString(vals, "title", title)
	addString(vals, "description", description)
	if duedate != nil {
		dd := duedate.Format(dateLayout)
		vals.Set("due_date", dd)
	}
	if state != nil {
		vals.Set("state_event", string(*state))
	}
	var m Milestone
	e := g.put(u, vals, &m)
	return &m, e
}
