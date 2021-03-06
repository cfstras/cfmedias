package util

import (
	"strconv"

	"github.com/cfstras/cfmedias/core"
	"github.com/cfstras/cfmedias/errrs"
)

// This file contains utils for fetching arguments used in the api.

// Fetches an argument from an ArgMap, used for single args
// Breaks and passes along the error given, if it is not nil.
// If the argument does not exist, the return value is nil.
// Parameter force can be used to return an error if the argument does not exist.
func GetArg(args core.ArgMap, arg string, force bool, err error) (*string, error) {
	if err != nil {
		return nil, err
	}
	value, ok := args[arg]
	if !ok || len(value) == 0 {
		if force {
			return nil, errrs.New("Argument '" + arg + "' missing!")
		}
		return nil, nil
	}
	if len(value) > 1 {
		return nil, errrs.New("argument " + arg + " cannot be supplied more than once!")
	}
	return &value[0], nil
}

// Converts a *string to a boolean.
// Passes along errrs, if not nil.
func CastBool(arg *string, err error) (*bool, error) {
	if err != nil {
		return nil, err
	}
	if arg == nil {
		return nil, nil
	}
	casted, err := strconv.ParseBool(*arg)
	if err != nil {
		return nil, errrs.New(*arg + " is not boolean")
	}
	return &casted, nil
}

// Converts a *string to a float32.
// Passes along errrs, if not nil.
func CastFloat(arg *string, err error) (*float32, error) {
	if err != nil {
		return nil, err
	}
	if arg == nil {
		return nil, nil
	}
	casted, err := strconv.ParseFloat(*arg, 32)
	if err != nil {
		return nil, errrs.New(*arg + " is not a floating point")
	}
	smaller := float32(casted)
	return &smaller, nil
}

// Converts a *string to a uint.
// Passes along errrs, if not nil.
func CastUint(arg *string, err error) (*uint, error) {
	if err != nil {
		return nil, err
	}
	if arg == nil {
		return nil, nil
	}
	casted, err := strconv.ParseUint(*arg, 10, 32)
	if err != nil {
		return nil, errrs.New(*arg + " is not an integer")
	}
	smaller := uint(casted)
	return &smaller, nil
}

// Converts a *string to a uint64.
// Passes along errrs, if not nil.
func CastInt64(arg *string, err error) (*int64, error) {
	if err != nil {
		return nil, err
	}
	if arg == nil {
		return nil, nil
	}
	casted, err := strconv.ParseInt(*arg, 10, 64)
	if err != nil {
		return nil, errrs.New(*arg + " is not an integer")
	}
	return &casted, nil
}

func Abs(i int64) int64 {
	if i < 0 {
		return -i
	}
	return i
}
