package main

import (
	"fmt"
)

type Job struct {
	ID    string
	File  string
	Codec string
	state JobState
}

func NewJob(id, file, codec string) *Job {
	return &Job{
		ID:    id,
		File:  file,
		Codec: codec,
		state: CreatedState{},
	}
}

func (j *Job) setState(state JobState) {
	j.state = state
}

func (j *Job) Status() string {
	return j.state.Name()
}

func (j *Job) Queue() error {
	return j.state.Queue(j)
}

func (j *Job) Start() error {
	return j.state.Start(j)
}

func (j *Job) Complete() error {
	return j.state.Complete(j)
}

func (j *Job) Fail() error {
	return j.state.Fail(j)
}

type JobRepository interface {
	Save(job *Job) error
	Get(id string) (*Job, error)
}

type InMemoryJobRepository struct {
	jobs map[string]*Job
}

func NewInMemoryJobRepository() *InMemoryJobRepository {
	return &InMemoryJobRepository{
		jobs: make(map[string]*Job),
	}
}

func (r *InMemoryJobRepository) Save(job *Job) error {
	r.jobs[job.ID] = job
	return nil
}

func (r *InMemoryJobRepository) Get(id string) (*Job, error) {
	job, exists := r.jobs[id]
	if !exists {
		return nil, fmt.Errorf("[InMemoryJobRepository] The requested job (id: %s) doesn't exist", id)
	}
	return job, nil
}

type JobState interface {
	Queue(job *Job) error
	Start(job *Job) error
	Complete(job *Job) error
	Fail(job *Job) error
	Name() string
}

type CreatedState struct{}

func (s CreatedState) Queue(job *Job) error {
	job.setState(QueuedState{})
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
	job.setState(ProcessingState{})
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

type ProcessingState struct{}

func (s ProcessingState) Queue(job *Job) error {
	return fmt.Errorf("cannot QUEUE job from PROCESSING")
}

func (s ProcessingState) Start(job *Job) error {
	return fmt.Errorf("cannot START job from PROCESSING")
}

func (s ProcessingState) Complete(job *Job) error {
	job.setState(CompletedState{})
	return nil
}

func (s ProcessingState) Fail(job *Job) error {
	job.setState(FailedState{})
	return nil
}

func (s ProcessingState) Name() string {
	return "PROCESSING"
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

type Encoder interface {
	Encode(file string) error
}

type H264Encoder struct{}

func (H264Encoder) Encode(file string) error {
	fmt.Println("Encoding", file, "with H264")
	return nil
}

type H265Encoder struct{}

func (H265Encoder) Encode(file string) error {
	fmt.Println("Encoding", file, "with H265")
	return nil
}

type VP9Encoder struct{}

func (VP9Encoder) Encode(file string) error {
	fmt.Println("Encoding", file, "with VP9")
	return nil
}

func EncoderFactory(codec string) (Encoder, error) {
	switch codec {
	case "h264":
		return H264Encoder{}, nil
	case "h265":
		return H265Encoder{}, nil
	case "vp9":
		return VP9Encoder{}, nil
	default:
		return nil, fmt.Errorf("unsupported codec: %s", codec)
	}
}

type Observer interface {
	Update(job *Job)
}

type NotificationObserver struct{}

func (NotificationObserver) Update(job *Job) {
	fmt.Println("[Notification] Job changed to:", job.Status())
}

type MetricsObserver struct{}

func (MetricsObserver) Update(job *Job) {
	fmt.Println("[Metric] Job changed to:", job.Status())
}

type MediaService struct {
	repo      JobRepository
	observers []Observer
}

func NewMediaService(repo JobRepository) *MediaService {
	return &MediaService{
		repo: repo,
	}
}

func (m *MediaService) Subscribe(o Observer) {
	m.observers = append(m.observers, o)
}

func (m *MediaService) notify(job *Job) {
	for _, observer := range m.observers {
		observer.Update(job)
	}
}

func (m *MediaService) CreateJob(id, file, codec string) (*Job, error) {
	job := NewJob(id, file, codec)

	if err := m.repo.Save(job); err != nil {
		return nil, err
	}

	return job, nil
}

func (m *MediaService) QueueJob(id string) error {
	job, err := m.repo.Get(id)
	if err != nil {
		return err
	}

	if err := job.Queue(); err != nil {
		return err
	}

	return m.repo.Save(job)
}

func (m *MediaService) StartJob(id string) error {
	job, err := m.repo.Get(id)
	if err != nil {
		return err
	}

	if err := job.Start(); err != nil {
		return err
	}

	if err := m.repo.Save(job); err != nil {
		return err
	}

	encoder, err := EncoderFactory(job.Codec)
	if err != nil {
		_ = job.Fail()
		_ = m.repo.Save(job)
		m.notify(job)
		return err
	}

	if err = encoder.Encode(job.File); err != nil {
		_ = job.Fail()
		_ = m.repo.Save(job)
		m.notify(job)
		return err
	}

	if err = job.Complete(); err != nil {
		return err
	}

	if err = m.repo.Save(job); err != nil {
		return err
	}

	m.notify(job)
	return nil
}

func (m *MediaService) GetJob(id string) (*Job, error) {
	return m.repo.Get(id)
}

func main() {
	repo := NewInMemoryJobRepository()
	service := NewMediaService(repo)

	service.Subscribe(NotificationObserver{})
	service.Subscribe(MetricsObserver{})

	job, err := service.CreateJob("job1", "movie.mp4", "vp9")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Created:", job.Status())

	if err := service.QueueJob(job.ID); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Queued:", job.Status())

	if err := service.StartJob(job.ID); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Final:", job.Status())
}
