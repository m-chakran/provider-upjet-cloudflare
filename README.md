# Upjet-based Crossplane Provider for Cloudflare

<div align="center">

[![CI](https://github.com/m-chakran/provider-upjet-cloudflare/actions/workflows/ci.yml/badge.svg)](https://github.com/m-chakran/provider-upjet-cloudflare/actions/workflows/ci.yml)
[![Tag](https://img.shields.io/github/v/tag/m-chakran/provider-upjet-cloudflare?sort=semver)](https://github.com/m-chakran/provider-upjet-cloudflare/tags)
[![Go Version](https://img.shields.io/github/go-mod/go-version/m-chakran/provider-upjet-cloudflare)](https://github.com/m-chakran/provider-upjet-cloudflare/blob/main/go.mod)

</div>

`provider-upjet-cloudflare` is a [Crossplane](https://crossplane.io/) provider
built with [Upjet](https://github.com/crossplane/upjet) on top of the official
[Cloudflare Terraform provider](https://github.com/cloudflare/terraform-provider-cloudflare).
It exposes ~190 Cloudflare resources as Kubernetes managed resources covering
DNS, zones, certificates, Workers, R2, Zero Trust, WAF, Load Balancing, and
more.

## Getting Started

### Prerequisites

- Kubernetes cluster running [Crossplane](https://crossplane.io/) v2.1+
- A Cloudflare account
- A scoped [Cloudflare API token](https://developers.cloudflare.com/fundamentals/api/get-started/create-token/)

### Install the Provider

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-upjet-cloudflare
spec:
  package: ghcr.io/m-chakran/provider-upjet-cloudflare:v0.1.1
```

```bash
kubectl wait provider.pkg/provider-upjet-cloudflare \
  --for=condition=Healthy --timeout=5m
```

### Create Credentials Secret

Create a Cloudflare API token from the dashboard
(**My Profile → API Tokens → Create Token**) with the permissions you need
(typically `Zone:Read`, `DNS:Edit`, etc.), then store it in a secret:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cloudflare-creds
  namespace: crossplane-system
type: Opaque
stringData:
  credentials: |
    { "api_token": "YOUR_CLOUDFLARE_API_TOKEN" }
```

The `credentials` value is a JSON object whose keys map to Cloudflare provider
attributes. Supported keys:

| Key                    | Description                                                       |
|------------------------|-------------------------------------------------------------------|
| `api_token`            | **Recommended.** Scoped API token.                                |
| `api_key` + `email`    | Legacy global API key. Both keys must be present together.        |
| `api_user_service_key` | Restricted-use service key for a small set of endpoints.          |

Provide exactly one of `api_token`, `api_key`, or `api_user_service_key`.

### Create ProviderConfig

```yaml
apiVersion: cloudflare.cloudflare.com/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: cloudflare-creds
      key: credentials
```

### Smoke Test

Create a DNS record on a zone you own:

```yaml
apiVersion: dns.cloudflare.crossplane.io/v1alpha1
kind: DNSRecord
metadata:
  name: test-record
spec:
  forProvider:
    zoneId: YOUR_ZONE_ID
    name: test
    type: A
    content: 192.0.2.1
    ttl: 300
  providerConfigRef:
    name: default
```

Then watch it reconcile:

```bash
kubectl get managed
kubectl describe dnsrecord test-record
```

## Supported Resources

This provider exposes ~190 resources across the following areas. For the full
list, see [`config/external_name.go`](config/external_name.go); for per-field
schemas see the CRD reference at
[doc.crds.dev](https://doc.crds.dev/github.com/m-chakran/provider-upjet-cloudflare).

| Area                | Examples |
|---------------------|----------|
| **Account & API**   | Account, AccountMember, AccountToken, APIToken |
| **Zone & DNS**      | Zone, ZoneSetting, ZoneDNSSEC, DNSRecord, DNSSettings |
| **Certificates & TLS** | CustomHostname, CustomSSL, OriginCACertificate, CertificatePack, KeylessCertificate, MTLSCertificate, TotalTLS, HostnameTLSSetting, AuthenticatedOriginPulls |
| **WAF & Security**  | Filter, FirewallRule, RateLimit, WAFRule, WAFOverride, Ruleset, ManagedHeaders, BotManagement, LeakedCredentialCheck, PageRule |
| **Access (Zero Trust)** | AccessApplication, AccessPolicy, AccessGroup, AccessIdentityProvider, AccessServiceToken, AccessOrganization |
| **Zero Trust**      | ZeroTrustList, ZeroTrustGatewayPolicy, ZeroTrustTunnelCloudflared, ZeroTrustDevicePostureRule, ZeroTrustDLPProfile, TeamsAccount/List/Location/Rule |
| **Load Balancing**  | LoadBalancer, LoadBalancerPool, LoadBalancerMonitor |
| **Workers & Pages** | WorkersScript, WorkerRoute, WorkerCronTrigger, WorkerDomain, WorkerSecret, WorkersKV, WorkersKVNamespace, PagesProject, PagesDomain |
| **R2 (Object Storage)** | R2Bucket, R2CustomDomain, R2BucketLifecycle, R2BucketCORS, R2BucketSippy, R2BucketLock |
| **Argo & Routing**  | Argo, ArgoTunnel, TieredCache, RegionalTieredCache, StaticRoute, GRETunnel, IPsecTunnel, MagicWANGRETunnel |
| **Email Routing**   | EmailRoutingAddress, EmailRoutingCatchAll, EmailRoutingDNS, EmailRoutingRule, EmailRoutingSettings |
| **Stream**          | Stream, StreamKey, StreamLiveInput, StreamWebhook, StreamCaption, StreamWatermark |
| **API Shield**      | APIShield, APIShieldOperation, APIShieldSchema, APIShieldOperationSchemaValidationSettings |
| **Other**           | Spectrum, LogpushJob, NotificationPolicy, Healthcheck, WaitingRoom, AddressMap, BYOIPPrefix, RegionalHostname, Snippet, TurnstileWidget, WebAnalyticsSite, Queue, HyperdriveConfig, AIGateway, D1Database, VectorizeIndex, BrowserIsolation*, ContentScanning, Zaraz |

Each resource is generated in two scopes:
- **Cluster-scoped** under `<group>.cloudflare.crossplane.io`
- **Namespaced** under `<group>.cloudflare.m.cloudflare.com` (Crossplane v2)

## Documentation

- [Examples](examples/) — generated YAML examples for every resource
- [CRD Reference](https://doc.crds.dev/github.com/m-chakran/provider-upjet-cloudflare) — rendered schemas
- [Cloudflare Terraform provider docs](https://registry.terraform.io/providers/cloudflare/cloudflare/5.15.0/docs) — argument reference for every resource

## Build Locally

```bash
# Initialize the build submodule (first time only)
git submodule update --init --recursive

# Run code generation
go run cmd/generator/main.go "$PWD"

# Build the provider binary and xpkg
make build

# Or build everything (multi-arch images + xpkg)
make -j2 build.all BUILD_ARGS=--load
```

The xpkg lands under `_output/xpkg/<platform>/`.

### Iterate without GitHub Actions (build → push → deploy)

For fast iteration against an already-installed Crossplane cluster (kubectl
context must point at the right cluster):

```bash
# One-time: log Docker into ghcr.io with a PAT that has read:packages, write:packages
echo "$GHCR_PAT" | docker login ghcr.io -u <your-gh-username> --password-stdin

# Build the xpkg, push under a dev tag (defaults to dev-<short-sha>),
# patch the Provider in the cluster to use it, and wait for Healthy.
make dev-deploy

# Or, if you've already built once:
make dev-push          # push current _output/xpkg with the dev tag
make dev-deploy        # rebuild + push + apply (idempotent)

# Tail the provider logs (filtering controller-runtime noise):
make dev-tail
```

Override the tag or owner if needed:

```bash
make dev-deploy DEV_TAG=test-foo GHCR_OWNER=my-org
```

## Troubleshooting

### Check provider status

```bash
kubectl get providers
kubectl get providerrevision -l pkg.crossplane.io/provider=provider-upjet-cloudflare -o wide
kubectl logs -n crossplane-system -l pkg.crossplane.io/provider=provider-upjet-cloudflare
```

### Check resource status

```bash
kubectl get managed
kubectl describe <kind> <name>
```

The `Synced` and `Ready` conditions surface upstream Cloudflare API errors
verbatim — usually permission/scope issues on the API token.

### Common issues

1. **`Healthy=False, reason=UnhealthyPackageRevision`** — A CRD failed to
   install. Run `kubectl describe providerrevision -l pkg.crossplane.io/provider=provider-upjet-cloudflare`
   to see which CRD and why.
2. **`401 Unauthorized` on resource creation** — The API token doesn't have
   the required permissions. Re-create the token with the right scopes.
3. **Reference resolution errors (`zoneIdRef`, `accountIdRef`, …)** — Ensure
   the referenced resource exists and is `Ready`.
4. **Resource stuck in `Creating`** — Check the Cloudflare dashboard; a stale
   resource with the same name may already exist.

## Adding new resources

If a Cloudflare resource exists in the upstream Terraform provider but isn't
exposed here yet:

1. Find it in the [Terraform provider docs](https://registry.terraform.io/providers/cloudflare/cloudflare/5.15.0/docs/resources).
2. Add a line in `config/external_name.go` (and a configurator under `config/cluster/<area>/` and `config/namespaced/<area>/` if the resource needs custom kind/group/references).
3. Run `make generate`.
4. Commit `apis/`, `package/crds/`, `examples-generated/`, and the config changes.
5. Tag and publish.

## Reporting Issues

Open an issue at
[m-chakran/provider-upjet-cloudflare/issues](https://github.com/m-chakran/provider-upjet-cloudflare/issues).

## License

Apache 2.0. See [LICENSE](LICENSE).
