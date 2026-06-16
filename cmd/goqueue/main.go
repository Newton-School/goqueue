package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/Newton-School/goqueue"
)

const (
	defaultNamespace   = "goqueue"
	defaultTaskQueue   = "default"
	defaultDeadLetters = 25
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

func printUsage() {
	fmt.Println("goqueue CLI")
	fmt.Println("Usage:")
	fmt.Println("  goqueue inspect <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  inspect ping")
	fmt.Println("  inspect stats --queue <queue>")
	fmt.Println("  inspect deadletters --queue <queue> [--count <n>]")
	fmt.Println("  inspect task --id <task_id> [--json]")
	fmt.Println("  inspect state --id <task_id> [--json]")
	fmt.Println("  inspect result --id <task_id> [--json]")
	fmt.Println("  inspect forget-result --id <task_id>")
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

func ctx() context.Context {
	return context.Background()
}
