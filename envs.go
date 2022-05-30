package envs

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type EnvConfig struct {
	Filepath string
	IsDebug  bool

	envs map[string]string
}

// Get returns an environment variable value as a string
func (c *EnvConfig) Get(key string, defaultValue ...string) string {
	value, ok := c.envs[key]
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

// GetBool returns an environment variable value as a boolean field
func (c *EnvConfig) GetBool(key string, defaultValue ...bool) bool {
	value, ok := c.envs[key]
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	res, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to parse key %s to boolean - ", key), err)
	}
	return res
}

// GetFloat returns an environment variable value as a float
func (c *EnvConfig) GetFloat(key string, defaultValue ...float32) float32 {
	value, ok := c.envs[key]
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	res, err := strconv.ParseFloat(value, 32)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to parse key %s to float - ", key), err)
	}
	return float32(res)
}

// GetInt returns an environment variable value as an integer
func (c *EnvConfig) GetInt(key string, defaultValue ...int) int {
	value, ok := c.envs[key]
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	res, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to parse key %s to integer - ", key), err)
	}
	return res
}

// GetMap returns an environment variable value as a map of strings
func (c *EnvConfig) GetMap(key string, defaultValue ...map[string]string) map[string]string {
	value, ok := c.envs[key]
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	pairs := strings.Split(value, ";")
	result := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		keyVal := strings.Split(pair, ":")
		result[keyVal[0]] = keyVal[1]
	}
	return result
}

// GetSlice returns an environment variable value as a slice of strings
func (c *EnvConfig) GetSlice(key string, defaultValue ...[]string) []string {
	value, ok := c.envs[key]
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return strings.Split(value, ",")
}

// GetSliceFloat returns an environment variable value as a slice of floats
func (c *EnvConfig) GetSliceFloat(key string, defaultValue ...[]float32) []float32 {
	value, ok := c.envs[key]
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	nums := strings.Split(value, ",")
	result := make([]float32, len(nums))
	for i, num := range nums {
		flo64, err := strconv.ParseFloat(num, 32)
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to parse key %s to slice of floats - ", key), err)
		}
		result[i] = float32(flo64)
	}
	return result
}

// GetSliceInt returns an environment variable value as a slice of integers
func (c *EnvConfig) GetSliceInt(key string, defaultValue ...[]int) []int {
	value, ok := c.envs[key]
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	nums := strings.Split(value, ",")
	var (
		err    error
		result = make([]int, len(nums))
	)
	for i, num := range nums {
		result[i], err = strconv.Atoi(num)
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to parse key %s to slice of integers - ", key), err)
		}
	}
	return result
}

// ReadEnvs obtains firstly from the file which path can be set as a struct field 'Filepath'
// and then from environment variables to rewrite got from the file
func (c *EnvConfig) ReadEnvs() {
	c.envs = map[string]string{}
	if c.Filepath == "" {
		c.Filepath = ".env"
	}
	file, err := os.Open(c.Filepath)
	if err != nil {
		if c.IsDebug {
			log.Println(fmt.Sprintf("Skip to read file %s, read all environment variables - ", c.Filepath), err)
		}
		for _, env := range os.Environ() {
			keyValuePair := strings.SplitN(env, "=", 2)
			c.envs[keyValuePair[0]] = keyValuePair[1]
		}
		return
	}
	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			log.Fatal("Failed to close file -", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if !strings.HasPrefix(text, "#") && text != "" {
			keyValuePair := strings.SplitN(text, "=", 2)
			c.envs[keyValuePair[0]] = keyValuePair[1]
		}
	}
	if err = scanner.Err(); err != nil {
		log.Fatal("Failed to scan file -", err)
	}

	for key := range c.envs {
		if value, exist := os.LookupEnv(key); exist {
			c.envs[key] = value
		}
	}
}
