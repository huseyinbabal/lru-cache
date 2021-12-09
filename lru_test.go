package lru

import (
	"sync"
	"testing"
	"time"
)

var cache = New(9)

var students = map[int]string{
	0: "mert",
	1: "john",
	2: "jack",
	3: "ahmet",
	4: "mehmet",
	5: "veli",
	6: "matthew",
	7: "jessie",
	8: "james",
}

func getStudentWithCache(id int) string {
	val := cache.Get([]byte(string(rune(id))))
	if val != nil {

		return string(val)
	}
	student := getStudent(id)
	cache.Put([]byte(string(rune(id))), []byte(student))
	return student
}

func getStudent(id int) string {
	addLoad()
	return students[id]
}

func addLoad() {
	time.Sleep(time.Millisecond * 10)
}
func BenchmarkGetWithoutCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getStudent(i % 9)
	}
}

func BenchmarkGetWithCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getStudentWithCache(i % 9)
	}
}

func BenchmarkGetWithCacheWithoutDataRace(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(a int) {
			getStudentWithCache(a % 9)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
