package main

import (
	"flag"
	"fmt"
	"github.com/ulrichSchreiner/gl"
)

var token = flag.String("token", "", "the private token")
var host = flag.String("host", "", "the gitlab host")

func main() {
	flag.Parse()
	gitlab, err := gl.OpenV3(*host)
	if err != nil {
		panic(err)
	}
	git := gitlab.Child()
	git.Token(*token)

	git.CreateProject("test", nil, nil, nil, false, false, false, false, false, nil, nil)
	prjs, err := git.AllVisibleProjects()
	if err != nil {
		panic(err)
	}
	for _, p := range prjs {
		fmt.Printf("%#v\n", p)
	}
	/*
			fmt.Println("-----------------")
			prj, err := git.Project(20)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%#v\n", prj)
			fmt.Println("-----------------")
		evts, err := git.AllIssues(20)
		if err != nil {
			panic(err)
		}
		for _, e := range evts {
			fmt.Printf("%#v\n", e)
		}
	*/
}
