package provider

import (
	"context"
	"fmt"

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
				"region": {
					Type:     schema.TypeString,
					Required: true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{
						"AWS_REGION",
						"AWS_DEFAULT_REGION",
					}, nil),
					Description: "AWS Region of Bref PHP runtime layers. Can be specified with the `AWS_REGION` " +
						"or `AWS_DEFAULT_REGION` environment variable.",
					InputDefault: "us-east-1",
				},
				"bref_version": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("BREF_VERSION", "1.5.0"),
					Description: "The Bref PHP runtime version to work with. Can be specified with the " +
						"`BREF_VERSION` environment variable.",
				},
				"bref_extra_version": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("BREF_EXTRA_VERSION", "0.11.34"),
					Description: "The Bref Extra PHP runtime version to work with. Can be specified with the " +
						"`BREF_EXTRA_VERSION` environment variable.",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"bref_lambda_layer":       dataSourceLambdaLayer(),
				"bref_extra_lambda_layer": extraDataSourceLambdaLayer(),
			},
			ResourcesMap: map[string]*schema.Resource{},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	Region       string
	Version      string
	ExtraVersion string
	AccountIds   map[string]string
	URLs         map[string]string
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		apiClient := apiClient{
			Region:       d.Get("region").(string),
			Version:      d.Get("bref_version").(string),
			ExtraVersion: d.Get("bref_extra_version").(string),
			AccountIds: map[string]string{
				"bref_lambda_layer":       "209497400698",
				"bref_extra_lambda_layer": "403367587399",
			},
			URLs: map[string]string{
				"bref_lambda_layer":       fmt.Sprintf("https://raw.githubusercontent.com/brefphp/bref/%s/layers.json", d.Get("bref_version").(string)),
				"bref_extra_lambda_layer": fmt.Sprintf("https://raw.githubusercontent.com/brefphp/extra-php-extensions/%s/layers.json", d.Get("bref_extra_version").(string)),
			},
		}

		return &apiClient, nil
	}
}
