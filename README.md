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

## Author

[Travis Jeffery](http://twitter.com/travisjeffery)

## License

MIT
