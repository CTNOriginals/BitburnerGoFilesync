/*
	BSD 3-Clause License

	Copyright (c) 2025, CTNOriginals

	Redistribution and use in source and binary forms, with or without
	modification, are permitted provided that the following conditions are met:

	1. Redistributions of source code must retain the above copyright notice, this
	list of conditions and the following disclaimer.

	2. Redistributions in binary form must reproduce the above copyright notice,
	this list of conditions and the following disclaimer in the documentation
	and/or other materials provided with the distribution.

	3. Neither the name of the copyright holder nor the names of its
	contributors may be used to endorse or promote products derived from
	this software without specific prior written permission.

	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
	AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
	IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
	DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
	FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
	DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
	SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
	CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
	OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	"log"
	"os"
	"slices"
	"time"

	"github.com/CTNOriginals/BitburnerGoFilesync/arguments"
	"github.com/CTNOriginals/BitburnerGoFilesync/cli"
	"github.com/CTNOriginals/BitburnerGoFilesync/clogger"
	"github.com/CTNOriginals/BitburnerGoFilesync/config"
	"github.com/CTNOriginals/BitburnerGoFilesync/constants"
	"github.com/CTNOriginals/BitburnerGoFilesync/watcher"
	"github.com/CTNOriginals/BitburnerGoFilesync/websocket"
)

var clog = clogger.Default.Clone(clogger.SClog{
	Name: "main",
})

func main() {
	log.SetFlags(0)

	var startTime = time.Now()
	clog.Messagef("\n\n---- FileSync START %s ----\n", startTime.Format(time.TimeOnly))
	defer func() {
		var now = time.Now()
		clog.Messagef(
			"---- FileSync END %s (%s) ----\n",
			now.Format(time.TimeOnly),
			time.Since(startTime).String(),
		)

		// Make sure that all gorotines also terminate when runtime.Goexit is called
		os.Exit(0)
	}()

	var args = os.Args

	constants.Debug = slices.Contains(args, "--test") || slices.Contains(args, "--debug")

	// Make sure that the config file path is the correct one before initializing the config
	arguments.ParseSpecificArgs(args, true, "--config")
	config.Initialize()

	clog.Debugf("Watching Directory: %s", config.Values.Directory)

	arguments.ParseSpecificArgs(args, false, "--config")

	if !constants.NoCli {
		go cli.CommandWatcher()
	}
	if !constants.NoWatcher {
		watcher.Initialize()
		go watcher.FileScanner()
	}

	if !constants.NoServer {
		go websocket.Client.Start(config.Values.Port)
		defer websocket.Client.Close()
	}

	for {
		time.Sleep(time.Second)
	}
}
