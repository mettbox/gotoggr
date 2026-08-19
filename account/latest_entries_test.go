package account

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func timeEntry(date string, pid int, description string, duration int64) TimeEntry {
	start, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}

	return TimeEntry{Pid: pid, Description: description, Duration: duration, Start: &start}
}

func TestBuildDayReportsOrdersDaysChronologically(t *testing.T) {
	account := Account{TimeEntries: []TimeEntry{
		timeEntry("2026-08-19", 1, "a", 3600),
		timeEntry("2026-08-10", 1, "b", 3600),
		timeEntry("2026-08-14", 1, "c", 3600),
		timeEntry("2026-08-11", 1, "d", 3600),
		timeEntry("2026-08-18", 1, "e", 3600),
		timeEntry("2026-08-12", 1, "f", 3600),
		timeEntry("2026-08-17", 1, "g", 3600),
		timeEntry("2026-08-13", 1, "h", 3600),
	}}

	want := []string{
		"2026-08-10", "2026-08-11", "2026-08-12", "2026-08-13",
		"2026-08-14", "2026-08-17", "2026-08-18", "2026-08-19",
	}

	reports := buildDayReports(account)

	if len(reports) != len(want) {
		t.Fatalf("got %d days, want %d", len(reports), len(want))
	}

	for index, date := range want {
		if reports[index].Date != date {
			t.Errorf("day %d: got %q, want %q", index, reports[index].Date, date)
		}
	}
}

func TestBuildDayReportsMergesEntriesOfSameProject(t *testing.T) {
	account := Account{TimeEntries: []TimeEntry{
		timeEntry("2026-08-10", 7, "ALL-7915", 7200),
		timeEntry("2026-08-10", 7, "Deployment umfangreicher Fixes", 3600),
	}}

	reports := buildDayReports(account)

	if len(reports) != 1 || len(reports[0].Entries) != 1 {
		t.Fatalf("got %d days with %v entries, want 1 day with 1 entry", len(reports), reports)
	}

	entry := reports[0].Entries[0]
	if entry.Duration != 10800 {
		t.Errorf("got duration %d, want 10800", entry.Duration)
	}

	want := "ALL-7915 // Deployment umfangreicher Fixes"
	if entry.Description != want {
		t.Errorf("got description %q, want %q", entry.Description, want)
	}

	if reports[0].Total != 10800 {
		t.Errorf("got total %d, want 10800", reports[0].Total)
	}
}

func TestBuildDayReportsKeepsDuplicateDescriptionOnce(t *testing.T) {
	account := Account{TimeEntries: []TimeEntry{
		timeEntry("2026-08-10", 4, "Daily", 900),
		timeEntry("2026-08-10", 4, "Daily", 900),
	}}

	reports := buildDayReports(account)

	entry := reports[0].Entries[0]
	if entry.Description != "Daily" {
		t.Errorf("got description %q, want %q", entry.Description, "Daily")
	}

	if entry.Duration != 1800 {
		t.Errorf("got duration %d, want 1800", entry.Duration)
	}
}

func TestBuildDayReportsSeparatesEntriesOfDifferentDays(t *testing.T) {
	account := Account{TimeEntries: []TimeEntry{
		timeEntry("2026-08-10", 7, "ALL-7915", 7200),
		timeEntry("2026-08-11", 7, "ALL-7915", 3600),
	}}

	reports := buildDayReports(account)

	if len(reports) != 2 {
		t.Fatalf("got %d days, want 2", len(reports))
	}

	if reports[0].Total != 7200 || reports[1].Total != 3600 {
		t.Errorf("got totals %d and %d, want 7200 and 3600", reports[0].Total, reports[1].Total)
	}
}

func TestBuildDayReportsKeepsProjectlessEntriesSeparate(t *testing.T) {
	account := Account{TimeEntries: []TimeEntry{
		timeEntry("2026-08-13", 0, "Developer Experience Exchange", 1800),
		timeEntry("2026-08-13", 0, "Monthly All-Hands", 3600),
	}}

	reports := buildDayReports(account)

	if len(reports[0].Entries) != 2 {
		t.Fatalf("got %d entries, want 2 (projectless entries must not merge)", len(reports[0].Entries))
	}
}

func TestBuildDayReportsOrdersEntriesByDurationDescending(t *testing.T) {
	account := Account{TimeEntries: []TimeEntry{
		timeEntry("2026-08-10", 1, "shortest", 900),
		timeEntry("2026-08-10", 2, "longest", 28800),
		timeEntry("2026-08-10", 3, "middle", 7200),
		timeEntry("2026-08-10", 4, "tied first", 3600),
		timeEntry("2026-08-10", 5, "tied second", 3600),
	}}

	want := []string{"longest", "middle", "tied first", "tied second", "shortest"}

	reports := buildDayReports(account)

	if len(reports[0].Entries) != len(want) {
		t.Fatalf("got %d entries, want %d", len(reports[0].Entries), len(want))
	}

	for index, description := range want {
		if reports[0].Entries[index].Description != description {
			t.Errorf("entry %d: got %q, want %q", index, reports[0].Entries[index].Description, description)
		}
	}
}

func TestBuildDayReportsSkipsEntriesWithoutStart(t *testing.T) {
	account := Account{TimeEntries: []TimeEntry{
		{Pid: 1, Description: "no start", Duration: 3600},
		timeEntry("2026-08-10", 1, "with start", 3600),
	}}

	reports := buildDayReports(account)

	if len(reports) != 1 {
		t.Fatalf("got %d days, want 1", len(reports))
	}

	if len(reports[0].Entries) != 1 || reports[0].Entries[0].Description != "with start" {
		t.Errorf("got entries %v, want only the entry with a start time", reports[0].Entries)
	}
}

func TestBuildDayReportsResolvesProjectAndClientNames(t *testing.T) {
	account := Account{
		Workspaces: []Workspace{{ID: 9, Name: "RMH"}},
		Clients:    []Client{{ID: 3, Name: "Merck"}},
		Projects:   []Project{{ID: 7, Cid: 3, Name: "P25-244: Merck Quizz App Retainer 12.25 - 26"}},
		TimeEntries: []TimeEntry{
			func() TimeEntry {
				entry := timeEntry("2026-08-10", 7, "ALL-7915", 7200)
				entry.Wid = 9

				return entry
			}(),
		},
	}

	entry := buildDayReports(account)[0].Entries[0]

	if entry.Project != "P25-244: Merck Quizz App Retainer 12.25 - 26" {
		t.Errorf("got project %q, want the resolved project name", entry.Project)
	}

	if entry.Client != "Merck" {
		t.Errorf("got client %q, want %q", entry.Client, "Merck")
	}

	if entry.Workspace != "RMH" {
		t.Errorf("got workspace %q, want %q", entry.Workspace, "RMH")
	}
}

func TestRenderDayReportsAlignsColumns(t *testing.T) {
	reports := []DayReport{{
		Date: "2026-08-17",
		Entries: []DetailedEntry{
			{Project: "P26-108: BIO-DE Apothekenfolder", Description: "ALL-7938", Duration: 28800},
			{Project: "", Description: "Developer Experience Exchange", Duration: 1800},
		},
		Total: 30600,
	}}

	want := "\n2026-08-17\n" +
		"  08:00  P26-108: BIO-DE Apothekenfolder  ALL-7938\n" +
		"  00:30  " + strings.Repeat(" ", 33) + "Developer Experience Exchange\n" +
		"  Total  08:30\n" +
		"\n"

	output := &bytes.Buffer{}
	renderDayReports(output, reports)

	if output.String() != want {
		t.Errorf("got:\n%s\nwant:\n%s", output.String(), want)
	}
}

func TestRenderDayReportsSizesColumnsPerDay(t *testing.T) {
	reports := []DayReport{
		{
			Date:    "2026-08-10",
			Entries: []DetailedEntry{{Project: "P22-544/T12: BIHP-ALL: Escalation of IS-2243", Description: "ALL-7836", Duration: 3600}},
			Total:   3600,
		},
		{
			Date:    "2026-08-11",
			Entries: []DetailedEntry{{Project: "P1", Description: "ALL-1", Duration: 3600}},
			Total:   3600,
		},
	}

	output := &bytes.Buffer{}
	renderDayReports(output, reports)

	want := "  01:00  P1  ALL-1\n"
	if !strings.Contains(output.String(), want) {
		t.Errorf("narrow day was padded to the wide day's width, got:\n%s", output.String())
	}
}

func TestAppendDescription(t *testing.T) {
	tests := []struct {
		name        string
		merged      string
		description string
		want        string
	}{
		{name: "first description", merged: "", description: "ALL-1", want: "ALL-1"},
		{name: "second description", merged: "ALL-1", description: "Fixes", want: "ALL-1 // Fixes"},
		{name: "duplicate description", merged: "ALL-1", description: "ALL-1", want: "ALL-1"},
		{name: "duplicate of merged part", merged: "ALL-1 // Fixes", description: "Fixes", want: "ALL-1 // Fixes"},
		{name: "empty description", merged: "ALL-1", description: "", want: "ALL-1"},
		{name: "empty description on empty", merged: "", description: "", want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := appendDescription(test.merged, test.description)
			if got != test.want {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}
