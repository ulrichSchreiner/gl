package glint

import (
	"fmt"
	"github.com/ulrichSchreiner/gl"
	"testing"
)

const (
	numGroups = 5
)

func testGroups(t *testing.T, git *gl.Client, p gl.Project) {
	groups := createGroups(t, git)
	t.Logf("Groups: %#v", groups)
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
