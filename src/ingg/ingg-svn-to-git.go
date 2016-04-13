package main

import (
    "github.com/codegangsta/cli"
    "os/exec"
    "fmt"
    "strings"
    "os"
)

var verbose bool = false;

var maxDepth int = 5;

var SvnToGit = cli.Command{
    Name: "svn-to-git",
    Category: "Vcs Operations",
    ArgsUsage:"[-p http_base_path] [-e exclude_repos]",
    Description: "Migrate a SVN repository to Git",
    Flags: [] cli.Flag{
        cli.StringFlag{
            Name: "http, p",
            Usage: "The SVN http path",
        },
        cli.IntFlag{
            Name: "max-depth, d",
            Value: 5,
            Usage: "The max depth allowed while traversing the tree. Default is 5",
            Destination: &maxDepth,
        },
        cli.StringFlag{
            Name: "exclude, e",
            Usage: "The sub folders to not migrate",
        },
        cli.BoolFlag{
            Name: "verbose, v",
            Usage: "Print more information of the process",
            Destination: &verbose,
        },
    },
    Action: func(c *cli.Context) {
        PrintMsg("Starting process of migrating SVN ) to Git")
        PrintMsg("Please wait !!! This may take a while... ")
        PrintMsg("")

        var binary string;
        var lookErr error;
        var isValid bool = true;

        binary, lookErr = exec.LookPath("svn")
        if lookErr != nil {
            danger("* svn: executable file not found in PATH")
            isValid = false
        }else {
            PrintInfo(fmt.Sprintf("* Found SVN at %s", binary))
        }

        binary, lookErr = exec.LookPath("git")
        if lookErr != nil {
            danger("* git: executable file not found in PATH")
            isValid = false
        }else {
            PrintInfo(fmt.Sprintf("* Found GIT at %s", binary))
        }

        PrintMsg("")

        if !isValid {
            danger("Requirements not met! Please make available SVN and GIT")
            os.Exit(1)
        }

        PrintMsg("Svn repositories found:")
        svnRepos := _findRepositories(c.String("http"), 0)

        if len(svnRepos) == 0 {
            PrintWarning("No repositories found")
            os.Exit(0)
        }
    },
}

// search through the repository to find all folders
// that are svn repositories and return the full http
// paths that will be used in the git svn clone.
func _findRepositories(httpBasePath string, depth int) []string {

    var svnRepos []string = make([]string, 0)

    if verbose {
        PrintMsg(fmt.Sprintf("Getting subfolders of svn http path %s", httpBasePath));
    }

    var subFolders []string = _getSubFolders(httpBasePath);

    // A repository is identified if there is
    // a subfolder with name "trunk"
    isRepo := func(subFolders[]string) bool {
        var joinedFolders = strings.Join(subFolders, ";");

        if (strings.Contains(joinedFolders, "trunk")) {
            return true

        }else {
            return false
        }
    }(subFolders)

    if !isRepo {
        depth = depth + 1

        if (depth > maxDepth) {
            if verbose {
                PrintWarning("Reached max depth. Stopping here.");
            }

            return []string{};
        }

        if verbose {
            PrintMsg(fmt.Sprintf("svn http path is not a repository: %s ", httpBasePath));
        }

        for _, path := range subFolders {
            reposFound := _findRepositories(httpBasePath + path, depth)

            if len(reposFound) > 0 {
                svnRepos = append(svnRepos, reposFound...);
            }
        }
    }else {
        PrintMsg(httpBasePath);
        svnRepos = append(svnRepos, httpBasePath);
    }

    return svnRepos;
}

func _getSubFolders(httpPath string) (subFolders []string) {
    cmd := exec.Command("svn", "list", httpPath)
    output, err := cmd.CombinedOutput()
    if err != nil {
        PrintDanger(fmt.Sprintf("Error getting subfolders of svn http path %s", httpPath))
        panic(err)
        os.Exit(2)
    }

    folders := make([]string, 0);
    for _, path := range strings.Split(string(output), "\n") {
        p := strings.Trim(path, "\r")

        if strings.HasSuffix(p, "/") {
            folders = append(folders, p)
        }
    }

    return folders
}