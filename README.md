# aws-appconfig-helper
A helper cli for AWS appconfig



### Available commands
```
appconfig get --app-id="" --env-id="" --watch="false|true" --files="<src_file:dest_path,src_file:dest_path>|all"
appconfig nuke app --app-id=""
appconfig nuke env --app-id="" --env-id=""
appconfig deploy --app-id="" --env-id="" --file="<file_name>|all" --version="<specific>|latest"
appconfig list --app-id="" --entity=""


###################
./aws-appconfig-helper
AWS appconfig helper

Usage:
  aws-appconfig-helper [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  deploy      deploy files
  get         Fetch files
  help        Help about any command
  list        list given type
  nuke        Wholesale deletion

Flags:
  -h, --help   help for aws-appconfig-helper

Use "aws-appconfig-helper [command] --help" for more information about a command.
```

### Examples

Fetch files test and test1 to respective paths at /tmp

```
aws-appconfig-helper get --app-id=5xmcnu7 --env-id=test --files=test:/tmp/test,test1:/tmp/test1
```

Fetch all config files from given environment to current folder
```
aws-appconfig-helper get --app-id=5xmcnu7 --env-id=test
```

get list of deployment strategies
```
aws-appconfig-helper list --app-id=5xmcnu7 --entity=deployment-strategy
{
  Items: [
    {
      DeploymentDurationInMinutes: 0,
      Description: "Quick",
      FinalBakeTimeInMinutes: 10,
      GrowthFactor: 100,
      GrowthType: "LINEAR",
      Id: "AppConfig.AllAtOnce",
      Name: "AppConfig.AllAtOnce",
      ReplicateTo: "NONE"
    },
    {
      DeploymentDurationInMinutes: 1,
      Description: "Test/Demo",
      FinalBakeTimeInMinutes: 1,
      GrowthFactor: 50,
      GrowthType: "LINEAR",
      Id: "AppConfig.Linear50PercentEvery30Seconds",
      Name: "AppConfig.Linear50PercentEvery30Seconds",
      ReplicateTo: "NONE"
    },
    {
      DeploymentDurationInMinutes: 20,
      Description: "AWS Recommended",
      FinalBakeTimeInMinutes: 10,
      GrowthFactor: 10,
      GrowthType: "EXPONENTIAL",
      Id: "AppConfig.Canary10Percent20Minutes",
      Name: "AppConfig.Canary10Percent20Minutes",
      ReplicateTo: "NONE"
    }
  ]
}
```

nuke entire appconfig application with id - 5xmcnu7. deletes all the environments and config profiles along with versions in the app.
```
aws-appconfig-helper nuke app --app-id 5xmcnu7
```

nuke environment in the appconfig application
```
aws-appconfig-helper nuke env --app-id 5xmcnu7 --env-id=abcdef
```

Deploy the latest versions of all config profiles to given environment
```
aws-appconfig-helper deploy --app-id=5xmcnu7 --env-id=abcdef --strategy="AppConfig.AllAtOnce" --profile-id=all
```

Deploy given version of config profile to given environment
```
aws-appconfig-helper deploy --app-id=5xmcnu7 --env-id=abcdef --strategy="AppConfig.AllAtOnce" --profile-id=hsnreou --version=2
```
