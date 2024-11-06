// https://linguinecode.com/post/how-to-import-local-files-packages-in-golang
package queue  

import "fmt"

type Queue struct { 
	slice []interface{} 
}

func (q *Queue) Enqueue(val interface{}) {
	q.slice = append(q.slice,val)
	// fmt.Println(len(q.slice))
}

func (q *Queue) Dequeue() {

	// var lastIndex = len(q.slice)-1
	q.slice = q.slice[1:]
	// fmt.Println(len(q.slice))
}

func (q *Queue) PeekFront() interface{} {
	return q.slice[0]
}

func (q *Queue) IsEmpty() bool {
	if len(q.slice) == 0 {
		return true 
	} else {
		return false 
	}
}

func (q *Queue) Length() int {
	return len(q.slice)
}

func (q *Queue) Print(){
	fmt.Println(q.slice)
}

func init(){
	// fmt.Println("hello from the queue")
}


// func main (){
// 	fmt.Println("hello main")

// 	var q = new(Queue)

// 	q.Enqueue(10)
// 	q.Enqueue(20)
// 	q.Enqueue(30)
// 	q.Enqueue(40)
// 	q.Enqueue(50)
// 	q.Enqueue(60)
// 	q.Print()
// 	q.Dequeue()
// 	q.Print()
// 	q.Dequeue()
// 	q.Print()
	


// }
