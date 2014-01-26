package passwd

type Keyring interface {
	// Gets the password for the service and username if exists.
	Get(service, username string) (string, error)

	// Sets a password for the service and username.
	Set(service, username, password string) error

	// Deletes a password belongs to the service and username.
	Delete(service, username string) error
}

// NewKeyring returns the Keyring client available on the platform.
// Currently only Factotum
// http://plan9.bell-labs.com/magic/man2html/4/factotum is implemented.
func NewKeyring() Keyring { return new(factotum) }
