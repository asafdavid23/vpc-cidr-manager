# VPC CIDR Manager
A Simple CLI tool to manage AWS VPC CIDR block reservations stored in DynamoDB.

## Features
- **Create Table**: Create a DynamoDB Table with IaC (Cloudformation).
- **Create IAM Role** Create an Assumable IAM Role for cross-account with Iac (Cloudformation).
- **Reserve CIDR**: Add a new CIDR block to the DynamoDB table.  
- **Release CIDR**: Remove a CIDR block from the table.  
- **Import CIDR**: Import live AWS VPC CIDRs into DynamoDB.  
- **Conflict Prevention**: Prevent CIDR overlap and maintain consistency across your infrastructure.
- **List CIDR**: List for existing CIDRs and print as Table/JSON.

## Installation
You can install vpc-cidr-manager by downloading the latest release from the [releases page](https://github.com/asafdavid23/vpc-cidr-manager/releases)

```bash
curl -LO https://github.com/asafdavid23/vpc-cidr-manager/releases/latest/download/vpc-cidr-manager
chmod +x vpc-cidr-manager
sudo mv vpc-cidr-manager /usr/local/bin/
```
For Windows, download the binary and add it to your system PATH.

Alternatively, you use brew for (MacOS / Linux)
```bash
brew tap asafdavid23/tap
brew update
brew install vpc-cidr-manager
```

## Usage
```bash
vpc-cidr-manager <command> [flags]
```

## Configuration
```
Available Commands:
  cloudformation Manage Supported Infrastructure as Code with CloudFormation
  completion     Generate the autocompletion script for the specified shell
  dynamodb       Manage DynamoDB Operations
  help           Help about any command

Flags:
      --config string      config file (default is $HOME/.vpc-cidr-manager.yaml)
  -h, --help               help for vpc-cidr-manager
      --log-level string   Set the log level (debug, info, warn, error, fatal) (default "info")
      --output string      Output type table/json/yaml (default "table")
      --version            Display the version of this CLI tool

Use "vpc-cidr-manager [command] --help" for more information about a command.
```

## License
This project is licensed under the MIT License.

