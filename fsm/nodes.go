package fsm

import (
	"context"

	"go.uber.org/zap"
)

var (
	l, _ = zap.NewProduction()
)

var CopyPost = Node{
	State: "CopyPost",
	Transitions: map[Event]*Transition{
		"HasAttachments": &Transition{
			Node: &HasAttachments,
			Action: func(ctx context.Context) error {
				return nil
			},
		},
		"NoAttachments": &Transition{
			Node: &NoAttachments,
			Action: func(ctx context.Context) error {
				return nil
			},
		},
	},
}

var NoteFound = Node{
	State: "NoteFound",
	Transitions: map[Event]*Transition{
		"CopyPost": &Transition{
			Node: &CopyPost,
			Action: func(ctx context.Context) error {
				//TODO: Here's the place to do validation
				// I had an idea about parsing frontmatter into go structs,
				// and then writing them out to a NullWriter via protobuf encoding.
				// this would do schema checking which is nice
				l := l.Named("NoteFound")
				note := ctx.Value("note").(string)
				l.Info("creating post from note", zap.String("filename", note))
				// scan for attachments here
				// if len() attachments > 0
				//
				return nil
			},
		},
	},
}

var NoAttachments = Node{
	State: "NoAttachments",
	Transitions: map[Event]*Transition{
		"Terminate": &Transition{
			Node: &Terminate,
			Action: func(ctx context.Context) error {
				return nil
			},
		},
	},
}

var HasAttachments = Node{
	State: "HasAttachments",
	Transitions: map[Event]*Transition{
		"CopyAttachments": &Transition{
			Node: &CopyAttachments,
			Action: func(ctx context.Context) error {
				return nil
			},
		},
	},
}

var CopyAttachments = Node{
	State: "CopyAttachments",
	Transitions: map[Event]*Transition{
		"SanitizeLinks": &Transition{
			Node: &SanitizeLinks,
			Action: func(ctx context.Context) error {
				return nil
			},
		},
	},
}

var SanitizeLinks = Node{
	State: "SanitizeLinks",
	Transitions: map[Event]*Transition{
		"Terminate": &Transition{
			Node: &Terminate,
			Action: func(ctx context.Context) error {
				return nil
			},
		},
	},
}

var Terminate = Node{
	State:       "Terminate",
	Transitions: map[Event]*Transition{},
}
