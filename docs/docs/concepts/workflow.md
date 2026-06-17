---
title: Workflow
---

Workflow tools compose multiple task signatures.

## Signature

`workflow.Signature` is a reusable task declaration:

- task name
- queue
- args/kwargs
- metadata
- timing / priority / retry policy

## Build flow

- `Canvas` created from app: `app.NewCanvas(...)`
- methods:
  - `ApplySignature`
  - `ApplyChain`
  - `ApplyGroup`
  - `ApplyChord`

## Chain

Runs signatures in order. Each step runs only after previous succeeds.

## Group

Runs signatures in parallel. Result returns all created task IDs.

## Chord

Runs a header group first. After all header tasks succeed, callback runs.

## Example

```go
sig1 := workflow.Signature{Name: "parse_report"}
sig2 := workflow.Signature{Name: "upload_report"}

result, err := canvas.ApplyChain(ctx, workflow.Chain{
    Signatures: []workflow.Signature{sig1, sig2},
})
_ = result.WorkflowID
```

## Result objects

- `ChainResult{WorkflowID, FirstTask}`
- `GroupResult{GroupID, TaskIDs}`
- `ChordResult{GroupID, TaskIDs}`

Workflow metadata keys include:

- `goqueue.WorkflowKindChain`, `goqueue.WorkflowKindGroup`, `goqueue.WorkflowKindChord`
- chain IDs, step index, and chord callback metadata
