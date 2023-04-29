package obspipeline

import (
	"context"
	"fmt"
)

type State string
type Event string

type Action func(ctx context.Context) error

type Node struct {
	State
	Transitions map[Event]*Transition
}

type Transition struct {
	*Node
	Action
}

type StateMachine struct {
	initialNode *Node
	CurrentNode *Node
}

func (m *StateMachine) getCurrentNode() *Node {
	return m.CurrentNode
}

func (m *StateMachine) getNextNode(event Event) (*Node, error) {
	if m.CurrentNode == nil {
		return nil, fmt.Errorf("nowhere to go anymore!\n")
	}

	transition, ok := m.CurrentNode.Transitions[event]
	if !ok {
		return nil, fmt.Errorf("invalid event: %v", event)
	}

	return transition.Node, nil
}

func (m *StateMachine) Transition(ctx context.Context, event Event) (*Node, error) {
	node, err := m.getNextNode(event) // gets current node
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	err = m.CurrentNode.Transitions[event].Action(ctx)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	m.CurrentNode = node

	return m.CurrentNode, nil
}
