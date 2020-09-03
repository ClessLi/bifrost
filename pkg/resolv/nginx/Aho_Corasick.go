package nginx

// acNode, AC 自动机节点对象
type acNode struct {
	data         rune
	children     []*acNode
	isEndingChar bool
	fail         *acNode
	length       int
}

// newACNode, 实例化 AC 自动机节点对象指针函数
// 参数:
//     data: 字符rune对象
//     isEnd: 是否为字符串结束字符
// 返回值:
//     AC 自动机节点对象实例指针
func newACNode(data rune, isEnd bool) *acNode {
	return &acNode{
		data:         data,
		isEndingChar: isEnd,
	}
}

// addChild, AC 自动机节点对象，添加子节点方法
// 参数:
//     data: 字符rune对象
//     isEnd: 是否为字符串结束字符
// 返回值:
//     AC 自动机子节点对象指针
func (acn *acNode) addChild(data rune, isEnd bool) *acNode {
	d := int(data)
	if acn.children == nil {
		acn.children = make([]*acNode, d)
		acn.children = append(acn.children, newACNode(data, isEnd))
	} else {
		n := len(acn.children)

		if n < d+1 {
			for i := 0; i < d+1-n; i++ {
				acn.children = append(acn.children, nil)
			}
		}

		if acn.children[d] == nil {
			acn.children[d] = newACNode(data, isEnd)
		}
	}
	acn.children[d].length = acn.length + 1
	return acn.children[d]
}

// AC, AC 自动机对象
type AC struct {
	root  acNode
	store []string
}

// NewAC, 实例化 AC 自动机对象指针函数
// 返回值:
//     AC 自动机对象实例指针
func NewAC() *AC {
	root := newACNode(-1, false)
	root.length = 0
	return &AC{root: *root}
}

// Insert, AC 自动机对象，字符串插入的方法
// 参数:
//     strings: 字符串集合
func (ac *AC) Insert(strings ...string) {
	for _, str := range strings {
		//fmt.Println(str)
		if !ac.Find(str) {
			p := &ac.root
			for _, data := range []rune(str) {
				//fmt.Println("rune:", data)
				pc := p.addChild(data, false)
				if p == &ac.root {
					pc.fail = &ac.root
				} else {
					q := p.fail
					for q != nil {
						var qc *acNode = nil
						if len(q.children) > int(pc.data) {
							//fmt.Println(len(q.children), int(pc.data), pc.data)
							qc = q.children[int(pc.data)]
						}
						if qc != nil {
							pc.fail = qc
							break
						}
						q = q.fail
					}
					if q == nil {
						pc.fail = &ac.root
					}
				}
				p = pc
			}
			p.isEndingChar = true
			ac.store = append(ac.store, str)
		}
	}
}

// Find, AC 自动机对象，查询字符串是否存在的方法
// 参数:
//     str: 需查询的字符串
// 返回值:
//     true: 存在该字符串; false: 不存在字符串
func (ac *AC) Find(str string) bool {
	p := &ac.root
	if p.children == nil {
		return false
	}
	for _, data := range []rune(str) {
		if int(data) >= len(p.children) || p.children[data] == nil {
			return false
		}
		p = p.children[data]
	}
	if !p.isEndingChar {
		return false
	} else {
		return true
	}
}

//func (ac *AC) Match(str string) {
//	if ac.root.children == nil {
//		return
//	}
//	p := &ac.root
//	chars := []rune(str)
//	n := len(chars)
//	for i := 0; i < n; i++ {
//		idx := int(chars[i])
//		if len(p.children) > idx {
//			if p.children[idx] == nil && p != &ac.root {
//				p = p.fail
//			}
//
//			p = p.children[idx]
//		} else {
//			p = nil
//		}
//
//		if p == nil {
//			p = &ac.root
//		}
//		tmp := p
//		for tmp != &ac.root {
//			if tmp.isEndingChar {
//				pos := i - tmp.length + 1
//				fmt.Println("匹配起始下标", pos, "；长度", tmp.length)
//				pos++
//				for j := 0; j < tmp.length-2; j++ {
//					chars[pos+j] = rune(42)
//				}
//			}
//			tmp = tmp.fail
//		}
//	}
//	fmt.Println(string(chars))
//}
