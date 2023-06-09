package fsm

import (
	"context"
	"errors"
	"fmt"
)

var NoMoreTransitions = errors.New("no more transitions")
var InvalidEvent = errors.New("invalid event")

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
		return nil, NoMoreTransitions
	}

	transition, ok := m.CurrentNode.Transitions[event]
	if !ok {
		return nil, InvalidEvent
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

func NewStateMachine(initialNode *Node) *StateMachine {
	if initialNode == nil {
		return &StateMachine{}
	}

	return &StateMachine{
		initialNode: initialNode,
		CurrentNode: initialNode,
	}
}
