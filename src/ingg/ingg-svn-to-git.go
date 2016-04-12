package main

import "github.com/codegangsta/cli"

var SvnToGit = cli.Command{
    Name: "svn-to-git",
    Category: "Vcs Operations",
    ArgsUsage:"[-p http_base_path] [-e exclude_repos]",
    Description: "Migrate a SVN repository to Git",
    Action: func(c *cli.Context) {
        println("Starting process of migrating SVN repositories to Git")
    },
}

func main(){
    SvnToGit
}
