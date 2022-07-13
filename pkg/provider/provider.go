package provider

type Provider interface {
	Post(msg string) error
}
