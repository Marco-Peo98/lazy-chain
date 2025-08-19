// network_clients.go
package settings

import (
	"fmt"
	"strings"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
)

func trimSlash(s string) string {
	return strings.TrimRight(strings.TrimSpace(s), "/")
}

func hasScheme(s string) bool {
	return strings.Contains(s, "://")
}

func hasExplicitPort(u string) bool {
	i := strings.LastIndex(u, ":")
	j := strings.Index(u, "://")
	return i > j+2
}

// createAlgodClient creates an algod client from network info
func createAlgodClient(networkInfo NetworkInfo) (*algod.Client, error) {
	base := strings.TrimSpace(networkInfo.AlgodURL)

	if hasScheme(base) {
		base = trimSlash(base)
		// Append a non-default, explicit port if provided and not already present
		if networkInfo.AlgodPort != "" && !hasExplicitPort(base) &&
			networkInfo.AlgodPort != "443" && networkInfo.AlgodPort != "80" {
			base = fmt.Sprintf("%s:%s", base, networkInfo.AlgodPort)
		}
	} else {
		base = fmt.Sprintf("%s:%s", trimSlash(base), networkInfo.AlgodPort)
	}

	return algod.MakeClient(base, networkInfo.AlgodToken)
}

// createIndexerClient creates an indexer client from network info
func createIndexerClient(networkInfo NetworkInfo) (*indexer.Client, error) {
	if networkInfo.IndexerURL == "" {
		return nil, fmt.Errorf("indexer URL not provided")
	}

	base := strings.TrimSpace(networkInfo.IndexerURL)

	if hasScheme(base) {
		base = trimSlash(base)
		if networkInfo.IndexerPort != "" && !hasExplicitPort(base) &&
			networkInfo.IndexerPort != "443" && networkInfo.IndexerPort != "80" {
			base = fmt.Sprintf("%s:%s", base, networkInfo.IndexerPort)
		}
	} else {
		base = fmt.Sprintf("%s:%s", trimSlash(base), networkInfo.IndexerPort)
	}

	return indexer.MakeClient(base, networkInfo.IndexerToken)
}
