package account

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/mettbox/gotoggr/util"
)

// descriptionSeparator joins the descriptions of entries merged into one project line.
const descriptionSeparator = " // "

func getWorkspace(id int, workspaces []Workspace) string {
	for _, workspace := range workspaces {
		if workspace.ID == id {
			return workspace.Name
		}
	}

	return ""
}

func getProject(id int, projects []Project) (projectName string, clientID int) {
	for _, project := range projects {
		if project.ID == id {
			projectName = project.Name
			clientID = project.Cid
		}
	}

	return
}

func getClient(id int, clients []Client) string {
	for _, client := range clients {
		if client.ID == id {
			return client.Name
		}
	}

	return ""
}

// indexOfProject finds the entry already collected for a project, or -1.
// Entries without a project (pid 0) never merge, as they are unrelated to each other.
func indexOfProject(entries []DetailedEntry, pid int) int {
	if pid == 0 {
		return -1
	}

	for index, entry := range entries {
		if entry.Pid == pid {
			return index
		}
	}

	return -1
}

// appendDescription joins a description onto those already merged for a project,
// skipping descriptions that are present already.
func appendDescription(merged string, description string) string {
	if description == "" {
		return merged
	}

	for _, existing := range strings.Split(merged, descriptionSeparator) {
		if existing == description {
			return merged
		}
	}

	if merged == "" {
		return description
	}

	return merged + descriptionSeparator + description
}

// buildDayReports merges time entries into one entry per project per day, ordered
// by date ascending and, within a day, by duration descending.
func buildDayReports(account Account) []DayReport {
	entriesByDate := map[string][]DetailedEntry{}

	for _, entry := range account.TimeEntries {
		if entry.Start == nil {
			continue
		}

		date := entry.Start.Format("2006-01-02")

		if index := indexOfProject(entriesByDate[date], entry.Pid); index >= 0 {
			merged := &entriesByDate[date][index]
			merged.Duration += entry.Duration
			merged.Description = appendDescription(merged.Description, entry.Description)

			continue
		}

		projectName, clientID := getProject(entry.Pid, account.Projects)
		entriesByDate[date] = append(entriesByDate[date], DetailedEntry{
			Workspace:   getWorkspace(entry.Wid, account.Workspaces),
			Pid:         entry.Pid,
			Project:     projectName,
			Cid:         clientID,
			Client:      getClient(clientID, account.Clients),
			Description: entry.Description,
			Duration:    entry.Duration,
		})
	}

	dates := make([]string, 0, len(entriesByDate))
	for date := range entriesByDate {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	reports := make([]DayReport, 0, len(dates))
	for _, date := range dates {
		entries := entriesByDate[date]
		sort.SliceStable(entries, func(i int, j int) bool {
			return entries[i].Duration > entries[j].Duration
		})

		total := int64(0)
		for _, entry := range entries {
			total += entry.Duration
		}

		reports = append(reports, DayReport{Date: date, Entries: entries, Total: total})
	}

	return reports
}

// renderDayReports writes one aligned table per day. Every day gets its own
// tabwriter, so column widths follow that day's content rather than the widest
// row of the whole report.
func renderDayReports(w io.Writer, reports []DayReport) {
	for _, report := range reports {
		fmt.Fprintf(w, "\n%s\n", report.Date)

		table := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		for _, entry := range report.Entries {
			fmt.Fprintf(table, "  %s\t%s\t%s\n", util.SecondsToHHMM(entry.Duration), entry.Project, entry.Description)
		}
		fmt.Fprintf(table, "  Total\t%s\n", util.SecondsToHHMM(report.Total))
		table.Flush()
	}

	fmt.Fprintln(w)
}

func LatestEntries() {
	account, err := Get()
	if err != nil {
		panic(err)
	}

	renderDayReports(os.Stdout, buildDayReports(account))
}
