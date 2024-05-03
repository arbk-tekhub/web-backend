package validator

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) HasErrors() bool {
	return len(v.Errors) != 0
}

func (v *Validator) AddError(key, msg string) {

	if _, exist := v.Errors[key]; !exist {
		v.Errors[key] = msg
	}
}

func (v *Validator) Check(ok bool, key string, msg string) {
	if !ok {
		v.AddError(key, msg)
	}
}
