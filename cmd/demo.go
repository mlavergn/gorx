package main

import (
	"log"
	"sync"
	"time"

	rx "github.com/mlavergn/gorx"
)

func demoInterval(observer *rx.Observer) *rx.Observable {
	interval := rx.NewInterval(200)
	interval.UID = "demoInterval." + interval.UID
	return interval.Take(10)
}

func demoSubject(observer *rx.Observer) *rx.Observable {
	interval := rx.NewInterval(200)
	interval.UID = "demoSubjectInterval." + interval.UID
	subject := rx.NewSubject()
	subject.UID = "demoSubjectObservable." + subject.UID
	interval.Pipe(subject)
	subject.Subscribe <- observer
	subject.Next <- 99
	return subject.Take(10)
}

func demoBehavior(observer *rx.Observer) *rx.Observable {
	return rx.NewBehaviorSubject(99).Take(1)
}

func demoReplay(observer *rx.Observer, count int) *rx.Observable {
	subject := rx.NewReplaySubject(count)
	subject.UID = "demoReplayObservable." + subject.UID

	subject.Next <- 11
	subject.Next <- 22
	subject.Next <- 33
	subject.Next <- 44
	subject.Next <- 55
	subject.Next <- 66
	subject.Next <- 77
	subject.Next <- 88
	subject.Next <- 99
	return subject.Take(5)
}

func demoRetry(observer *rx.Observer) *rx.Observable {
	subject, err := rx.NewHTTPTextSubject("http://httpbin.org/get", nil)
	if err != nil {
		log.Println("demoRetry", err)
		return nil
	}
	subject.UID = "demoRetryObservable." + subject.UID
	retry := 2
	subject.RepeatWhen(func() bool {
		retry--
		<-time.After(1 * time.Second)
		return (retry != 0)
	}).Take(1)
	return subject
}

func demoSSE(observer *rx.Observer) *rx.Observable {
	subject, err := rx.NewHTTPSSESubject("http://demo.howopensource.com/sse/stocks.php", nil)
	if err != nil {
		log.Println("demoSSE", err)
		return nil
	}
	subject.UID = "demoSSE." + subject.UID
	subject.Map(func(event interface{}) interface{} {
		result := event.(map[string]string)
		return result
	}).Take(5)
	return subject
}

func main() {
	rx.Config(false)
	closeCh := make(chan bool)
	observer := rx.NewObserver()
	observer.UID = "demoSubscription." + observer.UID

	parse := true
	// subject := demoInterval(observer)
	// subject := demoSubject(observer)
	// subject := demoBehavior(observer)
	subject := demoReplay(observer, 4)

	// subject := demoRetry(observer)
	// subject := demoSSE(observer)
	// parse = false

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		wg.Done()
		for {
			select {
			case event := <-observer.Next:
				log.Println(observer.UID, event)
				if parse {
					v := event.(int)
					if v == 99 || v == -1 {
						log.Println("Done")
						subject.Unsubscribe <- observer
					}
				}
			}
		}
	}()

	wg.Wait()
	subject.Subscribe <- observer

	<-subject.Finalize
	close(closeCh)
}
