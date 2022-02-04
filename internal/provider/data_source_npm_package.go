package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mohamed-gara/terraform-provider-npm/internal/client"
	"os"
	"path/filepath"
	"sort"

	"github.com/gabriel-vasile/mimetype"
)

func dataSourceNpmPackage() *schema.Resource {
	return &schema.Resource{
		Description: "A data source to download a package from an npm registry and access to its file list (using the files attribute).",
		ReadContext: dataSourceNpmPackageRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Package name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"version": {
				Description: "Package version",
				Type:        schema.TypeString,
				Required:    true,
			},
			"files": {
				Description: "Package files",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"absolute_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mime_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNpmPackageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	apiClient := meta.(*client.NpmRegistry)

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	e := apiClient.DownloadPackage(name, version, ".")
	if e != nil {
		return diag.FromErr(e)
	}

	filesNames, e := downloadedFiles()
	if e != nil {
		return diag.FromErr(e)
	}

	d.SetId(name + "-" + version)
	err := d.Set("files", filesNames)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func downloadedFiles() ([]map[string]string, error) {
	var filesName []map[string]string
	rootPath := "./package" //TODO: should be extracted to a constant
	e := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			absolutePath, pathErr := filepath.Abs(path)
			if pathErr != nil {
				return pathErr
			}

			mtype, err := mimetype.DetectFile(absolutePath)
			if err != nil {
				return err
			}

			fileObject := map[string]string{
				"absolute_path": absolutePath,
				"mime_type":     mtype.String(),
			}
			filesName = append(filesName, fileObject)
		}
		return nil
	})

	sort.Slice(filesName, func(i, j int) bool {
		return filesName[j]["absolte_path"] < filesName[i]["absolte_path"]
	})

	return filesName, e
}
