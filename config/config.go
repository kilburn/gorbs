package config

import "github.com/kilburn/gorbs/log"

type Interval struct {
	Name string
	Keep int
}

type Backup struct {
	Src     string
	Command string
	Dest    string
}

type Configuration struct {
	Snapshot_root    string
	No_create_root   bool
	Cmd_cp           string
	Cmd_rm           string
	Cmd_rsync        string
	Cmd_ssh          string
	Intervals        []Interval
	Verbose          log.Level
	Test             bool
	Loglevel         int
	Logfile          string
	Lockfile         string
	Rsync_short_args string
	Rsync_long_args  string
	Ssh_args         string
	Exclude_file     string
	Backups          []Backup
}

func New() *Configuration {
	return &Configuration{
		No_create_root:   false,
		Cmd_cp:           "/bin/cp",
		Cmd_rm:           "/bin/rm",
		Cmd_rsync:        "/usr/bin/rsync",
		Cmd_ssh:          "/usr/bin/ssh",
		Intervals:        make([]Interval, 0),
		Verbose:          2,
		Test:             false,
		Loglevel:         3,
		Lockfile:         "/var/run/rsnapshot.pid",
		Rsync_short_args: "-a",
		Rsync_long_args:  "--delete --numeric-ids --relative --delete-excluded",
		Backups:          make([]Backup, 0),
	}
}
