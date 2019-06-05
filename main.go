// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the URIs of this project regarding your
// rights to use or distribute this software.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/pflag"
	"github.com/sylabs/singularity/internal/pkg/sylog"
	pluginapi "github.com/sylabs/singularity/pkg/plugin"
	singularity "github.com/sylabs/singularity/pkg/runtime/engines/singularity/config"
)

var Plugin = pluginapi.Plugin{
	Manifest: pluginapi.Manifest{
		Name:        "github.com/mem/singularity-fusemnt",
		Author:      "Marcelo E. Magallon",
		Version:     "0.0.1",
		Description: "Singularity plugin allowing the user to mount a FUSE example filesystem",
	},

	Initializer: impl,
}

type pluginImplementation struct {
}

var impl = pluginImplementation{}

type FuseConfig struct {
	DevFuseFd  int
	MountPoint string
	Program    []string
}

type pluginConfig struct {
	Fuse FuseConfig
}

func fusecmdCallback(f *pflag.Flag, cfg *singularity.EngineConfig) {
	mnt := f.Value.String()

	// This will be called even if the flag was not used.
	// Assume that an empty mount point means the user did not pass
	// the flag and return silently.
	if mnt == "" {
		return
	}

	if !strings.HasPrefix(mnt, "/") {
		sylog.Fatalf("Invalid mount point %s.\n", mnt)
	}

	sylog.Verbosef("Mounting FUSE example filesystem at %s.\n", mnt)

	// Since this is a demonstration, it uses the existing
	// fuse-example driver in the host.  Normally the user would
	// create a container with the desired driver in it. For this
	// example, the expectation is that fuse-example works inside
	// the container, which means it's probably necessary to compile
	// it statically.
	//
	// The /bin/fusermount bind-mount is overkill.  fuse-example
	// looks for fusermount in order to do the mounting, but since
	// singularity takes care of obtaining a file descriptor for
	// /dev/fuse, fusermount is not really needed.

	fuseExampleDriver, err := exec.LookPath("fuse-example")
	if err != nil {
		sylog.Fatalf("Cannot find the fuse-example driver in $PATH. Abort.\n")
	}

	bindPaths := cfg.GetBindPath()
	bindPaths = append(bindPaths, "/bin/true:/bin/fusermount", fuseExampleDriver)
	cfg.SetBindPath(bindPaths)

	config := pluginConfig{
		Fuse: FuseConfig{
			MountPoint: mnt,
			Program:    []string{fuseExampleDriver},
		},
	}
	if err := cfg.SetPluginConfig(Plugin.Manifest.Name, config); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot set plugin configuration: %+v\n", err)
		return
	}
}

func (p pluginImplementation) Initialize(r pluginapi.HookRegistration) {
	flag := pluginapi.StringFlagHook{
		Flag: pflag.Flag{
			Name:  "fusemnt",
			Usage: "Mount FUSE example filesystem here",
		},
		Callback: fusecmdCallback,
	}

	r.RegisterStringFlag(flag)
}
