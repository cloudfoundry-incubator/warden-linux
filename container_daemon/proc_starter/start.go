package proc_starter

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"syscall"

	"flag"

	"code.cloudfoundry.org/garden-linux/container_daemon"
	"code.cloudfoundry.org/garden-linux/system"
	"github.com/docker/docker/pkg/reexec"
)

func init() {
	runtime.LockOSThread()
	reexec.Register("proc_starter", start)
}

// proc_starter starts a user process with the correct rlimits and after
// closing any open FDs.
func start() {
	rlimits := flag.String("rlimits", "", "encoded rlimits")
	dropCapabilities := flag.Bool("dropCapabilities", true, "drop capabilities before starting process")
	uid := flag.Int("uid", -1, "user id to run the process as")
	gid := flag.Int("gid", -1, "group id to run the process as")
	workDir := flag.String("workDir", "", "working dir for the process")
	extendedWhitelist := flag.Bool("extendedWhitelist", false, "whitelist CAP_SYS_ADMIN in addition to the default set. Use only with -dropCapabilities=true")
	flag.Parse()

	closeFds()

	mgr := &container_daemon.RlimitsManager{}
	must(mgr.Apply(mgr.DecodeLimits(*rlimits)))

	args := flag.Args()

	if *dropCapabilities {
		caps := &system.ProcessCapabilities{Pid: os.Getpid()}
		must(caps.Limit(*extendedWhitelist))
	}

	execer := system.UserExecer{}
	if err := execer.ExecAsUser(*uid, *gid, *workDir, args[0], args[1:]...); err != nil {
		fmt.Fprintf(os.Stderr, "proc_starter: ExecAsUser: %s\n", err)
		os.Exit(255)
	}
}

func closeFds() {
	fds, err := ioutil.ReadDir("/proc/self/fd")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: read /proc/self/fd: %s", err)
		os.Exit(255)
	}

	for _, fd := range fds {
		if fd.IsDir() {
			continue
		}

		fdI, err := strconv.Atoi(fd.Name())
		if err != nil {
			panic(err) // cant happen
		}

		if fdI <= 2 {
			continue
		}

		syscall.CloseOnExec(fdI)
	}
}

var must = mustNot

func mustNot(err error) {
	if err != nil {
		panic(err)
	}
}
