package azuread

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// Session info
type Session struct {
	TenantID   string
	Authorizer auth.Authorizer
}

// GetNewSession creates an session configured from environment variables/CLI in the order:
// 1. Client credentials
// 2. Client certificate
// 3. Username password
// 4. MSI
// 5. CLI
func GetNewSession(ctx context.Context, d *plugin.QueryData) (sess *Session, err error) {
	logger := plugin.Logger(ctx)

	// have we already created and cached the session?
	sessionCacheKey := "GetNewSession"
	if cachedData, ok := d.ConnectionManager.Cache.Get(sessionCacheKey); ok {
		return cachedData.(*Session), nil
	}

	azureADConfig := GetConfig(d.Connection)
	var tenantID string
	authMethod, authConfig, err := getApplicableAuthorizationDetails(ctx, azureADConfig)
	if err != nil {
		logger.Debug("GetNewSession__", "getApplicableAuthorizationDetails error", err)
		return nil, err
	}

	if authConfig.TenantID != "" {
		tenantID = authConfig.TenantID
	}

	// logger.Debug("GetNewSession", "TenantID", authConfig.TenantID)
	// logger.Debug("GetNewSession", "ClientCertData", authConfig.ClientCertData)
	// logger.Debug("GetNewSession", "ClientCertPassword", authConfig.ClientCertPassword)
	// logger.Debug("GetNewSession", "ClientCertPath", authConfig.ClientCertPath)
	// logger.Debug("GetNewSession", "ClientID", authConfig.ClientID)
	// logger.Debug("GetNewSession", "ClientSecret", authConfig.ClientSecret)
	// logger.Debug("GetNewSession", "EnableAzureCliToken", authConfig.EnableAzureCliToken)
	// logger.Debug("GetNewSession", "EnableClientCertAuth", authConfig.EnableClientCertAuth)
	// logger.Debug("GetNewSession", "EnableClientSecretAuth", authConfig.EnableClientSecretAuth)
	// logger.Debug("GetNewSession", "AzureADEndpoint", authConfig.Environment.AzureADEndpoint)
	// logger.Debug("GetNewSession", "EnableMsiAuth", authConfig.EnableMsiAuth)
	// logger.Debug("GetNewSession", "MsiEndpoint", authConfig.MsiEndpoint)

	authorizer, err := authConfig.NewAuthorizer(ctx, auth.MsGraph)
	if err != nil {
		log.Fatal(err)
	}

	if authMethod == "CLI" {
		tenantID, err = getTenantFromCLI()
		if err != nil {
			logger.Debug("GetNewSession__", "getTenantFromCLI error", err)
			return nil, err
		}
	}

	sess = &Session{
		Authorizer: authorizer,
		TenantID:   tenantID,
	}

	// Save session into cache
	d.ConnectionManager.Cache.Set(sessionCacheKey, sess)

	return sess, err
}

func getApplicableAuthorizationDetails(ctx context.Context, config azureADConfig) (authMethod string, authConfig auth.Config, err error) {

	var environment, tenantID, clientID, clientSecret, certificatePath, certificatePassword string
	// username, password string
	if config.TenantID != nil {
		tenantID = *config.TenantID
	} else {
		tenantID = os.Getenv("AZURE_TENANT_ID")
	}

	if config.Environment != nil {
		environment = *config.Environment
	} else {
		environment = os.Getenv("AZURE_ENVIRONMENT")
	}

	// Can be	"AZURECHINACLOUD", "AZUREGERMANCLOUD", "AZUREPUBLICCLOUD", "AZUREUSGOVERNMENTCLOUD"
	switch environment {
	case "AZURECHINACLOUD":
		authConfig.Environment = environments.China
	case "AZUREUSGOVERNMENTCLOUD":
		authConfig.Environment = environments.USGovernmentL4
	case "AZUREGERMANCLOUD":
		authConfig.Environment = environments.Germany
	default:
		authConfig.Environment = environments.Global
	}

	// 1. Client credentials
	if config.ClientID != nil {
		clientID = *config.ClientID
	} else {
		clientID = os.Getenv("AZURE_CLIENT_ID")
	}

	if config.ClientSecret != nil {
		clientSecret = *config.ClientSecret
	} else {
		clientSecret = os.Getenv("AZURE_CLIENT_SECRET")
	}

	// 2. Client certificate
	if config.CertificatePath != nil {
		certificatePath = *config.CertificatePath
	} else {
		certificatePath = os.Getenv("AZURE_CERTIFICATE_PATH")
	}

	if config.CertificatePassword != nil {
		certificatePassword = *config.CertificatePassword
	} else {
		certificatePassword = os.Getenv("AZURE_CERTIFICATE_PASSWORD")
	}

	// 3. Username password
	// if config.Username != nil {
	// 	username = *config.Username
	// } else {
	// 	username = os.Getenv("AZURE_USERNAME")
	// }

	// if config.Password != nil {
	// 	password = *config.Password
	// } else {
	// 	password = os.Getenv("AZURE_PASSWORD")
	// }

	// hamiltonauth.
	// authConfig := auth.Config{}
	// Environment: environments.Global,
	// TenantID:               tenantId,
	// ClientID:               clientId,
	// ClientSecret:           clientSecret,
	// EnableClientSecretAuth: true,
	// }

	authMethod = "CLI"
	if tenantID == "" {
		authMethod = "CLI"
		authConfig.EnableAzureCliToken = true
	} else if tenantID != "" && clientID != "" && clientSecret != "" {
		authConfig.TenantID = tenantID
		authConfig.ClientID = clientID
		authConfig.ClientSecret = clientSecret
		authConfig.EnableClientSecretAuth = true
		authMethod = "EnableClientSecretAuth"
	} else if tenantID != "" && certificatePath != "" && certificatePassword != "" {
		authConfig.TenantID = tenantID
		authConfig.ClientCertPath = certificatePath
		authConfig.ClientCertPassword = certificatePassword
		authConfig.EnableClientCertAuth = true
		authMethod = "EnableClientCertificateAuth"
	}
	// else if username != "" && password != "" {
	// 	authConfig. = true
	// 	authMethod = "EnableUserPasswordAuth"
	// }
	return
}

// https://github.com/Azure/go-autorest/blob/3fb5326fea196cd5af02cf105ca246a0fba59021/autorest/azure/cli/token.go#L126
// NewAuthorizerFromCLIWithResource creates an Authorizer configured from Azure CLI 2.0 for local development scenarios.
func getTenantFromCLI() (string, error) {
	// This is the path that a developer can set to tell this class what the install path for Azure CLI is.
	const azureCLIPath = "AzureCLIPath"

	// The default install paths are used to find Azure CLI. This is for security, so that any path in the calling program's Path environment is not used to execute Azure CLI.
	azureCLIDefaultPathWindows := fmt.Sprintf("%s\\Microsoft SDKs\\Azure\\CLI2\\wbin; %s\\Microsoft SDKs\\Azure\\CLI2\\wbin", os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramFiles"))

	// Default path for non-Windows.
	const azureCLIDefaultPath = "/bin:/sbin:/usr/bin:/usr/local/bin"

	// Execute Azure CLI to get token
	var cliCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cliCmd = exec.Command(fmt.Sprintf("%s\\system32\\cmd.exe", os.Getenv("windir")))
		cliCmd.Env = os.Environ()
		cliCmd.Env = append(cliCmd.Env, fmt.Sprintf("PATH=%s;%s", os.Getenv(azureCLIPath), azureCLIDefaultPathWindows))
		cliCmd.Args = append(cliCmd.Args, "/c", "az")
	} else {
		cliCmd = exec.Command("az")
		cliCmd.Env = os.Environ()
		cliCmd.Env = append(cliCmd.Env, fmt.Sprintf("PATH=%s:%s", os.Getenv(azureCLIPath), azureCLIDefaultPath))
	}
	cliCmd.Args = append(cliCmd.Args, "account", "get-access-token", "--resource-type=ms-graph", "-o", "json")

	var stderr bytes.Buffer
	cliCmd.Stderr = &stderr

	output, err := cliCmd.Output()
	if err != nil {
		return "", fmt.Errorf("Invoking Azure CLI failed with the following error: %v", err)
	}

	var tokenResponse struct {
		AccessToken string `json:"accessToken"`
		ExpiresOn   string `json:"expiresOn"`
		Tenant      string `json:"tenant"`
		TokenType   string `json:"tokenType"`
	}
	err = json.Unmarshal(output, &tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.Tenant, nil
}
