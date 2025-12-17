package config

import (
    "os"
    "gopkg.in/yaml.v3"
)

type Config struct {
    Server struct {
        Addr string `yaml:"addr"`
    } `yaml:"server"`
    DB struct {
        DSN string `yaml:"dsn"`
    } `yaml:"db"`
    Redis struct {
        Addr     string `yaml:"addr"`
        Password string `yaml:"password"`
        DB       int    `yaml:"db"`
    } `yaml:"redis"`
    Logs struct {
        Path string `yaml:"path"`
    } `yaml:"logs"`
}

func Load(path string) (*Config, error) {
    b, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var c Config
    if err := yaml.Unmarshal(b, &c); err != nil {
        return nil, err
    }
    return &c, nil
}
