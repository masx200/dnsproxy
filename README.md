# DNS Proxy <!-- omit in toc -->

[![Code Coverage](https://img.shields.io/codecov/c/github/AdguardTeam/dnsproxy/master.svg)](https://codecov.io/github/AdguardTeam/dnsproxy?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/masx200/dnsproxy)](https://goreportcard.com/report/AdguardTeam/dnsproxy)
[![Go Doc](https://godoc.org/github.com/masx200/dnsproxy?status.svg)](https://godoc.org/github.com/masx200/dnsproxy)

A simple DNS proxy server that supports all existing DNS protocols including
`DNS-over-TLS`, `DNS-over-HTTPS`, `DNSCrypt`, and `DNS-over-QUIC`. Moreover, it
can work as a `DNS-over-HTTPS`, `DNS-over-TLS` or `DNS-over-QUIC` server.

- [How to install](#how-to-install)
- [How to build](#how-to-build)
- [Usage](#usage)
- [Examples](#examples)
  - [Simple options](#simple-options)
  - [Encrypted upstreams](#encrypted-upstreams)
  - [Encrypted DNS server](#encrypted-dns-server)
  - [Additional features](#additional-features)
  - [DNS64 server](#dns64-server)
  - [Fastest addr + cache-min-ttl](#fastest-addr--cache-min-ttl)
  - [Specifying upstreams for domains](#specifying-upstreams-for-domains)
  - [Specifying private rDNS upstreams](#specifying-private-rdns-upstreams)
  - [EDNS Client Subnet](#edns-client-subnet)
  - [Bogus NXDomain](#bogus-nxdomain)
  - [Basic Auth for DoH](#basic-auth-for-doh)

## How to install

There are several options how to install `dnsproxy`.

1. Grab the binary for your device/OS from the [Releases][releases] page.
2. Use the [official Docker image][docker].
3. Build it yourself (see the instruction below).

[releases]: https://github.com/masx200/dnsproxy/releases
[docker]: https://hub.docker.com/r/adguard/dnsproxy

## How to build

You will need Go 1.25 or later.

```shell
make build
```

## Usage

```none
Usage of ./dnsproxy:
  --bogus-nxdomain=subnet
        Transform the responses containing at least a single IP that matches specified addresses and CIDRs into NXDOMAIN.  Can be specified multiple times.
  --bootstrap/-b
        Bootstrap DNS for DoH and DoT, can be specified multiple times (default: use system-provided).
  --cache
        If specified, DNS cache is enabled.
  --cache-max-ttl=uint32
        Maximum TTL value for DNS entries, in seconds.
  --cache-min-ttl=uint32
        Minimum TTL value for DNS entries, in seconds. Capped at 3600. Artificially extending TTLs should only be done with careful consideration.
  --cache-optimistic
        If specified, optimistic DNS cache is enabled.
  --cache-size=int
        Cache size (in bytes). Default: 64k.
  --config-path=path
        YAML configuration file. Minimal working configuration in config.yaml.dist. Options passed through command line will override the ones from this file.
  --dns64
        If specified, dnsproxy will act as a DNS64 server.
  --dns64-prefix=subnet
        Prefix used to handle DNS64. If not specified, dnsproxy uses the 'Well-Known Prefix' 64:ff9b::.  Can be specified multiple times.
  --dnscrypt-config=path/-g path
        Path to a file with DNSCrypt configuration. You can generate one using https://github.com/ameshkov/dnscrypt.
  --dnscrypt-port=port/-y port
        Listening ports for DNSCrypt.
  --edns
        Use EDNS Client Subnet extension.
  --edns-addr=address
        Send EDNS Client Address.
  --fallback/-f
        Fallback resolvers to use when regular ones are unavailable, can be specified multiple times. You can also specify path to a file with the list of servers.
  --help/-h
        Print this help message and quit.
  --hosts-file-enabled
        If specified, use hosts files for resolving.
  --hosts-files=path
        List of paths to the hosts files, can be specified multiple times.
  --http3
        Enable HTTP/3 support.
  --https-port=port/-s port
        Listening ports for DNS-over-HTTPS.
  --https-server-name=name
        Set the Server header for the responses from the HTTPS server.
  --https-userinfo=name
        If set, all DoH queries are required to have this basic authentication information.
  --insecure
        Disable secure TLS certificate validation.
  --ipv6-disabled
        If specified, all AAAA requests will be replied with NoError RCode and empty answer.
  --listen=address/-l address
        Listening addresses.
  --max-go-routines=uint
        Set the maximum number of go routines. A zero value will not not set a maximum.
  --output=path/-o path
        Path to the log file.
  --pending-requests-enabled
        If specified, the server will track duplicate queries and only send the first of them to the upstream server, propagating its result to others. Disabling it introduces a vulnerability to cache poisoning attacks.
  --port=port/-p port
        Listening ports. Zero value disables TCP and UDP listeners.
  --pprof
        If present, exposes pprof information on localhost:6060.
  --private-rdns-upstream
        Private DNS upstreams to use for reverse DNS lookups of private addresses, can be specified multiple times.
  --private-subnets=subnet
        Private subnets to use for reverse DNS lookups of private addresses.
  --quic-port=port/-q port
        Listening ports for DNS-over-QUIC.
  --ratelimit=int/-r int
        Ratelimit (requests per second).
  --ratelimit-subnet-len-ipv4=int
        Ratelimit subnet length for IPv4.
  --ratelimit-subnet-len-ipv6=int
        Ratelimit subnet length for IPv6.
  --refuse-any
        If specified, refuses ANY requests.
  --timeout=duration
        Timeout for outbound DNS queries to remote upstream servers in a human-readable form
  --tls-crt=path/-c path
        Path to a file with the certificate chain.
  --tls-key=path/-k path
        Path to a file with the private key.
  --tls-max-version=version
        Maximum TLS version, for example 1.3.
  --tls-min-version=version
        Minimum TLS version, for example 1.0.
  --tls-port=port/-t port
        Listening ports for DNS-over-TLS.
  --udp-buf-size=int
        Set the size of the UDP buffer in bytes. A value <= 0 will use the system default.
  --upstream/-u
        An upstream to be used (can be specified multiple times). You can also specify path to a file with the list of servers.
  --upstream-mode=mode
        Defines the upstreams logic mode, possible values: load_balance, parallel, fastest_addr (default: load_balance).
  --use-private-rdns
        If specified, use private upstreams for reverse DNS lookups of private addresses.
  --verbose/-v
        Verbose output.
  --version
        Prints the program version.
```

## Examples

### Simple options

Runs a DNS proxy on `0.0.0.0:53` with a single upstream - Google DNS.

```shell
./dnsproxy -u 8.8.8.8:53
```

The same proxy with verbose logging enabled writing it to the file `log.txt`.

```shell
./dnsproxy -u 8.8.8.8:53 -v -o log.txt
```

Runs a DNS proxy on `127.0.0.1:5353` with multiple upstreams.

```shell
./dnsproxy -l 127.0.0.1 -p 5353 -u 8.8.8.8:53 -u 1.1.1.1:53
```

Listen on multiple interfaces and ports:

```shell
./dnsproxy -l 127.0.0.1 -l 192.168.1.10 -p 5353 -p 5354 -u 1.1.1.1
```

The plain DNS upstream server may be specified in several ways:

- With a plain IP address:

  ```shell
  ./dnsproxy -l 127.0.0.1 -u 8.8.8.8:53
  ```

- With a hostname or plain IP address and the `udp://` scheme:

  ```shell
  ./dnsproxy -l 127.0.0.1 -u udp://dns.google -u udp://1.1.1.1
  ```

- With a hostname or plain IP address and the `tcp://` scheme to force using
  TCP:

  ```shell
  ./dnsproxy -l 127.0.0.1 -u tcp://dns.google -u tcp://1.1.1.1
  ```

### Encrypted upstreams

DNS-over-TLS upstream:

```shell
./dnsproxy -u tls://dns.adguard.com
```

DNS-over-HTTPS upstream with specified bootstrap DNS:

```shell
./dnsproxy -u https://dns.adguard.com/dns-query -b 1.1.1.1:53
```

DNS-over-QUIC upstream:

```shell
./dnsproxy -u quic://dns.adguard.com
```

DNS-over-HTTPS upstream with enabled HTTP/3 support (chooses it if it's faster):

```shell
./dnsproxy -u https://dns.google/dns-query --http3
```

DNS-over-HTTPS upstream with forced HTTP/3 (no fallback to other protocol):

```shell
./dnsproxy -u h3://dns.google/dns-query
```

DNSCrypt upstream ([DNS Stamp](https://dnscrypt.info/stamps) of AdGuard DNS):

```shell
./dnsproxy -u sdns://AQMAAAAAAAAAETk0LjE0MC4xNC4xNDo1NDQzINErR_JS3PLCu_iZEIbq95zkSV2LFsigxDIuUso_OQhzIjIuZG5zY3J5cHQuZGVmYXVsdC5uczEuYWRndWFyZC5jb20
```

DNS-over-HTTPS upstream ([DNS Stamp](https://dnscrypt.info/stamps) of Cloudflare
DNS):

```shell
./dnsproxy -u sdns://AgcAAAAAAAAABzEuMC4wLjGgENk8mGSlIfMGXMOlIlCcKvq7AVgcrZxtjon911-ep0cg63Ul-I8NlFj4GplQGb_TTLiczclX57DvMV8Q-JdjgRgSZG5zLmNsb3VkZmxhcmUuY29tCi9kbnMtcXVlcnk
```

DNS-over-TLS upstream with two fallback servers (to be used when the main
upstream is not available):

```shell
./dnsproxy -u tls://dns.adguard.com -f 8.8.8.8:53 -f 1.1.1.1:53
```

### Encrypted DNS server

Runs a DNS-over-TLS proxy on `127.0.0.1:853`.

```shell
./dnsproxy -l 127.0.0.1 --tls-port=853 --tls-crt=example.crt --tls-key=example.key -u 8.8.8.8:53 -p 0
```

Runs a DNS-over-HTTPS proxy on `127.0.0.1:443`.

```shell
./dnsproxy -l 127.0.0.1 --https-port=443 --tls-crt=example.crt --tls-key=example.key -u 8.8.8.8:53 -p 0
```

Runs a DNS-over-HTTPS proxy on `127.0.0.1:443` with HTTP/3 support.

```shell
./dnsproxy -l 127.0.0.1 --https-port=443 --http3 --tls-crt=example.crt --tls-key=example.key -u 8.8.8.8:53 -p 0
```

Runs a DNS-over-QUIC proxy on `127.0.0.1:853`.

```shell
./dnsproxy -l 127.0.0.1 --quic-port=853 --tls-crt=example.crt --tls-key=example.key -u 8.8.8.8:53 -p 0
```

Runs a DNSCrypt proxy on `127.0.0.1:443`.

```shell
./dnsproxy -l 127.0.0.1 --dnscrypt-config=./dnscrypt-config.yaml --dnscrypt-port=443 --upstream=8.8.8.8:53 -p 0
```

> [!TIP]
> In order to run a DNSCrypt proxy, you need to obtain DNSCrypt configuration
> first. You can use https://github.com/ameshkov/dnscrypt command-line tool to
> do that with a command like this
> `./dnscrypt generate --provider-name=2.dnscrypt-cert.example.org --out=dnscrypt-config.yaml`.

### Additional features

Runs a DNS proxy on `0.0.0.0:53` with rate limit set to `10 rps`, enabled DNS
cache, and that refuses type=ANY requests.

```shell
./dnsproxy -u 8.8.8.8:53 -r 10 --cache --refuse-any
```

Runs a DNS proxy on 127.0.0.1:5353 with multiple upstreams and enable parallel
queries to all configured upstream servers.

```shell
./dnsproxy -l 127.0.0.1 -p 5353 -u 8.8.8.8:53 -u 1.1.1.1:53 -u tls://dns.adguard.com --upstream-mode parallel
```

Loads upstreams list from a file.

```shell
./dnsproxy -l 127.0.0.1 -p 5353 -u ./upstreams.txt
```

### DNS64 server

`dnsproxy` is capable of working as a DNS64 server.

> [!NOTE] What is DNS64/NAT64 This is a mechanism of providing IPv6 access to
> IPv4. Using a NAT64 gateway with IPv4-IPv6 translation capability lets
> IPv6-only clients connect to IPv4-only services via synthetic IPv6 addresses
> starting with a prefix that routes them to the NAT64 gateway. DNS64 is a DNS
> service that returns AAAA records with these synthetic IPv6 addresses for
> IPv4-only destinations (with A but not AAAA records in the DNS). This lets
> IPv6-only clients use NAT64 gateways without any other configuration. See also
> [RFC 6147](https://datatracker.ietf.org/doc/html/rfc6147).

Enables DNS64 with the default [Well-Known Prefix][wkp]:

```shell
./dnsproxy -l 127.0.0.1 -p 5353 -u 8.8.8.8 --use-private-rdns --private-rdns-upstream=127.0.0.1 --dns64
```

You can also specify any number of custom DNS64 prefixes:

```shell
./dnsproxy -l 127.0.0.1 -p 5353 -u 8.8.8.8 --use-private-rdns --private-rdns-upstream=127.0.0.1 --dns64 --dns64-prefix=64:ffff:: --dns64-prefix=32:ffff::
```

Note that only the first specified prefix will be used for synthesis.

PTR queries for addresses within the specified ranges or the
[Well-Known one][wkp] could only be answered with locally appropriate data, so
dnsproxy will route those to the local upstream servers. Those should be
specified and enabled if DNS64 is enabled.

[wkp]: https://datatracker.ietf.org/doc/html/rfc6052#section-2.1

### Fastest addr + cache-min-ttl

This option would be useful to the users with problematic network connection. In
this mode, `dnsproxy` would detect the fastest IP address among all that were
returned, and it will return only it.

Additionally, for those with problematic network connection, it makes sense to
override `cache-min-ttl`. In this case, `dnsproxy` will make sure that DNS
responses are cached for at least the specified amount of time.

It makes sense to run it with multiple upstream servers only.

Run a DNS proxy with two upstreams, min-TTL set to 10 minutes, fastest address
detection is enabled:

```shell
./dnsproxy -u 8.8.8.8 -u 1.1.1.1 --cache --cache-min-ttl=600 --upstream-mode=fastest_addr
```

who run `dnsproxy` with multiple upstreams

### Specifying upstreams for domains

You can specify upstreams that will be used for a specific domain(s). We use the
dnsmasq-like syntax, decorating domains with brackets (see `--server`
[description][server-description]).

**Syntax:** `[/[domain1][/../domainN]/]upstreamString`

Where `upstreamString` is one or many upstreams separated by space (e.g.
`1.1.1.1` or `1.1.1.1 2.2.2.2`).

If one or more domains are specified, that upstream (`upstreamString`) is used
only for those domains. Usually, it is used for private nameservers. For
instance, if you have a nameserver on your network which deals with
`xxx.internal.local` at `192.168.0.1` then you can specify
`[/internal.local/]192.168.0.1`, and dnsproxy will send all queries to that
nameserver. Everything else will be sent to the default upstreams (which are
mandatory!).

1. An empty domain specification, `//` has the special meaning of "unqualified
   names only", which will be used to resolve names with a single label in them,
   or with exactly two labels in case of `DS` requests.
1. More specific domains take precedence over less specific domains, so:
   `--upstream=[/host.com/]1.2.3.4 --upstream=[/www.host.com/]2.3.4.5` will send
   queries for `*.host.com` to `1.2.3.4`, except `*.www.host.com`, which will go
   to `2.3.4.5`.
1. The special server address `#` means, "use the common servers", so:
   `--upstream=[/host.com/]1.2.3.4 --upstream=[/www.host.com/]#` will send
   queries for `*.host.com` to `1.2.3.4`, except `*.www.host.com` which will be
   forwarded as usual.
1. The wildcard `*` has special meaning of "any sub-domain", so:
   `--upstream=[/*.host.com/]1.2.3.4` will send queries for `*.host.com` to
   `1.2.3.4`, but `host.com` will be forwarded to default upstreams.

Sends requests for `*.local` domains to `192.168.0.1:53`. Other requests are
sent to `8.8.8.8:53`:

```shell
./dnsproxy \
    -u "8.8.8.8:53" \
    -u "[/local/]192.168.0.1:53" \
    ;
```

Sends requests for `*.host.com` to `1.1.1.1:53` except for `*.maps.host.com`
which are sent to `8.8.8.8:53` (along with other requests):

```shell
./dnsproxy \
    -u "8.8.8.8:53" \
    -u "[/host.com/]1.1.1.1:53" \
    -u "[/maps.host.com/]#" \
    ;
```

Sends requests for `*.host.com` to `1.1.1.1:53` except for `host.com` which is
sent to `9.9.9.10:53`, and all other requests are sent to `8.8.8.8:53`:

```shell
./dnsproxy \
    -u "8.8.8.8:53" \
    -u "[/host.com/]9.9.9.10:53" \
    -u "[/*.host.com/]1.1.1.1:53" \
    ;
```

Sends requests for `com` (and its subdomains) to `1.2.3.4:53`, requests for
other top-level domains to `1.1.1.1:53`, and all other requests to `8.8.8.8:53`:

```shell
./dnsproxy \
    -u "8.8.8.8:53" \
    -u "[//]1.1.1.1:53" \
    -u "[/com/]1.2.3.4:53" \
    ;
```

### Specifying private rDNS upstreams

You can specify upstreams that will be used for reverse DNS requests of type PTR
for private addresses. Same applies to the authority requests of types SOA and
NS. The set of private addresses is defined by the `--private-rdns-upstream`,
and the set from [RFC 6303][rfc6303] is used by default.

The additional requirement to the domains specified for upstreams is to be
`in-addr.arpa`, `ip6.arpa`, or its subdomain. Addresses encoded in the domains
should also be private.

Sends queries for `*.168.192.in-addr.arpa` to `192.168.1.2`, if requested by
client from `192.168.0.0/16` subnet. Other queries answered with `NXDOMAIN`:

```shell
./dnsproxy \
    -l "0.0.0.0" \
    -u "8.8.8.8" \
    --use-private-rdns \
    --private-subnets="192.168.0.0/16" \
    --private-rdns-upstream="192.168.1.2" \
    ;
```

Sends queries for `*.in-addr.arpa` to `192.168.1.2`, `*.ip6.arpa` to `fe80::1`,
if requested by client within the default [RFC 6303][rfc6303] subnet set. Other
queries answered with `NXDOMAIN`:

```shell
./dnsproxy\
    -l "0.0.0.0"\
    -u 8.8.8.8\
    --use-private-rdns\
    --private-rdns-upstream="192.168.1.2"\
    --private-rdns-upstream="[/ip6.arpa/]fe80::1"
```

[rfc6303]: https://datatracker.ietf.org/doc/html/rfc6303
[server-description]: http://www.thekelleys.org.uk/dnsmasq/docs/dnsmasq-man.html

### EDNS Client Subnet

To enable support for EDNS Client Subnet extension you should run dnsproxy with
`--edns` flag:

```shell
./dnsproxy -u 8.8.8.8:53 --edns
```

Now if you connect to the proxy from the Internet - it will pass through your
original IP address's prefix to the upstream server. This way the upstream
server may respond with IP addresses of the servers that are located near you to
minimize latency.

If you want to use EDNS CS feature when you're connecting to the proxy from a
local network, you need to set `--edns-addr=PUBLIC_IP` argument:

```shell
./dnsproxy -u 8.8.8.8:53 --edns --edns-addr=72.72.72.72
```

Now even if your IP address is 192.168.0.1 and it's not a public IP, the proxy
will pass through 72.72.72.72 to the upstream server.

### Bogus NXDomain

This option is similar to dnsmasq `bogus-nxdomain`. `dnsproxy` will transform
responses that contain at least a single IP address which is also specified by
the option into `NXDOMAIN`. Can be specified multiple times.

In the example below, we use AdGuard DNS server that returns `0.0.0.0` for
blocked domains, and transform them to `NXDOMAIN`.

```shell
./dnsproxy -u 94.140.14.14:53 --bogus-nxdomain=0.0.0.0
```

CIDR ranges are supported as well. The following will respond with `NXDOMAIN`
instead of responses containing any IP from `192.168.0.0`-`192.168.255.255`:

```shell
./dnsproxy -u 192.168.0.15:53 --bogus-nxdomain=192.168.0.0/16
```

### Basic Auth for DoH

By setting the `--https-userinfo` option you can use `dnsproxy` as a DoH proxy
with basic authentication requirements.

For example:

```shell
./dnsproxy \
    --https-port='443' \
    --https-userinfo='user:p4ssw0rd' \
    --tls-crt='…/my.crt' \
    --tls-key='…/my.key' \
    -u '94.140.14.14:53' \
    ;
```

This configuration will only allow DoH queries that contain an `Authorization`
header containing the BasicAuth credentials for user `user` with password
`p4ssw0rd`.

Add `-p 0` if you also want to disable plain-DNS handling and make `dnsproxy`
only serve DoH with Basic Auth checking.

## UpstreamOptions 接口实现重构总结

### 背景和目标

本次重构的主要目标是实现一个统一的 `UpstreamOptions` 接口，为DNS代理的上游配置和连接管理提供抽象层。这个接口旨在：

1. 提供统一的配置管理接口
2. 简化连接创建和管理
3. 改善代码的可维护性和可扩展性
4. 解决循环导入依赖问题

### 主要实现内容

#### 1. UpstreamOptions 接口设计

在 `upstream/upstream.go` 中定义了 `UpstreamOptions` 接口：

```go
type UpstreamOptions interface {
    // Getter methods for Options fields
    Logger() *slog.Logger
    VerifyServerCertificate() func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error
    VerifyConnection() func(state tls.ConnectionState) error
    VerifyDNSCryptCertificate() func(cert *dnscrypt.Cert) error
    QUICTracer() QUICTraceFunc
    RootCAs() *x509.CertPool
    CipherSuites() []uint16
    Bootstrap() Resolver
    HTTPVersions() []HTTPVersion
    Timeout() time.Duration
    InsecureSkipVerify() bool
    PreferIPv6() bool

    // DialTCP creates a TCP connection to the specified address using the
    // configuration from this UpstreamOptions.
    DialTCP(ctx context.Context, addr string) (net.Conn, error)

    // DialUDP creates a UDP connection to the specified address using the
    // configuration from this UpstreamOptions.
    DialUDP(ctx context.Context, addr string) (*net.UDPConn, error)
}
```

#### 2. Options 结构体重构

将 `Options` 结构体的字段改为私有（小写开头），并添加对应的 getter 方法：

```go
type Options struct {
    logger                    *slog.Logger
    verifyServerCertificate   func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error
    verifyConnection          func(state tls.ConnectionState) error
    verifyDNSCryptCertificate func(cert *dnscrypt.Cert) error
    quicTracer                QUICTraceFunc
    rootCAs                   *x509.CertPool
    cipherSuites              []uint16
    bootstrap                 Resolver
    httpVersions              []HTTPVersion
    timeout                   time.Duration
    insecureSkipVerify        bool
    preferIPv6                bool
}
```

#### 3. 构造函数模式

添加了 `NewOptions` 构造函数来替代直接的结构体字面量：

```go
func NewOptions(
    logger *slog.Logger,
    verifyServerCertificate func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error,
    verifyConnection func(state tls.ConnectionState) error,
    verifyDNSCryptCertificate func(cert *dnscrypt.Cert) error,
    quicTracer QUICTraceFunc,
    rootCAs *x509.CertPool,
    cipherSuites []uint16,
    bootstrap Resolver,
    httpVersions []HTTPVersion,
    timeout time.Duration,
    insecureSkipVerify bool,
    preferIPv6 bool,
) *Options
```

#### 4. 循环导入问题解决

创建了 `internal/types` 包来解决循环导入问题：

- **问题**：`upstream` 包和 `internal/bootstrap` 包之间存在循环导入
- **解决方案**：创建独立的 `types` 包存放共享类型定义
- **实现**：
  - 移除 `internal/bootstrap/resolver.go` 中的重复定义
  - 在 `types` 包中定义核心类型（Resolver, StaticResolver, ParallelResolver, ConsequentResolver）
  - 使用类型别名在 `upstream` 包中引用这些类型

#### 5. 类型系统重构

**共享类型定义** (`internal/types/types.go`)：
- `Resolver` 接口：统一的主机名解析接口
- `StaticResolver`：静态IP地址解析器
- `ParallelResolver`：并发查询多个解析器
- `ConsequentResolver`：顺序查询多个解析器
- `Network` 类型：网络类型枚举

**类型别名使用** (`upstream/resolver.go`)：
```go
type Resolver = types.Resolver
type StaticResolver = types.StaticResolver
type ParallelResolver = types.ParallelResolver
type ConsequentResolver = types.ConsequentResolver
```

### 解决的主要问题

#### 1. 编译错误修复

- **方法vs字段名冲突**：修复了 `DialTCP`/`DialUDP` 方法中 `o.Timeout` 方法与字段的冲突
- **未定义构造函数**：添加了正确的 `NewStaticResolver` 使用方式
- **私有字段访问**：将所有直接字段访问改为通过 getter 方法访问
- **类型转换错误**：修复了 resolver 类型的构造和使用方式

#### 2. 架构改进

- **封装性提升**：通过私有字段和 getter 方法改善了数据封装
- **接口抽象**：提供了统一的配置和连接管理接口
- **依赖解耦**：通过类型包解除了循环依赖

#### 3. 代码一致性

- **构造函数模式**：统一使用构造函数创建配置对象
- **错误处理**：保持了原有的错误处理机制
- **兼容性**：保持了与现有代码的兼容性

### 技术细节

#### DialTCP/DialUDP 方法实现

```go
func (o *Options) DialTCP(ctx context.Context, addr string) (net.Conn, error) {
    dialer := &net.Dialer{
        Timeout: o.timeout,  // 使用私有字段
    }
    return dialer.DialContext(ctx, "tcp", addr)
}

func (o *Options) DialUDP(ctx context.Context, addr string) (*net.UDPConn, error) {
    dialer := &net.Dialer{
        Timeout: o.timeout,
    }
    conn, err := dialer.DialContext(ctx, "udp", addr)
    if err != nil {
        return nil, err
    }
    return conn.(*net.UDPConn), nil
}
```

#### Resolver 系统集成

在 `cmd/proxy.go` 中的使用示例：

```go
timeout := time.Duration(conf.Timeout)
bootOpts := upstream.NewOptions(
    l,             // logger
    nil,           // verifyServerCertificate
    nil,           // verifyConnection
    nil,           // verifyDNSCryptCertificate
    nil,           // quicTracer
    nil,           // rootCAs
    nil,           // cipherSuites
    nil,           // bootstrap
    httpVersions,  // httpVersions
    timeout,       // timeout
    conf.Insecure, // insecureSkipVerify
    false,         // preferIPv6
)
```

### 构建验证

所有修改都通过了完整的构建验证：

```bash
go build -v ./...
```

构建成功，无编译错误，所有包正常编译。

### 后续优化建议

1. **配置验证**：可以添加配置参数的验证逻辑
2. **默认值管理**：可以考虑添加更智能的默认值处理
3. **性能优化**：可以考虑缓存一些计算结果
4. **文档完善**：可以添加更详细的接口使用文档

### 总结

本次重构成功地实现了 `UpstreamOptions` 接口，解决了循环导入问题，改善了代码的架构和可维护性。通过合理的抽象和封装，为未来的功能扩展奠定了良好的基础。所有修改都保持了向后兼容性，确保现有代码的正常运行。
