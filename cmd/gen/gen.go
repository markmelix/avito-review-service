package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

const api = "http://localhost:8080"

func main() {
	for i := 1; i <= 20; i++ {
		teamName := fmt.Sprintf("team-%d", i)
		members := []map[string]any{}

		for u := 1; u <= 10; u++ {
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
			fmt.Fprintf(os.Stderr, "Error marshalling data generated: %v\n", err)
			return
		}

		bodyDir := "reqbodies"
		os.Mkdir(bodyDir, 0755)
		filename := fmt.Sprintf("%s/team_add_%03d.json", bodyDir, i)
		err = os.WriteFile(filename, b, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing json %s: %v\n", filename, err)
			return
		}

		endpoint := api + "/team/add"

		fmt.Printf("POST %s\nContent-Type: application/json\n@%s\n\n", endpoint, filename)
	}
}
