package models

import (
	"testing"
)

func TestClassifyResult(t *testing.T) {
	cases := []struct {
		dns, http       ResultStatus
		dnsErr, httpErr error
		expect          ResultStatus
	}{
		{StatusBlocked, StatusBlocked, nil, nil, StatusBlocked},
		{StatusResolved, StatusBlocked, nil, nil, StatusBlocked},
		{StatusBlocked, StatusResolved, nil, nil, StatusBlocked},
		{StatusResolved, StatusResolved, nil, nil, StatusResolved},
		{StatusError, StatusBlocked, nil, nil, StatusBlocked},
		{StatusError, StatusError, nil, nil, StatusError},
	}
	for _, c := range cases {
		got := ClassifyResult(c.dns, c.http, c.dnsErr, c.httpErr)
		if got != c.expect {
			t.Errorf("ClassifyResult(%v, %v) = %v, want %v", c.dns, c.http, got, c.expect)
		}
	}
}

func TestComputeStats(t *testing.T) {
	results := []ClassifiedResult{
		{"a", StatusBlocked, 0, nil},
		{"b", StatusResolved, 0, nil},
		{"c", StatusError, 0, nil},
		{"d", StatusBlocked, 0, nil},
	}
	s := ComputeStats(results)
	if s.Total != 4 || s.Blocked != 2 || s.Resolved != 1 || s.Errored != 1 {
		t.Errorf("unexpected stats: %+v", s)
	}
	if s.PercentBlocked != 50.0 || s.PercentResolved != 25.0 {
		t.Errorf("unexpected percentages: blocked=%.1f resolved=%.1f", s.PercentBlocked, s.PercentResolved)
	}
}
