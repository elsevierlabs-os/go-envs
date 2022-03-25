package envs

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test_Get(t *testing.T) {
	envs := EnvConfig{}
	_ = os.Setenv("TEST_GET_ENVS", "get_envs")
	envs.ReadEnvs()
	if res := envs.Get("TEST_GET_ENVS"); res != "get_envs" {
		t.Error(fmt.Sprintf("expected 'get_envs', got '%s'", res))
	}
}

func Test_GetBool(t *testing.T) {
	envs := EnvConfig{}
	_ = os.Setenv("TEST_GET_ENVS_BOOL", "true")
	envs.ReadEnvs()
	if res := envs.GetBool("TEST_GET_ENVS_BOOL"); !res {
		t.Error(fmt.Sprintf("expected true, got %t", res))
	}
}

func Test_GetFloat(t *testing.T) {
	envs := EnvConfig{}
	_ = os.Setenv("TEST_GET_ENVS_FLOAT", "123.7645")
	envs.ReadEnvs()
	if res := envs.GetFloat("TEST_GET_ENVS_FLOAT"); res != float32(123.7645) {
		t.Error(fmt.Sprintf("expected 123.7645, got %f", res))
	}
}

func Test_GetInt(t *testing.T) {
	envs := EnvConfig{}
	_ = os.Setenv("TEST_GET_ENVS_INT", "123")
	envs.ReadEnvs()
	if res := envs.GetInt("TEST_GET_ENVS_INT"); res != 123 {
		t.Error(fmt.Sprintf("expected 123, got %d", res))
	}
}

func Test_GetMap(t *testing.T) {
	envs := EnvConfig{}
	_ = os.Setenv("TEST_GET_ENVS_MAP", "key1:value1;key2:value2")
	envs.ReadEnvs()
	res := envs.GetMap("TEST_GET_ENVS_MAP")
	if key := res["key1"]; key != "value1" {
		t.Error(fmt.Sprintf("expected 'value1', got '%s'", key))
	}
	if key := res["key2"]; key != "value2" {
		t.Error(fmt.Sprintf("expected 'value2', got '%s'", key))
	}
}

func Test_GetSlice(t *testing.T) {
	envs := EnvConfig{}
	_ = os.Setenv("TEST_GET_ENVS_SLICE", "first,second,third")
	envs.ReadEnvs()
	expected := []string{"first", "second", "third"}
	for i, value := range envs.GetSlice("TEST_GET_ENVS_SLICE") {
		if value != expected[i] {
			t.Error(fmt.Sprintf("expected '%s', got '%s'", expected[i], value))
		}
	}
}

func Test_GetSliceFloat(t *testing.T) {
	envs := EnvConfig{}
	_ = os.Setenv("TEST_GET_ENVS_SLICE_FLOAT", "123.554,64.23,78.09876")
	envs.ReadEnvs()
	expected := []float32{123.554, 64.23, 78.09876}
	for i, value := range envs.GetSliceFloat("TEST_GET_ENVS_SLICE_FLOAT") {
		if value != expected[i] {
			t.Error(fmt.Sprintf("expected %f, got %f", expected[i], value))
		}
	}
}

func Test_GetSliceInt(t *testing.T) {
	envs := EnvConfig{}
	_ = os.Setenv("TEST_GET_ENVS_SLICE_INT", "453,8,124567")
	envs.ReadEnvs()
	expected := []int{453, 8, 124567}
	for i, value := range envs.GetSliceInt("TEST_GET_ENVS_SLICE_INT") {
		if value != expected[i] {
			t.Error(fmt.Sprintf("expected %d, got %d", expected[i], value))
		}
	}
}

func Test_readEnvs(t *testing.T) {
	envs := EnvConfig{}

	{ // the .env file exists - envs are empty
		_ = ioutil.WriteFile(".env", []byte{}, 0644)
		envs.ReadEnvs()
		if res := envs.Get("TEST_READ_ENVS_1"); res != "" {
			t.Error(fmt.Sprintf("expected empty value, got %s", res))
		}
		_ = os.Remove(".env")
	}

	{ // the .env file exists - envs are read
		_ = ioutil.WriteFile(".env", []byte("# Comment\nTEST_READ_ENVS_2=read_envs_2"), 0644)
		envs.ReadEnvs()
		if res := envs.Get("TEST_READ_ENVS_2"); res != "read_envs_2" {
			t.Error(fmt.Sprintf("expected 'read_envs_2', got %s", res))
		}
		_ = os.Remove(".env")
	}

	{ // the .env file exists and the same env var was set - envs are read
		_ = os.Setenv("TEST_READ_ENVS_3", "read_envs_3_from_env?somePara=value")
		_ = ioutil.WriteFile(".env", []byte("# Comment\nTEST_READ_ENVS_3=read_envs_3"), 0644)
		envs.ReadEnvs()
		if res := envs.Get("TEST_READ_ENVS_3"); res != "read_envs_3_from_env?somePara=value" {
			t.Error(fmt.Sprintf("expected 'read_envs_3_from_env?somePara=value', got %s", res))
		}
		_ = os.Remove(".env")
	}

	{ // the .env file doesn't exists - envs are empty
		envs.ReadEnvs()
		if res := envs.Get("TEST_READ_ENVS_4"); res != "" {
			t.Error(fmt.Sprintf("expected empty value, got %s", res))
		}
	}

	{ // the .env file doesn't exists - envs are read
		_ = os.Setenv("TEST_READ_ENVS_5", "read_envs_5")
		envs.ReadEnvs()
		if res := envs.Get("TEST_READ_ENVS_5"); res != "read_envs_5" {
			t.Error(fmt.Sprintf("expected 'read_envs_5', got %s", res))
		}
	}
}
