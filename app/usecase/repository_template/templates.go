package usecase

import _ "embed"

//go:embed tenant_repo_template.md
var TenantRepoReadmeTemplate string

//go:embed tenant_config_template.json
var TenantConfigTemplate string

//go:embed gitignore_template
var GitignoreTemplate string
