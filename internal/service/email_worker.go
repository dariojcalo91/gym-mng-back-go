package service

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// EmailJob carries everything needed to send a real email — the old design
// only carried an address, which was enough for a generic "welcome" email
// but can't carry the member's name or expiration date for a reminder.
type EmailJob struct {
	To      string
	Subject string
	Body    string
}

// EmailWorkerPool is transport-agnostic infrastructure: it knows nothing about
// members or gyms, only how to process EmailJobs concurrently. Any service
// that needs to send email depends on it through the Notifier interface below.
type EmailWorkerPool struct {
	jobs chan EmailJob
	wg   sync.WaitGroup
}

type Notifier interface {
	Enqueue(job EmailJob)
}

func NewEmailWorkerPool(bufferSize int) *EmailWorkerPool {
	return &EmailWorkerPool{jobs: make(chan EmailJob, bufferSize)}
}

func (p *EmailWorkerPool) Start(ctx context.Context, workers int) {
	p.wg.Add(workers)
	for i := range workers {
		go p.worker(ctx, i+1)
	}
}

func (p *EmailWorkerPool) Enqueue(job EmailJob) {
	select {
	case p.jobs <- job:
	default:
		slog.Warn("email channel full, dropping email", "to", job.To, "subject", job.Subject)
	}
}

func (p *EmailWorkerPool) worker(ctx context.Context, id int) {
	defer p.wg.Done()
	slog.Info("worker ready to process emails", "worker_id", id)
	for {
		select {
		case job, ok := <-p.jobs:
			if !ok {
				return
			}
			_ = p.send(ctx, job, id)
		case <-ctx.Done():
			slog.Info("worker shutting down", "worker_id", id, "reason", ctx.Err())
			return
		}
	}
}

func (p *EmailWorkerPool) send(ctx context.Context, job EmailJob, workerID int) error {
	// TODO: replace with a real SMTP/SES call. Simulated for now.
	select {
	case <-time.After(2 * time.Second):
		slog.Info("email sent", "worker_id", workerID, "to", job.To, "subject", job.Subject)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *EmailWorkerPool) Shutdown() {
	slog.Info("closing email channel and waiting for workers")
	close(p.jobs)
	p.wg.Wait()
	slog.Info("all workers completed")
}
