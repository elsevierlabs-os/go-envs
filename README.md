# go-envs
A package to parse environment variables.

### Getting started
To read environment variables initialize a local variable using `EnvConfig`, that is a type of map, and run the method `ReadEnvs` to fill the map.
```go
envVars := EnvConfig{}
envVars.ReadEnvs()
```
<ins>**NOTE**</ins>: All errors while parsing environment variables will call the program termination, that's why the variables should be read and used at launch of the application.

### How it works
- The top priority is for environment variables
- Firstly, the handler tries to read the file `.env` (key-value pairs)
- If the file `.env` exists, all obtained variables will be also reread from the environments and rewritten (if are set) - for example, if the file has the variable `ENV=dev` and the environments has the variable `ENV=prod`, the last value (`prod`) will be relevant
- If the file `.env` does not exist, the handler will read all set environment variables

### General features
To get value of an environment variable, use the following methods:
- `Get`, `GetSlice`, `GetMap` - to get the value as a string, slice of strings or map of strings (`map[string]string`) appropriately
- `GetInt`, `GetSliceInt` - to get the value as an integer or slice of integers appropriately
- `GetFloat`, `GetSliceFloat` - to get the value as a float or slice of floats appropriately
- `GetBool` - to get the value as a boolean

For example:
```go
envStr := envVars.Get("ENV_STRING")
envSliceInt := envVars.GetSliceInt("ENV_SLICE_OF_INT")
envBool := envVars.GetBool("ENV_BOOL")
```
Also, the default value for an environment variable can be set. This value will be used if the environment variable doesn't exist:
```go
envStr := envVars.Get("ENV_STRING", "someDefaultValue")
```

### Some requirements
To read some kind of variables correctly, it should be written (in the `.env` file or environments) according to the following rules:
- if the variable is an array, use comma `,` as a separator without any extra whitespaces, like `LIST=first,second,third`
- if the variable is a map, use semicolon `;` as a separator for pair key-value and colon `:` between the key and the value, also without any extra whitespaces like `MAP=key1:value1;key2:value2`