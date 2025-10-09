// Package types contains shared types used by both upstream and bootstrap packages.
package types

import (
	"context"
	"net/netip"

	"github.com/AdguardTeam/golibs/errors"
)

// Resolver resolves the hostnames to IP addresses.
type Resolver interface {
	LookupNetIP(ctx context.Context, network Network, host string) ([]netip.Addr, error)
}

// Network is a network type for use in Resolver's methods.
type Network = string

const (
	// NetworkIP is a network type for both address families.
	NetworkIP Network = "ip"

	// NetworkIP4 is a network type for IPv4 address family.
	NetworkIP4 Network = "ip4"

	// NetworkIP6 is a network type for IPv6 address family.
	NetworkIP6 Network = "ip6"

	// NetworkTCP is a network type for TCP connections.
	NetworkTCP Network = "tcp"

	// NetworkUDP is a network type for UDP connections.
	NetworkUDP Network = "udp"
)

// StaticResolver is a resolver which always responds with an underlying slice
// of IP addresses.
type StaticResolver struct {
	addrs []netip.Addr
}

// NewStaticResolver creates a new StaticResolver with the specified addresses.
func NewStaticResolver(addrs []netip.Addr) *StaticResolver {
	return &StaticResolver{addrs: addrs}
}

// LookupNetIP implements the Resolver interface for StaticResolver.
func (r *StaticResolver) LookupNetIP(ctx context.Context, network Network, host string) ([]netip.Addr, error) {
	if network == NetworkIP || network == NetworkIP4 {
		return r.addrs, nil
	}
	if network == NetworkIP6 {
		var ipv6Addrs []netip.Addr
		for _, addr := range r.addrs {
			if addr.Is6() {
				ipv6Addrs = append(ipv6Addrs, addr)
			}
		}
		return ipv6Addrs, nil
	}
	return nil, nil
}

// ParallelResolver is a slice of resolvers that are queried concurrently until
// the first successful response is returned, as opposed to all resolvers being
// queried in order in ConsequentResolver.
type ParallelResolver struct {
	resolvers []Resolver
}

// NewParallelResolver creates a new ParallelResolver with the specified resolvers.
func NewParallelResolver(resolvers ...Resolver) *ParallelResolver {
	return &ParallelResolver{resolvers: resolvers}
}

// LookupNetIP implements the Resolver interface for ParallelResolver.
func (r *ParallelResolver) LookupNetIP(ctx context.Context, network Network, host string) ([]netip.Addr, error) {
	if len(r.resolvers) == 0 {
		return nil, errors.Error("no resolvers specified")
	}

	type result struct {
		addrs []netip.Addr
		err   error
	}

	resCh := make(chan result, len(r.resolvers))
	for _, resolver := range r.resolvers {
		go func(resolver Resolver) {
			addrs, err := resolver.LookupNetIP(ctx, network, host)
			resCh <- result{addrs: addrs, err: err}
		}(resolver)
	}

	var errs []error
	for range r.resolvers {
		res := <-resCh
		if res.err == nil && len(res.addrs) > 0 {
			return res.addrs, nil
		}
		if res.err != nil {
			errs = append(errs, res.err)
		}
	}

	return nil, joinErrors(errs)
}

// ConsequentResolver is a slice of resolvers that are queried in order until
// the first successful non-empty response, as opposed to just successful
// response requirement in ParallelResolver.
type ConsequentResolver struct {
	resolvers []Resolver
}

// NewConsequentResolver creates a new ConsequentResolver with the specified resolvers.
func NewConsequentResolver(resolvers ...Resolver) *ConsequentResolver {
	return &ConsequentResolver{resolvers: resolvers}
}

// LookupNetIP implements the Resolver interface for ConsequentResolver.
func (r *ConsequentResolver) LookupNetIP(ctx context.Context, network Network, host string) ([]netip.Addr, error) {
	var errs []error
	for _, resolver := range r.resolvers {
		addrs, err := resolver.LookupNetIP(ctx, network, host)
		if err == nil && len(addrs) > 0 {
			return addrs, nil
		}
		if err != nil {
			errs = append(errs, err)
		}
	}

	return nil, joinErrors(errs)
}

// joinErrors joins multiple errors into a single error.
func joinErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}

	var errStrs []string
	for _, err := range errs {
		errStrs = append(errStrs, err.Error())
	}

	return &multiError{errs: errStrs}
}

// multiError represents multiple errors as a single error.
type multiError struct {
	errs []string
}

func (e *multiError) Error() string {
	return joinStrings(e.errs, "\n")
}

func (e *multiError) Unwrap() []error {
	var errs []error
	for _, errStr := range e.errs {
		errs = append(errs, &errorString{s: errStr})
	}
	return errs
}

// errorString implements the error interface for a string.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// joinStrings joins a slice of strings with a separator.
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for _, s := range strs[1:] {
		result += sep + s
	}
	return result
}
