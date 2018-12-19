package main

import (
	"errors"
	"fmt"
	"math"
)

//////https://www.cnblogs.com/abc-begin/p/7953894.html

var (
	alli int
	allf float32
)

const (
	E1 = 0
	E2 = 1
)

func sumInt(a int, b int) int {
	return a + b
}

func swap(aa, bb string) (string, string) {
	return bb, aa
}

func change(aaa *int) {
	* aaa += 100
}


func main222() {

	//ClientConnection()

	const (
		a = iota
		b
		c
		d
		e = "ee"
		f
		g = 100
		h
		ii = iota
		j
	)

	fmt.Print(a, b, c, d, e, f, g, h, ii, j)

	var i = 123

	fmt.Print(i)

	fmt.Print("hello world")
	fmt.Print(math.Pi)

	if t := a * b; t < 10 {
		fmt.Print("t", t)
	}

	var grade string = "B"
	var marks int = 90
	switch marks {
	case 90:
		grade = "A"
	case 80:
		grade = "B"
	case 50, 60, 70:
		grade = "C"
	default:
		grade = "D"
	}
	switch {
	case grade == "A":
		fmt.Printf("优秀!\n")
	case grade == "B", grade == "C":
		fmt.Printf("良好\n")
	case grade == "D":
		fmt.Printf("及格\n")
	case grade == "F":
		fmt.Printf("不及格\n")
	default:
		fmt.Printf("差\n");
	}
	fmt.Printf("你的等级是 %s\n", grade)

	/////////////
	var x interface{} = "go"

	switch i := x.(type) {
	case nil:
		fmt.Printf(" x 的类型 :%T", i)
	case int:
		fmt.Printf("x 是 int 型")
	case float64:
		fmt.Printf("x 是 float64 型")
	case func(int) float64:
		fmt.Printf("x 是 func(int) 型")
	//case string:
	//	fmt.Print(" x is string")
	case bool, string:
		fmt.Printf("x 是 bool 或 string 型")
	default:
		fmt.Printf("未知型")
	}

	///////////

	fmt.Println("counting")

	for i := 0; i < 3; i++ {
		defer fmt.Println("===== ", i)
	}

	fmt.Println("done")

	//////////

	var tb int = 3
	var ta int

	/* for 循环 */
	for ta := 0; ta < 2; ta++ {
		fmt.Printf("111 ta 的值为: %d\n", ta)
	}

	fmt.Printf("333 ta 的值为: %d\n", ta)

	for ta < tb {
		ta++
		fmt.Printf("222 ta 的值为: %d\n", ta)
	}

	numbers := [6]int{1, 2, 3, 5}
	for i, x := range numbers {
		fmt.Printf("第 %d 位 x 的值 = %d\n", i, x)
	}
	//////
	sum := 1
	for sum < 1000 {
		sum += sum
		fmt.Printf(" sum cur %d \n", sum)
	}
	fmt.Println(sum)
	//////////////

	fmt.Println(sumInt(1, 2))

	var aa string = "hello"
	var bb string = "hi"
	aa, bb = swap(aa, bb)
	fmt.Println(swap("wang", "bo"))

	var aaa = 200
	change(&aaa)
	fmt.Println(aaa)

	/////////
	fmt.Println(split(100))

	//////////

	hypot := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}
	fmt.Println(hypot(5, 12))

	fmt.Println(compute(hypot))
	fmt.Println(compute(math.Pow))

	//////////

	/* nextNumber 为一个函数，函数 i 为 0 */
	nextNumber := getSequence()

	/* 调用 nextNumber 函数，i 变量自增 1 并返回 */
	fmt.Println(nextNumber())
	fmt.Println(nextNumber())
	fmt.Println(nextNumber())

	/* 创建新的函数 nextNumber1，并查看结果 */
	nextNumber1 := getSequence()
	fmt.Println(nextNumber1())
	fmt.Println(nextNumber1())

	///////////
	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
		fmt.Println(
			pos(i),
			neg(-2*i),
		)
	}

	//////////
	var c1 Circle
	c1.radius = 10.00
	fmt.Println("Area of Circle(c1) = ", c1.getArea())

	c2 := Circle{10.00}
	fmt.Println("c2  ", c2.getArea())
	////////
	v := &Vertex{3, 4}
	fmt.Println(v.Abs())

	////////

	f2 := MyFloat(-math.Sqrt2)
	fmt.Println(f2.Abs())

	///////
	v2 := &Vertex{3, 4}
	fmt.Printf("Before scaling: %+v, Abs: %v\n", v2, v2.Abs())
	v2.Scale(5)
	fmt.Printf("After scaling: %+v, Abs: %v\n", v2, v2.Abs())

	///////////

	var array = [3]int{11, 22, 33}

	forArray(array)

	var arr2 = []int{1, 3, 4, 5, 6}
	forArray2(arr2)

	//////////////
	Book1 := Books{"book1", "a", "aa", 111}

	var Book2 Books
	Book2.title = "book2"
	Book2.author = "b"
	Book2.subject = "bb"
	Book2.book_id = 222

	printBook(&Book1)
	printBook(&Book2)

	/////////////
	fmt.Println(v1, p, vv2, v3)

	/////////////
	printSpace()

	{
		var numbers []int
		printSlice(numbers)

		numbers = append(numbers, 1)
		printSlice(numbers)

		numbers = append(numbers, 2, 3, 4, 5, 6, 7)
		printSlice(numbers)

		/* 打印子切片从索引1(包含) 到索引4(不包含)*/
		fmt.Println("numbers[1:4] ==", numbers[1:4])

		/* 默认下限为 0*/
		fmt.Println("numbers[:3] ==", numbers[:3])

		/* 默认上限为 len(s)*/
		fmt.Println("numbers[4:] ==", numbers[4:])

		number1 := make([]int, 0, 5)
		printSlice(number1)

		/* 打印子切片从索引  0(包含) 到索引 2(不包含) */
		number2 := numbers[:2]
		printSlice(number2)

		var number3 []int
		copy(number3, numbers[1:4])
		printSlice(number3)

		number4 := make([]int, len(numbers[1:4]), len(numbers[1:4]))
		copy(number4, numbers[1:4])
		printSlice(number4)

		number4 = append(number4, 111, 222)
		printSlice(number4) ///// ????
		for eee, fff := range number4 {
			fmt.Println(eee, fff)
		}

		////////////////////////
		printSpace()

		countryCapitalMap := map[string]string{"china": "Bei Jing"}
		for country := range countryCapitalMap {
			fmt.Println("Capital of", country, "is", countryCapitalMap[country])
		}

		countryCapitalMap = make(map[string]string)

		countryCapitalMap["France"] = "Paris"
		countryCapitalMap["Italy"] = "Rome"
		countryCapitalMap["Japan"] = "Tokyo"
		countryCapitalMap["India"] = "New Delhi"

		for country := range countryCapitalMap {
			fmt.Println("Capital of", country, "is", countryCapitalMap[country])
		}

		for k, v := range countryCapitalMap {
			fmt.Printf(" country %s  capital %s \n", k, v)
		}

		captial1, ok1 := countryCapitalMap["France"]
		if (ok1) {
			fmt.Printf("Capital of %s is %\n", captial1, ok1)
		} else {
			fmt.Println("Capital of France is not present")
		}

		delete(countryCapitalMap, "France");
		captial2, ok2 := countryCapitalMap["France"]
		if (ok2) {
			fmt.Println("Capital of France is", captial2)
		} else {
			fmt.Println("Capital of France is not present")
		}

		//////////////
		printSpace()
		{
			//这是我们使用range去求一个slice的和。使用数组跟这个很类似
			nums := []int{2, 3, 4}
			sum := 0
			for _, num := range nums {
				sum += num
			}
			fmt.Println("sum:", sum)
			//在数组上使用range将传入index和值两个变量。上面那个例子我们不需要使用该元素的序号，所以我们使用空白符"_"省略了。有时侯我们确实需要知道它的索引。
			for i, num := range nums {
				if num == 3 {
					fmt.Println("index:", i)
				}
			}
			//range也可以用在map的键值对上。
			kvs := map[string]string{"a": "apple", "b": "banana"}
			for k, v := range kvs {
				fmt.Printf("%s -> %s\n", k, v)
			}
			//range也可以用来枚举Unicode字符串。第一个参数是字符的索引，第二个是字符（Unicode的值）本身。
			for i, c := range "ab" {
				fmt.Println(i, c)
			}
		}
	}

	//////////
	printSpace()

	var phone Phone

	phone = new(NokiaPhone)
	phone.call()

	phone = new(IPhone)
	phone.call()

	//////////

	result, err := Sqrt(-1)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

	///////////
	printSpace()
	{
		// 正常情况
		if result, errorMsg := Divide(100, 10); errorMsg == "" {
			fmt.Println("100/10 = ", result)
		}
		// 当被除数为零的时候会返回错误信息
		if _, errorMsg := Divide(100, 0); errorMsg != "" {
			fmt.Println("errorMsg is: ", errorMsg)
		}
	}
}

// 定义一个 DivideError 结构
type DivideError struct {
	dividee int
	divider int
}

// 实现     `error` 接口
func (de *DivideError) Error() string {
	strFormat := `
    Cannot proceed, the divider is zero.
    dividee: %d
    divider: 0
`
	return fmt.Sprintf(strFormat, de.dividee)
}

// 定义 `int` 类型除法运算的函数
func Divide(varDividee int, varDivider int) (result int, errorMsg string) {
	if varDivider == 0 {
		dData := DivideError{
			dividee: varDividee,
			divider: varDivider,
		}
		errorMsg = dData.Error()
		return
	} else {
		return varDividee / varDivider, ""
	}

}

func Sqrt(f float64) (float64, error) {
	if f < 0 {
		return 0, errors.New("math: square root of negative number")
	} else {
		return math.Sqrt(f), nil
	}
}

type Phone interface {
	call()
}

type NokiaPhone struct {
}

func (nokiaPhone NokiaPhone) call() {
	fmt.Println("I am Nokia, I can call you!")
}

type IPhone struct {
}

func (iPhone IPhone) call() {
	fmt.Println("I am iPhone, I can call you!")
}

func printSpace() {
	fmt.Println("================")
	fmt.Println("================")
	fmt.Println("================")

}

func printSlice(x []int) {
	fmt.Printf("len=%d cap=%d slice=%v\n", len(x), cap(x), x)
}

var (
	v1  = Vertex{1, 2}  // 类型为 Vertex
	vv2 = Vertex{X: 1}  // Y:0 被省略
	v3  = Vertex{}      // X:0 和 Y:0
	p   = &Vertex{1, 2} // 类型为 *Vertex
)

func printBook(book *Books) {
	fmt.Printf("Book title : %s\n", book.title);
	fmt.Printf("Book author : %s\n", book.author);
	fmt.Printf("Book subject : %s\n", book.subject);
	fmt.Printf("Book book_id : %d\n", book.book_id);
}

type Books struct {
	title   string
	author  string
	subject string
	book_id int
}

func forArray2(array []int) {

	fmt.Print(" ===== forArray2")
	for y := range array {
		fmt.Printf("第 %d 位的值 = %d\n", y, array[y])
	}

	for j, x := range array {
		fmt.Printf("第 %d 位的值 = %d\n", j, x)
	}

	var index int
	for index = 0; index < 3; index++ {
		fmt.Printf("第 %d 位的值 = %d\n", index, array[index])
	}

}

func forArray(array [3]int) {
	fmt.Print(" ===== forArray")

	for y := range array {
		fmt.Printf("第 %d 位的值 = %d\n", y, array[y])
	}

	for j, x := range array {
		fmt.Printf("第 %d 位的值 = %d\n", j, x)
	}

	var index int
	for index = 0; index < 3; index++ {
		fmt.Printf("第 %d 位的值 = %d\n", index, array[index])
	}

}

type MyFloat float64

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

/* 定义函数 */
type Circle struct {
	radius float64
}

//该 method 属于 Circle 类型对象中的方法
func (c Circle) getArea() float64 {
	//c.radius 即为 Circle 类型对象中的属性
	return 3.14 * c.radius * c.radius
}

type Vertex struct {
	X, Y float64
}

func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func compute(fn func(float64, float64) float64) float64 {
	return fn(3, 4)
}

func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
}

func getSequence() func() int {
	i := 0
	return func() int {
		i += 1
		return i
	}
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}
