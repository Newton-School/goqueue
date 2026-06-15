# Workflow Metadata

Workflow metadata is attached to task messages so workers can advance canvas
primitives without Redis-specific task payloads.

## Reserved Keys

The SDK owns these metadata keys:

- `goqueue.workflow.kind`
- `goqueue.workflow.chain_id`
- `goqueue.workflow.chain_step`
- `goqueue.workflow.group_id`
- `goqueue.workflow.group_index`
- `goqueue.workflow.chord_id`
- `goqueue.workflow.chord_callback`

User metadata is preserved unless it conflicts with a reserved workflow key.
When a conflict exists, goqueue writes the authoritative workflow value.

## Compatibility

Backends and workers should treat unknown metadata keys as user-owned and leave
them unchanged. New workflow metadata keys should be additive so existing tasks
remain readable.
