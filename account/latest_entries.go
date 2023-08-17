package account

import (
	"fmt"
	"sort"

	"github.com/mettbox/gotoggr/util"
)

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

func LatestEntries() {
	account, err := Get()
	if err != nil {
		panic(err)
	}

	days := map[string][]DetailedEntry{}

	for _, entry := range account.TimeEntries {
		date := string(entry.Start.Format("2006-01-02"))
		projectName, clientID := getProject(entry.Pid, account.Projects)

		exists := false
		for index, dayEntry := range days[date] {
			if dayEntry.Pid == entry.Pid && dayEntry.Description == entry.Description {
				days[date][index].Duration = days[date][index].Duration + entry.Duration
				exists = true
			}
		}

		if !exists {
			days[date] = append(days[date], DetailedEntry{
				Workspace:   getWorkspace(entry.Wid, account.Workspaces),
				Pid:         entry.Pid,
				Project:     projectName,
				Client:      getClient(clientID, account.Clients),
				Description: entry.Description,
				Duration:    entry.Duration,
			})
		}
	}

	keys := make([]string, 0, len(days))

	for k := range days {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	daysSorted := map[string][]DetailedEntry{}
	for _, k := range keys {
		daysSorted[k] = days[k]
	}

	for day, entries := range daysSorted {
		duration := 0
		fmt.Println()
		fmt.Println(day)
		for _, entry := range entries {
			duration = duration + int(entry.Duration)
			taskDuration := util.SecondsToHHMM(int64(entry.Duration))
			fmt.Printf("\t\t%s: %s > %s :: %s [%s] \n", entry.Workspace, entry.Client, entry.Project, entry.Description, taskDuration)
		}
		total := util.SecondsToHHMM(int64(duration))
		fmt.Printf("\t\tTotal: %s\n", total)
	}
	fmt.Println("")
}
