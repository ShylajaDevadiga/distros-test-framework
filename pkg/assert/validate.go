package assert

import (
	"fmt"
	"strings"
	"time"

	"github.com/rancher/distros-test-framework/shared"
)

// validate calls runAssertion for each cmd/assert pair
//
// the first caller - process tests will spawn a go routine per ip the cluster
//
// need to send KubeconfigFile
func validate(exec func(string) (string, error), args ...string) error {
	if len(args) < 2 || len(args)%2 != 0 {
		return fmt.Errorf("must receive an even number of arguments as cmd/assert pairs")
	}

	errorsChan := make(chan error, len(args)/2)
	timeout := time.After(220 * time.Second)
	ticker := time.NewTicker(3 * time.Second)

	for i := 0; i < len(args); i++ {
		cmd := args[i]
		if i+1 < len(args) {
			assert := args[i+1]
			i++

			err := runAssertion(cmd, assert, exec, ticker.C, timeout, errorsChan)
			if err != nil {
				fmt.Printf("error from runAssertion():\n %s\n", err)
				close(errorsChan)
				return err
			}
		}
	}
	close(errorsChan)

	return nil
}

// runAssertion runs a command and asserts that the value received against his respective command
func runAssertion(
	cmd, assert string,
	exec func(string) (string, error),
	ticker <-chan time.Time,
	timeout <-chan time.Time,
	errorsChan chan<- error,
) error {
	for {
		select {
		case <-timeout:
			timeoutErr := fmt.Errorf("timeout reached for command:\n%s\n "+
				"Trying to assert with:\n %s",
				cmd, assert)
			errorsChan <- timeoutErr
			return timeoutErr

		case <-ticker:
			res, err := exec(cmd)
			if err != nil {
				errorsChan <- err
				return fmt.Errorf("error from runCmd:\n %s\n %s", res, err)
			}
			fmt.Printf("\n---------------------\nCommand:\n"+
				"%s\n"+
				"\n---------------------\nResult:\n"+
				"%s\n"+
				"\n---------------------\nAssertion:\n"+
				"%s\n", cmd, res, assert)
			if strings.Contains(res, assert) {
				fmt.Printf("Matched with: \n%s\n", res)
				errorsChan <- nil
				return nil
			}
		}
	}
}

// ValidateOnHost runs an exec function on RunCommandHost and assert given is fulfilled.
// The last argument should be the assertion.
// Need to send kubeconfig file.
func ValidateOnHost(args ...string) error {
	exec := func(cmd string) (string, error) {
		return shared.RunCommandHost(cmd)
	}
	return validate(exec, args...)
}

// ValidateOnNode runs an exec function on RunCommandHost and assert given is fulfilled.
// The last argument should be the assertion.
func ValidateOnNode(ip string, args ...string) error {
	exec := func(cmd string) (string, error) {
		return shared.RunCommandOnNode(cmd, ip)
	}
	return validate(exec, args...)
}
