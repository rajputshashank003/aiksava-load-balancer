package models

import "time"

type Session struct {
	ID        string
	Backend   string
	LastSeen time.Time
}

type Backend struct {
	URL         string
	ActiveUsers int
}

type RoundRobin struct {
	Count  int
}