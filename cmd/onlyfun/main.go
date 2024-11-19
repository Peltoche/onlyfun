package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/Peltoche/onlyfun/internal/server"
	"github.com/Peltoche/onlyfun/internal/tools/buildinfos"
	"github.com/adrg/xdg"
)

const (
	binaryName        = "onlyfun"
	binaryDescription = `Only fun.

Usage:
  ` + binaryName + ` [flags...]

Flags:
`
)

var configDirs = append(xdg.DataDirs, xdg.DataHome)

type exitCode int

const (
	exitOK        exitCode = 0
	exitError     exitCode = 1
	exitInitError exitCode = 2
)

func main() {
	output := os.Stdout

	code := mainRun(os.Args, output)
	os.Exit(int(code))
}

func mainRun(args []string, output io.Writer) exitCode {
	ctx := context.Background()

	defaultFolder := getDefaultFolder()
	flags, err := parseFlags(args, defaultFolder, output)
	if err != nil {
		return exitInitError
	}

	if flags == nil {
		return exitOK
	}

	if flags.PrintVersion {
		fmt.Fprintf(output, "%s version %s\n", binaryName, buildinfos.Version())
		return exitOK
	}

	cfg, err := NewConfigFromFlags(flags)
	if err != nil {
		io.WriteString(output, err.Error())
		return exitInitError
	}

	_, err = server.Run(ctx, cfg)
	if err != nil {
		return exitError
	}

	return exitOK
}

func getDefaultFolder() string {
	var defaultFolder string

	for _, dir := range configDirs {
		_, err := os.Stat(path.Join(dir, binaryName))
		if err == nil {
			defaultFolder = path.Join(dir, binaryName)
			break
		}
	}

	if defaultFolder == "" {
		defaultFolder = path.Join(xdg.DataHome, binaryName)
	}

	return defaultFolder
}

func parseFlags(args []string, defaultFolder string, output io.Writer) (*flags, error) {
	flags := flags{}

	fs := flag.NewFlagSet("flags", flag.ContinueOnError)
	fs.SetOutput(output)

	if !buildinfos.IsRelease() {
		// Those flags are  only available outside the releases for security reasons.
		fs.BoolVar(&flags.Dev, "dev", false, "Run in dev mode and make json prettier")
		fs.BoolVar(&flags.HotReload, "hot-reload", false, "Enable the asset hot reload")
	}

	fs.BoolVar(&flags.Debug, "debug", false, "Force the debug level")
	fs.StringVar(&flags.LogLevel, "log-level", "info", "Log message verbosity LEVEL (debug, info, warning, error)")

	fs.StringVar(&flags.Folder, "folder", defaultFolder, "Specify you data directory location")
	fs.BoolVar(&flags.MemoryFS, "memory-fs", false, "Replace the OS filesystem by a in-memory stub. *Every data will disapear after each restart*.")

	fs.StringVar(&flags.TLSCert, "tls-cert", "", "Public HTTPS certificate file (.crt)")
	fs.StringVar(&flags.TLSKey, "tls-key", "", "Private HTTPS key file (.key)")
	fs.BoolVar(&flags.SelfSignedCert, "self-signed-cert", false, "Generate and use a self-signed HTTPS/TLS certificate.")

	fs.IntVar(&flags.HTTPPort, "http-port", 5764, "Web server port number.")
	fs.StringVar(&flags.HTTPHost, "http-host", "0.0.0.0", "Web server IP address")

	fs.BoolVar(&flags.PrintVersion, "version", false, "version for onlyfun")
	fs.BoolVar(&flags.PrintHelp, "help", false, "help for onlyfun")

	err := fs.Parse(args[1:])
	if err != nil {
		return nil, err
	}

	if flags.PrintHelp {
		io.WriteString(output, binaryDescription)
		fs.PrintDefaults()
		return nil, nil
	}

	return &flags, nil
}
