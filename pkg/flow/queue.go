package flow

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type JobInterface interface {
	Execute()
}

type JobWithID struct {
	Job JobInterface
	ID  string
}

type Worker struct {
	ID         int
	JobChannel chan JobWithID
}

type Queue struct {
	Worker      Worker
	JobChannel  chan JobWithID
	WaitGroup   sync.WaitGroup
	Jobs        map[string]JobWithID
	JobsMutex   sync.Mutex
	LastJobID   int
	flowSvc     *FlowService
	running     bool
	stopChannel chan struct{}
}

func newQueue(svc *FlowService) *Queue {
	return &Queue{
		JobChannel:  make(chan JobWithID),
		Jobs:        make(map[string]JobWithID),
		LastJobID:   0,
		flowSvc:     svc,
		stopChannel: make(chan struct{}),
	}
}

func (q *Queue) Start(rateLimit int) {
	if q.running {
		return
	}
	q.running = true

	q.Worker = Worker{
		JobChannel: q.JobChannel, // Assign the JobChannel of the worker to the JobChannel of the queue
	}
	go q.worker()

	ticker := time.NewTicker(time.Second / time.Duration(rateLimit))

	for {
		select {
		case <-q.stopChannel:
			ticker.Stop()
			return
		case <-ticker.C:
			q.JobsMutex.Lock()
			jobIDs := make([]string, 0, len(q.Jobs))
			for id := range q.Jobs {
				jobIDs = append(jobIDs, id)
			}
			q.JobsMutex.Unlock()

			for _, id := range jobIDs {
				q.WaitGroup.Add(1)
				jobWithID := q.Jobs[id]
				delete(q.Jobs, id)
				go func() {
					q.processJob(jobWithID)
					q.WaitGroup.Done()
				}()
			}
		}
	}
}

func (q *Queue) worker() {
	for jobWithID := range q.Worker.JobChannel {
		q.WaitGroup.Add(1)
		q.processJob(jobWithID)
		q.WaitGroup.Done()
		q.JobsMutex.Lock()
		delete(q.Jobs, jobWithID.ID)
		q.JobsMutex.Unlock()
	}
}

func (q *Queue) processJob(jobWithID JobWithID) {
	fmt.Printf("Processing job %s:", jobWithID.ID)
	jobWithID.Job.Execute()
	fmt.Printf("Job %s completed\n", jobWithID.ID)
}

func (q *Queue) Stop() {
	if q.running {
		close(q.stopChannel)
	}
}

func (q *Queue) RemoveJobByID(id string) error {
	q.JobsMutex.Lock()
	_, ok := q.Jobs[id]
	if !ok {
		q.JobsMutex.Unlock()
		return errors.New("Job not found: " + id)
	}
	delete(q.Jobs, id)
	q.JobsMutex.Unlock()

	return nil
}

func (q *Queue) CreateJob(job JobInterface) (string, error) {
	q.JobsMutex.Lock()
	defer q.JobsMutex.Unlock()

	newUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	jobID := newUUID.String()

	jobWithID := JobWithID{
		Job: job,
		ID:  jobID,
	}

	q.Jobs[jobID] = jobWithID
	go func() {
		q.JobChannel <- jobWithID
	}()

	return jobID, nil
}
