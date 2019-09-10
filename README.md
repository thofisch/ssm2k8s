## Installation

### Bash

The following will download the latest version:

```bash
curl -sSL https://raw.githubusercontent.com/thofisch/ssm2k8s/master/install.sh | sh -s -
```

## Usage

### Creating secrets:


```bash
$ ./secrets put foo/kafka/prod BOOTSTRAP_SERVERS=xxx GROUP_ID=yyy
```

Will create two parameters:

* `/foo/kafka/prod/BOOTSTRAP_SERVERS`
* `/foo/kafka/prod/GROUP_ID`

Which will be synchronized to the secrete `foo-kafka-prod` in Kubernetes.

### Updating secrets:

```bash
$ ./secrets put foo/kafka/prod BOOTSTRAP_SERVERS=xxx --overwrite
```

Will overwrite the `/foo/kafka/prod/BOOTSTRAP_SERVERS` parameter

### Deleting secrets:

```bash
# delete a single secret
$ ./secrets delete foo/kafka/prod BOOTSTRAP_SERVERS

# delete the entire secret
$ ./secrets delete foo/kafka/prod --force
```

### Listing secrets:

```bash
# list all secrets
$ ./secrets list
Getting AWS SSM Parameters
Found 11 secrets in "AWS SSM Parameter Store"

PATH                                           SECRET                                        KEYS  HASH     LAST MODIFIED
/a                                             a                                             3     e018c03  2019-09-09T13:24:21Z
/a/db-migrations                               a-db-migrations                               1     0759860  2019-09-05T11:24:55Z
/default-kafka/prod                            default-kafka-prod                            3     21f58dc  2019-09-04T20:56:58Z
/position-simulator-dbmigrations/prod          position-simulator-dbmigrations-prod          5     9267d66  2019-09-04T20:57:15Z
/position-simulator/prod                       position-simulator-prod                       1     b9a6be7  2019-09-04T20:57:26Z
/replay-service-dbmigrations/prod              replay-service-dbmigrations-prod              5     53dad8b  2019-09-04T20:58:12Z
/replay-service/prod                           replay-service-prod                           4     74297be  2019-09-05T09:42:52Z
/terminal-operation-service-dbmigrations/prod  terminal-operation-service-dbmigrations-prod  5     7459b36  2019-09-03T07:34:51Z
/terminal-operation-service/prod               terminal-operation-service-prod               1     331995f  2019-09-03T07:34:53Z
/tracking-service-dbmigrations/prod            tracking-service-dbmigrations-prod            5     7a80eb4  2019-09-04T21:01:18Z
/tracking-service/prod                         tracking-service-prod                         1     aa372bd  2019-09-04T21:02:21Z

# filter secrets by path
$ ./secrets list a
Getting AWS SSM Parameters
Found 2 secrets in "AWS SSM Parameter Store"

PATH              SECRET           KEYS  HASH     LAST MODIFIED
/a                a                3     e018c03  2019-09-09T13:24:21Z
/a/db-migrations  a-db-migrations  1     0759860  2019-09-05T11:24:55Z

# verbose list
$ ./secrets list a -v
Getting AWS SSM Parameters
Found 2 secrets in "AWS SSM Parameter Store"

PATH                     SECRET           VERSION  LAST MODIFIED         VALUE
/a/b                     a                1        2019-09-09T13:24:21Z  ***
/a/d                     a                1        2019-09-09T13:24:21Z  ***
/a/f                     a                1        2019-09-09T13:24:21Z  ***
/a/db-migrations/pghost  a-db-migrations  1        2019-09-05T11:24:55Z  ***

# verbose list with values decoded
$ ./secrets list a -vd
Getting AWS SSM Parameters
Found 2 secrets in "AWS SSM Parameter Store"

PATH                     SECRET           VERSION  LAST MODIFIED         VALUE
/a/b                     a                1        2019-09-09T13:24:21Z  c
/a/d                     a                1        2019-09-09T13:24:21Z  e
/a/f                     a                1        2019-09-09T13:24:21Z  g
/a/db-migrations/pghost  a-db-migrations  1        2019-09-05T11:24:55Z  sdfsdf
```
