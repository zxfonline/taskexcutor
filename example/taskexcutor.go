package main

import (
	"fmt"

	"math/rand"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zxfonline/golog"
	"github.com/zxfonline/taskexcutor"
	"github.com/zxfonline/timer"
)

func randInt(min int32, max int32) int32 {
	return min + rand.Int31n(max-min)
}

var logger = golog.New("ExcutorTest")

func main() {
	exc := taskexcutor.NewTaskExcutor(65536)
	num := 10
	for i := 0; i < num; i++ {
		go func(i int) {
			obj := &Obj{
				name: fmt.Sprintf("name_%d", i+1),
				t1:   T1{name: fmt.Sprintf("t1_%d", i+1)},
			}
			ev := timer.AddTimerEvent(taskexcutor.NewTaskService(obj.t1.Processing), fmt.Sprintf("event1_%d", i+1), 0, time.Duration((100+int(randInt(1, 50))*i)*1e6), num+1-i, true)
			obj.t1.event = ev
		}(i)
	}

	go func() {
		for {
			select {
			case event := <-exc:
				event.Call(logger)
				//				exc.Close()
			}
		}
	}()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-ch
}

type Obj struct {
	name string
	t1   T1
}
type T1 struct {
	name  string
	event *timer.TimerEvent
}

func (t *T1) Processing(param ...interface{}) {
	fmt.Printf("%p,e=%+v\n", t, t)
}
