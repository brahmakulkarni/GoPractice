package main

import (
	"fmt"
	"time"
)

type Encoder interface {
	Encode(file string) error
}

type H264Encoder struct{}

func (H264Encoder) Encode(file string) error {
	fmt.Println("Encoding", file, "using H264")
	time.Sleep(500 * time.Millisecond)
	return nil
}

type H265Encoder struct{}

func (H265Encoder) Encode(file string) error {
	fmt.Println("Encoding", file, "using H265")
	time.Sleep(500 * time.Millisecond)
	return nil
}

type AV1Encoder struct{}

func (AV1Encoder) Encode(file string) error {
	fmt.Println("Encoding", file, "using AV1")
	time.Sleep(500 * time.Millisecond)
	return nil
}

type VP9Encoder struct{}

func (VP9Encoder) Encode(file string) error {
	fmt.Println("Encoding", file, "using VP9")
	time.Sleep(500 * time.Millisecond)
	return nil
}

type FailingEncoder struct {}

func (FailingEncoder) Encode(file string) error {
	return fmt.Errorf("this encoder always fails")
}

// =================== Decorators

type LoggingEncoder struct {
	encoder Encoder
}

func NewLoggingEncoder(encoder Encoder) *LoggingEncoder {
	return &LoggingEncoder{
		encoder: encoder,
	}
}

func (l *LoggingEncoder) Encode(file string) error {
	fmt.Println("[LOG] Encoding started")

	err := l.encoder.Encode(file)

	fmt.Println("[LOG] Encoding Finished")
	return err
}

type MetricsEncoder struct {
	encoder Encoder
}

func NewMetricsEncoder(encoder Encoder) *MetricsEncoder {
	return &MetricsEncoder{
		encoder: encoder,
	}
}

func (m *MetricsEncoder) Encode(file string) error {
	start := time.Now()

	err := m.encoder.Encode(file)

	duration := time.Since(start)
	fmt.Println("[METRIC] Encoding took:", duration)
	return err
}

type RetryEncoder struct {
	encoder Encoder
	maxRetries int
}

func NewRetryEncoder(encoder Encoder, maxRetries int) *RetryEncoder {
	return &RetryEncoder{
		encoder: encoder,
		maxRetries: maxRetries,
	}
}

func (r *RetryEncoder) Encode(file string) error {
	var err error
	for i := range r.maxRetries {
		err = r.encoder.Encode(file)
		if err == nil {
			fmt.Printf("[RETRY] Encoder succeded on the attempt %d\n", i+1)
			break
		} else {
			fmt.Printf("[RETRY] Encoder failed on the attempt %d\n", i+1)
		}
	}
	return err
}

// ==========================

type VideoProcessor struct {
	encoder Encoder
}

func NewVideoProcessor(encoder Encoder) *VideoProcessor {
	return &VideoProcessor{
		encoder: encoder,
	}
}

func (v *VideoProcessor) Process(file string) error {
	return v.encoder.Encode(file)
}

func EncoderFactory(codec string) (Encoder, error) {
	switch codec {
	case "h264":
		return H264Encoder{}, nil
	case "h265":
		return H265Encoder{}, nil
	case "av1":
		return AV1Encoder{}, nil
	case "vp9":
		return VP9Encoder{}, nil
	case "failing":
		return FailingEncoder{}, nil
	default:
		return nil, fmt.Errorf("unsupported codec: %s", codec)
	}
}

func main() {
	codec := "failing"
	file := "movie.mp4"

	encoder, err := EncoderFactory(codec)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	encoder = NewLoggingEncoder(encoder)
	encoder = NewMetricsEncoder(encoder)
	encoder = NewRetryEncoder(encoder, 3)

	if err := encoder.Encode(file); err != nil {
		fmt.Println("encoding failed:", err)
	}

	// videoProcessor := NewVideoProcessor(encoder)
	// if err = videoProcessor.Process(file); err != nil {
	// 	fmt.Println("failed to process file:", err)
	// }
}
