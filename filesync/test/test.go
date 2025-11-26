package test

func DoTest() {
	// var testFile watcher.FileInfo

	// handler.ForEachFileInDirRecursive(constants.BitburnerRelativePath, func(file os.FileInfo, dir string) {
	// 	// if file.Name() == "fileWatch.ts" {
	// 	// 	testFile = watcher.FileInfo{Path: dir + "/" + file.Name(), Info: file}
	// 	// }

	// 	h, m, s := file.ModTime().Clock()
	// 	fmt.Printf("File %s/%s: %d:%d:%d\n", dir, file.Name(), h, m, s)
	// })

	// saveModTime := testFile.Info.ModTime()

	// for {
	// 	time.Sleep(time.Second)

	// 	curFile, _ := os.Stat(testFile.Path)
	// 	fmt.Printf("Old: %s\nNew: %s\nCur: %s\n\n", saveModTime, testFile.Info.ModTime(), curFile.ModTime())
	// }

	// println("")

	// watcher.Initialize()
	// println(watcher.FileStateMap.String())

	// watcher.FileScanner()

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
