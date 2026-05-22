package process

import (
	"dnslookup/entity"
	"dnslookup/intel"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func StartSchedulerAsync(store *entity.Storage, config entity.Configuration) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		slog.Error("unable to create scheduler", "error", err)
		return
	}

	_, err = scheduler.NewJob(gocron.DurationJob(time.Hour), gocron.NewTask(intel.UpdateIntel, store, config))
	if err != nil {
		slog.Error("unable to create schedule update intel", "error", err)
	}

	scheduler.Start()
}
