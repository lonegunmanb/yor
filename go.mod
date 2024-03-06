module github.com/bridgecrewio/yor

go 1.19

require (
	github.com/awslabs/goformation/v5 v5.2.7
	github.com/bridgecrewio/goformation/v5 v5.0.0-20210823083242-84a6d242099f
	github.com/go-git/go-git/v5 v5.11.0
	github.com/google/uuid v1.3.0
	github.com/hashicorp/hcl/v2 v2.8.2
	github.com/hashicorp/terraform-config-inspect v0.0.0-20211115214459-90acf1ca460f
	github.com/hashicorp/terraform-json v0.21.0
	github.com/lonegunmanb/terraform-alicloud-schema v1.218.0
	github.com/lonegunmanb/terraform-aws-schema/v5 v5.39.1
	github.com/lonegunmanb/terraform-azurerm-schema/v3 v3.94.0
	github.com/lonegunmanb/terraform-google-schema/v4 v4.84.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pmezard/go-difflib v1.0.0
	github.com/sanathkr/yaml v1.0.0
	github.com/stretchr/testify v1.8.4
	github.com/urfave/cli/v2 v2.3.0
	github.com/zclconf/go-cty v1.14.1
	go.opencensus.io v0.24.0
	gopkg.in/validator.v2 v2.0.0-20200605151824-2b28d334fa05
	gopkg.in/yaml.v2 v2.4.0
)

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/ProtonMail/go-crypto v1.1.0-alpha.0 // indirect
	github.com/agext/levenshtein v1.2.2 // indirect
	github.com/apparentlymart/go-textseg/v12 v12.0.0 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0-20190314233015-f79a8a8ca69d // indirect
	github.com/cyphar/filepath-securejoin v0.2.4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/sanathkr/go-yaml v0.0.0-20170819195128-ed9d249f429b // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/skeema/knownhosts v1.2.1 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/mod v0.15.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.13.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/awslabs/goformation/v5 => github.com/bridgecrewio/goformation/v5 v5.0.0-20210823081757-99ed9bf3c0e5
	github.com/go-git/go-git/v5 => github.com/lonegunmanb/go-git/v5 v5.11.2
)
