package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type MyMap struct {
	mu       sync.Mutex
	elements map[int]int
	chMap    map[int]chan int //每个key对应的channel
}

func NewMyMap() *MyMap {
	return &MyMap{
		mu:       sync.Mutex{},
		elements: make(map[int]int),
		chMap:    make(map[int]chan int),
	}
}

func (m *MyMap) Get(k int, maxWaitingTime time.Duration) (int, error) {
	m.mu.Lock()
	v, ok := m.elements[k] //若读到内容就返回值
	if ok {
		return v, nil
	}
	m.mu.Unlock()        //未读到内容就解开锁，让其他线程可进入临界区
	ch := make(chan int) //创建一个channel
	m.chMap[k] = ch      //放到对应位置
	select {             //开始监听此channel，若在maxWaitingTime内得到值，则返回值
	case val := <-ch:
		return val, nil
	case <-time.After(maxWaitingTime): //若超时则返回error
		m.mu.Lock()
		delete(m.chMap, k) //从chMap删掉并关闭对应channel防止影响到后面
		close(ch)
		m.mu.Unlock()
		return 0, errors.New("timeout")
	}
}
func (m *MyMap) Put(k, v int) {
	m.mu.Lock()
	m.elements[k] = v
	if ch, ok := m.chMap[k]; ok { //判断是否有对应channel等待接受值
		ch <- v
		delete(m.chMap, k)
		close(ch)
	}
	m.mu.Unlock()
}
func main() {

	m := NewMyMap()
	//一个线程停3秒 一个线程停2秒
	go func() {
		val, err := m.Get(1, 3*time.Second)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(val)
	}()
	go func() {
		val, err := m.Get(2, 2*time.Second)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(val)
	}()
	time.Sleep(2 * time.Second) //两秒后再输入数据
	m.Put(1, 1000)
	m.Put(2, 200)
	time.Sleep(5 * time.Second) //等待程序完成
}
