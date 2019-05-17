package adapter

import (
	"log"

	"github.com/pivotal-cf/on-demand-services-sdk/bosh"

	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

type TestServiceManifestGenerator struct {
	Logger *log.Logger
}

const stemcellAlias = "ubuntu"

func (mg TestServiceManifestGenerator) GenerateManifest(params serviceadapter.GenerateManifestParams) (serviceadapter.GenerateManifestOutput, error) {

	mg.Logger.Printf("Generating Manifest for Test Service...")

	igToJobMap := map[string][]string{
		"nginx": {"nginx"},
	}

	instanceGroups, err := serviceadapter.GenerateInstanceGroupsWithNoProperties(
		params.Plan.InstanceGroups,
		params.ServiceDeployment.Releases,
		stemcellAlias,
		igToJobMap,
	)

	if err != nil {
		mg.Logger.Printf("Error \"%s\" occured", err)
		return serviceadapter.GenerateManifestOutput{}, err
	}

	//Edit InstanceGroups
	instanceGroups[0].Jobs[0].Properties = map[string]interface{}{
		"nginx_conf": `user nobody vcap; # group vcap can read most directories
            worker_processes  1;
            error_log /var/vcap/sys/log/nginx/error.log   info;
            events {
              worker_connections  1024;
            }
            http {
              include /var/vcap/packages/nginx/conf/mime.types;
              default_type  application/octet-stream;
              sendfile        on;
              ssi on;
              keepalive_timeout  65;
              server_names_hash_bucket_size 64;
              server {
                server_name _; # invalid value which will never trigger on a real hostname.
                listen *:80;
                # FIXME: replace all occurrences of 'example.com' with your server's FQDN
                access_log /var/vcap/sys/log/nginx/example.com-access.log;
                error_log /var/vcap/sys/log/nginx/example.com-error.log;
                root /var/vcap/data/nginx/document_root;
                index index.shtml;
              }
            }`,
		"pre_start": `#!/bin/bash -ex
            NGINX_DIR=/var/vcap/data/nginx/document_root
            if [ ! -d $NGINX_DIR ]; then
              mkdir -p $NGINX_DIR
              cd $NGINX_DIR
              cat > index.shtml <<EOF
                <html><head><title>BOSH on IPv6</title>
                </head><body>
                <h2>Welcome to BOSH's nginx Release</h2>
                <h2>
                My hostname/IP: <b><!--# echo var="HTTP_HOST" --></b><br />
                Your IP: <b><!--# echo var="REMOTE_ADDR" --></b>
                </h2>
                </body></html>
            EOF
            fi`,
	}

	//Compose Output
	output := serviceadapter.GenerateManifestOutput{
		Manifest: bosh.BoshManifest{
			//	Addons         []bosh.Addon
			Name: params.ServiceDeployment.DeploymentName,
			Releases: []bosh.Release{{
				Name:    params.ServiceDeployment.Releases[0].Name,
				Version: params.ServiceDeployment.Releases[0].Version,
			}},
			Stemcells: []bosh.Stemcell{{
				Alias:   stemcellAlias,
				Version: params.ServiceDeployment.Stemcells[0].Version,
				OS:      params.ServiceDeployment.Stemcells[0].OS,
			}},
			InstanceGroups: instanceGroups,

			Update: &bosh.Update{
				Canaries:        params.Plan.Update.Canaries,
				MaxInFlight:     params.Plan.Update.MaxInFlight,
				Serial:          params.Plan.Update.Serial,
				CanaryWatchTime: params.Plan.Update.CanaryWatchTime,
				UpdateWatchTime: params.Plan.Update.UpdateWatchTime,
			},
			//	Update         *bosh.Update
			//	Properties     map[string]interface{}
			//	Variables      []bosh.Variable
			//	Tags           map[string]interface{}
			//	Features       bosh.BoshFeatures
		},
		// {Addons: , Name: , Releases: , Stemcells: , InstanceGroups: , Update: , Properties: , Variables: , Tags: , Features: },
	}

	return output, nil
}
