/* {{{ Copyright (c) Paul R. Tagliamonte <paultag@debian.org>, 2015
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE. }}} */

package dependency

import (
	"errors"
)

type topsortNetwork struct {
	Nodes map[string]*topsortNode
}

func sortSingleNodes(nodes []*topsortNode) ([]*topsortNode, error) {
	ret := make([]*topsortNode, 0)
	hasUnprunedNodes := false

	for _, node := range nodes {
		if node.Marked {
			continue /* Already output. */
		}

		hasUnprunedNodes = true

		/* Otherwise, let's see if we can prune it */
		if node.IsCanidate() {
			/* So, it has no deps and hasn't been marked; let's mark and
			 * output */
			node.Marked = true
			ret = append(ret, node)
		}
	}

	if hasUnprunedNodes && len(ret) == 0 {
		return nil, errors.New("Cycle detected :(")
	}

	return ret, nil
}

func sortNodes(nodes []*topsortNode) (ret []*topsortNode, err error) {
	for {
		generation, err := sortSingleNodes(nodes)
		if err != nil {
			return nil, err
		}
		if len(generation) == 0 {
			break
		}
		ret = append(ret, generation...)
	}
	return
}

func (tn *topsortNetwork) Sort() ([]*topsortNode, error) {
	nodes := make([]*topsortNode, 0)
	for _, v := range tn.Nodes {
		nodes = append(nodes, v)
	}
	return sortNodes(nodes)
}

func (tn *topsortNetwork) Get(name string) *topsortNode {
	return tn.Nodes[name]
}

func (tn *topsortNetwork) AddNode(name string) *topsortNode {
	node := topsortNode{
		Name:          name,
		InboundEdges:  make([]*topsortNode, 0),
		OutboundEdges: make([]*topsortNode, 0),
		Marked:        false,
	}

	tn.Nodes[name] = &node
	return &node
}

type topsortNode struct {
	Name          string
	OutboundEdges []*topsortNode
	InboundEdges  []*topsortNode
	Marked        bool
}

func (node *topsortNode) IsCanidate() bool {
	for _, edge := range node.InboundEdges {
		/* for each node, let's check if they're all marked */
		if !edge.Marked {
			return false
		}
	}
	return true
}

func (tn *topsortNetwork) CreateEdge(from string, to string) {
	/* Add the edge to `to`, as an "inbound" edge from `from`. */
	fromNode := tn.Get(from)
	toNode := tn.Get(to)

	if fromNode == nil || toNode == nil {
		return /* return err! */
	}

	toNode.InboundEdges = append(toNode.InboundEdges, fromNode)
	fromNode.OutboundEdges = append(fromNode.OutboundEdges, toNode)
}

func SortDependencies(els map[string][]Possibility) ([]string, error) {
	/* Right, so we have an incoming map of names -> other names via a
	 * Possibility; we need to construct the digraph and walk the DAG
	 * until we resolve the order. */

	/* First, let's create the network, empty. Then, we'll iterate over the
	 * keys to identify and allocate the nodes in the network. Finally, we'll
	 * take each Possi and construct an edge if the target is also within
	 * the network. Finally, we'll deconsuct that out by doing the actual
	 * topsort against the network in place. */

	network := topsortNetwork{Nodes: map[string]*topsortNode{}}

	/* So, let's create the known nodes here */
	for k, _ := range els {
		network.AddNode(k)
	}

	/* Now, let's establish links between the nodes */
	for name, possis := range els {
		for _, possi := range possis {
			network.CreateEdge(name, possi.Name)
		}
	}

	/* Now, we can offload that to the topsorter */

	order, err := network.Sort()
	if err != nil {
		return nil, err
	}

	keys := []string{}
	for _, el := range order {
		keys = append(keys, el.Name)
	}

	return keys, nil
}

// vim: foldmethod=marker
