//go:build !solution

package hogwarts

func learn(current string, graph map[string][]string, visited map[string]bool, result *[]string, visitedAll map[string]bool) {

	childs, ok := graph[current]

	if !ok && !visitedAll[current] {
		*result = append(*result, current)
		visitedAll[current] = true
		return
	}

	for _, child := range childs {

		if visited[child] {
			panic("it's circle dependency!")
		}

		visited[child] = true
		learn(child, graph, visited, result, visitedAll)
		if !visitedAll[child] {
			visitedAll[child] = true
			*result = append(*result, child)
		}

		visited[child] = false
	}
}

func GetCourseList(prereqs map[string][]string) []string {

	visited := make(map[string]bool)
	visitedAll := make(map[string]bool)
	var result []string

	for lesson, _ := range prereqs {
		visited[lesson] = true
		learn(lesson, prereqs, visited, &result, visitedAll)
		if !visitedAll[lesson] {
			visitedAll[lesson] = true
			result = append(result, lesson)
		}
		visited[lesson] = false
	}

	return result
}
