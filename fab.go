package main

func Fab(n int) int {

	if n < 2 {
		return n
	}
	return Fab(n-1) + Fab(n-2)

}

//func main() {
//	for i := 0; i < 10; i++ {
//		fmt.Printf(" %d \t ", Fab(i))
//	}
//	fmt.Println("type ", reflect.TypeOf(0xe0000000))
//}
