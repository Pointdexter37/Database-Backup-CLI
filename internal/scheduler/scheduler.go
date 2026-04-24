package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type JobConfig struct {
	ID         string `json:"id"`
	CronExpr   string `json:"cronExpr"`
	DBType     string `json:"dbType"`
	DBName     string `json:"dbName"`
	DBHost     string `json:"dbHost"`
	DBPort     int    `json:"dbPort"`
	DBUser     string `json:"dbUser"`
	DBPassword string `json:"dbPassword"`
	Storage    string `json:"storage"`
	OutputPath string `json:"outputPath"`
	S3Bucket   string `json:"s3Bucket"`
	S3Region   string `json:"s3Region"`
	Compress   bool   `json:"compress"`
}

type Manager struct {
	configFile string
	mu         sync.RWMutex
}

func NewManager(configFile string) *Manager {
	return &Manager{
		configFile: configFile,
	}
}

func (m *Manager) LoadJobs() ([]JobConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, err := os.ReadFile(m.configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []JobConfig{}, nil
		}
		return nil, fmt.Errorf("failed to read schedule config: %w", err)
	}

	var jobs []JobConfig
	if err := json.Unmarshal(data, &jobs); err != nil {
		return nil, fmt.Errorf("failed to parse schedule config: %w", err)
	}
	return jobs, nil
}

func (m *Manager) SaveJobs(jobs []JobConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode schedule config: %w", err)
	}

	return os.WriteFile(m.configFile, data, 0600)
}

func (m *Manager) AddJob(job JobConfig) error {
	jobs, err := m.LoadJobs()
	if err != nil {
		return err
	}
	jobs = append(jobs, job)
	return m.SaveJobs(jobs)
}

func (m *Manager) RemoveJob(id string) error {
	jobs, err := m.LoadJobs()
	if err != nil {
		return err
	}

	var newJobs []JobConfig
	for _, j := range jobs {
		if j.ID != id {
			newJobs = append(newJobs, j)
		}
	}

	if len(newJobs) == len(jobs) {
		return fmt.Errorf("job not found: %s", id)
	}

	return m.SaveJobs(newJobs)
}
