package middleware

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kataras/iris"
)

// Queue is the struct to store information
type Item struct {
	Domain   string
	Weight   int
	Priority int

	Index int // The index of the item in the heap.
}

// Que stack declaration
type PriorityQueue []*Item

// FinalQueue
var FinalQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest based on expiration number as the priority
	return pq[i].Weight+pq[i].Priority < pq[j].Weight+pq[j].Priority
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

// Repository should implement common methods
type Repository interface {
	Read() []*Item
}

// Read Repository interface implementation
func (q *Item) Read() []*Item {
	path, _ := filepath.Abs("")
	file, err := os.Open(path + "/api/middleware/domain.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	var final []*Item
	tmp := &Item{}
	count := 0
	for scanner.Scan() {
		count++
		if scanner.Text() == "" {
			count = 0
			continue
		}
		switch count {
		case 1:
			tmp.Domain = scanner.Text()
		case 2:
			r := strings.Split(scanner.Text(), ":")[1]
			res, _ := strconv.Atoi(r)
			tmp.Weight = res
		case 3:
			r := strings.Split(scanner.Text(), ":")[1]
			res, _ := strconv.Atoi(r)
			tmp.Priority = res
			// persist tmp struct
			final = append(final, tmp)
			// clean tmp struct
			tmp = &Item{}
		}
	}
	return final
}

// ProxyMiddleware should queue our incoming requests
func ProxyMiddleware(c iris.Context) {
	FinalQueue := make(PriorityQueue, len(FinalQueue))

	domain := c.GetHeader("domain")
	if len(domain) == 0 {
		c.JSON(iris.Map{"status": 400, "result": "domain error"})
		return
	}
	var repo Repository
	repo = &Item{}
	for _, row := range repo.Read() {
		// fmt.Printf("Domain: %s Weight: %d\n", row.Domain, row.Weight)

		if domain == row.Domain {
			// fmt.Printf(domain)
			heap.Push(&FinalQueue, row)
			// fmt.Printf("Domain: %s Weight: %d\n", row.Domain, row.Weight)
		}
		// fmt.Printf("Domain: %s Weight: %d\n", row.Domain, row.Weight)
	}

	for i, item := range FinalQueue {
		FinalQueue[i] = item
		// fmt.Printf("Domain: %s Weight: %d\n", item.Domain, item.Weight)
		FinalQueue[i].Index = i
	}
	heap.Init(&FinalQueue)

	fmt.Printf("Length of the queu: %d \n", FinalQueue.Len())
	// for FinalQueue.Len() > 0 {
	// 	item := heap.Pop(&FinalQueue).(*Item)
	// 	fmt.Printf("Domain: %s Weight: %d\n", item.Domain, item.Weight)
	// }

	c.Next()
}
