package main

import (
	"fmt"
	"strings"
	"time"

	swarming_api "go.chromium.org/luci/common/api/swarming/swarming/v1"
	"go.skia.org/infra/go/git/repograph"
	"go.skia.org/infra/go/sklog"
	"go.skia.org/infra/go/swarming"
	"go.skia.org/infra/task_scheduler/go/specs"
)

const (
	BOT_ID_PREFIX = "bot-"
)

// mockSwarmingBotsForAllTasksForTesting returns a list containing one swarming
// bot for each TaskSpec in the given repos, or nil on error.
func mockSwarmingBotsForAllTasksForTesting(repos repograph.Map) []*swarming_api.SwarmingRpcsBotInfo {
	botId := 0
	rv := []*swarming_api.SwarmingRpcsBotInfo{}
	for _, repo := range repos {
		branches, err := repo.Repo().Branches()
		if err != nil {
			sklog.Error(err)
			continue
		}
		for _, branch := range branches {
			if branch.Name != "master" {
				continue
			}
			contents, err := repo.Repo().GetFile(specs.TASKS_CFG_FILE, branch.Head)
			if err != nil {
				sklog.Error(err)
				continue
			}
			cfg, err := specs.ParseTasksCfg(contents)
			if err != nil {
				sklog.Error(err)
				continue
			}
			for _, spec := range cfg.Tasks {
				dimensions := make([]*swarming_api.SwarmingRpcsStringListPair, 0, len(spec.Dimensions))
				for _, d := range spec.Dimensions {
					split := strings.SplitN(d, ":", 2)
					dimensions = append(dimensions, &swarming_api.SwarmingRpcsStringListPair{
						Key:   split[0],
						Value: []string{split[1]},
					})
				}
				rv = append(rv, &swarming_api.SwarmingRpcsBotInfo{
					BotId:      fmt.Sprintf("%s%03d", BOT_ID_PREFIX, botId),
					Dimensions: dimensions,
				})
				botId++
			}
		}
	}
	return rv
}

// periodicallyUpdateMockTasksForTesting simulates running the mocked tasks in
// TestClient by updating the status, started/completed times, isolated output,
// etc. Does not return.
func periodicallyUpdateMockTasksForTesting(swarm *swarming.TestClient) {
	for range time.Tick(time.Minute) {
		swarm.DoMockTasks(func(task *swarming_api.SwarmingRpcsTaskRequestMetadata) {
			created, err := swarming.Created(task)
			if err != nil {
				return
			}
			if task.TaskResult.State == swarming.TASK_STATE_PENDING {
				task.TaskResult.State = swarming.TASK_STATE_RUNNING
				task.TaskResult.StartedTs = time.Now().Format(swarming.TIMESTAMP_FORMAT)
				task.TaskResult.BotId = fmt.Sprintf("A-Bot-To-Run-%s", task.TaskResult.Name)
			} else if task.TaskResult.State == swarming.TASK_STATE_RUNNING && created.Add(5*time.Minute).Before(time.Now()) {
				task.TaskResult.State = swarming.TASK_STATE_COMPLETED
				task.TaskResult.OutputsRef = &swarming_api.SwarmingRpcsFilesRef{
					Isolated: fmt.Sprintf("Isolated-%s", task.TaskId),
				}
				task.TaskResult.CompletedTs = time.Now().Format(swarming.TIMESTAMP_FORMAT)
			}
		})
	}
}
