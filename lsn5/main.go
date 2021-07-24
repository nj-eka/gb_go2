package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {


	// 0.2
	fmt.Println("0.2: умножение двух\nматриц.")
	//m1 := [][]int{
	//	{1, 1, 1},
	//	{2, 2, 2},
	//	{3, 3, 3},
	//}
	//m2 := [][]int{
	//	{4, 4, 4},
	//	{5, 5, 5},
	//	{6, 6, 6},
	//}
	m1 := generateMatrix(5,4)
	fmt.Printf("m1:\n%v\n", m1)
	m2 := generateMatrix(4,3)
	fmt.Printf("m2:\n%v\n", m2)
	mr, _ := multiplyMatrices(m1, m2)
	fmt.Printf("m1 * m2:\n%v\n", mr)
	fmt.Println("done")

	// 1:
	fmt.Println("1: запускает n(=100) потоков и дожидается завершения их всех.")
	completeGoroutines(100)
	fmt.Println("done")

	// 2
	fmt.Println("2: разблокировка мьютекса с помощью defer")
	deferMutex()
	fmt.Println("done")
}

func completeGoroutines(count int) {
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Second)
		}()
	}
	wg.Wait()
	fmt.Printf("All %d goroutines completed\n", count)
}

func deferMutex() {
	var m sync.Mutex
	m.Lock()
	defer mutexUnlock(&m)
}
func mutexUnlock(m *sync.Mutex) {
	m.Unlock()
}

func generateMatrix(rowNum, colNum int) [][]int{
	m :=  make([][]int, rowNum)
	rand.Seed(time.Now().UnixNano())
	min, max := -1000, 1000
	for i := 0; i < rowNum; i++ {
		m[i] = make([]int, colNum)
		for j := 0; j < colNum; j++ {
			m[i][j] = rand.Intn(max-min) + min
		}
	}
	return m
}

// C[i,j] = Sum(A[i,k]*B[k,j])
func multiplyMatrices (m1 [][]int, m2 [][]int) ([][]int, error) {
	if len(m2) != len(m1[0]) {
		return nil, errors.New("inapplicable dimensions")
	}
	res := make([][]int, len(m1))
	wg := sync.WaitGroup{}
	for i := 0; i < len(m1); i++ {
		res[i] = make([]int, len(m2[0]))
		for j := 0; j < len(m2[0]); j++ {
			for k := 0; k < len(m2); k++ {
				wg.Add(1)
				go func(i, j, k int) {
					defer wg.Done()
					res[i][j] += m1[i][k] * m2[k][j]
				}(i, j, k)
			}
		}
	}
	wg.Wait()
	return res, nil
}

type Set struct {
	sync.Mutex
	mm map[int]int
}

func NewSet() *Set {
	return &Set{
		mm: map[int]int{},
	}
}

func (s *Set) Add(i int) {
	s.Lock()
	defer s.Unlock()
	s.mm[i] = i
}

func (s *Set) Has(i int) bool {
	s.Lock()
	defer s.Unlock()
	_, ok := s.mm[i]
	return ok
}

type SetRW struct {
	sync.RWMutex
	mm map[int]int
}

func NewSetRW() *SetRW {
	return &SetRW{
		mm: map[int]int{},
	}
}

func (s *SetRW) AddRW(i int) {
	s.Lock()
	defer s.Unlock()
	s.mm[i] = i
}

func (s *SetRW) HasRW(i int) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.mm[i]
	return ok
}
