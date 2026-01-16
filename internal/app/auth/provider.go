package auth

type Provider interface {
	CheckCredentials(username string, password []byte) (bool, error)
}
