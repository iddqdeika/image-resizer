package cfg

import (
	"os"
	"testing"
)

func TestEnvCfg_StringWithDefaults(t *testing.T) {
	os.Setenv("key", "1")
	if EnvCfg.StringWithDefaults("key", "2") != "1"{
		t.Errorf("must return variable if it is set")
	}
	if EnvCfg.StringWithDefaults("d323gh-3gd4k-ddsa3", "1") != "1"{
		t.Errorf("must return default value if variable does not exists")
	}
}

func TestEnvCfg_IntWithDefaults(t *testing.T) {
	os.Setenv("key", "1")
	if EnvCfg.IntWithDefaults("key", 2) != 1{
		t.Errorf("must return variable if it is set")
	}
	if EnvCfg.IntWithDefaults("d323gh-3gd4k-ddsa3", 1) != 1{
		t.Errorf("must return default value if variable does not exists")
	}
	os.Setenv("key1", "1s")
	if EnvCfg.IntWithDefaults("key1", 2) != 2{
		t.Errorf("must return default value if variable is incorrect")
	}
}

func TestMapCfg_StringWithDefaults(t *testing.T) {
	cfg := MapCfg(map[string]string{"key":"1"})
	if cfg.StringWithDefaults("key", "2") != "1"{
		t.Errorf("must return variable if it is set")
	}
	if cfg.StringWithDefaults("d323gh-3gd4k-ddsa3", "1") != "1"{
		t.Errorf("must return default value if variable does not exists")
	}
}

func TestMapCfg_IntWithDefaults(t *testing.T) {
	cfg := MapCfg(map[string]string{"key":"1"})
	if cfg.IntWithDefaults("key", 2) != 1{
		t.Errorf("must return variable if it is set")
	}
	if cfg.IntWithDefaults("d323gh-3gd4k-ddsa3", 1) != 1{
		t.Errorf("must return default value if variable does not exists")
	}
	cfg = MapCfg(map[string]string{"key":"1"})
	if cfg.IntWithDefaults("key1", 2) != 2{
		t.Errorf("must return default value if variable is incorrect")
	}
}