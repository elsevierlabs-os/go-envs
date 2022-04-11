package envs

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type EnvConfig map[string]string

// Get returns an environment variable value as a string
func (c EnvConfig) Get(key string, defaultValue ...string) string {
	value, ok := c[key]
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

// GetBool returns an environment variable value as a boolean field
func (c EnvConfig) GetBool(key string, defaultValue ...bool) bool {
	value, ok := c[key]
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
func (c EnvConfig) GetFloat(key string, defaultValue ...float32) float32 {
	value, ok := c[key]
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
func (c EnvConfig) GetInt(key string, defaultValue ...int) int {
	value, ok := c[key]
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
func (c EnvConfig) GetMap(key string, defaultValue ...map[string]string) map[string]string {
	value, ok := c[key]
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
func (c EnvConfig) GetSlice(key string, defaultValue ...[]string) []string {
	value, ok := c[key]
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return strings.Split(value, ",")
}

// GetSliceFloat returns an environment variable value as a slice of floats
func (c EnvConfig) GetSliceFloat(key string, defaultValue ...[]float32) []float32 {
	value, ok := c[key]
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
func (c EnvConfig) GetSliceInt(key string, defaultValue ...[]int) []int {
	value, ok := c[key]
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

// ReadEnvs obtains firstly from the file .env and then from environment variables to rewrite got from the file
func (c EnvConfig) ReadEnvs() {
	file, err := os.Open(".env")
	if err != nil {
		log.Println("Failed to read file, read from all environment variables -", err)
		for _, env := range os.Environ() {
			keyValuePair := strings.SplitN(env, "=", 2)
			c[keyValuePair[0]] = keyValuePair[1]
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
			c[keyValuePair[0]] = keyValuePair[1]
		}
	}
	if err = scanner.Err(); err != nil {
		log.Fatal("Failed to scan file -", err)
	}

	for key := range c {
		if value, exist := os.LookupEnv(key); exist {
			c[key] = value
		}
	}
}
