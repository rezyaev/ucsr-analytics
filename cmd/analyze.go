package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"sync"

	dem "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/events"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
)

type LostCountByBombsite map[string]int
type Stats map[string]LostCountByBombsite

func AnalyzeDemos() error {
	var waitGroup sync.WaitGroup
	errorGroup, _ := errgroup.WithContext(context.Background())

	mapFolders, err := os.ReadDir("../demos")
	if err != nil {
		return err
	}

	channel := make(chan Stats)
	for _, folder := range mapFolders {
		mapFiles, err := os.ReadDir("../demos/" + folder.Name())
		if err != nil {
			return err
		}

		for _, file := range mapFiles {
			func(folder fs.DirEntry, file fs.DirEntry) {
				waitGroup.Add(1)
				errorGroup.Go(func() error {
					return parseDemo("../demos/"+folder.Name()+"/"+file.Name(), folder.Name(), channel, &waitGroup)
				})
			}(folder, file)
		}
	}

	go func() {
		waitGroup.Wait()
		close(channel)
	}()

	combinedStats := createEmptyStats()
	for stats := range channel {
		for mapName, lostCountByBombsite := range stats {
			combinedStats[mapName]["A"] += lostCountByBombsite["A"]
			combinedStats[mapName]["B"] += lostCountByBombsite["B"]
		}
	}
	fmt.Printf("\nCombined Stats: %+v \n", combinedStats)

	return errorGroup.Wait()
}

func parseDemo(path string, mapName string, channel chan Stats, waitGroup *sync.WaitGroup) error {
	logProgress("start")

	demoFile, error := os.Open(path)
	if error != nil {
		return error
	}
	defer demoFile.Close()

	parser := dem.NewParser(demoFile)
	defer parser.Close()

	stats := createEmptyStats()
	parser.RegisterEventHandler(func(event events.BombPlanted) {
		isFriendly := slices.ContainsFunc(event.Player.TeamState.Members(), func(player *common.Player) bool {
			return player.Name == "Rezyaev"
		})

		if !isFriendly {
			stats[mapName][string(event.Site)] += 1
		}
	})

	parser.RegisterEventHandler(func(event events.AnnouncementWinPanelMatch) {
		channel <- stats
		waitGroup.Done()
		logProgress("finish")
	})

	return parser.ParseToEnd()
}

func createEmptyStats() Stats {
	return Stats{
		"de_ancient":  {"A": 0, "B": 0},
		"de_anubis":   {"A": 0, "B": 0},
		"de_dust2":    {"A": 0, "B": 0},
		"de_inferno":  {"A": 0, "B": 0},
		"de_mirage":   {"A": 0, "B": 0},
		"de_nuke":     {"A": 0, "B": 0},
		"de_overpass": {"A": 0, "B": 0},
		"de_vertigo":  {"A": 0, "B": 0},
	}
}

var parseProgress = map[string]int{"started": 0, "finished": 0}

func logProgress(event string) {
	switch event {
	case "start":
		parseProgress["started"] += 1
	case "finish":
		parseProgress["finished"] += 1
	}

	fmt.Printf("\rStarted parsing demo %v, Finished parsing demo %v", parseProgress["started"], parseProgress["finished"])
}
