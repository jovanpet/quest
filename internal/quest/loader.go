package quest

import (
	"encoding/json"
	"os"

	"github.com/jovanpet/quest/internal/format"
	"github.com/jovanpet/quest/internal/types"
)

func LoadState() (*types.State, error) {
	data, err := os.ReadFile(StateFilePath)
	if err != nil {
		return nil, err
	}

	var state types.State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

func LoadPlan() (*types.Plan, error) {
	data, err := os.ReadFile(PlanFilePath)
	if err != nil {
		return nil, err
	}

	var plan types.Plan
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, err
	}

	return &plan, nil
}

func LoadStateAndPlan() (*types.State, *types.Plan, error) {
	state, err := LoadState()
	if err != nil {
		format.ErrorWithTip("Failed to load state", err, "Run 'quest begin' to start a new quest")
		return nil, nil, err
	}

	plan, err := LoadPlan()
	if err != nil {
		format.ErrorWithTip("Failed to load plan", err, "Run 'quest begin' to start a new quest")
		return nil, nil, err
	}

	return state, plan, nil
}
