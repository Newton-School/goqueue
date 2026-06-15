# Worker Advancement Order

Workers advance workflows after task result persistence and before queue
acknowledgement.

The order is:

1. Decode the task message.
2. Move the task to a running state.
3. Execute the handler.
4. Persist the final task state and result.
5. Advance workflow state or record group progress.
6. Acknowledge the queue message.

This order keeps workflow state aligned with durable task results. If workflow
advancement fails, the message is not acknowledged and can be retried by the
backend's reliability path.
