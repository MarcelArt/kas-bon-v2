package models

type AccessControlEval struct {
	Sub string `json:"sub"`
	App string `json:"app"`
	Dom string `json:"dom"`
	Obj string `json:"obj"`
	Act string `json:"act"`
}
