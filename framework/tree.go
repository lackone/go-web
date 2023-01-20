package framework

import (
	"errors"
	"strings"
)

// 代表trie树
type Tree struct {
	root *Node `json:"root"` //根节点
}

// 代表节点
type Node struct {
	isLast   bool                `json:"is_last"` //表示这个节点是否为最终路由规则
	segment  string              `json:"segment"` //uri字符串中的某段
	handlers []ControllerHandler `json:"-"`       //处理的handler
	childs   []*Node             `json:"childs"`  //节点下的所有子节点
	parent   *Node               `json:"parent"`  //父给节点
}

func NewTree() *Tree {
	root := NewNode()
	return &Tree{root: root}
}

func NewNode() *Node {
	return &Node{
		isLast:   false,
		segment:  "",
		handlers: []ControllerHandler{},
		childs:   []*Node{},
	}
}

func (this *Tree) AddRouter(uri string, handlers []ControllerHandler) error {
	root := this.root
	uri = strings.TrimPrefix(uri, "/")

	if root.MatchNode(uri) != nil {
		return errors.New("route exists: " + uri)
	}

	segments := strings.Split(uri, "/")

	//遍历每一个段
	for index, segment := range segments {
		var objNode *Node //有匹配的子节点

		if !IsWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}

		isLast := index == len(segments)-1

		nodes := root.FilterChildNodes(segment)
		//如果有匹配的子节点
		if len(nodes) > 0 {
			for _, node := range nodes {
				if node.segment == segment {
					objNode = node
					break
				}
			}
		}

		if objNode == nil {
			cnode := NewNode()
			cnode.segment = segment
			if isLast {
				cnode.isLast = true
				cnode.handlers = handlers
			}
			//修改父节点指针
			cnode.parent = root
			root.childs = append(root.childs, cnode)
			objNode = cnode
		}

		root = objNode
	}
	return nil
}

func (this *Tree) FindNode(uri string) *Node {
	uri = strings.TrimPrefix(uri, "/")
	node := this.root.MatchNode(uri)
	if node == nil {
		return nil
	}
	return node
}

// 判断一个segment是否是通用，即以:开头
func IsWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

// 获取所有满足segment规则的子节点
func (this *Node) FilterChildNodes(segment string) []*Node {
	if len(this.childs) == 0 {
		return nil
	}
	//如果是通配符，则所有下一层子节点都满足条件
	if IsWildSegment(segment) {
		return this.childs
	}
	nodes := make([]*Node, 0, len(this.childs))

	//遍历子节点，获取满足规则的
	for _, v := range this.childs {
		if IsWildSegment(v.segment) {
			//如果子节点有通配符，则满足条件
			nodes = append(nodes, v)
		} else if v.segment == segment {
			//如果文本完全匹配，也满足条件
			nodes = append(nodes, v)
		}
	}
	return nodes
}

func (this *Node) MatchNode(uri string) *Node {
	//把uri分割成两部分
	segments := strings.SplitN(uri, "/", 2)

	segment := segments[0]
	if !IsWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}

	nodes := this.FilterChildNodes(segment)
	//如果没有找到符合规则的子节点，则直接返回
	if nodes == nil || len(nodes) == 0 {
		return nil
	}

	//如果只有最后一个segment，说明是最后的标记
	if len(segments) == 1 {
		for _, v := range nodes {
			if v.isLast {
				return v
			}
		}
		return nil
	}

	//如果有2个以上segment，递归每个子节点继续查找
	for _, v := range nodes {
		node := v.MatchNode(segments[1])
		if node != nil {
			return node
		}
	}

	return nil
}

// 解析uri中的参数
func (this *Node) ParseParamsFromEndNode(uri string) map[string]string {
	ret := map[string]string{}
	segments := strings.Split(uri, "/")
	len := len(segments)
	cur := this
	for i := len - 1; i >= 0; i-- {
		if cur.segment == "" {
			break
		}
		if IsWildSegment(cur.segment) {
			ret[cur.segment[1:]] = segments[i]
		}
		cur = cur.parent
	}
	return ret
}
