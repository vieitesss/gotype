package main
 
type GameStatusUpdatedMsg struct {
	status GameStatus
}
type TextToWriteMsg []string
type UpdatedWordToRenderMsg string
