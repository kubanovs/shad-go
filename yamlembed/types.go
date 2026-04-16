package yamlembed

type Foo struct {
	A string `yaml:"aa"`
	p int64  `yaml:"-"`
}

type Bar struct {
	I      int64    `yaml:"-"`
	B      string   `yaml:"b"`
	UpperB string   `yaml:"-"`
	OI     []string `yaml:"oi,omitempty,flow"`
	F      []any    `yaml:"f,omitempty,flow"`
}

type Baz struct {
	Foo `yaml:",inline"`
	Bar `yaml:",inline"`
}
