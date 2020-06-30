package scheduler

import (
	"container/list"
	"fmt"
	"github.com/deckarep/golang-set"
)

type TaskGraph struct {
	graph map[string]*TaskNode // 任务别名 + 任务节点组成的map构成
	todo  mapset.Set           // 待执行的任务
}

// init taskGraph 构造器
func (taskGraph TaskGraph) New() TaskGraph {
	if taskGraph.graph == nil {
		taskGraph.graph = make(map[string]*TaskNode)
	}
	if taskGraph.todo == nil {
		taskGraph.todo = mapset.NewSet()
	}
	return taskGraph
}

func (taskGraph *TaskGraph) AddTask(taskName string, deps []string) bool {
	if _, v := taskGraph.graph[taskName]; v {
		// 任务节点已经 exist 跳过
		return false
	}
	taskNode := NewNode(deps)            // 创建任务节点，指定任务的前置依赖
	taskGraph.graph[taskName] = taskNode // 完成任务名称和任务节点之间的映射关系
	return true
}

// 亮点，任务执行过程中，同时维护任务的前后依赖状态！！！
func (taskGraph *TaskGraph) InitGraph() bool {
	graph := taskGraph.graph
	topStack := list.New()
	tmpOutCounter := make(map[string]int) // 前置依赖计数
	// 依次遍历任务节点
	for taskName, node := range graph {
		// 遍历每一个TaskNode的前置依赖
		for depIter := range node.OutEdge { // 得到前置依赖的taskName
			destIter := graph[depIter] // 得到TaskNode
			if nil == destIter {
				return false
			}
			destIter.InEdge[taskName] = false
			// 反过来，更新前置任务的后置依赖 （即当前Node）数 + 1
			destIter.InCounter += 1
		}
		tmpOutCounter[taskName] = node.OutCounter
		// 如果前置依赖为0，加入到顶层栈里面
		if node.OutCounter == 0 {
			topStack.PushBack(taskName)
		}
	}

	topCount := 0
	// dfs
	for topStack.Len() > 0 {
		topCount++

		item := topStack.Front()
		topStack.Remove(item)           // 出栈
		taskName := item.Value.(string) // interface  转成具体类型

		node := graph[taskName]
		// 后置依赖 - 1，去除掉边
		for iter := range node.InEdge {
			tmpOutCounter[iter] -= 1
			if tmpOutCounter[iter] == 0 {
				topStack.PushBack(iter)
			}
		}
	}
	// circle !!!
	if topCount != len(graph) {
		return false
	}

	for iter, node := range graph {
		if node.OutCounter == 0 {
			taskGraph.todo.Add(iter)
		}
	}

	return true
}

// 需要执行的任务列表！！！
func (taskGraph *TaskGraph) GetTodoTasks() []string {
	var todo []string
	for taskName := range taskGraph.todo.Iter() {
		todo = append(todo, taskName.(string))
	}
	return todo
}

// 标记taskName代表的任务成功！！
func (taskGraph *TaskGraph) MarkTaskDone(taskName string) bool {
	if !taskGraph.todo.Contains(taskName) {
		// todo表征还没有执行完成的任务
		return false
	}

	taskGraph.todo.Remove(taskName)

	node := taskGraph.graph[taskName]

	node.Done = true
	// 任务的后置依赖更新
	for k := range node.InEdge {
		from := taskGraph.graph[k]    // fromNode
		from.OutEdge[taskName] = true // 标记后置任务，其一个前置任务（当前节点）执行完成！！！
		from.OutCounter -= 1          // 后置依赖更新其前置依赖数据-1
		if from.OutCounter == 0 {     // 如果后置依赖的入度为0, 没有需要等待的前置任务，将k 加入等待执行的任务列表！！
			taskGraph.todo.Add(k) //
		}
	}
	// 遍历当前任务的前置依赖
	for k := range node.OutEdge {
		dest := taskGraph.graph[k]   // 前置依赖的任务节点
		dest.InEdge[taskName] = true // 前置依赖都标记完成
		dest.InCounter -= 1          // 前置依赖的出度数目减去1
	}

	return true
}

func (taskGraph *TaskGraph) PrintGraph() {
	fmt.Println("-----------------------------------")
	for k, node := range taskGraph.graph {
		fmt.Println("任务名：", k)
		if node.Done {
			fmt.Println("是否完成：", "YES")
		} else {
			fmt.Println("是否完成：", "NO")
		}
		fmt.Print("（当前）依赖这些任务：")

		for taskName, v := range node.OutEdge {
			if !v {
				fmt.Print(" ", taskName, " ")
			}
		}
		// fmt.Println()

		fmt.Print("\t（当前）被这些任务依赖：")
		for taskName, v := range node.InEdge {
			if !v {
				fmt.Print(" ", taskName, " ")
			}
		}
		fmt.Println()
	}
	fmt.Println("-----------------------------------")
}
