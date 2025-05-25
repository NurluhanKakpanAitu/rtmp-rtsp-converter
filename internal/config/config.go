package config

import (
    "fmt"
    "os"
    "time"
    
    "gopkg.in/yaml.v2"
)

type Config struct {
    Server struct {
        Port int    `yaml:"port"`
        Host string `yaml:"host"`
    } `yaml:"server"`
    
    RTSP struct {
        Port int    `yaml:"port"`
        Host string `yaml:"host"`
    } `yaml:"rtsp"`
    
    Logging struct {
        Level  string `yaml:"level"`
        Format string `yaml:"format"`
    } `yaml:"logging"`
    
    Streams struct {
        MaxConcurrent int           `yaml:"max_concurrent"`
        Timeout       time.Duration `yaml:"timeout"`
    } `yaml:"streams"`
}

func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    
    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }
    
    return &config, nil
}