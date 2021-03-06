package embetcd

import (
	"context"
	"fmt"
	"net/url"
	"time"

	cli "github.com/coreos/etcd/clientv3"
)

var (
	// ErrNameConflict is an error indicating that the server name is in conflict with an existing member of the cluster
	ErrNameConflict = fmt.Errorf("server name is in conflict with an existing cluster member")
	// ErrAlreadyRunning is an error indicating that the server is already running
	ErrAlreadyRunning = fmt.Errorf("server is already running")
	// ErrAlreadyStopped is an error indicating that the server is already stopped
	ErrAlreadyStopped = fmt.Errorf("server is already stopped")
	// ErrClusterNameConflict is an error indicating that the configured cluster name conflicts with the target cluster
	ErrClusterNameConflict = fmt.Errorf("cluster name either does not exist in the cluster under '/_etcd-cluster/name' or is different from this server's cluster name")
)

// WaitForStructChOrErrCh waits for the struct channel, error channel or context to return a value
func WaitForStructChOrErrCh(ctx context.Context, structCh <-chan struct{}, errCh <-chan error) error {
	// wait for the server to start or error out
	select {
	case <-structCh:
		return nil
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// DurationOrDefault returns the pointed duration or the specified default
func DurationOrDefault(in *time.Duration, def time.Duration) time.Duration {
	if in != nil {
		return *in
	}
	return def
}

// URLSToStringSlice converts urls slices to string slices
func URLSToStringSlice(urls []url.URL) []string {
	strs := make([]string, 0, len(urls))
	for _, u := range urls {
		strs = append(strs, u.String())
	}
	return strs
}

// CloseServer closes an embetcd server with nil checks
func CloseServer(s *Server) {
	// close etcd server if it was improperly created in previous loop iterations
	if s != nil && s.Etcd != nil {
		s.Etcd.Close()
	}
}

// CloseClient closes an embetcd client with nil checks
func CloseClient(client *Client) {
	if client != nil && client.Client != nil {
		client.Close()
	}
}

// CancelContext is a helper function for canceling context using the context.CancelFunc if it isn't nil
func CancelContext(cancel context.CancelFunc) {
	if cancel != nil {
		cancel()
	}
}

// RevokeLease revokes a lease using a client as long as the client and lease are not nil
func RevokeLease(ctx context.Context, client *Client, lease *cli.LeaseGrantResponse) {
	if client != nil && lease != nil {
		client.Revoke(ctx, lease.ID)
	}
}

// StringIsInStringSlice returns true if the given string is in the slice of strings
func StringIsInStringSlice(s string, strs []string) (resp bool) {
	for _, i := range strs {
		if s == i {
			resp = true
			break
		}
	}
	return
}
