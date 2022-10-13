<AWSACCOUNT>n

import (
	"flag"
	"fmt"
	"github.com/TwiN/go-color"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"strings"
	"text/template"
	"github.com/rollerd/commonconfig/filelib"
)

const version string = "1.1.0"
var home = getEnv("HOME") 

type MasterConfig  struct {
	BinaryVersion,
	AwsAccessKeyId,
	AwsSecretAccessKey,
	Username,
	ShortName,
	DevRole,
	StagingRole,
	DevEksRole,
	StagingEksRole,
	ProdEksRole string
}

func main() {
    reset := flag.Bool("r", false, "Reset config files with new values")
	bversion := flag.Bool("v", false, "Version info")
    flag.Parse()

	if (*reset) {
		fmt.Println(color.Cyan + "Resetting config files" + color.Reset)
	}

	if (*bversion) {
		fmt.Printf(color.Blue + "Common config version: %s\n" + color.Reset, version)
		os.Exit(0)
	}

	masterConfig := readMasterConfig(*reset)
	downloadFiles(masterConfig)
	configureShellrc()
	configureRdsauth(masterConfig, *reset)
	configureKubeConfig(*reset)
	configureAwsCredentials(masterConfig, *reset)
	configureAwsConfig(*reset)
	copyBinaries(*reset)
}

func downloadFiles(masterConfig *MasterConfig) {
	filelib.S3Download(masterConfig.AwsAccessKeyId, masterConfig.AwsSecretAccessKey, "icm-commonconfig", "rdsauth", "./bin/rdsauth")
	filelib.S3Download(masterConfig.AwsAccessKeyId, masterConfig.AwsSecretAccessKey, "icm-commonconfig", "awsenv", "./bin/awsenv")
	filelib.S3Download(masterConfig.AwsAccessKeyId, masterConfig.AwsSecretAccessKey, "icm-commonconfig", "kubectl", "./bin/kubectl")
	filelib.S3Download(masterConfig.AwsAccessKeyId, masterConfig.AwsSecretAccessKey, "icm-commonconfig", "kubectx", "./bin/kubectx")
	filelib.S3Download(masterConfig.AwsAccessKeyId, masterConfig.AwsSecretAccessKey, "icm-commonconfig", "kubens", "./bin/kubens")
	filelib.S3Download(masterConfig.AwsAccessKeyId, masterConfig.AwsSecretAccessKey, "icm-commonconfig", "commonconfig", "./bin/commonconfig")
}

func getEnv(varName string) string {
	varValue := os.Getenv(varName)
	return varValue
}

func copyBinaries(reset bool) {
	success := filelib.CopyFile("./bin/rdsauth", "/usr/local/bin/rdsauth")
	if (!success) {
		fmt.Printf(color.Bold + color.Purple + "Try to manually run: 'sudo cp ./bin/rdsauth /usr/local/bin/ && sudo chmod 755 /usr/local/bin/rdsauth'\n" + color.Reset)
	}else{
		_ = os.Chmod("/usr/local/bin/rdsauth", 0755)
	}
	success = filelib.CopyFile("./bin/awsenv", "/usr/local/bin/awsenv")
	if (!success) {
		fmt.Printf(color.Bold + color.Purple + "Try to manually run: 'sudo cp ./bin/awsenv /usr/local/bin/ && sudo chmod 755 /usr/local/bin/awsenv'\n" + color.Reset)
		_ = os.Chmod("/usr/local/bin/awsenv", 0755)
	}else{
		_ = os.Chmod("/usr/local/bin/awsenv", 0755)
	}
	success = filelib.CopyFile("./bin/kubectx", "/usr/local/bin/kubectx")
	if (!success) {
		fmt.Printf(color.Bold + color.Purple + "Try to manually run: 'sudo cp ./bin/kubectx /usr/local/bin/ && sudo chmod 755 /usr/local/bin/kubectx'\n" + color.Reset)
		_ = os.Chmod("/usr/local/bin/kubectx", 0755)
	}else{
		_ = os.Chmod("/usr/local/bin/kubectx", 0755)
	}
	success = filelib.CopyFile("./bin/kubens", "/usr/local/bin/kubens")
	if (!success) {
		fmt.Printf(color.Bold + color.Purple + "Try to manually run: 'sudo cp ./bin/kubens /usr/local/bin/ && sudo chmod 755 /usr/local/bin/kubens'\n" + color.Reset)
		_ = os.Chmod("/usr/local/bin/kubens", 0755)
	}else{
		_ = os.Chmod("/usr/local/bin/kubens", 0755)
	}
	success = filelib.CopyFile("./bin/kubectl", "/usr/local/bin/kubectl")
	if (!success) {
		fmt.Printf(color.Bold + color.Purple + "Try to manually run: 'sudo cp ./bin/kubectl /usr/local/bin/ && sudo chmod 755 /usr/local/bin/kubectl'\n" + color.Reset)
		_ = os.Chmod("/usr/local/bin/kubectl", 0755)
	}else{
		_ = os.Chmod("/usr/local/bin/kubectl", 0755)
	}
	success = filelib.CopyFile("./bin/commonconfig", "/usr/local/bin/commonconfig")
	if (!success) {
		fmt.Printf(color.Bold + color.Purple + "Try to manually run: 'sudo cp ./bin/commonconfig /usr/local/bin/ && sudo chmod 755 /usr/local/bin/commonconfig'\n" + color.Reset)
		_ = os.Chmod("/usr/local/bin/commonconfig", 0755)
	}else{
		_ = os.Chmod("/usr/local/bin/commonconfig", 0755)
	}
}

func configureAwsConfig(reset bool) {
	awsConfigFilename := fmt.Sprintf("%s/.aws/config", home)
	fileExists := filelib.CheckFileExists(awsConfigFilename)

	if (fileExists) {
		if (reset) {
			filelib.BackupFile(awsConfigFilename)
			filelib.CopyFile("configs/aws_config", awsConfigFilename)
		}else{
			log.Printf(color.Green + "Found ~/.aws/config file" + color.Reset)
		}
	}else{
		log.Printf(color.Green + "Creating file: %s" + color.Reset, awsConfigFilename)
		filelib.CopyFile("configs/aws_config", awsConfigFilename)
	}
}

func configureAwsCredentials(masterConfig *MasterConfig, reset bool) {
	awsCredentialsFilename := fmt.Sprintf("%s/.aws/credentials", home)
	fileExists := filelib.CheckFileExists(awsCredentialsFilename)

	if (fileExists) {
		if (reset) {
			filelib.BackupFile(awsCredentialsFilename)
			writeTemplate(masterConfig, awsCredentialsFilename, awsCredentialsTemplate)
		}else{
			log.Printf(color.Green + "Found ~/.aws/credentials file" + color.Reset)
		}
	}else{
		log.Printf(color.Green + "Creating file: %s" + color.Reset, awsCredentialsFilename)
		filelib.CreateDir(fmt.Sprintf("%s/.aws", home))
		writeTemplate(masterConfig, awsCredentialsFilename, awsCredentialsTemplate)
	}
}

func configureKubeConfig(reset bool) {
	kubeConfigFilename := fmt.Sprintf("%s/.kube/config", home)
	fileExists := filelib.CheckFileExists(kubeConfigFilename)

	if (fileExists) {
		if (reset) {
			filelib.BackupFile(kubeConfigFilename)
			filelib.CopyFile("configs/kube_config", kubeConfigFilename)
		}else{
			log.Print(color.Green + "Found kube config file" + color.Reset) 
		}
	}else{
		log.Printf(color.Green + "Creating file: %s" + color.Reset, kubeConfigFilename)
		filelib.CreateDir(fmt.Sprintf("%s/.kube", home))
		filelib.CopyFile("configs/kube_config", kubeConfigFilename)
	}
}

func configureShellrc() {
	shell := getEnv("SHELL")
	log.Printf(color.Green + "Found shell: %s" + color.Reset, shell)

	var shellConfig string
	if (strings.Contains(shell, "bash")) {
		shellConfig = fmt.Sprintf("%s/.bashrc", home)
	}
	if (strings.Contains(shell, "zsh")) {
		shellConfig = fmt.Sprintf("%s/.zshrc", home)
	}


	var foundProfile bool
	foundProfile = false
	aws_profile := "export AWS_PROFILE"
	aws_key		:= "export AWS_ACCESS_KEY_ID"
	aws_secret	:= "export AWS_SECRET_ACCESS_KEY"

	fileExists := filelib.CheckFileExists(shellConfig)
	if (fileExists) {
		filelib.BackupFile(shellConfig)
	}

	content := filelib.ReadFile(shellConfig)

	for i := 0; i < len(content); i++ {
		if (strings.HasPrefix(content[i], aws_profile)) {
			log.Printf(color.Green + "Found AWS_PROFILE in %s" + color.Reset, shellConfig)
			foundProfile = true
		}
		if (strings.HasPrefix(content[i], aws_key)) {
			content[i] = ""
		}
		if (strings.HasPrefix(content[i], aws_secret)) {
			content[i] = ""
		}
	}

	if (foundProfile == false) {
		content = append(content, "export AWS_PROFILE=default")
	}

	filelib.WriteFile(shellConfig, content)
}

func configureRdsauth(masterConfig *MasterConfig, reset bool) {
	rdsauthFilename := fmt.Sprintf("%s/.rdsauth.ini", home)

	fileExists := filelib.CheckFileExists(rdsauthFilename)
	if (fileExists) {
		if (reset) {
			filelib.BackupFile(rdsauthFilename)
			writeTemplate(masterConfig, rdsauthFilename, rdsauthTemplate)
		}else{
			log.Print(color.Green + "Found rdsauth.ini file" + color.Reset) 
		}
	}else{
		log.Printf(color.Green + "Creating file: %s" + color.Reset, rdsauthFilename)
		writeTemplate(masterConfig, rdsauthFilename, rdsauthTemplate)
	}
}

func createMasterConfig() *MasterConfig {
	var masterConfig MasterConfig

	fmt.Println(color.Blue + "Enter your AWS_ACCESS_KEY_ID: " + color.Reset)
	var awsAccessKeyId string 
	fmt.Scanln(&awsAccessKeyId)

	fmt.Println(color.Blue + "Enter your AWS_SECRET_ACCESS_KEY: " + color.Reset)
	var awsSecretAccessKey string 
	fmt.Scanln(&awsSecretAccessKey)

	fmt.Println(color.Blue + "Enter your ICM email: " + color.Reset)
	var username string 
	fmt.Scanln(&username)

	fmt.Println(color.Blue + "Enter your user short name (first part of email before '@'): " + color.Reset)
	var shortName string 
	fmt.Scanln(&shortName)

	fmt.Println(color.Blue + "Enter AWS dev role name: " + color.Reset)
	var devRole string 
	fmt.Scanln(&devRole)

	fmt.Println(color.Blue + "Enter AWS staging role name: " + color.Reset)
	var stagingRole string 
	fmt.Scanln(&stagingRole)

	fmt.Println(color.Blue + "Enter EKS dev role name: " + color.Reset)
	var devEksRole string 
	fmt.Scanln(&devEksRole)

	fmt.Println(color.Blue + "Enter EKS staging role name: " + color.Reset)
	var stagingEksRole string 
	fmt.Scanln(&stagingEksRole)

	fmt.Println(color.Blue + "Enter EKS prod role name: " + color.Reset)
	var prodEksRole string 
	fmt.Scanln(&prodEksRole)

	masterConfig.BinaryVersion		= version
	masterConfig.AwsAccessKeyId		= awsAccessKeyId
	masterConfig.AwsSecretAccessKey = awsSecretAccessKey
	masterConfig.Username			= username
	masterConfig.ShortName			= shortName
	masterConfig.DevRole			= devRole
	masterConfig.StagingRole		= stagingRole
	masterConfig.DevEksRole			= devEksRole
	masterConfig.StagingEksRole		= stagingEksRole
	masterConfig.ProdEksRole		= prodEksRole

	writeTemplate(&masterConfig, fmt.Sprintf("%s/.icm_commonconfig", home), icmCommonConfigTemplate)

	return &masterConfig
}

func readMasterConfig(reset bool) *MasterConfig {
	masterConfigFilename := fmt.Sprintf("%s/.icm_commonconfig", home)
	fileExists := filelib.CheckFileExists(masterConfigFilename)

	if (fileExists) {
		log.Print(color.Green + "Found existing icm-commonconfig master config!" + color.Reset)
		if (reset) {
			filelib.BackupFile(masterConfigFilename)
			masterConfig := createMasterConfig()
			return masterConfig
		}
		cfg, err := ini.Load(masterConfigFilename)
		if err != nil {
			log.Fatal(err)
		}
		var masterConfig MasterConfig

		masterConfig.BinaryVersion		= cfg.Section("main").Key("BINARY_VERSION").String() 
		masterConfig.AwsAccessKeyId		= cfg.Section("main").Key("AWS_ACCESS_KEY_ID").String()
		masterConfig.AwsSecretAccessKey = cfg.Section("main").Key("AWS_SECRET_ACCESS_KEY").String() 
		masterConfig.Username			= cfg.Section("main").Key("USERNAME").String() 
		masterConfig.ShortName			= cfg.Section("main").Key("SHORTNAME").String() 
		masterConfig.DevRole			= cfg.Section("main").Key("DEV_ROLE").String() 
		masterConfig.StagingRole		= cfg.Section("main").Key("STAGING_ROLE").String() 
		masterConfig.DevEksRole			= cfg.Section("main").Key("DEV_EKS_ROLE").String() 
		masterConfig.StagingEksRole		= cfg.Section("main").Key("STAGING_EKS_ROLE").String() 
		masterConfig.ProdEksRole		= cfg.Section("main").Key("PROD_EKS_ROLE").String() 

		return &masterConfig

	}else{
		masterConfig := createMasterConfig()
		return masterConfig
	}
}

func writeTemplate(configData *MasterConfig, filename string, templateData string) {
	tmpl, err := template.New("masterconfig yaml").Parse(templateData)
	if err != nil {
		log.Fatalf(color.Red + "Could not create %s template: %s" + color.Reset, filename, err)
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf(color.Red + "Could not create file: %s: %s" + color.Reset, filename, err)
	}

	err = tmpl.Execute(file, configData)
	if err != nil {
		log.Fatalf(color.Red + "Could not write template data to %s: %s" + color.Reset, filename, err)
	}

	err = os.Chmod(filename, 0600)
	if err != nil {
		log.Printf(color.Yellow + "Could not change filemode for file: %s" + color.Reset, filename)
	}
}

var rdsauthTemplate= `[dev]
db     = <DB>:3306
region = us-west-2
user   = {{ .ShortName }}

[staging]
db     = <DB>:3306
region = us-west-2
user   = {{ .ShortName }}

[prod]
db     = <DB>:3306
region = us-west-2
user   = {{ .ShortName }}
`

var icmCommonConfigTemplate= `[main]
BINARY_VERSION = {{ .BinaryVersion }}
AWS_ACCESS_KEY_ID = {{ .AwsAccessKeyId }}
AWS_SECRET_ACCESS_KEY = {{ .AwsSecretAccessKey }}
USERNAME = {{ .Username }}
SHORTNAME = {{ .ShortName }}
DEV_ROLE = {{ .DevRole }}
STAGING_ROLE = {{ .StagingRole }}
DEV_EKS_ROLE = {{ .DevEksRole }}
STAGING_EKS_ROLE = {{ .StagingEksRole }}
PROD_EKS_ROLE = {{ .ProdEksRole }}
`

var awsCredentialsTemplate= `[prod]
aws_access_key_id = {{ .AwsAccessKeyId }}
aws_secret_access_key = {{ .AwsSecretAccessKey }}
alias = prod

[staging]
role_arn = arn:aws:iam::<AWSACCOUNT>:role/{{ .StagingRole }}
role_session_name = {{ .Username }}
source_profile = prod
alias = staging

[dev]
role_arn = arn:aws:iam::<AWSACCOUNT>:role/{{ .DevRole }}
role_session_name = {{ .Username }}
source_profile = prod
alias = dev

[prod-eks]
role_arn = arn:aws:iam::<AWSACCOUNT>:role/{{ .ProdEksRole }}
role_session_name = {{ .Username }}
source_profile = prod
alias = prod-eks

[staging-eks]
role_arn = arn:aws:iam::<AWSACCOUNT>:role/{{ .StagingEksRole }}
role_session_name = {{ .Username }}
source_profile = prod
alias = staging-eks

[dev-eks]
role_arn = arn:aws:iam::<AWSACCOUNT>:role/{{ .DevEksRole }}
role_session_name = {{ .Username }}
source_profile = prod
alias = dev-eks

[default]
aws_access_key_id = {{ .AwsAccessKeyId }}
aws_secret_access_key = {{ .AwsSecretAccessKey }}
alias = prod
`
