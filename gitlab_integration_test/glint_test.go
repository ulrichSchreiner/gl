package glint

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/ulrichSchreiner/gl"
)

// run this tests with: go test -short=false

var endpoint = flag.String("socket", "unix:///var/run/docker.sock", "the docker socket to use")
var refresh = flag.Bool("refresh", false, "set to true to create a new container")

var testLog = log.New(os.Stdout, "TEST", log.LstdFlags)

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
		if len(parm) == 0 {
			t.Fatal(msg)
		} else {
			t.Fatalf(msg, parm)
		}
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

func createOrStartContainer(t *testing.T, client *docker.Client, contname string) (*gl.Client, *docker.Container) {
	conts, e := client.ListContainers(docker.ListContainersOptions{All: true})
	checkErrorCondition(t, e != nil, "cannot list containers")
	var id string
	for _, c := range conts {
		for _, n := range c.Names {
			if n[1:] == contname {
				id = c.ID
			}
		}
	}
	if *refresh && id != "" {
		// we have an old container but we should create a new one
		client.RemoveContainer(docker.RemoveContainerOptions{
			ID:            id,
			RemoveVolumes: true,
			Force:         true,
		})
		id = ""
	}
	if id == "" {
		cfg := docker.Config{
			Image: "ulrichschreiner/gitlabdev",
		}
		opts := docker.CreateContainerOptions{
			Config: &cfg,
			Name:   contname,
		}
		cnt, e := client.CreateContainer(opts)

		checkErrorCondition(t, e != nil, "cannot create docker container")
		e = client.StartContainer(cnt.ID, nil)
		checkErrorCondition(t, e != nil, "cannot start docker container")
		id = cnt.ID
	}
	cnt, e := client.InspectContainer(id)
	checkErrorCondition(t, e != nil, "cannot inspect docker container")

	gitlabURL := "http://" + cnt.NetworkSettings.IPAddress + ":8080"
	e = waitForGitlabDev(gitlabURL, 60)
	checkErrorCondition(t, e != nil, "cannot start gitlab server")
	gitlab, e := gl.OpenV3(gitlabURL)
	checkErrorCondition(t, e != nil, "cannot open gitlabV3 API url")
	return gitlab, cnt
}

func TestGitlab(t *testing.T) {
	if testing.Short() {
		t.Skip("this is a long running integration test, not for short mode")
	}
	t.Log("Fire up a gitlab server ...")
	flag.Parse()
	client, e := docker.NewClient(*endpoint)
	checkErrorCondition(t, e != nil, "cannot create docker client")
	gitlab, _ := createOrStartContainer(t, client, "gitlabtest_container")

	/*rm := docker.RemoveContainerOptions{
		ID:            cnt.ID,
		RemoveVolumes: true,
		Force:         true,
	}
	defer client.RemoveContainer(rm)*/

	usr, e := gitlab.Session("root", nil, "start123")
	checkErrorCondition(t, e != nil, "cannot open root session")
	git := gitlab.Child()
	git.Token(usr.PrivateToken)
	projects := createProjects(t, git, TESTPROJECT, 20)
	defer removeProjectsWithId(t, git, projects)

	listAllProjects(t, git, 20)
	fetchSingleProject(t, git, "root", "testproject_0_1_2_3")
	fetchProjectsPaged(t, git, 5, 20)

	testUsersAndMembers(t, git, projects[0])
	testGroups(t, git)
	testRepositories(t, git)

}

func listAllProjects(t *testing.T, g *gl.Client, numExp int) {
	prjs, e := g.AllProjects()
	checkErrorCondition(t, e != nil, "cannot query projects")
	checkErrorCondition(t, len(prjs) != numExp, "wrong number of projects")
}

func removeProjectsWithId(t *testing.T, git *gl.Client, prjs []gl.Project) {
	t.Logf("Removing Projects from Gitlab with their ID's")
	for _, p := range prjs {
		e := git.RemoveProject(p.Id)
		checkErrorCondition(t, e != nil, "cannot remove project '%s': %s", p.Name, e)
	}
}

func fetchSingleProject(t *testing.T, git *gl.Client, ns, name string) {
	pname := url.QueryEscape(fmt.Sprintf("%s/%s", ns, name))
	_, err := git.Project(pname)
	checkErrorCondition(t, err != nil, "cannot fetch project %s/%s: %s", ns, pname, err)
}

func fetchProjectsPaged(t *testing.T, git *gl.Client, pg, all int) {
	page := gl.Page{Page: 0, PerPage: pg}
	fetched := 0
	for {
		prj, pag, err := git.Projects(&page)
		checkErrorCondition(t, err != nil, "cannot fetch projects with page: %d: %s", page.Page, err)
		checkErrorCondition(t, len(prj) != page.PerPage, "returned projects differ from Pagecount %d <> %d", len(prj), page.PerPage)
		fetched += page.PerPage
		if fetched >= all {
			break
		}
		page = *pag.NextPage
	}
}

func createProjects(t *testing.T, git *gl.Client, templ gl.Project, num int) []gl.Project {
	var projects []gl.Project
	t.Logf("Create %d projects in gitlab", num)
	for i := 0; i < num; i++ {
		templ.Name = fmt.Sprintf("%s_%d", templ.Name, i)
		pr, e := git.CreateProject(
			templ.Name, nil, nil,
			&templ.Description,
			templ.IssuesEnabled,
			templ.MergeRequestsEnabled,
			templ.WikiEnabled,
			templ.SnippetsEnabled,
			templ.Public,
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
