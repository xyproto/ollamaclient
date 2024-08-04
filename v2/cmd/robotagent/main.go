package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Tool struct {
	Type     string      `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Parameters  ToolParameters `json:"parameters"`
}

type ToolParameters struct {
	Type       string                 `json:"type"`
	Properties map[string]ToolProperty `json:"properties"`
	Required   []string               `json:"required"`
}

type ToolProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

type Robot struct {
	X, Y int
	HasBox bool
}

func (r *Robot) Move(direction string) {
	switch direction {
	case "up":
		if r.Y > 0 {
			r.Y--
		}
	case "down":
		if r.Y < 18 {
			r.Y++
		}
	case "left":
		if r.X > 0 {
			r.X--
		}
	case "right":
		if r.X < 18 {
			r.X++
		}
	}
}

func (r *Robot) PickUpBox() {
	r.HasBox = true
}

func (r *Robot) SetDownBox() {
	r.HasBox = false
}

func main() {
	oc := New("llama3.1")
	oc.Verbose = true

	err := oc.PullIfNeeded(true)
	if err != nil {
		t.Fatalf("Failed to pull model: %v", err)
	}

	if found, err := oc.Has("llama3.1"); err != nil || !found {
		t.Error("Expected to have 'llama3.1' model downloaded, but it's not present")
	}

	oc.SetSystemPrompt("You are a robot moving around a 19x19 grid, picking up and setting down boxes.")
	oc.SetRandom()

	var toolMove, toolPickUpBox, toolSetDownBox Tool

	// Define the move tool
	json.Unmarshal(json.RawMessage(`{
		"type": "function",
		"function": {
		  "name": "move",
		  "description": "Move the robot in a direction",
		  "parameters": {
			"type": "object",
			"properties": {
			  "direction": {
				"type": "string",
				"description": "The direction to move the robot in, e.g. 'up', 'down', 'left', 'right'",
				"enum": ["up", "down", "left", "right"]
			  }
			},
			"required": ["direction"]
		  }
		}
	  }`), &toolMove)

	// Define the pick up box tool
	json.Unmarshal(json.RawMessage(`{
		"type": "function",
		"function": {
		  "name": "pick_up_box",
		  "description": "Pick up a box"
		}
	  }`), &toolPickUpBox)

	// Define the set down box tool
	json.Unmarshal(json.RawMessage(`{
		"type": "function",
		"function": {
		  "name": "set_down_box",
		  "description": "Set down the box"
		}
	  }`), &toolSetDownBox)

	oc.SetTool(toolMove)
	oc.SetTool(toolPickUpBox)
	oc.SetTool(toolSetDownBox)

	// Create a robot instance
	robot := &Robot{X: 0, Y: 0, HasBox: false}

	// Simulate a prompt and robot actions
	prompt := "Move the robot to the right, then pick up a box."
	generatedOutput := oc.MustOutputChat(prompt)
	if generatedOutput.Error != "" {
		t.Error(generatedOutput.Error)
	}

	for _, toolCall := range generatedOutput.ToolCalls {
		switch toolCall.Function.Name {
		case "move":
			direction := toolCall.Arguments["direction"].(string)
			robot.Move(direction)
			fmt.Printf("Robot moved %s to position (%d, %d)\n", direction, robot.X, robot.Y)
		case "pick_up_box":
			robot.PickUpBox()
			fmt.Println("Robot picked up a box")
		case "set_down_box":
			robot.SetDownBox()
			fmt.Println("Robot set down the box")
		}
	}

	// Display final state of the robot
	fmt.Printf("Final position of robot: (%d, %d), HasBox: %v\n", robot.X, robot.Y, robot.HasBox)
}
