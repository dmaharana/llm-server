package main

type Models struct {
	Models []Model `json:"models"`
}

type Model struct {
	Name    string      `json:"name"`
	Model   string      `json:"model"`
	Details ModelDetail `json:"details"`
}

type ModelDetail struct {
	ParentModel       string   `json:"parent_model"`
	Format            string   `json:"format"`
	Family            string   `json:"family"`
	Families          []string `json:"families"`
	ParameterSize     string   `json:"parameter_size"`
	QuantizationLevel string   `json:"quantization_level"`
}
