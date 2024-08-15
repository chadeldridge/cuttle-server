package tests

import "time"

type TestArg struct {
	Key   string
	Value any
}

func FindArg(args []TestArg, key string) any {
	if len(args) == 0 {
		return nil
	}

	for _, a := range args {
		if a.Key == key {
			return a.Value
		}
	}

	return nil
}

// Quiet returns a "quiet" TestArg with value set to true.
func Quiet() TestArg {
	return TestArg{
		Key:   "quiet",
		Value: true,
	}
}

// BeQuiet returns true if the "quiet" argument is set to true.
func BeQuiet(args []TestArg) bool {
	v := FindArg(args, "quiet")
	if v == nil {
		return false
	}

	return v.(bool)
}

func GetTimeout(args []TestArg, defaultTimeout time.Duration) time.Duration {
	v := FindArg(args, "timeout")
	if v == nil {
		return defaultTimeout
	}

	switch any(v).(type) {
	case int:
		if v.(int) == 0 {
			return defaultTimeout
		}

		return time.Second * time.Duration(v.(int))
	case int64:
		if v.(int64) == 0 {
			return defaultTimeout
		}

		return time.Second * time.Duration(v.(int64))
	case time.Duration:
		return v.(time.Duration)
	default:
		return defaultTimeout
	}
}
