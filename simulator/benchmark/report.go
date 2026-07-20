package benchmark

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type WorkloadMetadata struct {
	Scenario         Scenario `json:"scenario"`
	Clients          int      `json:"clients"`
	GroupSenders     int      `json:"group_senders,omitempty"`
	TargetRate       int      `json:"target_messages_per_second"`
	Warmup           string   `json:"warmup"`
	Duration         string   `json:"duration"`
	DeliveryGrace    string   `json:"delivery_grace"`
	PayloadBytes     int      `json:"payload_bytes"`
	StoreMessages    bool     `json:"store_messages"`
	CountMessages    bool     `json:"count_messages"`
	WebSocketURL     string   `json:"websocket_url"`
	APIURL           string   `json:"api_url"`
	EnvironmentLabel string   `json:"environment_label"`
	DatabaseLabel    string   `json:"database_label"`
}

type EnvironmentMetadata struct {
	ServerCommit     string `json:"server_commit"`
	WorkingTreeDirty bool   `json:"working_tree_dirty"`
	OS               string `json:"os"`
	Architecture     string `json:"architecture"`
	GoVersion        string `json:"go_version"`
	CPUModel         string `json:"cpu_model"`
	LogicalCPUs      int    `json:"logical_cpus"`
	TotalMemoryBytes uint64 `json:"total_memory_bytes"`
}

type Report struct {
	SchemaVersion           int                 `json:"schema_version"`
	RunID                   string              `json:"run_id"`
	GeneratedAt             time.Time           `json:"generated_at"`
	StartedAt               time.Time           `json:"started_at"`
	MeasurementStartedAt    time.Time           `json:"measurement_started_at"`
	MeasurementEndedAt      time.Time           `json:"measurement_ended_at"`
	Workload                WorkloadMetadata    `json:"workload"`
	Environment             EnvironmentMetadata `json:"environment"`
	Connections             MetricSnapshot      `json:"connections"`
	MessageAcknowledgements MetricSnapshot      `json:"message_acknowledgements"`
	Deliveries              MetricSnapshot      `json:"deliveries"`
	Limitations             []string            `json:"limitations"`
}

func newReport(config Config, runID string, startedAt, measurementStarted, measurementEnded time.Time) Report {
	groupSenders := 0
	if config.Scenario == ScenarioGroup {
		groupSenders = config.GroupSenders
	}
	commit, dirty := gitState()
	if config.ServerCommit != "" {
		commit = config.ServerCommit
	}
	return Report{
		SchemaVersion:        1,
		RunID:                runID,
		GeneratedAt:          time.Now().UTC(),
		StartedAt:            startedAt,
		MeasurementStartedAt: measurementStarted,
		MeasurementEndedAt:   measurementEnded,
		Workload: WorkloadMetadata{
			Scenario: config.Scenario, Clients: config.Clients, GroupSenders: groupSenders,
			TargetRate: config.Rate, Warmup: config.Warmup.String(), Duration: config.Duration.String(),
			DeliveryGrace: config.DeliveryGrace.String(), PayloadBytes: config.PayloadBytes,
			StoreMessages: config.StoreMessages, CountMessages: config.CountMessages,
			WebSocketURL: config.WSURL, APIURL: config.APIURL,
			EnvironmentLabel: config.EnvironmentLabel, DatabaseLabel: config.DatabaseLabel,
		},
		Environment: EnvironmentMetadata{
			ServerCommit: commit, WorkingTreeDirty: dirty, OS: runtime.GOOS,
			Architecture: runtime.GOARCH, GoVersion: runtime.Version(), CPUModel: cpuModel(),
			LogicalCPUs: runtime.NumCPU(), TotalMemoryBytes: totalMemoryBytes(),
		},
		Limitations: []string{
			"ACK latency uses the load generator's monotonic clock; delivery latency uses wall-clock timestamps within the same load-generator process.",
			"Acknowledgement latency measures client publish to server ACK; delivery latency measures client publish to recipient callback.",
			"The load generator and local Docker stack may share hardware unless the environment label says otherwise.",
			"Results from different hardware or workload settings are not directly comparable.",
		},
	}
}

func WriteReport(path string, report Report) (string, error) {
	if path == "" {
		path = filepath.Join("benchmark-results", report.RunID+"-"+string(report.Workload.Scenario)+".json")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("create report directory: %w", err)
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encode benchmark report: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("write benchmark report: %w", err)
	}
	return path, nil
}

func gitState() (string, bool) {
	commitOutput, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "unknown", false
	}
	dirtyOutput, _ := exec.Command("git", "status", "--porcelain").Output()
	return strings.TrimSpace(string(commitOutput)), len(strings.TrimSpace(string(dirtyOutput))) > 0
}

func cpuModel() string {
	if runtime.GOOS == "darwin" {
		if output, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output(); err == nil {
			return strings.TrimSpace(string(output))
		}
	}
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "unknown"
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "model name") || strings.HasPrefix(line, "Hardware") {
			if _, value, ok := strings.Cut(line, ":"); ok {
				return strings.TrimSpace(value)
			}
		}
	}
	return "unknown"
}

func totalMemoryBytes() uint64 {
	if runtime.GOOS == "darwin" {
		if output, err := exec.Command("sysctl", "-n", "hw.memsize").Output(); err == nil {
			value, _ := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 64)
			return value
		}
	}
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 && fields[0] == "MemTotal:" {
			value, _ := strconv.ParseUint(fields[1], 10, 64)
			return value * 1024
		}
	}
	return 0
}
