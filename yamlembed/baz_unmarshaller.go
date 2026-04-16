package yamlembed

func (bz *Baz) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&bz.Foo); err != nil {
		return err
	}
	if err := unmarshal(&bz.Bar); err != nil {
		return err
	}
	return nil
}
