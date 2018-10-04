# tyr

A command line security audit tool for Amazon Web Services

## About
Tyr is a command line tool that scans for vulnerabilities in your AWS Account. In easy way you will be able to
identify unsecure parts of your infrastructure and prepare your AWS account for security audit.

## Installation
Currently Tyr does not support any package managers, but the work is in progress. 
### Building from sources
First of all you need to download Tyr to your GO workspace:

```bash
$GOPATH $ go get github.com/Appliscale/tyr
$GOPATH $ cd tyr
```

Then build and install configuration for the application inside perun directory by executing:

```bash
tyr $ make all
```

## Usage
### Initialising Session
If you're using MFA you need to tell Tyr to authenticate you before trying to connect by using flag `--mfa`. 
Example:
```
$ tyr --service s3 --mfa --mfa-duration 3600
```

### EC2 Scan
#### How to use
To perform audit on all EC2 instances, type:
```
$ tyr --service ec2
```
You can narrow the audit to a region, by using the flag `-r` or `--region`. Tyr also supports AWS profiles -
to specify profile use the flag `-p` or `--profile`.

#### Example output

```bash
+---------------------+--------------------------------+-------------+----------+
|         EC2         |            VOLUMES             |  SECURITY   |          |
|                     |                                |             | EC2 TAGS |
|                     |     (NONE) - NOT ENCRYPTED     |   GROUPS    |          |
|                     |                                |             |          |
|                     |    (DKMS) - ENCRYPTED WITH     |             |          |
|                     |         DEFAULT KMSKEY         |             |          |
+---------------------+--------------------------------+-------------+----------+
| i-0fa455c90ace32283 | vol-0a8143f0b2e78424d[DKMS]    | sg-aaaaaaa  | App:some |
|                     | vol-0c4bacc1704c98f56[NONE]    |             | Key:Val  |
|                     |                                |             |          |
|                     |                                |             |          |
+---------------------+--------------------------------+-------------+----------+
```

#### How to read it

 1. First column `EC2` contains instance ID.
 2. Second column `Volumes` contains IDs of attached volumes(virtual disks) to given EC2. Suffixes meaning:
    * `[NONE]` - Volume not encrypted.
    * `[DKMS]` - Volume encrypted using AWS Default KMS Key. More about KMS you can find [here](https://aws.amazon.com/kms/faqs/)
 3. Third column `Security Groups` contains IDs of security groups that have too open permissions. e.g. CIDR block is equal to `0.0.0.0/0`(open to the whole world).
 4. Fourth column `EC2 TAGS` contains tags of a given EC2 instance to help you identify purpose of this instance.

#### Docs
You can find more information about encryption in the following documentation:
  1. https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSEncryption.html

### S3 Scan
#### How to use
To perform audit on all S3 buckets, type:
```
$ tyr --service s3
```
Tyr supports AWS profiles - to specify profile use the flag `-p` or `--profile`.

#### Example output

```bash
+------------------------------+-------------+-----------------+
|          BUCKET NAME         | DEFAULT SSE | LOGGING ENABLED |
+------------------------------+-------------+-----------------+
| bucket1                      | NONE        | true            |
+------------------------------+-------------+-----------------+
| bucket2                      | DKMS        | false           |
+------------------------------+-------------+-----------------+
| bucket3                      | AES256      | false           |
+------------------------------+-------------+-----------------+
```

#### How to read it

 1. First column `BUCKET NAME` contains names of the s3 buckets.
 2. Second column `DEFAULT SSE` gives you information on which default type of server side encryption was used in your S3 bucket:
   * `NONE` - Default SSE not enabled.
   * `DKMS` - Default SSE enabled, AWS KMS Key used to encryp data.
   * `AES256` - Default SSE enabled, [AES256](https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingServerSideEncryption.html).
 3. Third column `LOGGING ENABLED` contains information if logging was enabled in given S3 bucket.  

#### Docs
You can find more about securing your S3's in the following documentations:
 1. https://docs.aws.amazon.com/AmazonS3/latest/dev/serv-side-encryption.html
 2. https://docs.aws.amazon.com/AmazonS3/latest/dev/ServerLogs.html
 3. https://docs.aws.amazon.com/AmazonS3/latest/user-guide/server-access-logging.html
