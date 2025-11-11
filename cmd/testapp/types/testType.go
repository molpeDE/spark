package types

type TestStruct struct {
	Test        string
	Floatyfloat float32 `json:"iamfloating"`
}

type ExtendedField struct {
	TestStruct
	Banana  bool
	Flotarr []float32
	Bytes   []uint8
}
