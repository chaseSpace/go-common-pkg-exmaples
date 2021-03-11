package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
)

func TestXXX(t *testing.T) {
	startDB()
	defer closeDB()

	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		x := 80
		wg.Add(x)
		for i := 1; i <= x; i++ {
			// select
			// update

			// select
			// update
			go update(&wg, fmt.Sprintf("update concurrenct_update_test set gid=gid+2222 where uid=1"))
		}
		wg.Wait()
	}

}

func update(wg *sync.WaitGroup, sql string) {
	x := db.Exec(sql)
	log.Println(x.RowsAffected, x.Error)
	wg.Done()
}
