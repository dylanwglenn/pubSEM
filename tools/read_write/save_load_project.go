package read_write

import (
	"encoding/json"
	"main/model"
	"main/utils"
	"os"
)

func SaveProject(m *model.Model, path string) {

	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		panic(err)
	}
}

func LoadProject(path string) (*model.Model, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	m := new(model.Model)
	err = json.Unmarshal(data, &m)
	if err != nil {
		panic(err)
	}

	m.Font.Faces = utils.LoadAllFontFaces()

	return m, nil
}
