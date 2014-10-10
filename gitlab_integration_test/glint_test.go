package glint

import (
	"flag"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/ulrichSchreiner/gl"
	"net/http"
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

func checkErrorCondition(t *testing.T, cond bool, msg string, parm ...interface{}) {
	if cond {
		t.Fatalf(msg, parm)
	}
}

func checkProject(t *testing.T, p1, p2 *gl.Project, withpublic bool) {
	checkErrorCondition(t, p1.Name != p2.Name, "Name of project differs: '%s' <> '%s'", p1.Name, p2.Name)
	checkErrorCondition(t, p1.Description != p2.Description, "Description of project differs: '%s' <> '%s'", p1.Description, p2.Description)
	checkErrorCondition(t, p1.IssuesEnabled != p2.IssuesEnabled, "IssuesEnabled of project differs: '%v' <> '%v'", p1.IssuesEnabled, p2.IssuesEnabled)
	checkErrorCondition(t, p1.IssuesEnabled != p2.IssuesEnabled, "IssuesEnabled of project differs: '%v' <> '%v'", p1.IssuesEnabled, p2.IssuesEnabled)
	checkErrorCondition(t, p1.MergeRequestsEnabled != p2.MergeRequestsEnabled, "MergeRequestsEnabled of project differs: '%v' <> '%v'", p1.MergeRequestsEnabled, p2.MergeRequestsEnabled)
	checkErrorCondition(t, p1.SnippetsEnabled != p2.SnippetsEnabled, "SnippetsEnabled of project differs: '%v' <> '%v'", p1.SnippetsEnabled, p2.SnippetsEnabled)
	if withpublic {
		checkErrorCondition(t, p1.Public != p2.Public, "Public of project differs: '%v' <> '%v'", p1.Public, p2.Public)
	}
}

func TestGitlab(t *testing.T) {
	if testing.Short() {
		t.Skip("this is a long running integration test, not for short mode")
	}
	t.Log("Fire up a gitlab server ...")
	flag.Parse()
	client, e := docker.NewClient(*endpoint)
	checkErrorCondition(t, e != nil, "cannot create docker client")
	cfg := docker.Config{
		Image: "ulrichschreiner/gitlabdev",
	}
	opts := docker.CreateContainerOptions{
		Config: &cfg,
	}
	cnt, e := client.CreateContainer(opts)

	checkErrorCondition(t, e != nil, "cannot create docker container")
	e = client.StartContainer(cnt.ID, nil)
	checkErrorCondition(t, e != nil, "cannot start docker container")
	cnt, e = client.InspectContainer(cnt.ID)
	checkErrorCondition(t, e != nil, "cannot inspect docker container")

	rm := docker.RemoveContainerOptions{
		ID:            cnt.ID,
		RemoveVolumes: true,
		Force:         true,
	}
	defer client.RemoveContainer(rm)

	gitlabURL := "http://" + cnt.NetworkSettings.IPAddress + ":8080"
	//fmt.Printf("%#v\n", cnt.NetworkSettings)
	e = waitForGitlabDev(gitlabURL, 60)
	checkErrorCondition(t, e != nil, "cannot start gitlab server")
	gitlab, e := gl.OpenV3(gitlabURL)
	checkErrorCondition(t, e != nil, "cannot open gitlabV3 API url")
	usr, e := gitlab.Session("root", nil, "5iveL!fe")
	checkErrorCondition(t, e != nil, "cannot open root session")
	git := gitlab.Child()
	git.Token(usr.PrivateToken)
	projects := createProjects(t, git, TESTPROJECT, 5)
	listAllProjects(t, git, 5)
	removeProjectsWithId(t, git, projects)
}

func listAllProjects(t *testing.T, g *gl.Client, numExp int) {
	prjs, e := g.AllProjects()
	checkErrorCondition(t, e != nil, "cannot query projects")
	checkErrorCondition(t, len(prjs) != numExp, "wrong number of projects")
}

func removeProjectsWithId(t *testing.T, git *gl.Client, prjs []gl.Project) {
	t.Logf("Removing Projects from Gitlab with their ID's")
	for _, p := range prjs {
		t.Logf("Removing '%s' with id '%d'", p.Name, p.Id)
		rmp, e := git.RemoveProject(p.Id)
		checkErrorCondition(t, e != nil, "cannot remove project '%s'", p.Name)
		checkProject(t, &p, rmp, false)
	}
}

func createProjects(t *testing.T, git *gl.Client, templ gl.Project, num int) []gl.Project {
	var projects []gl.Project
	t.Logf("Create %d projects in gitlab", num)
	for i := 0; i < num; i++ {
		templ.Name = fmt.Sprintf("%s_%d", templ.Name, i)
		t.Logf("Creating %s ...", templ.Name)
		pr, e := git.CreateProject(
			templ.Name, nil, nil,
			&templ.Description,
			&templ.IssuesEnabled,
			&templ.MergeRequestsEnabled,
			&templ.WikiEnabled,
			&templ.SnippetsEnabled,
			&templ.Public,
			nil, nil)
		checkErrorCondition(t, e != nil, "cannot create project: '%s'", e)
		projects = append(projects, *pr)
		checkProject(t, &templ, pr, true)
	}
	return projects
}

func waitForGitlabDev(url string, maxWait int) error {
	wait := 0
	for {
		_, err := http.Get(url)
		if err == nil {
			return nil
		}
		if wait > maxWait {
			return err
		}
		time.Sleep(time.Second * 5)
		wait += 5
	}
}
