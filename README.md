# VPC CIDR Manager
A Simple CLI tool to manage AWS VPC CIDR block reservations stored in DynamoDB.

## Features
- **Create Table**: Create a DynamoDB Table
- **Reserve CIDR**: Add a new CIDR block to the DynamoDB table.  
- **Delete CIDR**: Remove a CIDR block from the table.  
- **Import CIDR**: Import live AWS VPC CIDRs into DynamoDB.  

## Installation
You can install eolctl by downloading the latest release from the [releases page](https://github.com/asafdavid23/vpc-cidr-manager/releases)

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
The tool uses environment variables for AWS credentials and DynamoDB table name:

- AWS_REGION
- DDB_TABLE_NAME

## License
This project is licensed under the MIT License.

