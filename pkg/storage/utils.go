package storage

import "fmt"

func GetStorageProvider(engine string, providers map[string]StorageProvider) (StorageProvider, error) {
	if _, ok := providers[engine]; !ok {
		return nil, fmt.Errorf("No matching storage provider found: %s.", engine)
	}
	storageProvider := providers[engine]
	installed, err := storageProvider.IsInstalled()
	if err != nil {
		return nil, err
	}
	if !installed {
		return nil, fmt.Errorf("Storage provider is not installed: %s.", engine)
	}
	return storageProvider, nil
}
