package quest

import (
	"encoding/json"
	"os"

	"github.com/jovanpet/quest/internal/format"
	"github.com/jovanpet/quest/internal/types"
)

func UploadState(state *types.State) error {
	jsonState, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		format.Error("Error marshaling state", err)
		return err
	}
	err = os.WriteFile(StateFilePath, jsonState, 0644)
	if err != nil {
		format.ErrorWithTip("Error writing state file", err, "Check folder permissions")
		return err
	}
	return nil
}

func UploadPlan(plan *types.Plan) error {
	jsonPlan, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		format.Error("Error marshaling plan", err)
		return err
	}
	err = os.WriteFile(PlanFilePath, jsonPlan, 0644)
	if err != nil {
		format.ErrorWithTip("Error writing plan file", err, "Check folder permissions")
		return err
	}
	return nil
}

func UploadStateAndPlan(state *types.State, plan *types.Plan) error {
	err := UploadState(state)
	if err != nil {
		format.Error("Failed to save state", err)
		return err
	}

	err = UploadPlan(plan)
	if err != nil {
		format.Error("Failed to save plan", err)
		return err
	}

	return nil
}
