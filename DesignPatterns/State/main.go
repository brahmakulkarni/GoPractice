package main

import (
	"fmt"
)

type JobState interface {
	Queue(job *Job) error
	Start(job *Job) error
	Complete(job *Job) error
	Fail(job *Job) error
	Name() string
}

type Job struct {
	state JobState
}

func NewJob() *Job {
	return &Job{
		state: CreatedState{},
	}
}

func (job *Job) SetState(state JobState) {
	job.state = state
}

func (job *Job) Queue() error {
	return job.state.Queue(job)
}

func (job *Job) Start() error {
	return job.state.Start(job)
}

func (job *Job) Complete() error {
	return job.state.Complete(job)
}

func (job *Job) Fail() error {
	return job.state.Fail(job)
}

func (job *Job) Status() string {
	return job.state.Name()
}

type CreatedState struct{}

func (s CreatedState) Queue(job *Job) error {
	job.SetState(QueuedState{})
	return nil
}

func (s CreatedState) Start(job *Job) error {
	return fmt.Errorf("cannot START job from CREATED")
}

func (s CreatedState) Complete(job *Job) error {
	return fmt.Errorf("cannot COMPLETE job from CREATED")
}

func (s CreatedState) Fail(job *Job) error {
	return fmt.Errorf("cannot FAIL job from CREATED")
}

func (s CreatedState) Name() string {
	return "CREATED"
}

type QueuedState struct{}

func (s QueuedState) Queue(job *Job) error {
	return fmt.Errorf("cannot QUEUE job from QUEUED")
}

func (s QueuedState) Start(job *Job) error {
	job.SetState(StartedState{})
	return nil
}

func (s QueuedState) Complete(job *Job) error {
	return fmt.Errorf("cannot COMPLETE job from QUEUED")
}

func (s QueuedState) Fail(job *Job) error {
	return fmt.Errorf("cannot FAIL job from QUEUED")
}

func (s QueuedState) Name() string {
	return "QUEUED"
}

type StartedState struct{}

func (s StartedState) Queue(job *Job) error {
	return fmt.Errorf("cannot QUEUE job from STARTED")
}

func (s StartedState) Start(job *Job) error {
	return fmt.Errorf("cannot START job from STARTED")
}

func (s StartedState) Complete(job *Job) error {
	job.SetState(CompletedState{})
	return nil
}

func (s StartedState) Fail(job *Job) error {
	job.SetState(FailedState{})
	return nil
}

func (s StartedState) Name() string {
	return "STARTED"
}

type CompletedState struct{}

func (s CompletedState) Queue(job *Job) error {
	return fmt.Errorf("cannot QUEUE job from COMPLETED")
}

func (s CompletedState) Start(job *Job) error {
	return fmt.Errorf("cannot START job from COMPLETED")
}

func (s CompletedState) Complete(job *Job) error {
	return fmt.Errorf("cannot COMPLETE job from COMPLETED")
}

func (s CompletedState) Fail(job *Job) error {
	return fmt.Errorf("cannot FAIL job from COMPLETED")
}

func (s CompletedState) Name() string {
	return "COMPLETED"
}

type FailedState struct{}

func (s FailedState) Queue(job *Job) error {
	return fmt.Errorf("cannot QUEUE job from FAILED")
}

func (s FailedState) Start(job *Job) error {
	return fmt.Errorf("cannot START job from FAILED")
}

func (s FailedState) Complete(job *Job) error {
	return fmt.Errorf("cannot COMPLETE job from FAILED")
}

func (s FailedState) Fail(job *Job) error {
	return fmt.Errorf("cannot FAIL job from FAILED")
}

func (s FailedState) Name() string {
	return "FAILED"
}

func main() {
	job := NewJob()

	fmt.Println(job.Status())

	if err := job.Queue(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(job.Status())

	if err := job.Start(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(job.Status())

	if err := job.Complete(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(job.Status())

	if err := job.Start(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(job.Status())
}
