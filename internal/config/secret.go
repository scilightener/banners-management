package config

// Secret is a type for secret values. It implements fmt.Stringer interface to prevent logging itself as a plain text.
type Secret string

func (Secret) String() string {
	return "*****"
}
