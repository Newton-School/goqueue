package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Newton-School/goqueue"
)

const (
	defaultNamespace              = "goqueue"
	defaultTaskQueue              = "default"
	defaultDeadLetters            = 25
	defaultDeadLetterIDsSeparator = ","
)

func main() {
	rootCmd := os.Args
	if len(rootCmd) < 2 {
		fmt.Fprintln(os.Stderr, "goqueue: expected command. try `goqueue help`")
		os.Exit(1)
	}

	switch rootCmd[1] {
	case "help", "-h", "--help":
		printUsage()
		return
	case "inspect":
		runInspect(rootCmd[2:])
	case "control":
		runControl(rootCmd[2:])
	default:
		fmt.Fprintf(os.Stderr, "goqueue: unknown command %q\n", rootCmd[1])
		printUsage()
		os.Exit(1)
	}
}

func runInspect(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "goqueue inspect: expected subcommand")
		printInspectUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "ping":
		runInspectPing(args[1:])
	case "stats":
		runInspectStats(args[1:])
	case "deadletters":
		runInspectDeadLetters(args[1:])
	case "task":
		runInspectTask(args[1:])
	case "state":
		runInspectTaskState(args[1:])
	case "result":
		runInspectTaskResult(args[1:])
	case "forget-result":
		runInspectForgetResult(args[1:])
	case "-h", "--help", "help":
		printInspectUsage()
	default:
		fmt.Fprintf(os.Stderr, "goqueue inspect: unknown subcommand %q\n", args[0])
		printInspectUsage()
		os.Exit(1)
	}
}

func runControl(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "goqueue control: expected subcommand")
		printControlUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "retry-task":
		runControlRetryTask(args[1:])
	case "revoke-task":
		runControlRevokeTask(args[1:])
	case "replay-dead-letter":
		runControlReplayDeadLetter(args[1:])
	case "delete-dead-letter":
		runControlDeleteDeadLetter(args[1:])
	case "purge-queue":
		runControlPurgeQueue(args[1:])
	case "-h", "--help", "help":
		printControlUsage()
	default:
		fmt.Fprintf(os.Stderr, "goqueue control: unknown subcommand %q\n", args[0])
		printControlUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("goqueue CLI")
	fmt.Println("Usage:")
	fmt.Println("  goqueue inspect <command> [options]")
	fmt.Println("  goqueue control <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  inspect ping")
	fmt.Println("  inspect stats --queue <queue>")
	fmt.Println("  inspect deadletters --queue <queue> [--count <n>]")
	fmt.Println("  inspect task --id <task_id> [--json]")
	fmt.Println("  inspect state --id <task_id> [--json]")
	fmt.Println("  inspect result --id <task_id> [--json]")
	fmt.Println("  inspect forget-result --id <task_id>")
	fmt.Println("  control retry-task --id <task_id> [--queue <queue>] [--scheduled-at <RFC3339>] [--countdown <duration>] [--preserve-attempt] [--clear-state] [--clear-result] [--json]")
	fmt.Println("  control revoke-task --id <task_id> [--reason <text>]")
	fmt.Println("  control replay-dead-letter --queue <queue> --stream-id <id> [--destination-queue <queue>] [--delete-source]")
	fmt.Println("  control delete-dead-letter --queue <queue> --stream-id <id>[,<id>...] [--json]")
	fmt.Println("  control purge-queue --queue <queue> --yes [--delete-messages] [--delete-states] [--delete-results]")
	fmt.Println()
	fmt.Println("Global flags (per command):")
	fmt.Println("  --redis-url string")
	fmt.Println("  --namespace string")
	fmt.Println("  --json")
}

func printInspectUsage() {
	fmt.Println("goqueue inspect")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  goqueue inspect ping")
	fmt.Println("  goqueue inspect task --id <task_id>")
	fmt.Println("  goqueue inspect state --id <task_id>")
	fmt.Println("  goqueue inspect result --id <task_id>")
	fmt.Println("  goqueue inspect forget-result --id <task_id>")
	fmt.Println("  goqueue inspect stats --queue <queue>")
	fmt.Println("  goqueue inspect deadletters --queue <queue> --count <n>")
}

func printControlUsage() {
	fmt.Println("goqueue control")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  goqueue control retry-task --id <task_id> [--json]")
	fmt.Println("  goqueue control revoke-task --id <task_id>")
	fmt.Println("  goqueue control replay-dead-letter --queue <queue> --stream-id <id>")
	fmt.Println("  goqueue control delete-dead-letter --queue <queue> --stream-id <id>[,<id>...]")
	fmt.Println("  goqueue control purge-queue --queue <queue> --yes")
}

func exitOnError(err error, context string) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "goqueue: %s: %v\n", context, err)
	os.Exit(1)
}

func printJSON(value any) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		exitOnError(err, "render json")
	}
	fmt.Println(string(data))
}

func newInspector(redURL, namespace string) (*goqueue.Inspector, error) {
	app, err := goqueue.New(
		goqueue.WithRedisURL(redURL),
		goqueue.WithNamespace(namespace),
	)
	if err != nil {
		return nil, err
	}

	return app.NewInspector()
}

func newAdmin(redURL, namespace string) (*goqueue.Admin, error) {
	app, err := goqueue.New(
		goqueue.WithRedisURL(redURL),
		goqueue.WithNamespace(namespace),
	)
	if err != nil {
		return nil, err
	}

	return app.NewAdmin()
}

func splitStreamIDs(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	parts := strings.Split(raw, defaultDeadLetterIDsSeparator)
	ids := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		ids = append(ids, trimmed)
	}

	if len(ids) == 0 {
		return nil, nil
	}

	return ids, nil
}

func resolveRedisURL() string {
	if value := os.Getenv("GOQUEUE_REDIS_URL"); value != "" {
		return value
	}
	return "redis://localhost:6379/0"
}

func resolveNamespace(defaultNamespace string) string {
	if value := os.Getenv("GOQUEUE_NAMESPACE"); value != "" {
		return value
	}
	return defaultNamespace
}

func resolveQueue(defaultQueue string) string {
	if value := os.Getenv("GOQUEUE_DEFAULT_QUEUE"); value != "" {
		return value
	}
	return defaultQueue
}

func runInspectPing(args []string) {
	fs := flag.NewFlagSet("inspect ping", flag.ExitOnError)
	redisURL := fs.String("redis-url", resolveRedisURL(), "Redis URL used by goqueue")
	namespace := fs.String("namespace", resolveNamespace(defaultNamespace), "Redis namespace used by goqueue")
	outputJSON := fs.Bool("json", false, "Print output as JSON")
	_ = fs.Parse(args)

	inspector, err := newInspector(*redisURL, *namespace)
	exitOnError(err, "inspect ping")

	err = inspector.Ping(ctx())
	exitOnError(err, "inspect ping")

	if *outputJSON {
		fmt.Println(`{"ok":true}`)
		return
	}

	fmt.Println("status: ok")
}

func runInspectStats(args []string) {
	fs := flag.NewFlagSet("inspect stats", flag.ExitOnError)
	redisURL := fs.String("redis-url", resolveRedisURL(), "Redis URL used by goqueue")
	namespace := fs.String("namespace", resolveNamespace(defaultNamespace), "Redis namespace used by goqueue")
	queue := fs.String("queue", resolveQueue(defaultTaskQueue), "Queue name")
	outputJSON := fs.Bool("json", false, "Print output as JSON")
	_ = fs.Parse(args)

	inspector, err := newInspector(*redisURL, *namespace)
	exitOnError(err, "inspect stats")

	stats, err := inspector.QueueStats(ctx(), goqueue.QueueName(*queue))
	exitOnError(err, "inspect stats")

	if *outputJSON {
		printJSON(map[string]any{
			"queue":             string(stats.Queue),
			"ready_count":       stats.ReadyCount,
			"scheduled_count":   stats.ScheduledCount,
			"dead_letter_count": stats.DeadLetterCount,
		})
		return
	}

	fmt.Printf("queue=%s\nready=%d\nscheduled=%d\ndead_letter=%d\n", stats.Queue, stats.ReadyCount, stats.ScheduledCount, stats.DeadLetterCount)
}

func runInspectDeadLetters(args []string) {
	fs := flag.NewFlagSet("inspect deadletters", flag.ExitOnError)
	redisURL := fs.String("redis-url", resolveRedisURL(), "Redis URL used by goqueue")
	namespace := fs.String("namespace", resolveNamespace(defaultNamespace), "Redis namespace used by goqueue")
	queue := fs.String("queue", resolveQueue(defaultTaskQueue), "Queue name")
	count := fs.Int("count", defaultDeadLetters, "Number of dead letters to read")
	outputJSON := fs.Bool("json", false, "Print output as JSON")
	_ = fs.Parse(args)

	inspector, err := newInspector(*redisURL, *namespace)
	exitOnError(err, "inspect deadletters")

	records, err := inspector.ReadDeadLetters(ctx(), goqueue.QueueName(*queue), int64(*count))
	exitOnError(err, "inspect deadletters")

	if *outputJSON {
		printJSON(records)
		return
	}

	fmt.Printf("queue=%s\ndead_letters=%d\n", *queue, len(records))
	for idx, record := range records {
		fmt.Printf("#%d stream=%s reason=%s error=%s source=%s\n", idx, record.StreamID, record.Reason, record.Error, record.SourceStreamID)
	}
}

func runInspectTask(args []string) {
	fs := flag.NewFlagSet("inspect task", flag.ExitOnError)
	redisURL := fs.String("redis-url", resolveRedisURL(), "Redis URL used by goqueue")
	namespace := fs.String("namespace", resolveNamespace(defaultNamespace), "Redis namespace used by goqueue")
	taskID := fs.String("id", "", "Task ID")
	outputJSON := fs.Bool("json", false, "Print output as JSON")
	_ = fs.Parse(args)

	if *taskID == "" {
		exitOnError(fmt.Errorf("inspect task: --id required"), "inspect task")
	}

	inspector, err := newInspector(*redisURL, *namespace)
	exitOnError(err, "inspect task")

	snapshot, err := inspector.TaskSnapshot(ctx(), goqueue.TaskID(*taskID))
	exitOnError(err, "inspect task")

	if *outputJSON {
		printJSON(snapshot)
		return
	}

	fmt.Printf("task_id=%s\n", snapshot.TaskID)
	fmt.Printf("state_found=%t\n", snapshot.StateFound)
	fmt.Printf("state=%s\n", snapshot.State.State)
	fmt.Printf("state_error=%s\n", snapshot.State.Error)
	if snapshot.ResultFound {
		fmt.Printf("result_state=%s\n", snapshot.Result.State)
		fmt.Printf("result_error=%s\n", snapshot.Result.Error)
	}
}

func runInspectTaskState(args []string) {
	fs := flag.NewFlagSet("inspect state", flag.ExitOnError)
	redisURL := fs.String("redis-url", resolveRedisURL(), "Redis URL used by goqueue")
	namespace := fs.String("namespace", resolveNamespace(defaultNamespace), "Redis namespace used by goqueue")
	taskID := fs.String("id", "", "Task ID")
	outputJSON := fs.Bool("json", false, "Print output as JSON")
	_ = fs.Parse(args)

	if *taskID == "" {
		exitOnError(fmt.Errorf("inspect state: --id required"), "inspect state")
	}

	inspector, err := newInspector(*redisURL, *namespace)
	exitOnError(err, "inspect state")

	state, err := inspector.TaskState(ctx(), goqueue.TaskID(*taskID))
	exitOnError(err, "inspect state")

	if *outputJSON {
		printJSON(state)
		return
	}

	fmt.Printf("task_id=%s\nstate=%s\nerror=%s\nupdated_at=%s\n", state.TaskID, state.State, state.Error, state.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
}

func runInspectTaskResult(args []string) {
	fs := flag.NewFlagSet("inspect result", flag.ExitOnError)
	redisURL := fs.String("redis-url", resolveRedisURL(), "Redis URL used by goqueue")
	namespace := fs.String("namespace", resolveNamespace(defaultNamespace), "Redis namespace used by goqueue")
	taskID := fs.String("id", "", "Task ID")
	outputJSON := fs.Bool("json", false, "Print output as JSON")
	_ = fs.Parse(args)

	if *taskID == "" {
		exitOnError(fmt.Errorf("inspect result: --id required"), "inspect result")
	}

	inspector, err := newInspector(*redisURL, *namespace)
	exitOnError(err, "inspect result")

	result, err := inspector.TaskResult(ctx(), goqueue.TaskID(*taskID))
	exitOnError(err, "inspect result")

	if *outputJSON {
		printJSON(result)
		return
	}

	fmt.Printf("task_id=%s\nstate=%s\nerror=%s\n", result.TaskID, result.Result.State, result.Result.Error)
}

func runInspectForgetResult(args []string) {
	fs := flag.NewFlagSet("inspect forget-result", flag.ExitOnError)
	redisURL := fs.String("redis-url", resolveRedisURL(), "Redis URL used by goqueue")
	namespace := fs.String("namespace", resolveNamespace(defaultNamespace), "Redis namespace used by goqueue")
	taskID := fs.String("id", "", "Task ID")
	_ = fs.Parse(args)

	if *taskID == "" {
		exitOnError(fmt.Errorf("inspect forget-result: --id required"), "inspect forget-result")
	}

	inspector, err := newInspector(*redisURL, *namespace)
	exitOnError(err, "inspect forget-result")

	err = inspector.ForgetTaskResult(ctx(), goqueue.TaskID(*taskID))
	exitOnError(err, "inspect forget-result")
	fmt.Printf("task_id=%s result_cleared=true\n", *taskID)
}

func runControlRetryTask(args []string) {
	fs := flag.NewFlagSet("control retry-task", flag.ExitOnError)
	redisURL := fs.String("redis-url", resolveRedisURL(), "Redis URL used by goqueue")
	namespace := fs.String("namespace", resolveNamespace(defaultNamespace), "Redis namespace used by goqueue")
	taskID := fs.String("id", "", "Task ID")
	queue := fs.String("queue", "", "Target retry queue")
	scheduledAt := fs.String("scheduled-at", "", "Retry ETA in RFC3339")
	countdown := fs.Duration("countdown", 0, "Retry delay")
	preserveAttempt := fs.Bool("preserve-attempt", false, "Preserve original retry attempt")
	clearState := fs.Bool("clear-state", false, "Clear task state before retry")
	clearResult := fs.Bool("clear-result", false, "Clear task result before retry")
	outputJSON := fs.Bool("json", false, "Print output as JSON")
	_ = fs.Parse(args)

	if *taskID == "" {
		exitOnError(fmt.Errorf("control retry-task: --id required"), "control retry-task")
	}

	admin, err := newAdmin(*redisURL, *namespace)
	exitOnError(err, "control retry-task")

	var eta time.Time
	if *scheduledAt != "" {
		parsed, parseErr := time.Parse(time.RFC3339, *scheduledAt)
		if parseErr != nil {
			exitOnError(parseErr, "control retry-task: invalid --scheduled-at")
		}
		eta = parsed
	}

	result, err := admin.RetryTask(ctx(), goqueue.TaskID(*taskID), goqueue.RetryTaskOptions{
		Queue:           goqueue.QueueName(*queue),
		ScheduledAt:     eta,
		CountDown:       *countdown,
		PreserveAttempt: *preserveAttempt,
		ClearState:      *clearState,
		ClearResult:     *clearResult,
		Now:             time.Now,
	})
	exitOnError(err, "control retry-task")

	if *outputJSON {
		printJSON(result)
		return
	}

	streamID := result.EnqueueResult.StreamID
	if result.EnqueueResult.StreamID == "" {
		streamID = ""
	}
	fmt.Printf("task_id=%s\n", result.TaskID)
	fmt.Printf("queue=%s\n", result.Queue)
	fmt.Printf("original_queue=%s\n", result.OriginalQueue)
	fmt.Printf("attempt=%d\n", result.Attempt)
	fmt.Printf("scheduled_at=%s\n", result.ScheduledAt.Format(time.RFC3339))
	fmt.Printf("stream_id=%s\n", streamID)
	fmt.Printf("scheduled=%t\n", result.EnqueueResult.Scheduled)

	if result.EnqueueResult.TaskID == "" {
		fmt.Printf("enqueue_task_id=%s\n", *taskID)
	} else {
		fmt.Printf("enqueue_task_id=%s\n", result.EnqueueResult.TaskID)
	}
}

func runControlRevokeTask(args []string) {
	fs := flag.NewFlagSet("control revoke-task", flag.ExitOnError)
	redisURL := fs.String("redis-url", resolveRedisURL(), "Redis URL used by goqueue")
	namespace := fs.String("namespace", resolveNamespace(defaultNamespace), "Redis namespace used by goqueue")
	taskID := fs.String("id", "", "Task ID")
	reason := fs.String("reason", "operator request", "Revoke reason")
	outputJSON := fs.Bool("json", false, "Print output as JSON")
	_ = fs.Parse(args)

	if *taskID == "" {
		exitOnError(fmt.Errorf("control revoke-task: --id required"), "control revoke-task")
	}

	admin, err := newAdmin(*redisURL, *namespace)
	exitOnError(err, "control revoke-task")

	result, err := admin.RevokeTask(ctx(), goqueue.TaskID(*taskID), *reason)
	exitOnError(err, "control revoke-task")

	if *outputJSON {
		printJSON(result)
		return
	}

	fmt.Printf("task_id=%s\n", result.TaskID)
	fmt.Printf("state=%s\n", result.State)
	fmt.Printf("reason=%s\n", *reason)
}

func runControlReplayDeadLetter(args []string) {
	fs := flag.NewFlagSet("control replay-dead-letter", flag.ExitOnError)
	redisURL := fs.String("redis-url", resolveRedisURL(), "Redis URL used by goqueue")
	namespace := fs.String("namespace", resolveNamespace(defaultNamespace), "Redis namespace used by goqueue")
	queue := fs.String("queue", resolveQueue(defaultTaskQueue), "Queue name")
	streamID := fs.String("stream-id", "", "Dead-letter stream ID")
	destinationQueue := fs.String("destination-queue", "", "Destination queue")
	deleteSource := fs.Bool("delete-source", false, "Delete source dead-letter record")
	outputJSON := fs.Bool("json", false, "Print output as JSON")
	_ = fs.Parse(args)

	if *queue == "" {
		exitOnError(fmt.Errorf("control replay-dead-letter: --queue required"), "control replay-dead-letter")
	}
	if *streamID == "" {
		exitOnError(fmt.Errorf("control replay-dead-letter: --stream-id required"), "control replay-dead-letter")
	}

	admin, err := newAdmin(*redisURL, *namespace)
	exitOnError(err, "control replay-dead-letter")

	result, err := admin.ReplayDeadLetter(ctx(), goqueue.QueueName(*queue), *streamID, goqueue.ReplayDeadLetterOptions{
		DestinationQueue: goqueue.QueueName(*destinationQueue),
		DeleteSource:     *deleteSource,
	})
	exitOnError(err, "control replay-dead-letter")

	if *outputJSON {
		printJSON(result)
		return
	}

	fmt.Printf("stream_id=%s\n", result.StreamID)
	fmt.Printf("queue=%s\n", result.Queue)
	fmt.Printf("destination=%s\n", result.Destination)
	fmt.Printf("source_deleted=%t\n", result.SourceDeleted)
	fmt.Printf("enqueue_task_id=%s\n", result.EnqueueResult.TaskID)
}

func runControlDeleteDeadLetter(args []string) {
	fs := flag.NewFlagSet("control delete-dead-letter", flag.ExitOnError)
	redisURL := fs.String("redis-url", resolveRedisURL(), "Redis URL used by goqueue")
	namespace := fs.String("namespace", resolveNamespace(defaultNamespace), "Redis namespace used by goqueue")
	queue := fs.String("queue", resolveQueue(defaultTaskQueue), "Queue name")
	streamIDsRaw := fs.String("stream-id", "", "Comma separated stream ids")
	outputJSON := fs.Bool("json", false, "Print output as JSON")
	_ = fs.Parse(args)

	if *queue == "" {
		exitOnError(fmt.Errorf("control delete-dead-letter: --queue required"), "control delete-dead-letter")
	}

	ids, err := splitStreamIDs(*streamIDsRaw)
	if err != nil {
		exitOnError(err, "control delete-dead-letter")
	}
	if len(ids) == 0 {
		exitOnError(fmt.Errorf("control delete-dead-letter: at least one --stream-id required"), "control delete-dead-letter")
	}

	admin, err := newAdmin(*redisURL, *namespace)
	exitOnError(err, "control delete-dead-letter")

	result, err := admin.DeleteDeadLetters(ctx(), goqueue.QueueName(*queue), ids...)
	exitOnError(err, "control delete-dead-letter")

	if *outputJSON {
		printJSON(result)
		return
	}

	fmt.Printf("queue=%s\n", result.Queue)
	fmt.Printf("deleted=%d\n", result.Deleted)
	fmt.Printf("stream_ids=%s\n", strings.Join(ids, defaultDeadLetterIDsSeparator))
}

func runControlPurgeQueue(args []string) {
	fs := flag.NewFlagSet("control purge-queue", flag.ExitOnError)
	redisURL := fs.String("redis-url", resolveRedisURL(), "Redis URL used by goqueue")
	namespace := fs.String("namespace", resolveNamespace(defaultNamespace), "Redis namespace used by goqueue")
	queue := fs.String("queue", resolveQueue(defaultTaskQueue), "Queue name")
	deleteMessage := fs.Bool("delete-messages", false, "Delete persisted task messages")
	deleteStates := fs.Bool("delete-states", false, "Delete persisted task states")
	deleteResults := fs.Bool("delete-results", false, "Delete persisted task results")
	confirmed := fs.Bool("yes", false, "Confirm destructive queue purge")
	outputJSON := fs.Bool("json", false, "Print output as JSON")
	_ = fs.Parse(args)

	if err := validatePurgeConfirmation(*confirmed); err != nil {
		exitOnError(err, "control purge-queue")
	}

	admin, err := newAdmin(*redisURL, *namespace)
	exitOnError(err, "control purge-queue")

	result, err := admin.PurgeQueue(ctx(), goqueue.PurgeQueueOptions{
		Queue:         goqueue.QueueName(*queue),
		DeleteMessage: *deleteMessage,
		DeleteState:   *deleteStates,
		DeleteResult:  *deleteResults,
	})
	exitOnError(err, "control purge-queue")

	if *outputJSON {
		printJSON(result)
		return
	}

	fmt.Printf("queue=%s\n", result.Queue)
	fmt.Printf("ready_stream_deleted=%d\n", result.ReadyStream)
	fmt.Printf("scheduled_set_deleted=%d\n", result.ScheduledSet)
	fmt.Printf("dead_letter_deleted=%d\n", result.DeadLetterStream)
	fmt.Printf("task_messages_deleted=%d\n", result.TaskMessages)
	fmt.Printf("task_states_deleted=%d\n", result.TaskStates)
	fmt.Printf("task_results_deleted=%d\n", result.TaskResults)
}

func validatePurgeConfirmation(confirmed bool) error {
	if !confirmed {
		return fmt.Errorf("purge queue requires --yes")
	}

	return nil
}

func ctx() context.Context {
	return context.Background()
}
