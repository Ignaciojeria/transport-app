package usecase

import _ "embed"

//go:embed tenant_repo_template.md
var TenantRepoReadmeTemplate string

//go:embed tenant_config_template.json
var TenantConfigTemplate string

//go:embed gitignore_template
var GitignoreTemplate string

//go:embed github_workflow_template.yml
var GitHubWorkflowTemplate string

//go:embed index_template.html
var IndexTemplate string

//go:embed firebase_template.json
var FirebaseTemplate string

//go:embed deploy_readme_template.md
var DeployReadmeTemplate string
