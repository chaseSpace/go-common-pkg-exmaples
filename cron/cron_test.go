package cron

import (
	"github.com/robfig/cron/v3"
	"log"
	"sync"
	"testing"
	"time"
)

var secConf = "* * * * * ?"

/*
cron有两种模式，一个支持秒，一个是默认不支持
秒模式下，配置项6个（任一模式下最后一位都是问号?）
`* * * * * ?` 表示每秒执行
`* 3 * * * ?` 表示分钟数=3时的每一秒都执行
`0 3 * * * ?` 表示分钟数=3时的第0秒都执行
`0 0/3 * * * ?` 表示每隔3分钟的第0秒执行
`1 0-30/3 * * * ?` 表示在分钟等于0-30的范围内每隔3分钟的第1秒执行
*/

type task struct {
	Name    string
	counter int
	wg      *sync.WaitGroup
}

func (t *task) Run() {
	t.counter++
	log.Println(t.Name, t.counter)
	t.wg.Done()
}

func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}

func TestCron(t *testing.T) {
	// 启动支持秒级别的定时器
	cronItem := cron.New(cron.WithSeconds())
	var wg = &sync.WaitGroup{}
	wg.Add(1011)

	// 两种方式添加job
	job1, _ := cronItem.AddJob(secConf, &task{"job1", 0, wg})
	job2 := cronItem.Schedule(cron.Every(1*time.Second), &task{"job2", 0, wg})

	// 每分钟的0秒执行（第一次执行是启动后的第一个0秒）
	_, _ = cronItem.AddFunc("0 0/2 * * * ?", func() {
		log.Printf("job3")
	})

	if name := cronItem.Entry(job1).Job.(*task).Name; name != "job1" {
		t.Errorf("job1's Name is wrong [%s]", name)
	}
	if name := cronItem.Entry(job2).Job.(*task).Name; name != "job2" {
		t.Errorf("job2's Name is wrong [%s]", name)
	}
	cronItem.Start()
	defer cronItem.Stop()
	//var ctx, cancel = context.WithCancel(context.Background())
	//go func() {
	//	time.Sleep(3 * time.Second)
	//	cronItem.Remove(job1)
	//	time.Sleep(1 * time.Minute)
	//	cronItem.Remove(job2)
	//	log.Printf("all job stopped!")
	//	cancel()
	//}()
	select {
	//case <-ctx.Done():
	case <-wait(wg):
	}
}

func TestMinute(t *testing.T) {
	cronItem := cron.New()
	var wg = &sync.WaitGroup{}
	wg.Add(11111)

	// 每分钟的0秒执行一次（第一次执行是启动后的第一个0秒）
	_, _ = cronItem.AddFunc("* * * * ?", func() {
		log.Printf("minute job1")
		wg.Done()
	})
	// 每2分钟的0秒执行一次（第一次执行是启动后的第【二】个0秒）
	_, _ = cronItem.AddFunc("0/2 * * * ?", func() {
		log.Printf("minute job2")
		wg.Done()
	})
	// “9/2”中的9不会识别(语法错误，如果只能是minute=9时执行使用9-9/2，但这个语法明显有歧义)
	// 仍然是每2分钟的0秒执行一次（第一次执行是启动后的第一个0秒）
	_, _ = cronItem.AddFunc("9/2 * * * ?", func() {
		log.Printf("minute job3")
		wg.Done()
	})
	// “0-1/2”中的0-1表示在minute为0到1的范围内每2分钟执行一次（如11:02-11:59这个范围都不会执行）
	_, _ = cronItem.AddFunc("0-1/2 * * * ?", func() {
		log.Printf("minute 0-1/2")
		wg.Done()
	})
	// “0-30/2”中的0-30表示在minute为0到30的范围内每2分钟执行一次（如11:31-11:59这个范围都不会执行）
	_, _ = cronItem.AddFunc("0-30/2 * * * ?", func() {
		log.Printf("minute 0-30/2")
		wg.Done()
	})
	_ = cronItem.Schedule(cron.Every(time.Second), &task{"sec job", 0, wg})
	cronItem.Start()
	defer cronItem.Stop()
	select {
	case <-wait(wg):
	}
}

func TestAny(t *testing.T) {
	cronItem := cron.New(cron.WithSeconds())
	var wg = &sync.WaitGroup{}
	wg.Add(1111)
	_, _ = cronItem.AddFunc("* 47 * * * ?", func() {
		log.Printf("minute job1")
		wg.Done()
	})

	_ = cronItem.Schedule(cron.Every(time.Second), &task{"sec job", 0, wg})
	cronItem.Start()
	defer cronItem.Stop()
	select {
	case <-wait(wg):
	}
}
