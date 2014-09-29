package glint

import (
	"flag"
	"github.com/fsouza/go-dockerclient"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/ulrichSchreiner/gl"
	"testing"
	"time"
)

// run this tests with: go test -short=false

var endpoint = flag.String("socket", "unix:///var/run/docker.sock", "the docker socket to use")

var TESTPROJECT = gl.Project{
	Name:                 "testproject",
	Description:          "my test description",
	IssuesEnabled:        true,
	WikiEnabled:          true,
	MergeRequestsEnabled: true,
	SnippetsEnabled:      true,
	Public:               true,
}

func checkProject(p1, p2 *gl.Project) {
	So(p1.Name, ShouldEqual, p2.Name)
	So(p1.Description, ShouldEqual, p2.Description)
	So(p1.IssuesEnabled, ShouldEqual, p2.IssuesEnabled)
	So(p1.WikiEnabled, ShouldEqual, p2.WikiEnabled)
	So(p1.MergeRequestsEnabled, ShouldEqual, p2.MergeRequestsEnabled)
	So(p1.SnippetsEnabled, ShouldEqual, p2.SnippetsEnabled)
}

func TestGitlab(t *testing.T) {
	Convey("fire up a gitlab server ...", t, func() {
		if testing.Short() {
			t.Skip("this is a long running integration test, not for short mode")
		}
		flag.Parse()
		client, e := docker.NewClient(*endpoint)
		So(e, ShouldBeNil)
		cfg := docker.Config{
			Image: "ulrichschreiner/gitlabdev",
		}
		opts := docker.CreateContainerOptions{
			Config: &cfg,
		}
		cnt, e := client.CreateContainer(opts)

		So(e, ShouldBeNil)
		e = client.StartContainer(cnt.ID, nil)
		So(e, ShouldBeNil)
		cnt, e = client.InspectContainer(cnt.ID)
		So(e, ShouldBeNil)
		/*
			rm := docker.RemoveContainerOptions{
				ID:            cnt.ID,
				RemoveVolumes: true,
				Force:         true,
			}
			defer client.RemoveContainer(rm)*/
		gitlabURL := "http://" + cnt.NetworkSettings.IPAddress + ":8080"
		//fmt.Printf("%#v\n", cnt.NetworkSettings)
		time.Sleep(time.Second * 15) // git gitlab time to start ...
		gitlab, e := gl.OpenV3(gitlabURL)
		So(e, ShouldBeNil)
		usr, e := gitlab.Session("root", nil, "5iveL!fe")
		So(e, ShouldBeNil)
		git := gitlab.Child()
		git.Token(usr.PrivateToken)
		pr, e := git.CreateProject(
			TESTPROJECT.Name, nil, nil,
			&TESTPROJECT.Description,
			&TESTPROJECT.IssuesEnabled,
			&TESTPROJECT.MergeRequestsEnabled,
			&TESTPROJECT.WikiEnabled,
			&TESTPROJECT.SnippetsEnabled,
			&TESTPROJECT.Public,
			nil, nil)
		checkProject(&TESTPROJECT, pr)
		So(pr.Public, ShouldEqual, TESTPROJECT.Public)
		listAllProjects(git, 1)
		rmp, e := git.RemoveProject(pr.Id)
		So(e, ShouldBeNil)
		checkProject(&TESTPROJECT, rmp)
	})
}

func listAllProjects(g *gl.Client, numExp int) {
	prjs, e := g.AllProjects()
	So(e, ShouldBeNil)
	So(len(prjs), ShouldEqual, numExp)
}
