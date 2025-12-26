package config

import (
    "os"
    "time"

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
    JWT struct {
        Secret string       `yaml:"secret"`
        Expire yamlDuration `yaml:"expire"`
    } `yaml:"jwt"`
}

var Conf *Config

type yamlDuration struct{ D int64 }

func (y *yamlDuration) UnmarshalYAML(value *yaml.Node) error {
    var s string
    if err := value.Decode(&s); err == nil {
        d, err := time.ParseDuration(s)
        if err != nil {
            return err
        }
        y.D = int64(d)
        return nil
    }
    var i int64
    if err := value.Decode(&i); err == nil {
        y.D = i
        return nil
    }
    return nil
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
    Conf = &c
    return &c, nil
}
