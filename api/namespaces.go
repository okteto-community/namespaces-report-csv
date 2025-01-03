package api

import (
	"fmt"
	"log/slog"

	"github.com/okteto-community/list-namespace-information/model"
)

const (
	developmentNamespaceType = "development"
)

// GetNamespaces retrieves all the namespaces
func GetDevelopmentNamespaces(baseURL, token string, logger *slog.Logger) ([]model.Namespace, error) {
	namespacesURL := fmt.Sprintf("https://%s/%s?type=%s", baseURL, namespacesAPIPath, developmentNamespaceType)
	var namespaces []model.Namespace
	if err := sendRequest(namespacesURL, token, &namespaces, logger); err != nil {
		return nil, err
	}
	return namespaces, nil
}
