### common_configs

This repository holds the commonconfig binary for creating and updating aws, kube, rdsauth, and other common config files

#### Configuration files:

- AWS credentials
- AWS config
- Kube config
- rdsauth.ini config

##### Requirements:

You will need the following before starting:

- AWS account

##### Usage:

Download the commonconfig binary from the 'Releases' section of the repository.

Clone this common-config repository and move the downloaded commonconfig binary to the root of the repo directory

From a Mac, simply make the commonconfig binary executable (`chmod 755 commonconfig`) and then run it:

```
./commonconfig
```

Have the following values available to complete template rendering:

AWS_ACCESS_KEY_ID

AWS_SECRET_ACCESS_KEY

USERNAME/EMAIL (your AWS username)

SHORTNAME (this is the name portion of your AWS account, before the '@'. May be different than your email)

DEV_ROLE name

STAGING_ROLE name

PROD_EKS_ROLE name

STAGING_EKS_ROLE name

DEV_EKS_ROLE name

*NOTE:* If you would like to re-enter all values and start from scratch (or if this is the first time running the script), add '-r' as a commandline arg. This will DELETE existing ~/.aws and ~/.kube directories, but will back them up with a 'date' extension.

```
$ ./commonconfig -r
```
