package main
 
// Thrown when updating the status of the game.
// It carries the new status.
type GameStatusUpdatedMsg struct {
	status GameStatus
}

// Thrown when the words that the player will have to type are ready.
// It represents the words.
type TextToWriteMsg []string

// Thrown when the current word being typed is ready to be rendered.
// It represents word to be typed styled depending on the input.
type UpdatedWordToRenderMsg string
