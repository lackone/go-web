package framework

import (
	"errors"
	"strings"
)

// 代表trie树
type Tree struct {
	Root *Node `json:"root"` //根节点
}

// 代表节点
type Node struct {
	IsLast  bool              `json:"is_last"` //表示这个节点是否为最终路由规则
	Segment string            `json:"segment"` //uri字符串中的某段
	Handler ControllerHandler `json:"-"`       //处理的handler
	Childs  []*Node           `json:"childs"`  //节点下的所有子节点
}

func NewTree() *Tree {
	root := NewNode()
	return &Tree{Root: root}
}

func NewNode() *Node {
	return &Node{
		IsLast:  false,
		Segment: "",
		Childs:  []*Node{},
	}
}

func (this *Tree) AddRouter(uri string, handler ControllerHandler) error {
	root := this.Root
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
				if node.Segment == segment {
					objNode = node
					break
				}
			}
		}

		if objNode == nil {
			cnode := NewNode()
			cnode.Segment = segment
			if isLast {
				cnode.IsLast = true
				cnode.Handler = handler
			}
			root.Childs = append(root.Childs, cnode)
			objNode = cnode
		}

		root = objNode
	}
	return nil
}

func (this *Tree) FindHandler(uri string) ControllerHandler {
	uri = strings.TrimPrefix(uri, "/")
	node := this.Root.MatchNode(uri)
	if node == nil {
		return nil
	}
	return node.Handler
}

// 判断一个segment是否是通用，即以:开头
func IsWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

// 获取所有满足segment规则的子节点
func (this *Node) FilterChildNodes(segment string) []*Node {
	if len(this.Childs) == 0 {
		return nil
	}
	//如果是通配符，则所有下一层子节点都满足条件
	if IsWildSegment(segment) {
		return this.Childs
	}
	nodes := make([]*Node, 0, len(this.Childs))

	//遍历子节点，获取满足规则的
	for _, v := range this.Childs {
		if IsWildSegment(v.Segment) {
			//如果子节点有通配符，则满足条件
			nodes = append(nodes, v)
		} else if v.Segment == segment {
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
			if v.IsLast {
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
