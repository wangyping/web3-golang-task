package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// 题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，
//然后在主函数中调用该函数并输出修改后的值。
// 考察点 ：指针的使用、值传递与引用传递的区别。
// 题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
// 考察点 ：指针运算、切片操作。

func goFunc(a *int) {
	*a += 10
}

func doubleSlice(a []int) {
	for i := range a {
		a[i] *= 2
	}
}

// 题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
// 考察点 ： go 关键字的使用、协程的并发执行。
// 题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
// 考察点 ：协程原理、并发任务调度。

func twoGo() {
	go func() {
		fmt.Println("1到10的奇数:")
		for i := 1; i <= 10; i += 2 {
			fmt.Println("奇数:", i)
		}
	}()

	go func() {
		fmt.Println("1到10的偶数:")
		for i := 2; i <= 10; i += 2 {
			fmt.Println("偶数:", i)
		}
	}()
}

func taskDo(funcs []func(int) int) {
	for index, dofunc := range funcs {
		go func(i int) {
			start := time.Now()
			dofunc(i + 1)
			fmt.Printf("任务；%d 耗时: %dms\n", i+1, time.Since(start).Milliseconds())
		}(index)
	}
}

func func1(i int) int {
	time.Sleep(time.Duration(i*100) * time.Millisecond)
	return i
}

// 题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
// 考察点 ：接口的定义与实现、面向对象编程风格。
// 题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
// 考察点 ：组合的使用、方法接收者。

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	radius float64
}

type Rectangle struct {
	length float64
	width  float64
}

func (c *Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

func (r *Rectangle) Area() float64 {
	return r.length * r.width
}

func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.length + r.width)
}

type Person struct {
	name string
	age  int
}

type Employee struct {
	person     Person
	employeeId int
}

func (e *Employee) PrintInfo() {
	//输出员工的信息。
	fmt.Printf("员工信息：姓名=%s, 年龄=%d, 员工ID=%d\n", e.person.name, e.person.age, e.employeeId)
}

// 题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
// 考察点 ：通道的基本使用、协程间通信。
// 题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
// 考察点 ：通道的缓冲机制。
func sendChannel(i chan<- int) {
	for j := 1; j <= 10; j++ {
		i <- j
	}
	close(i)

}

func goChannel(i <-chan int) {
	for i := range i {
		fmt.Println(i)
	}

}

// ✅锁机制
// 题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ： sync.Mutex 的使用、并发数据安全。
// 题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ：原子操作、并发数据安全。

func syncCounter() {
	var mutex sync.Mutex
	counter := 0

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				mutex.Lock()
				counter++
				mutex.Unlock()
			}
		}()
	}
	time.Sleep(2 * time.Second)
	fmt.Println("有锁计数器counter:", counter)

}

func noMutexCounter() {
	var counter atomic.Int64
	counter.Store(0)

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				counter.Add(1)
			}
		}()
	}
	time.Sleep(2 * time.Second)
	fmt.Println("无锁计数器counter:", counter.Load())

}

func main() {
	v := 1
	a := &v
	fmt.Printf("int指针a的值为: %d \n", *a)
	goFunc(a)
	fmt.Printf("int指针a的值增加10后的值为: %d\n", *a)

	a1 := []int{1, 2, 3, 4}
	b := a1[0:2]
	fmt.Println("切片:", b)
	doubleSlice(b)
	fmt.Println("切片中的每个元素乘以2后得到：", b)

	fmt.Println("两个go程同时执行:")
	twoGo()

	funcV := []func(int) int{func1, func1}
	taskDo(funcV)

	time.Sleep(1 * time.Second)
	c := Circle{radius: 3}
	d := Rectangle{length: 5, width: 4}

	fmt.Println("圆的周长和面积:", c.Perimeter(), c.Area())
	fmt.Println("长方形的周长和面积", d.Perimeter(), d.Area())

	person := Person{age: 18, name: "张三"}
	e := Employee{person: person, employeeId: 111}
	e.PrintInfo()

	channelParam := make(chan int)
	go sendChannel(channelParam)
	go goChannel(channelParam)

	time.Sleep(1 * time.Second)

	// 创建缓冲区大小为10的通道
	ch := make(chan int, 10)

	// 生产者协程
	go func() {
		for i := 1; i <= 100; i++ {
			ch <- i
			fmt.Printf("发送: %d\n", i)
		}
		close(ch)
	}()

	// 消费者协程
	go func() {
		for val := range ch {
			fmt.Printf("接收: %d\n", val)
			time.Sleep(10 * time.Millisecond) // 模拟处理时间
		}
	}()

	go syncCounter()

	go noMutexCounter()

	time.Sleep(3 * time.Second)
}
