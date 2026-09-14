package main

import (
	"fmt"
)

type Observer interface {
	Update(job *Job)
}

type NotificationObserver struct{}

func (n NotificationObserver) Update(job *Job) {
	fmt.Println("[Notification] Job changed to:", job.Status())
}

type MetricsObserver struct{}

func (n MetricsObserver) Update(job *Job) {
	fmt.Println("[Metrics] Job changed to:", job.Status())
}

type AuditObserver struct{}

func (n AuditObserver) Update(job *Job) {
	fmt.Println("[Audit] Job changed to:", job.Status())
}

type Job struct {
	status    string
	observers []Observer
}

func NewJob() *Job {
	return &Job{
		status: "CREATED",
	}
}

func (j *Job) Subscribe(o Observer) {
	j.observers = append(j.observers, o)
}

func (j *Job) Notify() {
	for _, observer := range j.observers {
		observer.Update(j)
	}
}

func (j *Job) SetStatus(status string) {
	j.status = status
	j.Notify()
}

func (j *Job) Status() string {
	return j.status
}

func main() {
	job := NewJob()

	// Create all three observers
	notificationObserver := NotificationObserver{}
	metricsObserver := MetricsObserver{}
	auditObserver := AuditObserver{}
	// Subscribe them
	job.Subscribe(notificationObserver)
	job.Subscribe(metricsObserver)
	job.Subscribe(auditObserver)

	job.SetStatus("QUEUED")
	job.SetStatus("PROCESSING")
	job.SetStatus("COMPLETED")
}
