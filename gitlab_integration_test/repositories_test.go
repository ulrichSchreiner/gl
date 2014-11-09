package glint

import (
	"fmt"
	"github.com/ulrichSchreiner/gl"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func uploadGitRepo(t *testing.T, git *gl.Client, uid int, remote string, host string) {
	tmpdir, e := ioutil.TempDir("/tmp", "gl_test")
	checkErrorCondition(t, e != nil, "Cannot create temp direcotry for git repo: %s", e)
	defer os.RemoveAll(tmpdir)

	privkey, e := gl.GeneratePrivateKey()
	checkErrorCondition(t, e != nil, "Cannot generate private key: %s", e)

	ioutil.WriteFile(filepath.Join(tmpdir, "privatekey"), []byte(gl.MarshalPrivateKey(privkey)), 0600)
	checkErrorCondition(t, e != nil, "Cannot write a private key: %s", e)

	ioutil.WriteFile(filepath.Join(tmpdir, "git.sh"), []byte(fmt.Sprintf("/usr/bin/ssh -o StrictHostKeyChecking=no -i %s/privatekey $@", tmpdir)), 0700)
	checkErrorCondition(t, e != nil, "Cannot write a git wrapper: %s", e)

	_, e = git.CreateSshKey(uid, "ssh pubkey", gl.MarshalPublicKey(&privkey.PublicKey))
	checkErrorCondition(t, e != nil, "cannot create remote user key: %s", e)

	// ugly, but gitlab needs a little time to store the key in its filesystem so gitlab-shell can use it
	time.Sleep(5 * time.Second)

	ioutil.WriteFile(filepath.Join(tmpdir, "README.md"), []byte("Hello World"), 0777)
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpdir
	e = cmd.Run()
	checkErrorCondition(t, e != nil, "cannot init git repo: %s", e)
	cmd = exec.Command("git", "add", "README.md")
	cmd.Dir = tmpdir
	e = cmd.Run()
	checkErrorCondition(t, e != nil, "cannot add data to git repo: %s", e)
	cmd = exec.Command("git", "commit", "-m", "initial commit")
	cmd.Dir = tmpdir
	e = cmd.Run()
	checkErrorCondition(t, e != nil, "cannot commit data to git repo: %s", e)
	cmd = exec.Command("git", "remote", "add", "origin", strings.Replace(remote, "localhost", host, 1))
	cmd.Dir = tmpdir
	e = cmd.Run()
	checkErrorCondition(t, e != nil, "cannot set remote: %s", e)
	cmd = exec.Command("git", "push", "origin", "master")
	cmd.Dir = tmpdir
	cmd.Env = []string{
		fmt.Sprintf("PATH=%s", os.ExpandEnv("PATH")),
		fmt.Sprintf("GIT_SSH=%s", filepath.Join(tmpdir, "git.sh")),
	}
	out, e := cmd.CombinedOutput()
	checkErrorCondition(t, e != nil, "cannot push to remote: %s: %s", e, string(out))
}

func testRepositories(t *testing.T, admingit *gl.Client) {
	t.Log("create a new user for testing repositories")
	u, e := admingit.CreateUser("test@example.com", "username2", "start123", "myname2", nil, nil, nil, nil, nil, nil, nil, nil, true, true)
	checkErrorCondition(t, e != nil, "cannot create user 'username': '%s'", e)
	defer func() {
		t.Log("remove testuser for group repositories")
		admingit.DeleteUser(u.Id)
	}()
	usr, e := admingit.Session("username2", nil, "start123")
	checkErrorCondition(t, e != nil, "cannot open username2 session")
	git := admingit.Child()
	git.Token(usr.PrivateToken)
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

	uploadGitRepo(t, git, u.Id, pr.SshRepoUrl, git.Host())

	_, e = git.CreateFile(pr.Sid(), "mytestfile.txt", "master", "first commit", "Hello world", "text")
	checkErrorCondition(t, e != nil, "cannot create new file in project: %s", e)
}
