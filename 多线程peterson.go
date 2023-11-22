package main

import (
	"fmt"
)

type MyLock struct {
	level []int
	line  []int
	n     int
}

func NewMyLock(n int) *MyLock {
	m := MyLock{}
	m.level = make([]int, n+1)
	m.line = make([]int, n+1)
	m.n = n
	return &m
}

func (m *MyLock) Lock(code int) {

	// i 是优先级
	for i := 0; i < m.n; i++ { //最低优先级为0，最高优先级为n
		m.level[code] = i          //把当前线程优先级设为i
		m.line[i] = code           //把当前线程放到优先级为i的line里等待
		for j := 0; j < m.n; j++ { //遍历所有其他线程
			if j != code {
				//如果 优先级为i的line是当前线程 && 存在某个线程等级高于当前线程优先级  就等待
				for m.line[i] == code && m.level[j] >= i {
					// 等待ing...
				}
			}

		}
	}
}

func (m *MyLock) UnLock(code int) {
	fmt.Println(code, "out")
	m.level[code] = -1 //把线程优先级设为-1,表示已完成
}

func main() {
	//假设有十个线程
	m := NewMyLock(10)
	for i := 0; i < 9; i++ {
		//GO没有暴露线程的唯一标识符，于是就用for循环赋值来代替标识符了
		go func(a int) {

			m.Lock(a)
			fmt.Println(a, "in")
			m.UnLock(a)

		}(i)
	}
	select {}
}
