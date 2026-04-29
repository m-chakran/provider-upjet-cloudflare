package certificate

import "github.com/crossplane/upjet/v2/pkg/config"

const shortGroup = "certificate"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("cloudflare_custom_hostname", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "CustomHostname"
		r.References["zone_id"] = config.Reference{
			TerraformName: "cloudflare_zone",
		}
	})

	p.AddResourceConfigurator("cloudflare_custom_hostname_fallback_origin", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "CustomHostnameFallbackOrigin"
		r.References["zone_id"] = config.Reference{
			TerraformName: "cloudflare_zone",
		}
	})

	p.AddResourceConfigurator("cloudflare_custom_ssl", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "CustomSSL"
		r.References["zone_id"] = config.Reference{
			TerraformName: "cloudflare_zone",
		}
	})

	p.AddResourceConfigurator("cloudflare_origin_ca_certificate", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "OriginCACertificate"
	})

	p.AddResourceConfigurator("cloudflare_certificate_pack", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "CertificatePack"
		r.References["zone_id"] = config.Reference{
			TerraformName: "cloudflare_zone",
		}
	})

	p.AddResourceConfigurator("cloudflare_keyless_certificate", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "KeylessCertificate"
		r.References["zone_id"] = config.Reference{
			TerraformName: "cloudflare_zone",
		}
	})

	p.AddResourceConfigurator("cloudflare_mtls_certificate", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "MTLSCertificate"
		r.References["account_id"] = config.Reference{
			TerraformName: "cloudflare_account",
		}
	})

	p.AddResourceConfigurator("cloudflare_total_tls", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "TotalTLS"
		r.References["zone_id"] = config.Reference{
			TerraformName: "cloudflare_zone",
		}
	})

	p.AddResourceConfigurator("cloudflare_hostname_tls_setting", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "HostnameTLSSetting"
		r.References["zone_id"] = config.Reference{
			TerraformName: "cloudflare_zone",
		}
		// 'value' is a Dynamic-typed field. Marking it required would cause upjet
		// to emit a CEL rule referencing self.forProvider.value, which Kubernetes
		// rejects because the field has no typed schema (it uses
		// x-kubernetes-preserve-unknown-fields). Force it optional so no rule is
		// emitted; the underlying API still validates the value server-side.
		if s, ok := r.TerraformResource.Schema["value"]; ok {
			s.Required = false
			s.Optional = true
		}
	})

	p.AddResourceConfigurator("cloudflare_hostname_tls_setting_ciphers", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "HostnameTLSSettingCiphers"
		r.References["zone_id"] = config.Reference{
			TerraformName: "cloudflare_zone",
		}
	})

	p.AddResourceConfigurator("cloudflare_authenticated_origin_pulls", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "AuthenticatedOriginPulls"
		r.References["zone_id"] = config.Reference{
			TerraformName: "cloudflare_zone",
		}
	})

	p.AddResourceConfigurator("cloudflare_authenticated_origin_pulls_certificate", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "AuthenticatedOriginPullsCertificate"
		r.References["zone_id"] = config.Reference{
			TerraformName: "cloudflare_zone",
		}
	})

	p.AddResourceConfigurator("cloudflare_ssl_recommendation", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "SSLRecommendation"
		r.References["zone_id"] = config.Reference{
			TerraformName: "cloudflare_zone",
		}
	})

	p.AddResourceConfigurator("cloudflare_custom_hostname_certificate_pack", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "CustomHostnameCertificatePack"
		r.References["zone_id"] = config.Reference{
			TerraformName: "cloudflare_zone",
		}
	})
}
