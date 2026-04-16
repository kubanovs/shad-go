package main

type Stack struct {
	slice []int
	size  int
}

func (stack *Stack) push(el int) {
	stack.slice = append(stack.slice, el)
	stack.size++
}

func (stack *Stack) pop() int {
	if stack.size > 0 {
		last := stack.slice[len(stack.slice)-1]
		stack.slice = stack.slice[:stack.size-1]
		stack.size--
		return last
	} else {
		panic("aye")
	}
}

func (stack *Stack) peek() int {
	if stack.size == 0 {
		panic("stack is empty")
	}
	return stack.slice[stack.size-1]
}
