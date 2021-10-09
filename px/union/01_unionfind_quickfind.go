package main

/*
	查找v所属的集合（根节点）
*/
func findQf(v int) int {
	return parents[v]
}

/*
	合并v1、v2所在的集合
*/
func unionQf(v1, v2 int) {
	p1 := findQf(v1)
	p2 := findQf(v2)
	if p1 == p2 {
		return
	}

	for i, parent := range parents {
		if parent == p1 {
			parents[i] = p2
		}
	}
}

/*
	检查v1、v2是否属于同一个集合
*/
func isSameQf(v1, v2 int) bool {
	return findQf(v1) == findQf(v2)
}
