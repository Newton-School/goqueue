# Canvas Workflows

Phase 7 adds Celery-style workflow primitives for task composition:

- `Signature`: a reusable task invocation.
- `Chain`: ordered signatures that run one after another.
- `Group`: signatures dispatched in parallel and tracked as one unit.
- `Chord`: a group plus a callback that runs once the group succeeds.

## Signatures

```go
signature := goqueue.Signature{
	Name:  "email.send",
	Queue: "default",
	Args:  []any{"u_123"},
	Kwargs: map[string]any{
		"template": "welcome",
	},
}

result, err := canvas.ApplySignature(ctx, signature)
```

## Chains

```go
chainResult, err := canvas.ApplyChain(ctx, goqueue.Chain{
	Signatures: []goqueue.Signature{
		{Name: "email.prepare", Args: []any{"u_123"}},
		{Name: "email.send", Args: []any{"u_123"}},
	},
})
```

Workers advance chains only after the current task result has been persisted.
If the next-step enqueue fails, the current stream message is not acknowledged.

## Groups

```go
groupResult, err := canvas.ApplyGroup(ctx, goqueue.Group{
	Signatures: []goqueue.Signature{
		{Name: "email.send", Args: []any{"u_1"}},
		{Name: "email.send", Args: []any{"u_2"}},
	},
})
```

Group child tasks carry workflow metadata so workers can record terminal
progress after each child completes.

## Chords

```go
chordResult, err := canvas.ApplyChord(ctx, goqueue.Chord{
	Header: goqueue.Group{
		Signatures: []goqueue.Signature{
			{Name: "email.send", Args: []any{"u_1"}},
			{Name: "email.send", Args: []any{"u_2"}},
		},
	},
	Callback: goqueue.Signature{Name: "email.report"},
})
```

The callback is dispatched once when all header tasks succeed. If any header
task fails, the group is marked done but the callback is not dispatched.
