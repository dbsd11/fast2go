package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

var jobQuene chan Job
var workPool chan chan Job

func main() {
	jobQueueSize := 100;
	workerNum := 4;

	StartDispatcher(jobQueueSize, workerNum)

	if error := fasthttp.ListenAndServe(":65532", func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("text/plain; charset=utf8")

		// Set arbitrary headers
		ctx.Response.Header.Set("X-My-Header", "my-header-value")

		// Set cookies
		var c fasthttp.Cookie
		c.SetKey("cookie-name")
		c.SetValue("cookie-value")
		ctx.Response.Header.SetCookie(&c)

		jobQuene <- Job{
			Run:func() {
				fmt.Printf("run handler 01 at connectionID: %s\n", ctx.ConnTime());
			},
		}
		jobQuene <- Job{
			Run:func() {
				fmt.Printf("run handler 02 at connectionID: %s\n", ctx.ConnTime());
			},
		}
		jobQuene <- Job{
			Run:func() {
				fmt.Printf("run handler 03 at connectionID: %s\n", ctx.ConnTime());
			},
		}
	}); error != nil {
		fmt.Println(error.Error())
	}

	fmt.Println("server started")

}

type Job struct {
	Run func()
}
type Worker struct {
	ID         int
	JobChannel chan Job
	WorkerPool chan chan Job
	QuitChan   chan bool
}

func (worker *Worker) start() {
	go func() {
		for {
			worker.WorkerPool <- worker.JobChannel
			select {
			case job := <-worker.JobChannel:
				job.Run()
			case <-worker.QuitChan:
				return
			}
		}
	}()
}

func (worker *Worker) stop() {
	go func() {
		worker.QuitChan <- true;
	}()
}

func StartDispatcher(jobQueueSize int, workerNum int) {
	jobQuene = make(chan Job, jobQueueSize)
	workPool = make(chan chan Job, workerNum)
	for i := 0; i < workerNum; i++ {
		worker := &Worker{
			ID:i,
			JobChannel:make(chan Job),
			WorkerPool:workPool,
			QuitChan:make(chan bool),
		};
		worker.start();
	}
	go func() {
		for {
			select {
			case job := <-jobQuene:
				go func() {
					<-workPool<- job
				}()
			}
		}
	}()
}