package megabrowser

const (
	login = `login`
	pass  = `password`
)

type StorageBrowser interface {
	GetObjectNode(file string) (string, error)
}

type MegaBrowser struct {
	megaClient StorageClient
}

func NewMegaBrowser(megaClient StorageClient) (*MegaBrowser, error) {
	browser := &MegaBrowser{
		megaClient: megaClient,
	}
	err := browser.megaClient.Login(login, pass)
	if err != nil {
		return nil, err
	}
	return browser, nil
}

func (b *MegaBrowser) GetObjectNode(file string) (string, error) {
	return "", nil
}
