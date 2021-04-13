package validators

type PreSharedKeyValidator struct {
	preSharedKey string
}

func NewPreSharedKeyValidator(preSharedKey string) *PreSharedKeyValidator {
	return &PreSharedKeyValidator{preSharedKey}
}

func (v *PreSharedKeyValidator) Validate(preSharedKey string) bool {
	return preSharedKey == v.preSharedKey
}
