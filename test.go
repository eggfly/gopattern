package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"strconv"
	//"strings"
	//"encoding/gob"
	"encoding/json"
	"github.com/eggfly/gopattern/behavior/nullobject"
	"github.com/eggfly/gopattern/behavior/templatemethod"
	"github.com/eggfly/gopattern/creation/abstractfactory"
	"github.com/eggfly/gopattern/creation/builder"
	"github.com/eggfly/gopattern/creation/factorymethod"
	"github.com/eggfly/gopattern/creation/prototype"
	"github.com/eggfly/gopattern/creation/singleton"
	"math"
	"math/big"
	"sync"
	"time"
)

const PORT = "8888"

func testGoroutine() {
	const MAXUINT = ^uint(0)
	const LOOP_COUNT = 1000 * 200
	var i uint
	for i = 0; i < LOOP_COUNT; i++ {
		go waitSeconds()
	}
	time.Sleep(time.Second * 10)
	runtime.GC()
	time.Sleep(time.Second * 10)
}

func waitSeconds() {
	time.Sleep(time.Second * 10)
}

func testSocket() {
	pServerMode := flag.Bool("server", false, "true for server mode, otherwise false")
	flag.Parse()
	var mode string
	if *pServerMode {
		mode = "server"
	} else {
		mode = "client"
	}
	fmt.Println("Run Mode:", mode)

	var err error
	if *pServerMode {
		err = startServer()
	} else {
		err = startClient()
	}

	fmt.Println("final:", err)
}
func startServer() error {
	ln, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
	return err
}

func checkPrintError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func handleConnection(conn net.Conn) error {
	log.Println(conn)

	_, err := io.WriteString(conn, "\033[1;30;41mWelcome to the chat by golang!\033[0m\n")
	checkPrintError(err)
	for {
		reader := bufio.NewReader(conn)
		line, err := reader.ReadString('\n')
		io.WriteString(conn, line)
		if err != nil {
			checkPrintError(err)
			conn.Close()
			break
		}
	}
	return err
}

func startClient() error {
	conn, err := net.Dial("tcp", "localhost:"+PORT)
	if err != nil {
		return err
	}
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	buf := bufio.NewReader(conn)
	line, err := buf.ReadString('\n')
	fmt.Println("now client close connection.")
	conn.Close()
	fmt.Println(line)
	return err
}

func printHelloWorld() {
	fmt.Println("Hello World!")
}

func testGOMAXPROCS() {
	val := runtime.GOMAXPROCS(8)
	fmt.Println("GOMAXPROCS:", val)
	val = runtime.GOMAXPROCS(0)
	fmt.Println("GOMAXPROCS:", val)
	val = runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("GOMAXPROCS:", val)
	a := func() {
		fmt.Println(runtime.CPUProfile())
		time.Sleep(time.Second * 10)
	}
	fmt.Println("NumCPU():", runtime.NumCPU())
	fmt.Println("NumCgoCall():", runtime.NumCgoCall())
	go a()
	go a()
	go a()
	go a()
	go a()
	go a()
	time.Sleep(time.Second * 11)
}

func testReaderWriter() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("MAXPROCS:", runtime.GOMAXPROCS(0))
	c := make(chan int)
	// lock := sync.RWMutex{}
	var rwLock sync.RWMutex
	if &rwLock == nil {
		fmt.Println()
	}
	read := func(id int) {
		rwLock.RLock()
		fmt.Println("goroutine:", id+1, color("read", FORE_GREEN, BACK_BLACK))
		time.Sleep(time.Millisecond * 20)
		rwLock.RUnlock()
		c <- 0
	}
	write := func(id int) {
		rwLock.Lock()
		fmt.Println("goroutine:", id+1, color("write", FORE_RED, BACK_BLACK))
		time.Sleep(time.Millisecond * 20)
		rwLock.Unlock()
		c <- 0
	}
	num_goroutines := 0
	for i := 0; i < 5; i++ {
		go write(i)
		num_goroutines++
		go read(i)
		num_goroutines++
	}
	for i := 0; i < num_goroutines; i++ {
		<-c
	}

	fmt.Println("\033[32;49;1m [DONE] \033[39;49;0m")
	fmt.Print(color("green on red", FORE_GREEN, BACK_RED))
	fmt.Print(color("black on white", FORE_BLACK, BACK_WHITE))
	fmt.Println()
}

type Color int

const (
	FORE_BLACK Color = 30 + iota
	FORE_RED
	FORE_GREEN
	FORE_YELLOW
	FORE_BLUE
	FORE_MAGENTA
	FORE_CYAN
	FORE_WHITE
)
const (
	BACK_BLACK Color = 40 + iota
	BACK_RED
	BACK_GREEN
	BACK_YELLOW
	BACK_BLUE
	BACK_MAGENTA
	BACK_CYAN
	BACK_WHITE
)

func color(origin string, foreground Color, background Color) string {
	return fmt.Sprintf("\033[%d;%d;1m%s\033[0m", foreground, background, origin)
}
func testConcurrentPrint() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	c := make(chan bool)
	f := func(id int) {
		a := strconv.Itoa(id)
		for i := 0; i < 10; i++ {
			fmt.Print(a)
		}
		fmt.Println()
		c <- true
	}
	for i := 0; i < 10; i++ {
		go f(i)
	}
	for i := 0; i < 10; i++ {
		<-c
	}
}

type Task struct {
	channel chan interface{}
}

func (this Task) New() *Task {
	return &Task{}
}
func NewTask(f func(args []interface{}) interface{}, args []interface{}) *Task {
	t := Task{}
	t.channel = make(chan interface{}, 1)
	wrapperFunc := func() {
		t.channel <- f(args)
	}
	go wrapperFunc()
	return &t
}

func (this *Task) Get() interface{} {
	return <-this.channel
}

func printMem() {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	fmt.Println(m.TotalAlloc, m.EnableGC, m.NumGC, m.Sys)
}

func testTask() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	printMem()
	f := func(args []interface{}) interface{} {
		fmt.Println("args:", args)
		return true
	}
	tasks := []*Task{}
	for i := 0; i < 17; i++ {
		task := NewTask(f, []interface{}{i})
		tasks = append(tasks, task)
	}
	for i, task := range tasks {
		fmt.Println("result", i, "is", task.Get())
		printMem()
	}
	fmt.Println("len of tasks:", len(tasks), "cap of tasks", cap(tasks))
	t := Task.New(Task{})
	fmt.Println(t)
}
func testSwitch() {
	ch := make(chan int)
	ch <- 1
	select {
	case v := <-ch:
		fmt.Println(v)
	}
}

func testOneWayChannel() {
	go func(<-chan interface{}) chan<- int {
		return make(chan<- int)
	}(make(<-chan interface{}))
}
func test2B() {
	a := 1
	a = a
}

type Request int

func Serve(queue chan *Request) {
	for req := range queue {
		//<-sem
		// Buggy
		req := req // Create new instance of req for the goroutine.
		go func() {
			//process(req)
			func(req *Request) {
			}(req)
			//sem <- 1
		}()
	}
}

func testFibonacci() {
	printTestHeader("testFibonacci")
	sqrt5 := math.Sqrt(5)
	const1 := sqrt5 / 5
	const2 := (sqrt5 + 1) / 2
	const3 := (1 - sqrt5) / 2
	round := func(f float64) int64 {
		return int64(math.Floor(f + 0.5))
	}
	fibonacciFormula := func(n int64) int64 {
		nn := float64(n)
		return round(const1 * (math.Pow(const2, nn) - math.Pow(const3, nn)))
	}
	t := time.Now()
	for i := 0; i < 10; i++ {
		fmt.Println(fibonacciFormula(int64(i)))
	}
	fmt.Println(time.Now().Sub(t))
}

func testDecimal() {
	printTestHeader("testDecimal")
	i := big.NewInt(math.MaxInt64)
	i.Mul(i, i)
	i.Mul(i, i)
	i.Mul(i, i)
	i.Mul(i, i)
	i.Mul(i, i)
	i.Mul(i, i)
	fmt.Println(i)
	ii := new(big.Int)
	fmt.Println(ii)
}

// var s interface{} = nil
// var singletonLock = sync.RWMutex{}

func testRightSingleton() {
	// global value auto thread safe?
	c := make(chan bool)
	f := func() {
		if singleton.S == nil {
			log.Println("init value")
			singleton.S = "value"
		}
		c <- true
	}
	const COUNT = 10
	for i := 0; i < COUNT; i++ {
		go f()
	}
	for i := 0; i < COUNT; i++ {
		<-c
	}
}
func testWrongSingleton() {
	c := make(chan bool)
	var s interface{} = nil
	f := func() {
		if s == nil {
			log.Println("init value")
			s = "value"
		}
		c <- true
	}
	const COUNT = 10
	for i := 0; i < COUNT; i++ {
		go f()
	}
	for i := 0; i < COUNT; i++ {
		<-c
	}
	// result may like this:
	// 2014/01/16 23:07:00 init value
	// 2014/01/16 23:07:00 init value
	// 2014/01/16 23:07:00 init value
	// 2014/01/16 23:07:00 init value
	// 2014/01/16 23:07:00 init value
	// 2014/01/16 23:07:00 init value
	// 2014/01/16 23:07:00 init value
}
func testSingleton() {
	printTestHeader("testSingleton")
	runtime.GOMAXPROCS(runtime.NumCPU())
	testRightSingleton()
	testWrongSingleton()
}
func testPackageVarAccessSpeed() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// a local var
	var s interface{} = nil

	// access package var
	accessPackageVar := func() {
		_ = singleton.S
	}
	// access local var
	accessLocalVar := func() {
		_ = s
	}

	ch := make(chan bool)
	// benchmark function
	benchmark := func(accessFunc func()) {
		t := time.Now()
		const COUNT = 100000
		for i := 0; i < COUNT; i++ {
			go func(accessFunc func()) {
				accessFunc()
				ch <- true
			}(accessFunc)
		}
		// wait for end
		for i := 0; i < COUNT; i++ {
			<-ch
		}

		fmt.Println("time:", time.Now().Sub(t))
	}

	benchmark(accessPackageVar)
	benchmark(accessLocalVar)
}

func testPackageVarThreadSafe() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ch := make(chan bool)

	const COUNT = 100000
	for i := 0; i < COUNT; i++ {
		go func() {
			singleton.I++
			ch <- true
		}()
	}
	// wait for end
	for i := 0; i < COUNT; i++ {
		<-ch
	}

	fmt.Println("singleton.I:", singleton.I)
	// finally it's not safe
}

func testAbstractFactory() {
	printTestHeader("testAbstractFactory")
	var ibm abstractfactory.IComputerFactory = abstractfactory.IBMFactory{}
	desktop := ibm.CreateDesktopProduct()
	laptop := ibm.CreateLaptopProduct()
	desktop.PlayDesktop()
	laptop.PlayLaptop()
	desktop = abstractfactory.AppleFactory{}.CreateDesktopProduct()
	desktop.PlayDesktop()
}

func testBuilder() {
	printTestHeader("testBuilder")
	director := builder.BuildHouseDirector{}
	director.SetHouseBuilder(&builder.YoungHouseBuilder{})
	director.BuildHouse()
	fmt.Println(director.GetHouse())
}

func testFactoryMethod() {
	printTestHeader("testFactoryMethod")
	var factory factorymethod.IButtonFactory = factorymethod.WindowsButtonFactory{}
	factory.CreateButton().Click()
	factory = factorymethod.MacButtonFactory{}
	factory.CreateButton().Click()
}

func testPrototype() {
	printTestHeader("testPrototype")
	templateBullet := prototype.Bullet{}
	templateBullet.Init()
	templateBullet.X = 2

	bullet1 := prototype.Bullet{}
	bullet1.CopyFrom(templateBullet)
	fmt.Println(bullet1)
	bullet2 := bullet1.Clone()
	fmt.Println(bullet2)
}

func testCreatePrivate() {
	_ = templatemethod.CreatePrivate()
}

func testTemplateMethod() {
	monopoly := templatemethod.Monopoly{}
	monopoly.IGame = monopoly // this
	monopoly.PlayOneGame(3)
	game := templatemethod.Chess{}
	game.IGame = game // this
	game.PlayOneGame(4)
}

func testNullObject() {
	printTestHeader("testNullObject")
	n := nullobject.Null
	n.Do()
	n = nullobject.ConcreteImpl{}
	n.Do()
}

func testJSON() {
	printTestHeader("testJSON")
	type Message struct {
		A, B int
	}
	m := Message{1, 999}
	b, err := json.Marshal(m)
	fmt.Println(string(b), err)
}

func testGob() {
	printTestHeader("testGob")

}

func printTestHeader(name string) {
	const SIGN = " **** "
	println(SIGN + name + SIGN)
}

func main() {
	//testSocket()
	//testGoroutine()
	//printHelloWorld()
	//testGOMAXPROCS()
	//testReaderWriter()
	//testConcurrentPrint()
	//testTask()
	//testSwitch()
	//testOneWayChannel()
	//test2B()
	//testPackageVarAccessSpeed()
	//testPackageVarThreadSafe()
	testFibonacci()
	testDecimal()
	testCreatePrivate()
	testJSON()
	testGob()

	testSingleton()
	testAbstractFactory()
	testBuilder()
	testFactoryMethod()
	testPrototype()
	testTemplateMethod()
	testNullObject()
}
