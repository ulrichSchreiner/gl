package glint

import (
	"fmt"
	"github.com/ulrichSchreiner/gl"
	"testing"
)

const (
	numGroups = 5
)

func testGroups(t *testing.T, admingit *gl.Client) {
	t.Log("create a new user for testing groups")
	u, e := admingit.CreateUser("test@example.com", "username2", "start123", "myname2", nil, nil, nil, nil, nil, nil, nil, nil, true, true)
	checkErrorCondition(t, e != nil, "cannot create user 'username': '%s'", e)
	defer func() {
		t.Log("remove testuser for group testing")
		admingit.DeleteUser(u.Id)
	}()
	usr, e := admingit.Session("username2", nil, "start123")
	checkErrorCondition(t, e != nil, "cannot open username2 session")
	git := admingit.Child()
	git.Token(usr.PrivateToken)
	t.Log("now creating groups")
	groups := createGroups(t, git)
	defer func() {
		t.Log("remove groups in group-test")
		for _, g := range groups {
			git.DeleteGroup(g.Id)
		}
	}()

	t.Log("reading groups")
	readGroups(t, git, groups)
	t.Log("set user to be a group member")
	testGroupMember(t, git, groups[0], u)
	t.Log("transfer a new project to the group")
	testTransfer(t, git, admingit, u)
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
	checkErrorCondition(t, gm.Name != "myname2", "name differs: %s", gm.Name)
	checkErrorCondition(t, gm.Username != "username2", "username differs: %s", gm.Username)
	gms, e := git.AllGroupMembers(g.Id)
	checkErrorCondition(t, e != nil, "cannot query all group members: %s", e)
	checkErrorCondition(t, len(gms) != 1, "there must be exact one group member")
	checkErrorCondition(t, gm.Name != gms[0].Name, "email differs: %s", gm.Email)
	checkErrorCondition(t, gm.Username != gms[0].Username, "username differs: %s", gm.Username)

	e = git.DeleteGroupMember(g.Id, u.Id)
	checkErrorCondition(t, e != nil, "cannot delete group member: %s", e)
}

func testTransfer(t *testing.T, git *gl.Client, admingit *gl.Client, u *gl.User) {
	tp := TESTPROJECT
	pr, e := git.CreateProject(
		"testuserproject", nil, nil,
		&tp.Description,
		tp.IssuesEnabled,
		tp.MergeRequestsEnabled,
		tp.WikiEnabled,
		tp.SnippetsEnabled,
		tp.Public,
		nil, nil)

	checkErrorCondition(t, e != nil, "cannot create project: '%s'", e)
	defer git.RemoveProject(pr.Id)

	g, e := git.AddGroup("transfer_testgroup", "transfer_testpath")
	checkErrorCondition(t, e != nil, "cannot create group: %s", e)
	defer git.DeleteGroup(g.Id)

	g, e = admingit.TransferProjectToGroup(g.Id, pr.Id)
	checkErrorCondition(t, e != nil, "cannot transfer project to group: %s", e)
}
