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
	"github.com/eggfly/gopattern/creation/singleton"
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

// var s interface{} = nil
// var singletonLock = sync.RWMutex{}

func testRightSingleton() {
	// global value auto thread safe?
	c := make(chan bool)
	f := func() {
		if singleton.S == nil {
			singleton.S = "value"
			log.Println("init value")
		}
		c <- true
	}
	const COUNT = 1000000
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
			s = "value"
			log.Println("init value")
		}
		c <- true
	}
	const COUNT = 1000000
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
	runtime.GOMAXPROCS(runtime.NumCPU())
	testRightSingleton()
	testWrongSingleton()
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
	testSingleton()
}
