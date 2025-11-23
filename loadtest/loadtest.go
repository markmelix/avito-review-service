package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tsenart/vegeta/v12/lib"
)

const (
	api            = "http://api:8080"
	targetFile     = "/target.list"
	resultDir      = "/result"
	resultFile     = resultDir + "/vegeta_results.bin"
	durationString = "300s"
	rps            = 5
	teamAmount     = 20
	usersPerTeam   = 10
	pullreqsAmount = 1000
)

func createTeams() {
	for i := 1; i <= teamAmount; i++ {
		teamName := fmt.Sprintf("team-%d", i)
		members := []map[string]any{}

		for u := 1; u <= usersPerTeam; u++ {
			members = append(members, map[string]any{
				"user_id":   fmt.Sprintf("u-%d-%d", i, u),
				"username":  fmt.Sprintf("User-%d-%d", i, u),
				"is_active": rand.Intn(5) != 0,
			})
		}

		body := map[string]any{
			"team_name": teamName,
			"members":   members,
		}

		b, err := json.Marshal(body)
		if err != nil {
			panic(fmt.Sprintf("Error marshalling data generated: %v\n", err))
		}

		_, err = http.Post(api+"/team/add", "application/json", bytes.NewReader(b))
		if err != nil {
			panic(fmt.Sprintf("Error making /team/data request: %v\n", err))
		}
	}
}

func createPullReqs() {
	for i := 1; i <= pullreqsAmount; i++ {
		body := map[string]any{
			"pull_request_id":   fmt.Sprintf("pr-%d", i),
			"pull_request_name": "LoadTest",
			"author_id":         fmt.Sprintf("u-%d-%d", rand.Intn(20)+1, rand.Intn(10)+1),
		}

		b, err := json.Marshal(body)
		if err != nil {
			panic(fmt.Sprintf("Error marshalling data generated: %v\n", err))
		}

		_, err = http.Post(api+"/pullRequest/create", "application/json", bytes.NewReader(b))
		if err != nil {
			panic(fmt.Sprintf("Error making /pullRequest/create request: %v\n", err))
		}
	}
}

func writeVegetTargetFile() error {
	target := make([]string, pullreqsAmount)

	bodyDir := "reqbodies"
	err := os.Mkdir(bodyDir, 0755)
	if err != nil {
		panic(fmt.Sprintf("Error making reqbodies dir: %v\n", err))
	}

	for i := 1; i <= pullreqsAmount; i++ {
		body := map[string]any{
			"pull_request_id": fmt.Sprintf("pr-%d", i),
		}

		b, err := json.Marshal(body)
		if err != nil {
			panic(fmt.Sprintf("Error marshalling data generated: %v\n", err))
		}

		filename := fmt.Sprintf("%s/merge_%06d.json", bodyDir, i)
		err = os.WriteFile(filename, b, 0644)
		if err != nil {
			panic(fmt.Sprintf("Error writing json %s: %v\n", filename, err))
		}

		endpoint := api + "/pullRequest/merge"

		target = append(target, fmt.Sprintf("POST %s\nContent-Type: application/json\n@%s\n\n", endpoint, filename))
	}

	os.WriteFile(targetFile, []byte(strings.Join(target, "")), 0644)

	return nil
}

func runVegeta() {
	file, err := os.Open(targetFile)
	if err != nil {
		panic(fmt.Sprintf("Error opening target.list file: %v", err))
	}
	defer file.Close()

	targeter := vegeta.NewHTTPTargeter(file, nil, nil)

	rate := vegeta.Rate{Freq: rps, Per: time.Second}
	duration, err := time.ParseDuration(durationString)
	if err != nil {
		panic(fmt.Sprintf("Failed parsing duration string: %v", duration))
	}

	attacker := vegeta.NewAttacker()

	err = os.Mkdir(resultDir, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		panic(fmt.Sprintf("Error making %s dir: %v", resultDir, err))
	}

	results, err := os.Create(resultFile)
	if err != nil {
		panic(err)
	}
	defer results.Close()

	encoder := vegeta.NewEncoder(results)

	var metrics vegeta.Metrics

	for res := range attacker.Attack(targeter, rate, duration, "team+pr-test") {
		metrics.Add(res)

		if err := encoder.Encode(res); err != nil {
			panic(fmt.Sprintf("failed encoding vegeta result: %v", err))
		}
	}

	fmt.Printf("Saved results binary to %s\n", resultFile)

	metrics.Close()

	fmt.Println("Report:")
	fmt.Printf("Duration: %ds\n", duration/time.Second)
	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Rate: %.2f rps\n", metrics.Rate)
	fmt.Printf("SLI Success Rate: %.2f%%\n", metrics.Success*100)
	fmt.Printf("SLI Latency (p95): %s\n", metrics.Latencies.P95)
	fmt.Printf("Errors: %v\n", metrics.Errors)
}

func main() {
	fmt.Printf("Creating teams (%d; %d user per each)\n", teamAmount, usersPerTeam)
	createTeams()

	fmt.Printf("Creating pull requests (%d)\n", pullreqsAmount)
	createPullReqs()

	fmt.Println("Writing request bodies")
	writeVegetTargetFile()

	fmt.Println("Running vegeta loadtest (duration: %s)", durationString)
	runVegeta()

	for {
	}
}
