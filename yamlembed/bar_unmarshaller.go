package yamlembed

import (
	"strings"
)

func (b *Bar) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type PlainBar Bar // алиас "обнуляет" метод UnmarshalYAML
	var temp PlainBar
	if err := unmarshal(&temp); err != nil {
		return err
	}
	b.B = temp.B
	b.UpperB = strings.ToUpper(temp.B)
	b.OI = temp.OI
	b.F = temp.F
	return nil
}
