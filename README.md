# Cloud Security Audit [![Release](https://img.shields.io/github/release/Appliscale/tyr.svg?style=flat-square)](https://github.com/Appliscale/tyr/releases/latest) [![CircleCI](https://circleci.com/gh/Appliscale/tyr.svg?style=svg)](https://circleci.com/gh/Appliscale/tyr) [![License](https://img.shields.io/badge/License-Apache%202.0-orange.svg)](https://github.com/Appliscale/tyr/blob/master/LICENSE.md)  [![Go_Report_Card](https://goreportcard.com/badge/github.com/Appliscale/tyr?style=flat-square&fuckgithubcache=1)](https://goreportcard.com/report/github.com/Appliscale/tyr) [![GoDoc](https://godoc.org/github.com/Appliscale/tyr?status.svg)](https://godoc.org/github.com/Appliscale/tyr)


A command line security audit tool for Amazon Web Services

## About
Cloud Security Audit is a command line tool that scans for vulnerabilities in your AWS Account. In easy way you will be able to
identify unsecure parts of your infrastructure and prepare your AWS account for security audit.

## Installation
Currently Cloud Security Audit does not support any package managers, but the work is in progress. 
### Building from sources
First of all you need to download Cloud Security Audit to your GO workspace:

```bash
$GOPATH $ go get github.com/Appliscale/cloud-security-audit
$GOPATH $ cd cloud-security-audit
```

Then build and install configuration for the application inside perun directory by executing:

```bash
cloud-security-audit $ make all
```

## Usage
### Initialising Session
If you're using MFA you need to tell Cloud Security Audit to authenticate you before trying to connect by using flag `--mfa`. 
Example:
```
$ cloud-security-audit --service s3 --mfa --mfa-duration 3600
```

### EC2 Scan
#### How to use
To perform audit on all EC2 instances, type:
```
$ cloud-security-audit --service ec2
```
You can narrow the audit to a region, by using the flag `-r` or `--region`. Cloud Security Audit also supports AWS profiles -
to specify profile use the flag `-p` or `--profile`.

#### Example output

```bash

+---------------+---------------------+--------------------------------+-----------------------------------+----------+
| AVAILABILITY  |         EC2         |            VOLUMES             |             SECURITY              |          |
|               |                     |                                |                                   | EC2 TAGS |
|     ZONE      |                     |     (NONE) - NOT ENCRYPTED     |              GROUPS               |          |
|               |                     |                                |                                   |          |
|               |                     |    (DKMS) - ENCRYPTED WITH     |    (INCOMING CIDR = 0.0.0.0/0)    |          |
|               |                     |         DEFAULT KMSKEY         |                                   |          |
|               |                     |                                |       ID : PROTOCOL : PORT        |          |
+---------------+---------------------+--------------------------------+-----------------------------------+----------+
| eu-central-1a | i-0fa345j6756nb3v23 | vol-0a81288qjd188424d[DKMS]    | sg-aaaaaaaa : tcp : 22            | App:some |
|               |                     | vol-0c2834re8dfsd8sdf[NONE]    | sg-aaaaaaaa : tcp : 22            | Key:Val  |
+---------------+---------------------+--------------------------------+-----------------------------------+----------+
```

#### How to read it

 1. First column `AVAILABILITY ZONE` contains information where the instance is placed
 2. Second column `EC2` contains instance ID.
 3. Third column `Volumes` contains IDs of attached volumes(virtual disks) to given EC2. Suffixes meaning:
    * `[NONE]` - Volume not encrypted.
    * `[DKMS]` - Volume encrypted using AWS Default KMS Key. More about KMS you can find [here](https://aws.amazon.com/kms/faqs/)
 4. Fourth column `Security Groups` contains IDs of security groups that have too open permissions. e.g. CIDR block is equal to `0.0.0.0/0`(open to the whole world).
 5. Fifth column `EC2 TAGS` contains tags of a given EC2 instance to help you identify purpose of this instance.

#### Docs
You can find more information about encryption in the following documentation:
  1. https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSEncryption.html

### S3 Scan
#### How to use
To perform audit on all S3 buckets, type:
```
$ cloud-security-audit --service s3
```
Cloud Security Audit supports AWS profiles - to specify profile use the flag `-p` or `--profile`.

#### Example output

```bash
+------------------------------+---------+---------+-------------+------------+
|          BUCKET NAME         | DEFAULT | LOGGING |     ACL     |  POLICY    |
|                              |         |         |             |            |
|                              |   SSE   | ENABLED |  IS PUBLIC  | IS PUBLIC  |
|                              |         |         |             |            |
|                              |         |         |  R - READ   |  R - READ  |
|                              |         |         |             |            |
|                              |         |         |  W - WRITE  | W - WRITE  |
|                              |         |         |             |            |
|                              |         |         | D - DELETE  | D - DELETE |
+------------------------------+---------+---------+-------------+------------+
| bucket1                      | NONE    | true    | false       | false      |
+------------------------------+---------+---------+-------------+------------+
| bucket2                      | DKMS    | false   | false       | true [R]   |
+------------------------------+---------+---------+-------------+------------+
| bucket3                      | AES256  | false   | true [RWD]  | false      |
+--------------------------- --+---------+---------+-------------+------------+
```

#### How to read it

 1. First column `BUCKET NAME` contains names of the s3 buckets.
 2. Second column `DEFAULT SSE` gives you information on which default type of server side encryption was used in your S3 bucket:
   * `NONE` - Default SSE not enabled.
   * `DKMS` - Default SSE enabled, AWS KMS Key used to encrypt data.
   * `AES256` - Default SSE enabled, [AES256](https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingServerSideEncryption.html).
 3. Third column `LOGGING ENABLED` contains information if Server access logging was enabled for a given S3 bucket. This provides detailed records for the requests that are made to an S3 bucket. More information about Server Access Logging can be found [here](https://docs.aws.amazon.com/AmazonS3/latest/user-guide/server-access-logging.html)
 4. Fourth column `ACL IS PUBLIC` provides information if ACL (Access Control List) contains permissions, that make the bucket public (allow read/writes for anyone). More information about ACLs [here](https://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html)
 5. Fifth column `POLICY IS PUBLIC` contains information if bucket's policy allows any action (read/write) for an anonymous user. More about bucket policies [here](https://docs.aws.amazon.com/AmazonS3/latest/dev/using-iam-policies.html)
R, W and D letters describe what type of action is available for everyone.
#### Docs
You can find more about securing your S3's in the following documentations:
 1. https://docs.aws.amazon.com/AmazonS3/latest/dev/serv-side-encryption.html
 2. https://docs.aws.amazon.com/AmazonS3/latest/dev/ServerLogs.html
 3. https://docs.aws.amazon.com/AmazonS3/latest/user-guide/server-access-logging.html
 
## License

[Apache License 2.0](LICENSE)

## Maintainers

- [Michał Połcik](https://github.com/mwpolcik)
- [Maksymilian Wojczuk](https://github.com/maxiwoj)
- [Piotr Figwer](https://github.com/pfigwer)
- [Sylwia Gargula](https://github.com/SylwiaGargula)
- [Mateusz Piwowarczyk](https://github.com/piwowarc)
