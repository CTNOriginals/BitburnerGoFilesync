package test

import (
	"filesync/constants"
	"filesync/handler"
	"fmt"
	"os"
)

func DoTest() {
	handler.ForEachFileInDirRecursive(constants.BitburnerRelativePath, func(file os.FileInfo, dir string) {
		h, m, s := file.ModTime().Clock()
		fmt.Printf("File %s/%s: %d:%d:%d\n", dir, file.Name(), h, m, s)
	})
	// info, _ := os.Stat(constants.BitburnerRelativePath + "fileWatch.ts")
	// println(info.ModTime().Clock())

	// // fileContent := handler.GetFileByPath("fileTest.ts")
	// // sanitized := handler.SanitizeFileContent(fileContent)
	// // println(string(sanitized))

	// def := definitions.RPCDefinitions.GetByMethod(definitions.GetFile)
	// defaultParameters, _ := debug.DefaultParameters[def.Method]

	// jsonRpc := rpc.NewRPC(def.Method, defaultParameters...).String()
	// jsonMap := handler.JSONToMap([]byte(jsonRpc))

	// println(jsonRpc)

	// for key, val := range *jsonMap {
	// 	fmt.Printf("%s: (%T) %v\n", key, val, val)
	// }
}
