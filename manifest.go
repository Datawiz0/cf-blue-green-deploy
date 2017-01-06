package main

import (
	"code.cloudfoundry.org/cli/plugin/models"
	"fmt"
	"github.com/bluemixgaragelondon/cf-blue-green-deploy/from-cf-codebase/manifest"
	"github.com/bluemixgaragelondon/cf-blue-green-deploy/from-cf-codebase/models"
)

type ManifestReader func(manifest.Repository, string) *models.AppParams

type ManifestAppFinder struct {
	Repo    manifest.Repository
	ManifestPath string
	AppName string
}

func (f *ManifestAppFinder) RoutesFromManifest(defaultDomain string) []plugin_models.GetApp_RouteSummary {
	if appParams := f.AppParams(); appParams != nil {

		manifestRoutes := make([]plugin_models.GetApp_RouteSummary, 0)

		for _, host := range appParams.Hosts {
			if appParams.Domains == nil {
				manifestRoutes = append(manifestRoutes, plugin_models.GetApp_RouteSummary{Host: host, Domain: plugin_models.GetApp_DomainFields{Name: defaultDomain}})
				continue
			}

			for _, domain := range appParams.Domains {
				manifestRoutes = append(manifestRoutes, plugin_models.GetApp_RouteSummary{Host: host, Domain: plugin_models.GetApp_DomainFields{Name: domain}})
			}
		}

		return manifestRoutes
	}
	return nil
}

func (f *ManifestAppFinder) AppParams() *models.AppParams {
	var manifest *manifest.Manifest
	var err error
	if f.ManifestPath == "" {
		manifest, err = f.Repo.ReadManifest("./")
	} else {
		manifest, err = f.Repo.ReadManifest(f.ManifestPath)
	}

	if err != nil {
		return nil
	}

	apps, err := manifest.Applications()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for index, app := range apps {
		if app.IsHostEmpty() {
			continue
		}

		if app.Name != nil && *app.Name != f.AppName {
			continue
		}

		return &apps[index]
	}

	return nil
}
