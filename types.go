package main

import (
	"github.com/drone/drone-go/drone"
)

type Params struct {
	Commands []string          `json:"commands"`
	Login    string            `json:"user"`
	Port     int               `json:"port"`
	Host     drone.StringSlice `json:"host"`
	Sleep    int               `json:"sleep"`
}
