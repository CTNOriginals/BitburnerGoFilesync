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
	"filesync/arguments"
	"fmt"
	"os"
	"time"
)

// var (
// 	debugMode = false
// 	testMode  = false
// 	noServer  = false
// 	noWatcher = false
// 	keepAlive = false
// )

func main() {
	startTime := time.Now()
	fmt.Printf("\n\n---- FileSync START %s ----\n", startTime.Format(time.TimeOnly))
	defer fmt.Printf("---- FileSync END %s ----\n", startTime.Format(time.TimeOnly))

	var args = os.Args
	arguments.ParseArgs(args)

	//! Comment this out if you dont have a test/test.go file with a DoTest() func in it
	// if testMode {
	// 	test.DoTest()
	// }

	// if debugMode {
	// 	go debug.DebugCommandListener()
	// }

	// if !noWatcher {
	// 	watcher.Initialize()
	// 	go watcher.FileScanner()
	// }

	// if !noServer {
	// 	communication.StartServer(constants.Port)
	// } else if keepAlive {
	// 	for {
	// 		time.Sleep(time.Millisecond)
	// 	}
	// }
}

// func parseArgs(args []string) {
// 	for _, arg := range args {
// 		switch arg {
// 		case "--test":
// 			testMode = true
// 		case "--debug":
// 			debugMode = true
// 		case "--no-watcher":
// 			noWatcher = true
// 		case "--no-server":
// 			noServer = true
// 		case "--keep-alive":
// 			keepAlive = true
// 		}
// 	}
// }
