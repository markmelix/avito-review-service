package team

import "review/internal/domain"

type TeamAddDto struct {
	Team domain.Team `json:"team"`
}
