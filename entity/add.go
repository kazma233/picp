package entity

type (
	AddMask struct {
		Origin string `json:"origin"`

		Infos []AddInfo `json:"infos"`
	}

	AddInfo struct {
		Type string `json:"type"`

		X int `json:"x"`
		Y int `json:"y"`

		Word  string  `json:"word"`
		Color string  `json:"color"`
		Size  float64 `json:"size"`
		Font  int     `json:"font"`
		Dpi   float64 `json:"dpi"`

		WaterMask string  `json:"waterMask"`
		Opacity   float64 `json:"opacity"`
	}
)
