package glint

import (
	"fmt"
	"github.com/ulrichSchreiner/gl"
	"strconv"
	"testing"
)

const numUsers = 10

func testUsersAndMembers(t *testing.T, git *gl.Client, p gl.Project) {
	users := testMembers(t, git, p)
	retrieveMembers(t, git, p, numUsers+1) // we created 10 members + root member
	removeMember(t, git, p, users[0])
	retrieveMembers(t, git, p, numUsers)
	removeUsers(t, git, p, users)
}

func testMembers(t *testing.T, git *gl.Client, p gl.Project) []gl.User {
	var res []gl.User
	t.Logf("creating users...")
	for i := 0; i < numUsers; i++ {
		username := fmt.Sprintf("user_%d", i)
		email := fmt.Sprintf("%s@example.com", username)
		u, e := git.CreateUser(email, username, "mypassword", "myname", nil, nil, nil, nil, nil, nil, nil, nil, false, true)
		checkErrorCondition(t, e != nil, "cannot create user '%s': '%s'", username, e)
		res = append(res, *u)
		git.AddTeamMember(strconv.Itoa(p.Id), u.Id, gl.Developer)
	}
	return res
}

func retrieveMembers(t *testing.T, git *gl.Client, p gl.Project, num int) {
	t.Logf("check if there are %d members", num)
	memb, e := git.AllTeamMembers(strconv.Itoa(p.Id), nil)
	checkErrorCondition(t, e != nil, "cannot get team members from project '%s': '%s'", p.Name, e)
	checkErrorCondition(t, len(memb) != num, "number of members seems incorrect %d != %d", len(memb), num)
}

func removeMember(t *testing.T, git *gl.Client, p gl.Project, u gl.User) {
	t.Logf("remove '%s' from the member list", u.Name)
	_, e := git.DeleteTeamMember(strconv.Itoa(p.Id), u.Id)
	checkErrorCondition(t, e != nil, "cannot delete team member '%s' from project '%s': '%s'", u.Name, p.Name, e)
}

func removeUsers(t *testing.T, git *gl.Client, p gl.Project, usr []gl.User) {
	t.Logf("removing all Users")
	for _, u := range usr {
		du, e := git.DeleteUser(u.Id)
		checkErrorCondition(t, e != nil, "cannot delete user '%s': '%s'", u.Name, e)
		checkUser(t, u, *du)
	}
}

func checkUser(t *testing.T, u1, u2 gl.User) {
	checkErrorCondition(t, u1.Name != u2.Name, "names of users differ: %s != %s", u1.Name, u2.Name)
}
