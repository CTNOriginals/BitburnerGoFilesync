package debug

import (
	"bufio"
	"filesync/communication"
	"filesync/communication/definitions"
	"filesync/utils"
	"fmt"
	"os"
	"strings"
)

var DefaultParameters = map[definitions.Method][]string{
	definitions.GetFileNames:    {"home"},
	definitions.GetFile:         {"proto.ts", "home"},
	definitions.GetFileMetadata: {"proto.ts", "home"},
	definitions.PushFile:        {"proto.ts", string(utils.SanitizeFileContent(utils.GetFileContentByPath("proto.ts"))), "home"},
}

func DebugCommandListener() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := strings.TrimSpace(scanner.Text())

		if cmd == "" {
			continue
		}

		switch cmd {
		case "--help", "-h", "?":
			var msg strings.Builder
			for _, method := range definitions.MethodsAsArray() {
				if _, exists := DefaultParameters[method]; !exists {
					continue
				}

				msg.WriteString(string(method) + " ")
			}
			println(msg.String())
			continue
		case "--test":
			for method, params := range DefaultParameters {
				var msg = communication.SendRequest(method, params...)
				println(msg.String())
			}
			continue
		}

		if communication.ActiveConnection == nil {
			fmt.Println("No active connection")
			continue
		}

		def, exists := definitions.RPCDefinitions[definitions.Method(cmd)]
		defaultParameters, defaultExists := DefaultParameters[def.Method]

		if !exists {
			fmt.Printf("Invalid method name: %s\nSending raw input instead.\n", cmd)
			continue
		}

		if !defaultExists {
			fmt.Printf("This method does not have default parameters defined")
			continue
		}

		var msg = communication.SendRequest(def.Method, defaultParameters...)
		println(msg.String())
	}
}
