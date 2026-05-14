package test

import "github.com/CTNOriginals/BitburnerGoFilesync/cli"

func DoTest() {
	// var arr = make([]int, 50)
	//
	// for i := range arr {
	// 	arr[i] = i
	// }
	//
	// for i := 0; i < len(arr); i++ {
	// 	// for i, num := range arr {
	// 	var num = arr[i]
	// 	fmt.Printf("%d: %d", i, num)
	//
	// 	if num%5 == 0 {
	// 		fmt.Printf("--")
	// 		arr = append(arr[0:i], arr[i+1:]...)
	// 	}
	//
	// 	fmt.Print("\n")
	// }
	//
	// for i, num := range arr {
	// 	fmt.Printf("%d: %d\n", i, num)
	// }

	cli.TestCli()
}
