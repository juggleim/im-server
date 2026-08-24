package benchmark

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/msgdefines"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"
)

type benchmarkPayload struct {
	BenchmarkID string `json:"benchmark_id"`
	Phase       string `json:"phase"`
	Sequence    uint64 `json:"sequence"`
	SentAtNS    int64  `json:"sent_at_unix_nano"`
	Padding     string `json:"padding"`
}

type connectedClient struct {
	user   registeredUser
	client *wsclients.WsImClient
}

type Runner struct {
	config       Config
	runID        string
	groupID      string
	api          *apiClient
	connections  *metricRecorder
	acknowledged *metricRecorder
	deliveries   *metricRecorder
	sequence     atomic.Uint64
}

func NewRunner(config Config) (*Runner, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	nonce, err := randomNonce()
	if err != nil {
		return nil, fmt.Errorf("create benchmark run ID: %w", err)
	}
	runID := fmt.Sprintf("%d-%s", time.Now().UTC().Unix(), nonce[:8])
	return &Runner{
		config:       config,
		runID:        runID,
		groupID:      "bench-group-" + runID,
		api:          newAPIClient(config),
		connections:  newMetricRecorder(),
		acknowledged: newMetricRecorder(),
		deliveries:   newMetricRecorder(),
	}, nil
}

func (r *Runner) Run(ctx context.Context) (Report, error) {
	startedAt := time.Now().UTC()
	users, err := r.prepareUsers(ctx)
	if err != nil {
		return Report{}, err
	}
	if r.config.Scenario == ScenarioGroup {
		memberIDs := make([]string, len(users))
		for index, user := range users {
			memberIDs[index] = user.ID
		}
		if err := r.api.createGroup(ctx, r.groupID, memberIDs); err != nil {
			return Report{}, fmt.Errorf("prepare benchmark group: %w", err)
		}
	}

	clients := r.connect(ctx, users)
	if len(clients) != len(users) {
		r.disconnect(clients)
		return Report{}, fmt.Errorf("connected %d of %d clients; refusing to publish a partial baseline", len(clients), len(users))
	}
	defer r.disconnect(clients)

	if r.config.Warmup > 0 {
		if err := r.runTraffic(ctx, clients, r.config.Warmup, false); err != nil {
			return Report{}, fmt.Errorf("warm-up: %w", err)
		}
	}
	measurementStarted := time.Now().UTC()
	if err := r.runTraffic(ctx, clients, r.config.Duration, true); err != nil {
		return Report{}, fmt.Errorf("measurement: %w", err)
	}
	measurementEnded := measurementStarted.Add(r.config.Duration)
	if r.config.DeliveryGrace > 0 {
		timer := time.NewTimer(r.config.DeliveryGrace)
		select {
		case <-ctx.Done():
			timer.Stop()
			return Report{}, ctx.Err()
		case <-timer.C:
		}
	}

	report := newReport(r.config, r.runID, startedAt, measurementStarted, measurementEnded)
	report.Connections = r.connections.snapshot(0)
	report.MessageAcknowledgements = r.acknowledged.snapshot(r.config.Duration)
	report.Deliveries = r.deliveries.snapshot(r.config.Duration)
	return report, nil
}

func (r *Runner) prepareUsers(ctx context.Context) ([]registeredUser, error) {
	users := make([]registeredUser, r.config.Clients)
	jobs := make(chan int)
	errCh := make(chan error, r.config.Clients)
	var workers sync.WaitGroup
	workerCount := min(r.config.SetupConcurrency, r.config.Clients)
	for worker := 0; worker < workerCount; worker++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for index := range jobs {
				userID := fmt.Sprintf("bench-%s-%06d", r.runID, index+1)
				user, err := r.api.registerUser(ctx, userID)
				if err != nil {
					errCh <- fmt.Errorf("register %s: %w", userID, err)
					continue
				}
				users[index] = user
			}
		}()
	}
	for index := range users {
		select {
		case <-ctx.Done():
			close(jobs)
			workers.Wait()
			return nil, ctx.Err()
		case jobs <- index:
		}
	}
	close(jobs)
	workers.Wait()
	close(errCh)
	for err := range errCh {
		return nil, err
	}
	return users, nil
}

func (r *Runner) connect(ctx context.Context, users []registeredUser) []connectedClient {
	clients := make([]connectedClient, len(users))
	semaphore := make(chan struct{}, r.config.ConnectConcurrency)
	var wait sync.WaitGroup
	for index, user := range users {
		select {
		case <-ctx.Done():
			return compactClients(clients)
		case semaphore <- struct{}{}:
		}
		wait.Add(1)
		go func(index int, user registeredUser) {
			defer wait.Done()
			defer func() { <-semaphore }()
			client := wsclients.NewWsImClient(r.config.WSURL, r.config.AppKey, user.Token, r.onMessage, nil, nil)
			client.Verbose = false
			started := time.Now()
			code, _ := client.Connect("benchmark", "local")
			latency := time.Since(started)
			if code != utils.ClientErrorCode_Success {
				r.connections.recordFailure(strconv.Itoa(int(code)), latency)
				return
			}
			r.connections.recordSuccess(latency)
			clients[index] = connectedClient{user: user, client: client}
		}(index, user)
	}
	wait.Wait()
	return compactClients(clients)
}

func compactClients(clients []connectedClient) []connectedClient {
	result := make([]connectedClient, 0, len(clients))
	for _, client := range clients {
		if client.client != nil {
			result = append(result, client)
		}
	}
	return result
}

func (r *Runner) disconnect(clients []connectedClient) {
	for _, client := range clients {
		client.client.Disconnect()
	}
}

func (r *Runner) runTraffic(parent context.Context, clients []connectedClient, duration time.Duration, measured bool) error {
	ctx, cancel := context.WithTimeout(parent, duration)
	defer cancel()
	senderCount := len(clients)
	if r.config.Scenario == ScenarioGroup {
		senderCount = min(r.config.GroupSenders, len(clients))
	}
	tokens := make(chan struct{})
	var workers sync.WaitGroup
	for senderIndex := 0; senderIndex < senderCount; senderIndex++ {
		workers.Add(1)
		go func(senderIndex int) {
			defer workers.Done()
			for range tokens {
				r.sendOne(clients, senderIndex, measured)
			}
		}(senderIndex)
	}

	interval := time.Second / time.Duration(r.config.Rate)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			close(tokens)
			workers.Wait()
			if parent.Err() != nil {
				return parent.Err()
			}
			return nil
		case <-ticker.C:
			select {
			case tokens <- struct{}{}:
			case <-ctx.Done():
			}
		}
	}
}

func (r *Runner) sendOne(clients []connectedClient, senderIndex int, measured bool) {
	sequence := r.sequence.Add(1)
	phase := "warmup"
	if measured {
		phase = "measure"
	}
	payload, err := makePayload(r.runID, phase, sequence, time.Now(), r.config.PayloadBytes)
	if err != nil {
		if measured {
			r.acknowledged.recordFailure("payload", 0)
		}
		return
	}
	flags := int32(0)
	if r.config.StoreMessages {
		flags = msgdefines.SetStoreMsg(flags)
	}
	if r.config.CountMessages {
		flags = msgdefines.SetCountMsg(flags)
	}
	message := &pbobjs.UpMsg{MsgType: msgdefines.InnerMsgType_Text, MsgContent: payload, Flags: flags}
	started := time.Now()
	var code utils.ClientErrorCode
	if r.config.Scenario == ScenarioPrivate {
		target := clients[(senderIndex+1)%len(clients)].user.ID
		code, _ = clients[senderIndex].client.SendPrivateMsg(target, message)
	} else {
		code, _ = clients[senderIndex].client.SendGroupMsg(r.groupID, message)
	}
	if !measured {
		return
	}
	latency := time.Since(started)
	if code == utils.ClientErrorCode_Success {
		r.acknowledged.recordSuccess(latency)
	} else {
		r.acknowledged.recordFailure(strconv.Itoa(int(code)), latency)
	}
}

func (r *Runner) onMessage(message *pbobjs.DownMsg) {
	var payload benchmarkPayload
	if err := json.Unmarshal(message.MsgContent, &payload); err != nil {
		return
	}
	if payload.BenchmarkID != r.runID || payload.Phase != "measure" || payload.SentAtNS <= 0 {
		return
	}
	latency := time.Since(time.Unix(0, payload.SentAtNS))
	if latency < 0 {
		return
	}
	r.deliveries.recordObservation(latency)
}

func makePayload(runID, phase string, sequence uint64, sentAt time.Time, size int) ([]byte, error) {
	payload := benchmarkPayload{
		BenchmarkID: runID,
		Phase:       phase,
		Sequence:    sequence,
		SentAtNS:    sentAt.UnixNano(),
	}
	base, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	if missing := size - len(base); missing > 0 {
		payload.Padding = strings.Repeat("x", missing)
	}
	return json.Marshal(payload)
}
