package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LinuxWebAppResource struct{}

func TestAccLinuxWebApp_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep(),
	})
}

// TODO - func TestAccLinuxWebApp_complete(t *testing.T)

// TODO - func TestAccLinuxWebApp_completeUpdate(t *testing.T)

func TestAccLinuxWebApp_backup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBackup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_backupUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBackup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_completeLogging(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withLoggingComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.detailed_error_logging").HasValue("true"),
				check.That(data.ResourceName).Key("logs.0.detailed_error_messages").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_completeLoggingUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withDetailedLogging(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.detailed_error_logging").HasValue("false"),
				check.That(data.ResourceName).Key("logs.0.detailed_error_messages").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withLoggingComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.detailed_error_logging").HasValue("true"),
				check.That(data.ResourceName).Key("logs.0.detailed_error_messages").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccLinuxWebApp_updateServicePlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.secondServicePlan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_logsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.logsEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_loadBalancing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.loadBalancing(data, "WeightedRoundRobin"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withIPRestrictions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withIPRestrictionsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withIPRestrictionsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAuthSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withAuthSettingsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withAuthSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withAuthSettingsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withAuthSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.virtual_application.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

// App stacks...
func TestAccLinuxWebApp_withDotNet21(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotNet(data, "2.1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withDotNet31(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotNet(data, "3.1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withDotNet50(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotNet(data, "5.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withPhp56(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.php(data, "5.6"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withPhp72(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.php(data, "7.2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withPhp73(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.php(data, "7.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withPhp74(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.php(data, "7.4"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withPython27(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.python(data, "2.7"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withPython36(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.python(data, "3.6"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withPython37(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.python(data, "3.7"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withPython38(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.python(data, "3.8"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withNode101(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.node(data, "10.1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withNode106(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.node(data, "10.6"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withNode1014(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.node(data, "10.14"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withNode10LTS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.node(data, "10-lts"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withNode12LTS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.node(data, "12-lts"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withNode14LTS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.node(data, "14-lts"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withJre8Java(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "jre8", "JAVA", "8"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withJre11Java(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "java11", "JAVA", "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("JAVA|11-java11"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withJava1109(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "11.0.9", "JAVA", ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("JAVA|11.0.9"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withJava8u242(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "8u242", "JAVA", ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("JAVA|8u242"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withJava11Tomcat9(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "java11", "TOMCAT", "9.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("TOMCAT|9.0-java11"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withJava11Tomcat9041(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "java11", "TOMCAT", "9.0.41"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("TOMCAT|9.0.41-java11"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withJava11Tomcat85(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "java11", "TOMCAT", "8.5"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("TOMCAT|8.5-java11"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withJava11Tomcat8561(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "java11", "TOMCAT", "8.5.61"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("TOMCAT|8.5.61-java11"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withJava8JBOSSEAP73(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.javaPremiumV3Plan(data, "java8", "JBOSSEAP", "7.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("JBOSSEAP|7.3-java8"),
			),
		},
		data.ImportStep(),
	})
}

// TODO - finish known Java matrix combination tests...?

func TestAccLinuxWebApp_withDocker(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.docker(data, "mcr.microsoft.com/appsvc/staticsite", "latest"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("DOCKER|mcr.microsoft.com/appsvc/staticsite:latest"),
			),
		},
		data.ImportStep(),
	})
}

// TODO - Needs more property tests for autoheal
func TestAccLinuxWebApp_withAutoHealRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoHealRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccLinuxWebApp_withAutoHealRulesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoHealRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoHealRulesUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_withAutoHealRulesStatusCodeRange(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoHealRulesStatusCodeRange(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("app_settings.secret").HasValue("sauce"),
			),
		},
	})
}

// Exists func

func (r LinuxWebAppResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.WebAppID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.WebAppsClient.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Linux Web App %s: %+v", id, err)
	}
	if utils.ResponseWasNotFound(resp.Response) {
		return utils.Bool(false), nil
	}
	return utils.Bool(true), nil
}

// Configs

func (r LinuxWebAppResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppResource) withBackup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  backup {
    name                = "acctest"
    storage_account_url = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
    schedule {
      frequency_interval = 1
      frequency_unit     = "Day"
    }
  }
}
`, r.templateWithStorageAccount(data), data.RandomInteger)
}

func (r LinuxWebAppResource) appSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  app_settings = {
    foo    = "bar"
    secret = "sauce"
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppResource) secondServicePlan(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_service_plan" "test2" {
  name                = "acctestASP2-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "B1"
}

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test2.id
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_linux_web_app" "import" {
  name                = azurerm_linux_web_app.test.name
  location            = azurerm_linux_web_app.test.location
  resource_group_name = azurerm_linux_web_app.test.resource_group_name
  service_plan_id     = azurerm_linux_web_app.test.service_plan_id
}
`, r.basic(data))
}

func (r LinuxWebAppResource) loadBalancing(data acceptance.TestData, loadBalancingMode string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    load_balancing_mode = "%s"
  }
}
`, r.baseTemplate(data), data.RandomInteger, loadBalancingMode)
}

func (r LinuxWebAppResource) logsEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  logs {
    application_logs {
      file_system_level = "Warning"
      azure_blob_storage {
        level             = "Information"
        sas_url           = "http://x.com/"
        retention_in_days = 2
      }
    }

    http_logs {
      azure_blob_storage {
        sas_url           = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
        retention_in_days = 3
      }
    }
  }
}

`, r.templateWithStorageAccount(data), data.RandomInteger)
}

func (r LinuxWebAppResource) withDetailedLogging(data acceptance.TestData, detailedErrorLogging bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  logs {
    detailed_error_messages = %t
  }
}
`, r.baseTemplate(data), data.RandomInteger, detailedErrorLogging)
}

func (r LinuxWebAppResource) withLoggingComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  logs {
    detailed_error_messages = true
    failed_request_tracing  = true

    application_logs {
      file_system_level = "Warning"

      azure_blob_storage {
        level             = "Information"
        sas_url           = "http://x.com/"
        retention_in_days = 7
      }
    }

    http_logs {
      file_system {
        retention_in_days = 4
        retention_in_mb   = 25
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppResource) withIPRestrictions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppResource) withIPRestrictionsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da", "6bde7211-57bc-4476-866a-c9676e22b9d7"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64", "9.9.9.8/32"]
        x_forwarded_host  = ["example.com", "anotherexample.com"]
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppResource) withAuthSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%s"

    additional_login_parameters = {
      test_key = "test_value"
    }

    allowed_external_redirect_urls = ["https://example.com"]

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger, data.Client().TenantID)
}

func (r LinuxWebAppResource) withAuthSettingsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%s"

    additional_login_parameters = {
      test_key = "test_value"
    }

    allowed_external_redirect_urls = ["https://example.com"]

    default_provider              = "AzureActiveDirectory"
    token_refresh_extension_hours = 24
    token_store_enabled           = true
    unauthenticated_client_action = "RedirectToLoginPage"

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }

    facebook {
      app_id     = "facebookappid"
      app_secret = "facebookappsecret"

      oauth_scopes = [
        "facebookscope",
      ]
    }

    google {
      client_id     = "googleclientid"
      client_secret = "googleclientsecret"

      oauth_scopes = [
        "googlescope",
      ]
    }

    microsoft {
      client_id     = "microsoftclientid"
      client_secret = "microsoftclientsecret"

      oauth_scopes = [
        "microsoftscope",
      ]
    }

    twitter {
      consumer_key    = "twitterconsumerkey"
      consumer_secret = "twitterconsumersecret"
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger, data.Client().TenantID)
}

func (r LinuxWebAppResource) dotNet(data acceptance.TestData, dotNetVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    application_stack {
      dotnet_framework_version = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, dotNetVersion)
}

func (r LinuxWebAppResource) php(data acceptance.TestData, phpVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    application_stack {
      php_version = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, phpVersion)
}

func (r LinuxWebAppResource) python(data acceptance.TestData, pythonVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    application_stack {
      python_version = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, pythonVersion)
}

func (r LinuxWebAppResource) node(data acceptance.TestData, nodeVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    application_stack {
      node_version = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, nodeVersion)
}

func (r LinuxWebAppResource) java(data acceptance.TestData, javaVersion, javaServer, javaServerVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    application_stack {
      java_version        = "%s"
      java_server         = "%s"
      java_server_version = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, javaVersion, javaServer, javaServerVersion)
}

func (r LinuxWebAppResource) javaPremiumV3Plan(data acceptance.TestData, javaVersion, javaServer, javaServerVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    application_stack {
      java_version        = "%s"
      java_server         = "%s"
      java_server_version = "%s"
    }
  }
}

`, r.premiumV3PlanTemplate(data), data.RandomInteger, javaVersion, javaServer, javaServerVersion)
}

func (r LinuxWebAppResource) docker(data acceptance.TestData, containerImage, containerTag string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  app_settings = {
    "DOCKER_REGISTRY_SERVER_URL"          = "https://mcr.microsoft.com"
    "DOCKER_REGISTRY_SERVER_USERNAME"     = ""
    "DOCKER_REGISTRY_SERVER_PASSWORD"     = ""
    "WEBSITES_ENABLE_APP_SERVICE_STORAGE" = "false"

  }

  site_config {
    application_stack {
      docker_image     = "%s"
      docker_image_tag = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, containerImage, containerTag)
}

func (r LinuxWebAppResource) autoHealRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    auto_heal = true

    auto_heal_setting {
      trigger {
        status_code {
          status_code_range = "500"
          interval          = "00:01:00"
          count             = 10
        }
      }

      action {
        action_type                    = "Recycle"
        minimum_process_execution_time = "00:05:00"
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppResource) autoHealRulesUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    auto_heal = true

    auto_heal_setting {
      trigger {
        status_code {
          status_code_range = "500"
          interval          = "00:01:00"
          count             = 10
        }
        status_code {
          status_code_range = "400-404"
          interval          = "00:10:00"
          count             = 10
        }
      }

      action {
        action_type                    = "LogEvent"
        minimum_process_execution_time = "00:10:00"
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppResource) autoHealRulesStatusCodeRange(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {
    auto_heal = true

    auto_heal_setting {
      trigger {
        status_code {
          status_code_range = "500-599"
          interval          = "00:01:00"
          count             = 10
        }
      }

      action {
        action_type                    = "Recycle"
        minimum_process_execution_time = "00:05:00"
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

// TODO - Test for new acr creds?

// Templates

func (LinuxWebAppResource) baseTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "B1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LinuxWebAppResource) premiumV3PlanTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "P1v3"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

// nolint: unused
func (r LinuxWebAppResource) templateWithStorageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_share" "test" {
  name                 = "test"
  storage_account_name = azurerm_storage_account.test.name
}

data "azurerm_storage_account_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  https_only        = true

  resource_types {
    service   = false
    container = false
    object    = true
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "2021-04-01"
  expiry = "2024-03-30"

  permissions {
    read    = false
    write   = true
    delete  = false
    list    = false
    add     = false
    create  = false
    update  = false
    process = false
  }
}
`, r.baseTemplate(data), data.RandomInteger, data.RandomString)
}
