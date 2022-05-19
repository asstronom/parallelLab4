package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const (
	arrayLen = 5
)

func SortFuture(array []int64) chan []int64 {
	future := make(chan []int64)
	go func() {
		sort.Slice(array, func(i, j int) bool {
			return array[i] < array[j]
		})
		future <- array
	}()
	return future
}

func MultiplyFuture(array []int64) chan []int64 {
	future := make(chan []int64)
	go func() {
		for i := range array {
			array[i] *= 2
		}
		future <- array
	}()
	return future
}

func DivideFuture(array []int64, divider int64) chan []int64 {
	future := make(chan []int64)
	go func() {
		for i := range array {
			array[i] /= divider
		}
		future <- array
	}()
	return future
}

func RemoveOddFuture(array []int64) chan []int64 {
	future := make(chan []int64)
	go func() {
		newarr := make([]int64, 0)
		for _, v := range array {
			if v%2 == 0 {
				newarr = append(newarr, v)
			}
		}
		future <- newarr
	}()
	return future
}

func RemoveRangeFuture(array []int64, lo, hi float64, max int64) chan []int64 {
	future := make(chan []int64)
	go func() {
		low := int64(float64(max) * lo)
		high := int64(float64(max) * hi)
		newarr := make([]int64, 0)
		for _, v := range array {
			if v >= low && v <= high {
				newarr = append(newarr, v)
			}
		}
		future <- newarr
	}()
	return future
}

func NumberGenerator(ch chan int64) {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	for {
		ch <- generator.Int63n(50)
	}
}

func FillArrayFuture(array []int64, rng chan int64) chan []int64 {
	future := make(chan []int64)
	go func() {
		for i := range array {
			array[i] = <-rng
		}
		future <- array
	}()
	return future
}

func DumpArrFuture(sourcearray, destinationarray []int64, startidx int) chan []int64 {
	future := make(chan []int64)
	go func() {
		for i, v := range sourcearray {
			destinationarray[startidx+i] = v
		}
		future <- destinationarray
	}()
	return future
}

func main() {
	rng := make(chan int64)
	go NumberGenerator(rng)

	ar1, ar2, ar3 := make([]int64, arrayLen), make([]int64, arrayLen), make([]int64, arrayLen)

	f1 := func(array []int64) chan []int64 {
		future := make(chan []int64)
		go func() {
			for i := range array {
				array[i] = <-rng
			}
			fmt.Println("Filled arr1:\n", ar1)

			sort.Slice(array, func(i, j int) bool {
				return array[i] < array[j]
			})
			fmt.Println("Sorted arr1:\n", ar1)

			for i := range array {
				array[i] *= 2
			}
			fmt.Println("Multiplied by 2 arr1:\n", ar1)

			for i := range array {
				array[i] /= 2
			}
			fmt.Println("Divided by 2 arr1:\n", ar1)

			future <- array
		}()
		return future
	}(ar1)

	f2 := func(array []int64) chan []int64 {
		future := make(chan []int64)
		go func() {
			for i := range array {
				array[i] = <-rng
			}
			fmt.Println("Filled arr2:\n", array)

			sort.Slice(array, func(i, j int) bool {
				return array[i] < array[j]
			})
			fmt.Println("Sorted arr2:\n", array)

			newarr := make([]int64, 0)
			for _, v := range array {
				if v%2 == 0 {
					newarr = append(newarr, v)
				}
			}
			array = newarr
			fmt.Println("Removed odd arr2:\n", array)

			for i := range array {
				array[i] /= 2
			}
			fmt.Println("Divided by 2 arr2:\n", array)

			future <- array
		}()
		return future
	}(ar2)

	f3 := func(array []int64) chan []int64 {
		future := make(chan []int64)
		go func() {
			for i := range array {
				array[i] = <-rng
			}
			fmt.Println("Filled arr3:\n", array)

			sort.Slice(array, func(i, j int) bool {
				return array[i] < array[j]
			})
			fmt.Println("Sorted arr3:\n", array)

			max := array[len(array)-1]

			low := int64(float64(max) * 0.3)
			high := int64(float64(max) * 0.6)
			newarr := make([]int64, 0)
			for _, v := range array {
				if v >= low && v <= high {
					newarr = append(newarr, v)
				}
			}
			array = newarr
			fmt.Println("Left range [", low, high, "] arr3:\n", array)

			for i := range array {
				array[i] /= 2
			}
			fmt.Println("Divided by 2 arr3:\n", array)

			future <- array
		}()
		return future
	}(ar3)

	ar1 = <-f1
	ar2 = <-f2
	ar3 = <-f3

	finalarr := make([]int64, len(ar1)+len(ar2)+len(ar3))

	fmt.Print("Arr1 after operations:\n", ar1, "\n")
	fmt.Print("Arr2 after operations:\n", ar2, "\n")
	fmt.Print("Arr3 after operations:\n", ar3, "\n")

	//fmt.Println(len(ar1), len(ar2), len(ar3))

	f1 = DumpArrFuture(ar1, finalarr, 0)
	f2 = DumpArrFuture(ar2, finalarr, len(ar1))
	f3 = DumpArrFuture(ar3, finalarr, len(ar1)+len(ar2))

	<-f1
	<-f2
	<-f3

	fmt.Println("Dumped into shared array")
	fmt.Println(finalarr)

}
