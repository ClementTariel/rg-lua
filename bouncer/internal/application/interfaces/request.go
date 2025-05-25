package interfaces

type AddPendingMatchRequest struct {
	BlueName string `json:"blueName"`
	RedName  string `json:"redName"`
}
