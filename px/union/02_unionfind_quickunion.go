package main

/*
	查找v所属的集合（根节点）
*/
func findQu(v int) int {
	for v != parents[v] {
		v = parents[v]
	}

	return v
}

/*
	合并v1、v2所在的集合
*/
func unionQu(v1, v2 int) {
	p1 := findQu(v1)
	p2 := findQu(v2)
	if p1 == p2 {
		return
	}

	parents[v1] = p2
}

/*
	检查v1、v2是否属于同一个集合
*/
func isSameQu(v1, v2 int) bool {
	return findQu(v1) == findQu(v2)
}
