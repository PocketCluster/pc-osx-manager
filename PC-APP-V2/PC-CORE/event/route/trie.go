package route

import (
    "strings"
)

// node represents a struct of each node in the tree.
type node struct {
    children     []*node
    component    string
    methods      map[string]Handle
}

// addNode - adds a node to our tree. Will add multiple nodes if path can be broken up into multiple components.
// Those nodes will have no handler implemented and will fall through to the default handler.
func (n *node) addNode(method, path string, handler Handle) {
    var (
        components = strings.Split(path, "/")[1:]
        count = len(components)
    )

    for {
        aNode, component := n.traverse(components)
        if aNode.component == component && count == 1 { // update an existing node.
            aNode.methods[method] = handler
            return
        }
        newNode := node{component: component, methods: make(map[string]Handle)}

        if count == 1 { // this is the last component of the url resource, so it gets the handler.
            newNode.methods[method] = handler
        }
        aNode.children = append(aNode.children, &newNode)
        count--
        if count == 0 {
            break
        }
    }
}

// traverse moves along the tree adding named params as it comes and across them.
// Returns the node and component found.
func (n *node) traverse(components []string) (*node, string) {
    component := components[0]
    if len(n.children) > 0 { // no children, then bail out.
        for _, child := range n.children {
            if component == child.component {
                next := components[1:]
                if len(next) > 0 { // http://xkcd.com/1270/
                    return child.traverse(next) // tail recursion is it's own reward.
                }
                return child, component
            }
        }
    }
    return n, component
}
