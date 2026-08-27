package onnx

import (
	"fmt"
	"os"
	"runtime"

	onnxruntime "github.com/yalue/onnxruntime_go"

	"github.com/photoprism/photoprism/internal/event"
)

var log = event.Log

// ExecutionProviderType defines the type of execution provider to use.
type ExecutionProviderType string

const (
	// ExecutionProviderCPU uses CPU execution (default).
	ExecutionProviderCPU ExecutionProviderType = "cpu"

	// ExecutionProviderCUDA uses NVIDIA CUDA execution provider.
	ExecutionProviderCUDA ExecutionProviderType = "cuda"

	// ExecutionProviderOpenVINO uses Intel OpenVINO execution provider.
	ExecutionProviderOpenVINO ExecutionProviderType = "openvino"

	// ExecutionProviderDirectML uses DirectML execution provider (Windows).
	ExecutionProviderDirectML ExecutionProviderType = "directml"

	// ExecutionProviderTensorRT uses NVIDIA TensorRT execution provider.
	ExecutionProviderTensorRT ExecutionProviderType = "tensorrt"

	// ExecutionProviderAuto automatically selects best available provider.
	ExecutionProviderAuto ExecutionProviderType = "auto"
)

// ExecutionProviderConfig holds configuration for ONNX execution providers.
type ExecutionProviderConfig struct {
	Type       ExecutionProviderType
	DeviceID   int
	Enabled    bool
	FallbackCPU bool
}

// DefaultExecutionProviderConfig returns default execution provider configuration.
func DefaultExecutionProviderConfig() ExecutionProviderConfig {
	return ExecutionProviderConfig{
		Type:       ExecutionProviderCPU,
		DeviceID:   0,
		Enabled:    false,
		FallbackCPU: true,
	}
}

// ConfigureExecutionProvider attempts to configure GPU acceleration for an ONNX session.
// It tries the specified provider and falls back to CPU if unavailable or disabled.
func ConfigureExecutionProvider(sessionOpts *onnxruntime.SessionOptions, config ExecutionProviderConfig) error {
	if !config.Enabled || config.Type == ExecutionProviderCPU {
		// CPU execution, no special configuration needed
		return nil
	}

	// Determine which provider to use
	provider := config.Type
	if provider == ExecutionProviderAuto {
		provider = detectBestExecutionProvider()
	}

	// Try to append the execution provider
	var err error
	switch provider {
	case ExecutionProviderCUDA:
		err = appendCUDAExecutionProvider(sessionOpts, config.DeviceID)
	case ExecutionProviderOpenVINO:
		err = appendOpenVINOExecutionProvider(sessionOpts)
	case ExecutionProviderDirectML:
		err = appendDirectMLExecutionProvider(sessionOpts, config.DeviceID)
	case ExecutionProviderTensorRT:
		err = appendTensorRTExecutionProvider(sessionOpts, config.DeviceID)
	default:
		return fmt.Errorf("unsupported execution provider: %s", provider)
	}

	if err != nil {
		if config.FallbackCPU {
			log.Warnf("onnx: %s execution provider unavailable, falling back to CPU: %v", provider, err)
			return nil
		}
		return fmt.Errorf("onnx: failed to configure %s execution provider: %w", provider, err)
	}

	log.Infof("onnx: using %s execution provider (device %d)", provider, config.DeviceID)
	return nil
}

// detectBestExecutionProvider automatically detects the best available execution provider.
func detectBestExecutionProvider() ExecutionProviderType {
	// Check for NVIDIA GPU
	if _, err := os.Stat("/dev/nvidia0"); err == nil {
		return ExecutionProviderCUDA
	}

	// Check for Intel GPU (DRI devices)
	if _, err := os.Stat("/dev/dri/renderD128"); err == nil {
		// Intel GPU available, check if OpenVINO is available
		if runtime.GOOS == "linux" {
			return ExecutionProviderOpenVINO
		}
	}

	// Check for DirectML on Windows
	if runtime.GOOS == "windows" {
		return ExecutionProviderDirectML
	}

	// Fallback to CPU
	return ExecutionProviderCPU
}

// appendCUDAExecutionProvider adds NVIDIA CUDA execution provider.
func appendCUDAExecutionProvider(sessionOpts *onnxruntime.SessionOptions, deviceID int) error {
	// Create CUDA provider options
	cudaOptions := &onnxruntime.CUDAProviderOptions{
		DeviceID: deviceID,
	}

	// Try to append CUDA execution provider
	if err := sessionOpts.AppendExecutionProviderCUDA(cudaOptions); err != nil {
		return fmt.Errorf("failed to append CUDA execution provider: %w (ensure ONNX Runtime is built with CUDA support)", err)
	}

	return nil
}

// appendOpenVINOExecutionProvider adds Intel OpenVINO execution provider.
func appendOpenVINOExecutionProvider(sessionOpts *onnxruntime.SessionOptions) error {
	// OpenVINO options - use defaults for now
	// Can be configured with device_type, precision, etc.
	options := map[string]string{
		"device_type": "CPU_FP32", // Can be CPU_FP32, GPU_FP32, GPU_FP16, etc.
	}

	// Try to append OpenVINO execution provider
	if err := sessionOpts.AppendExecutionProviderOpenVINO(options); err != nil {
		return fmt.Errorf("failed to append OpenVINO execution provider: %w (ensure ONNX Runtime is built with OpenVINO support)", err)
	}

	return nil
}

// appendDirectMLExecutionProvider adds DirectML execution provider (Windows).
func appendDirectMLExecutionProvider(sessionOpts *onnxruntime.SessionOptions, deviceID int) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("DirectML is only available on Windows")
	}

	// Try to append DirectML execution provider
	if err := sessionOpts.AppendExecutionProviderDirectML(deviceID); err != nil {
		return fmt.Errorf("failed to append DirectML execution provider: %w", err)
	}

	return nil
}

// appendTensorRTExecutionProvider adds NVIDIA TensorRT execution provider.
func appendTensorRTExecutionProvider(sessionOpts *onnxruntime.SessionOptions, deviceID int) error {
	// Create TensorRT provider options
	tensorRTOptions := &onnxruntime.TensorRTProviderOptions{
		DeviceID: deviceID,
	}

	// Try to append TensorRT execution provider
	if err := sessionOpts.AppendExecutionProviderTensorRT(tensorRTOptions); err != nil {
		return fmt.Errorf("failed to append TensorRT execution provider: %w (ensure ONNX Runtime is built with TensorRT support)", err)
	}

	return nil
}

// IsGPUAvailable checks if any GPU is available for ONNX inference.
func IsGPUAvailable() bool {
	// Check for NVIDIA GPU
	if _, err := os.Stat("/dev/nvidia0"); err == nil {
		return true
	}

	// Check for Intel/AMD GPU (DRI devices)
	if _, err := os.Stat("/dev/dri/renderD128"); err == nil {
		return true
	}

	// Check for other GPU indicators
	if runtime.GOOS == "windows" {
		// Windows typically has GPU available
		return true
	}

	return false
}

// GetAvailableExecutionProviders returns a list of execution providers that might be available.
func GetAvailableExecutionProviders() []ExecutionProviderType {
	providers := []ExecutionProviderType{ExecutionProviderCPU}

	if _, err := os.Stat("/dev/nvidia0"); err == nil {
		providers = append(providers, ExecutionProviderCUDA, ExecutionProviderTensorRT)
	}

	if _, err := os.Stat("/dev/dri/renderD128"); err == nil {
		providers = append(providers, ExecutionProviderOpenVINO)
	}

	if runtime.GOOS == "windows" {
		providers = append(providers, ExecutionProviderDirectML)
	}

	return providers
}
