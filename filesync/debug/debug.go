package debug

import (
	"bufio"
	"filesync/rpc"
	"filesync/rpc/definitions"
	"filesync/utils"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
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
		}

		if rpc.ActiveConnection == nil {
			fmt.Println("No active connection")
			continue
		}

		def := definitions.RPCDefinitions.GetByMethod(definitions.Method(cmd))
		defaultParameters, exists := DefaultParameters[def.Method]

		if !def.IsError() {
			if !exists {
				fmt.Printf("This method does not have default parameters defined")
				return
			}

			cmd = rpc.NewRPC(def.Method, defaultParameters...).String()
		} else {
			fmt.Printf("Invalid method name: %s\nSending raw input instead.\n", cmd)
		}

		fmt.Printf("✅ Sending: %s\n", cmd)

		err := rpc.ActiveConnection.WriteMessage(websocket.TextMessage, []byte(cmd))
		if err != nil {
			fmt.Println("Error sending:", err)
		}
	}
}
