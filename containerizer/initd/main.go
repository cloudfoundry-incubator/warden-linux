package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/cloudfoundry-incubator/cf-lager"
	"github.com/cloudfoundry-incubator/garden-linux/container_daemon"
	"github.com/cloudfoundry-incubator/garden-linux/container_daemon/unix_socket"
	"github.com/cloudfoundry-incubator/garden-linux/containerizer"
	"github.com/cloudfoundry-incubator/garden-linux/containerizer/system"
	"github.com/cloudfoundry-incubator/garden-linux/network"
	"github.com/cloudfoundry-incubator/garden-linux/process"
	"github.com/cloudfoundry/gunk/command_runner/linux_command_runner"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "initd: panicked: %s\n", r)
			os.Exit(4)
		}
	}()

	rootFsPath := flag.String("root", "", "Path for the root file system directory")
	configFilePath := flag.String("config", "./etc/config", "Path for the configuration file")
	cf_lager.AddFlags(flag.CommandLine)
	flag.Parse()

	if *rootFsPath == "" {
		missing("--root")
	}

	syncReader := os.NewFile(uintptr(3), "/dev/a")
	syncWriter := os.NewFile(uintptr(4), "/dev/d")
	sync := &containerizer.PipeSynchronizer{
		Reader: syncReader,
		Writer: syncWriter,
	}

	if err := sync.Wait(time.Second * 3); err != nil {
		fail(fmt.Sprintf("initd: wait for host: %s", err), 8)
	}
	syncReader.Close()

	env, err := process.EnvFromFile(*configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "initd: failed to get env from config file: %s\n", err)
		os.Exit(3)
	}

	initializer := &system.ContainerInitializer{
		Steps: []system.Initializer{
			&step{system.Mount{
				Type:  system.Tmpfs,
				Flags: syscall.MS_NODEV,
				Path:  "/dev/shm",
			}.Mount},
			&step{system.Mount{
				Type:  system.Proc,
				Flags: syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_NOEXEC,
				Path:  "/proc",
			}.Mount},
			&step{system.Unmount{
				Dir: "/tmp/garden-host",
			}.Unmount},
			&step{func() error {
				return setupNetwork(env)
			}},
			&containerizer.ShellRunnerStep{
				Runner: linux_command_runner.New(),
				Path:   "/etc/seed",
			},
		},
	}

	containerizer := containerizer.Containerizer{
		RootfsPath:  *rootFsPath,
		Initializer: initializer,
		Waiter:      sync,
		Signaller:   sync,
	}

	if err := containerizer.Init(); err != nil {
		fail(fmt.Sprintf("failed to init containerizer: %s", err), 2)
	}

	// EXEC INITD =======

	logger, _ := cf_lager.New("init")

	reaper := system.StartReaper(logger)
	defer reaper.Stop()

	daemon := &container_daemon.ContainerDaemon{
		CmdPreparer: &container_daemon.ProcessSpecPreparer{
			Users:           container_daemon.LibContainerUser{},
			Rlimits:         &container_daemon.RlimitsManager{},
			ProcStarterPath: "/sbin/proc_starter",
		},
		Spawner: &container_daemon.Spawn{
			Runner: reaper,
			PTY:    system.KrPty,
		},
	}

	socketFile := os.NewFile(uintptr(5), "/dev/host.sock")
	defer socketFile.Close()

	listener, err := unix_socket.NewListenerFromFile(socketFile)
	if err != nil {
		fail(fmt.Sprintf("initd: failed to create listener: %s\n", err), 5)
	}

	if err := sync.SignalSuccess(); err != nil {
		fail(fmt.Sprintf("signal host: %s", err), 6)
	}
	syncWriter.Close()

	if err := daemon.Run(listener); err != nil {
		fail(fmt.Sprintf("run daemon: %s", err), 7)
	}
}

func fail(err string, code int) {
	fmt.Fprintf(os.Stderr, "initd: %s\n", err)
	os.Exit(code)
}

func missing(flagName string) {
	fmt.Fprintf(os.Stderr, "initd: %s is required\n", flagName)
	flag.Usage()
	os.Exit(1)
}

func setupNetwork(env process.Env) error {
	_, ipNet, err := net.ParseCIDR(env["network_cidr"])
	if err != nil {
		return fmt.Errorf("initd: failed to parse network CIDR: %s", err)
	}

	mtu, err := strconv.ParseInt(env["container_iface_mtu"], 0, 64)
	if err != nil {
		return fmt.Errorf("initd: failed to parse container interface MTU: %s", err)
	}

	logger, _ := cf_lager.New("hook")
	configurer := network.NewConfigurer(logger.Session("initd: hook.CHILD_AFTER_PIVOT"))
	err = configurer.ConfigureContainer(&network.ContainerConfig{
		Hostname:      env["id"],
		ContainerIntf: env["network_container_iface"],
		ContainerIP:   net.ParseIP(env["network_container_ip"]),
		GatewayIP:     net.ParseIP(env["network_host_ip"]),
		Subnet:        ipNet,
		Mtu:           int(mtu),
	})
	if err != nil {
		return fmt.Errorf("initd: failed to configure container network: %s", err)
	}

	return nil
}

type step struct {
	fn func() error
}

func (s *step) Init() error {
	return s.fn()
}