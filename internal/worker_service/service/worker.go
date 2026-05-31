package service

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"my-go-server/internal/worker_service"
)

type WorkerPool struct {
	repo workerservice.Repository
	workersCount int
	jobs chan string
	wg sync.WaitGroup
}

func NewWorkerPool(repo workerservice.Repository, workersCount int) *WorkerPool {
	return &WorkerPool{
		repo: repo,
		workersCount: workersCount,
		jobs: make(chan string, 100),
	}
}

func (p *WorkerPool) Run(ctx context.Context) {
	for i := 0; i < p.workersCount; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i+1)
	}

	p.repo.StartNewOrdersConsumer(ctx, func(orderID string) error {
		select {
		case p.jobs <- orderID:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	<-ctx.Done()

	close(p.jobs)
	p.wg.Wait()
}

func (p *WorkerPool) worker(ctx context.Context, workerID int) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case orderID, ok := <-p.jobs:
			if !ok {
				return
			}

			if err := p.processOrder(ctx, workerID, orderID); err != nil {
				fmt.Printf("worker %d failed to process order %s: %v\n", workerID, orderID, err)
			}
		}
	}
}

func (p *WorkerPool) processOrder(ctx context.Context, workerID int, orderID string) error {
	fmt.Printf("worker %d started processing order %s\n", workerID, orderID)

	processingTime := time.Duration(rand.Intn(4)+1)*time.Second

	select {
	case <-time.After(processingTime):
	case <-ctx.Done():
		return ctx.Err()
	}

	status := "CONFIRMED"

	if err := p.repo.PublishNewOrderStatus(ctx, orderID, status); err != nil {
		return err
	}

	fmt.Printf("worker %d confirmed order %s\n", workerID, orderID)

	return nil
}