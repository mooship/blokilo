package models

import "time"

type ResultStatus string

const (
	StatusBlocked  ResultStatus = "BLOCKED"
	StatusResolved ResultStatus = "RESOLVED"
	StatusError    ResultStatus = "ERROR"
)

type ClassifiedResult struct {
	Domain       string
	Status       ResultStatus
	ResponseTime time.Duration
	Err          error
}

func ClassifyResult(dnsStatus, httpStatus ResultStatus, dnsErr, httpErr error) ResultStatus {
	if dnsStatus == StatusBlocked || httpStatus == StatusBlocked {
		return StatusBlocked
	}
	if dnsStatus == StatusError || httpStatus == StatusError {
		return StatusError
	}
	if dnsStatus == StatusResolved || httpStatus == StatusResolved {
		return StatusResolved
	}
	return StatusError
}

type Stats struct {
	Total           int
	Blocked         int
	Resolved        int
	Errored         int
	PercentBlocked  float64
	PercentResolved float64
}

func ComputeStats(results []ClassifiedResult) Stats {
	var s Stats
	s.Total = len(results)
	for _, r := range results {
		switch r.Status {
		case StatusBlocked:
			s.Blocked++
		case StatusResolved:
			s.Resolved++
		case StatusError:
			s.Errored++
		}
	}
	if s.Total > 0 {
		s.PercentBlocked = float64(s.Blocked) / float64(s.Total) * 100
		s.PercentResolved = float64(s.Resolved) / float64(s.Total) * 100
	}
	return s
}
