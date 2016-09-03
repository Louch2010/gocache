package main

import (
	"log"
	"strconv"
	"time"

	"github.com/louch2010/gocache/core"
)

func main() {
	log.Println("测试开始")
	table, _ := core.Cache("test")
	//table := core.DefaultCache()
	table.Add("name", "luociang", 3*time.Second)
	table.Add("age", 26, 5*time.Second)
	table.Add(110, "北京", 15*time.Second)
	table.Add("home", "江西", 10*time.Second)

	log.Println("-----------1秒后---------------")
	time.Sleep(1 * time.Second)
	table.Get("name")
	table.Get("age")
	table.Get(110)
	table.Get("home")

	log.Println("----------删除name后----------------")
	table.Delete("name")
	table.Get("name")
	table.Get("age")
	table.Get(110)
	table.Get("home")

	log.Println("----------age过期后----------------")
	time.Sleep(5 * time.Second)
	table.Get("name")
	table.Get("age")
	table.Get(110)
	table.Get("home")

	log.Println("----------home重设值后----------------")
	table.Add("home", "上饶", 0)
	table.Get("name")
	table.Get("age")
	table.Get(110)
	table.Get("home")

	log.Println("----------全部过期后----------------")
	time.Sleep(15 * time.Second)
	table.Get("name")
	table.Get("age")
	table.Get(110)
	table.Get("home")

	for i := 0; i < 10; i++ {
		log.Println("表容量：", strconv.Itoa(table.ItemCount()))
		time.Sleep(2 * time.Second)
	}

}
