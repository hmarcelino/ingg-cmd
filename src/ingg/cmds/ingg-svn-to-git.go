package cmds

import (
    "github.com/codegangsta/cli"
    "os/exec"
    "fmt"
    "strings"
    "os"
    "sync"
    Print "ingg/utils"
)

var (
    maxDepth int = 5
    destinationFolder string = "."
    totalWorkers = 5
)

var SvnToGit = cli.Command{
    Name: "svn-to-git",
    Category: "Vcs Operations",
    ArgsUsage:"--http <http_base_path> --destination <destination_folder>",
    Description: "Migrate a SVN repository to Git",
    Flags: [] cli.Flag{
        cli.StringFlag{
            Name: "http, p",
            Usage: "The SVN http path",
        },
        cli.IntFlag{
            Name: "max-depth",
            Value: 5,
            Usage: "The max depth allowed while traversing the tree. Default is 5",
            Destination: &maxDepth,
        },
        cli.StringFlag{
            Name: "destination, d",
            Value: ".",
            Usage: "The destination folder",
            Destination: &destinationFolder,
        },
    },
    Action: func(c *cli.Context) {
        Print.PrintMsg("Starting process of migrating SVN to Git")
        Print.PrintMsg("Please wait !!! This may take a while... ")
        Print.PrintMsg("")

        var binary string
        var lookErr error
        var isValid bool = true

        binary, lookErr = exec.LookPath("svn")
        if lookErr != nil {
            Print.Danger("* svn: executable file not found in PATH")
            isValid = false
        } else {
            Print.PrintInfo(fmt.Sprintf("* Found SVN at %s", binary))
        }

        binary, lookErr = exec.LookPath("git")
        if lookErr != nil {
            Print.Danger("* git: executable file not found in PATH")

        } else {
            Print.PrintInfo(fmt.Sprintf("* Found GIT at %s", binary))
        }

        baseHttp := c.String("http")
        if baseHttp == "" {
            Print.Danger("* No svn http server defined")
            isValid = false
        }

        Print.PrintMsg("")

        if !isValid {
            Print.Danger("Requirements not met! Please fix the identified errors")
            os.Exit(1)
        }

        if destinationFolder != "." {
            _, err := os.Stat(destinationFolder)
            if err != nil {
                os.MkdirAll(destinationFolder, 0755)
            }
        }

        Print.PrintMsg("Svn repositories found:")
        if !strings.HasSuffix(baseHttp, "/") {
            baseHttp = baseHttp + "/"
        }

        svnRepos := findRepositories(baseHttp, 0)

        if len(svnRepos) == 0 {
            Print.PrintWarning("No repositories found")
            os.Exit(0)
        }

        Print.PrintMsg("")
        migrateToGit(svnRepos)
    },
}

// search through the repository to find all folders
// that are svn repositories and return the full http
// paths that will be used in the git svn clone.

func findRepositories(httpBasePath string, depth int) []string {

    var svnRepos []string = make([]string, 0)

    if Verbose {
        Print.PrintMsg(fmt.Sprintf("Getting subfolders of svn http path %s", httpBasePath));
    }

    var subFolders []string = getSubFolders(httpBasePath);

    // A repository is identified if there is
    // a subfolder with name "trunk"
    isRepo := func(subFolders[]string) bool {
        var joinedFolders = strings.Join(subFolders, ";");

        if (strings.Contains(joinedFolders, "trunk")) {
            return true

        } else {
            return false
        }
    }(subFolders)

    if !isRepo {
        depth = depth + 1

        if (depth > maxDepth) {
            if Verbose {
                Print.PrintWarning("Reached max depth. Stopping here.");
            }

            return []string{};
        }

        if Verbose {
            Print.PrintMsg(fmt.Sprintf("svn http path is not a repository: %s ", httpBasePath));
        }

        for _, path := range subFolders {
            reposFound := findRepositories(httpBasePath + path, depth)

            if len(reposFound) > 0 {
                svnRepos = append(svnRepos, reposFound...);
            }
        }
    } else {
        Print.PrintMsg("* " + httpBasePath);
        svnRepos = append(svnRepos, httpBasePath);
    }

    return svnRepos;
}

func getSubFolders(httpPath string) ([]string) {
    cmd := exec.Command("svn", "list", httpPath)
    output, err := cmd.CombinedOutput()
    if err != nil {
        Print.PrintDanger(fmt.Sprintf("Error getting subfolders of svn http path %s", httpPath))
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

func migrateToGit(svnRepos []string) {
    Print.PrintMsg("Starting process of migrating SVN to Git")

    os.Chdir(destinationFolder);

    migrateRepoTasks := make(chan string)

    var wg sync.WaitGroup
    wg.Add(totalWorkers);

    workerIdx := 0;
    for workerIdx < totalWorkers {
        go cloneRepoWorker(workerIdx, migrateRepoTasks, &wg);
        workerIdx++;
    }

    go func() {
        for _, svnRepoPath := range svnRepos {
            migrateRepoTasks <- svnRepoPath
        }

        close(migrateRepoTasks)

        if Verbose {
            Print.PrintMsg("migrateRepoTasks channel is closed");
        }
    }()

    wg.Wait()
}

func cloneRepoWorker(workerIdx int, migrateRepoTasks <- chan string, wg *sync.WaitGroup) {

    if Verbose {
        Print.PrintMsg(fmt.Sprintf("[worker: %d] Starting worker", workerIdx));
    }

    for {

        svnRepoPath, more := <-migrateRepoTasks

        if more {
            if Verbose {
                Print.PrintMsg(fmt.Sprintf("[worker: %d] Cloning repository %s", workerIdx, svnRepoPath));
            }

            cmd := exec.Command("git", "svn", "clone", "-s", svnRepoPath)
            _, err := cmd.CombinedOutput()

            if err == nil {
                Print.PrintSuccess("* Svn repository migrated to GIT: " + svnRepoPath)

            } else {
                Print.PrintDanger(fmt.Sprintf("Error migrating to GIT SVN repo %s", svnRepoPath))
                Print.PrintError(err)
            }

        } else {
            if (Verbose) {
                Print.PrintMsg(fmt.Sprintf("[worker: %d] Worker terminated", workerIdx));
            }

            wg.Done()
            return
        }
    }

}
