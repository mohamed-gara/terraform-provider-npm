package provider

import (
	"context"
	client "github.com/mohamed-gara/terraform-provider-npm/internal/client"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"url": {
					Description: "The npm registry URL. The value can be defined using the TERRAFORM_PROVIDER_NPM_URL environment variable.",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("TERRAFORM_PROVIDER_NPM_URL", nil),
				},
				"username": {
					Description: "The username to use for authentication. The username should be defined with the TERRAFORM_PROVIDER_NPM_USERNAME environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("TERRAFORM_PROVIDER_NPM_USERNAME", nil),
				},
				"password": {
					Description: "The password to use for authentication. The password should be defined with the TERRAFORM_PROVIDER_NPM_PASSWORD environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("TERRAFORM_PROVIDER_NPM_PASSWORD", nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"npm_package": dataSourceNpmPackage(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

/*
type apiClient struct {
	// Add whatever fields, client or connection info, etc. here
	// you would need to setup to communicate with the upstream
	// API.
}
*/

func configure(_ string, _ *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// Setup a User-Agent for your API client (replace the provider name for yours):
		// userAgent := p.UserAgent("terraform-provider-npm", version)
		// TODO: myClient.UserAgent = userAgent
		url := d.Get("url").(string)
		username := d.Get("username").(string)
		password := d.Get("password").(string)
		apiClient, err := client.NewNpmRegistry(&url, &password, &username)
		if err != nil {
			log.Fatal("Can't load API client!")
		}

		return apiClient, nil
	}
}
