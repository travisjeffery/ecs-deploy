# ecs-deploy

Deploy ECS service to a Docker image.

## Install

[Download the binary](https://github.com/travisjeffery/ecs-deploy/releases/latest), or go get:

```
$ go get github.com/travisjeffery/ecs-deploy
```

## Example

```
$ ecs-deploy --service=app --image=travisjeffery/app --tag=1.0.0
default/app-stage 2015/12/05 02:10:38 [info] --> desired: 2, pending: 0, running: 0
default/app-stage 2015/12/05 02:10:43 [info] --> desired: 1, pending: 1, running: 0
default/app-stage 2015/12/05 02:10:43 [info] --> desired: 0, pending: 0, running: 2
default/app-stage 2015/12/05 02:10:48 [info] update service succes
```

## Usage

```
usage: ecs-deploy --service=SERVICE [<flags>]

Deploy ECS service.

Flags:
  --help                Show context-sensitive help (also try --help-long and --help-man).
  --service=SERVICE     Name of Service to update.
  --task=TASK-DEF       Name of Task Definition to update. Defaults to service.
  --image=IMAGE         Name of Docker image to run.
  --tag=TAG             Tag of Docker image to run.
  --cluster="default"   Name of ECS cluster.
  --region="us-east-1"  Name of AWS region.
  --count=-1            Desired count of instantiations to run. Defaults to existing running count.
  --nowait              Disable waiting for task definitions to start running.
  --version             Show application version.
```

You can also override the default region by setting the `AWS_DEFAULT_REGION` environmental variable.

## Author

[Travis Jeffery](http://twitter.com/travisjeffery)

## License

MIT
