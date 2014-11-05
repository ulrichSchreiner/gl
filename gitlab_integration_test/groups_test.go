package glint

import (
	"fmt"
	"github.com/ulrichSchreiner/gl"
	"testing"
)

const (
	numGroups = 5
)

func testGroups(t *testing.T, git *gl.Client) {
	u, e := git.CreateUser("test@example.com", "username", "start123", "myname", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	checkErrorCondition(t, e != nil, "cannot create user 'username': '%s'", e)
	defer func() {
		git.DeleteUser(u.Id)
	}()
	u, e = git.EditUser(u.Id, "test@example.com", "username", "start123", "myname", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	checkErrorCondition(t, e != nil, "cannot edit user 'username': '%s'", e)
	groups := createGroups(t, git)
	defer func() {
		for _, g := range groups {
			git.DeleteGroup(g.Id)
		}
	}()

	readGroups(t, git, groups)
	testGroupMember(t, git, groups[0], u)
	testTransfer(t, git, u)
}

func readGroups(t *testing.T, git *gl.Client, grps gl.Groups) {
	gr, e := git.AllGroups()
	checkErrorCondition(t, e != nil, "cannot read all groups: %s", e)
	checkErrorCondition(t, len(gr) != len(grps), "the groups differ from the expected groups")
}

func createGroups(t *testing.T, git *gl.Client) gl.Groups {
	var res gl.Groups
	for i := 0; i < numGroups; i++ {
		n := fmt.Sprintf("gname_%d", i)
		p := fmt.Sprintf("gpath_%d", i)
		g, e := git.AddGroup(n, p)
		checkErrorCondition(t, e != nil, "cannot create group '%s/%s': '%s'", n, p, e)
		res = append(res, *g)
	}
	return res
}

func testGroupMember(t *testing.T, git *gl.Client, g gl.Group, u *gl.User) {
	gm, e := git.AddGroupMember(g.Id, u.Id, gl.Master)
	checkErrorCondition(t, e != nil, "cannot add group member: %s", e)
	checkErrorCondition(t, gm.Name != "myname", "email differs: %s", gm.Email)
	checkErrorCondition(t, gm.Username != "username", "username differs: %s", gm.Username)
	gms, e := git.AllGroupMembers(g.Id)
	checkErrorCondition(t, e != nil, "cannot query all group members: %s", e)
	checkErrorCondition(t, len(gms) != 1, "there must be exact one group member")
	checkErrorCondition(t, gm.Name != gms[0].Name, "email differs: %s", gm.Email)
	checkErrorCondition(t, gm.Username != gms[0].Username, "username differs: %s", gm.Username)

	e = git.DeleteGroupMember(g.Id, u.Id)
	checkErrorCondition(t, e != nil, "cannot delete group member: %s", e)
}

func testTransfer(t *testing.T, git *gl.Client, u *gl.User) {
	tp := TESTPROJECT
	git2 := git.Child()
	git2.Sudo(u.Username)
	pr, e := git2.CreateProject(
		"testuserproject", nil, nil,
		&tp.Description,
		&tp.IssuesEnabled,
		&tp.MergeRequestsEnabled,
		&tp.WikiEnabled,
		&tp.SnippetsEnabled,
		&tp.Public,
		nil, nil)

	checkErrorCondition(t, e != nil, "cannot create project: '%s'", e)

	g, e := git2.AddGroup("testgroup", "testpath")
	checkErrorCondition(t, e != nil, "cannot create group: %s", e)

	git2.AddGroupMember(g.Id, u.Id, gl.Owner)
	git2.SetLogger(testLog)
	t.Logf("Group=%#v\nProject=%#v\n", g, pr)
	g, e = git2.TransferProjectToGroup(g.Id, pr.Id)
	checkErrorCondition(t, e != nil, "cannot transfer project to group: %s", e)
	t.Logf("g = %+v", g)
}
