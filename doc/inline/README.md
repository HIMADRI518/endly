## Inline Workflow

- [Introduction](#introduction)
- [Inline Workflow format](#format)
- [Workflow invocation](#workflow)
- [Action invocation](#action)
- [Data model vs request namespace](#model_vs_request)
- [State modification](#state)



<a name="introduction"></a>
### Introduction

Endly uses [Inline Workflow](../../model/inline_workflow.go) to define simple sequential tasks with yaml files.
For instance the following workflow runs SSH command (service: exec, action: run).


```bash
endly -r=run
```

@run.yaml
```yaml
pipeline:
  action: exec:run
  target:
    url:  ssh://127.0.0.1/
    credentials: ${env.HOME}/.secret/localhost.json
  commands:
    - mkdir -p /tmp/app/build 
    - chown ${os.user} /tmp/app/build 
```



<a name="format"></a>
### Inline Workflow format

The general inline workflow syntax: 

@xxx.yaml
```yaml
params:
  k1:v1
init: var/@init
defaults:
  d1:v1

pipeline:
  task1:
     action: serviceID:action
     requestField1: val1
     requestFieldN: valN
           
  taskN:
    subTaskA:
      workflow: workflowSelector
      tasks: task selector
      paramsKey1: val1
      paramKeyN: valN
      
    subTaskZ:
       action: serviceID:action
       requestField1: val1
       requestFieldN: valN
      

post: 
  - age = $response.user.age
```

- _params_ node defines command line arguments
- _init_ node defines initial workflow variables/state
- _default_ node defines attributes that will be merge with every actionable node.
- _pipeline_ node defines set of tasks with its actions, which are be executed sequentially unless endly -t: task switch is used
- _post_ node defines post execution current workflow state data extraction to wrokflow run response

```yaml
pipeline:
  service:
    mysql:
      workflow: service/mysql
      tasks: start
    aerospike:
      action: workflow:run
      name: service/aerospike
      tasks: start
  frontend:
    deploy:
      workflow: app/deploy
      sdk: node:8.1
      app: demp-ui
  backend:    
    deploy:
      workflow: app/deploy
      sdk: go:1.9
      app: demo
    
  test:    
```

A _task_ can either be a groping or actual actions node. The latter invokes selected endly service operation, with defined request attributes.


To see the [*model.Workflow model](../../model/workflow.go) tree converted from a inline workflow run the following

```bash
endly -r=PIPELINE_FILE.yaml -p  -f=yaml|json
endly -w=WORKFLOW_FILE.csv -p  -f=yaml|json
```


<a name="action"></a>
### Action invocation

The generic service action invocation syntax:

```yaml
pipeline:
  task1:
    action: [SERVICE.]ACTION
    param1: val1
    paramX: valX
```

If SERVICE is omitted, 'workflow' service is used.


Run the following to check available workflow actions:

```bash
endly -s=workflow
```


Run the following to check particular workflow actions:

```bash
endly -s=workflow -a=run
```


Run the following to check available services

```bash
endly -s='*'
```

for example the following workflow task1:

```bash
endly -r=test
```

@test.yaml    
 ```yaml
pipeline:
  task1:
    action: workflow.print
    message: hello world
```

or 

@test.yaml    
 ```yaml
pipeline:
  task1:
    action: workflow.print
    request: '@print_req.yaml'
```

@print_req.yaml
```yaml
message: hello world
```


To see pipeline converted workflow  [*model.Workflow](../../model/workflow.go) run the following

```bash
endly -r=test -p -f=yaml
endly -r=test -p
```



<a name="workflow"></a>
### Sub workflow invocation

The generic external workflow invocation syntax:

```yaml
pipeline:
  task1:
     workflow: WORKFLOW_NAME[:TASKS_TO_RUN]
     param1: val1
     paramX: valX
```

for example the following workflow task1: invokes [assert workflow](../../shared/workflow/assert/assert.csv) with task:'assert' and the following  [@run](#assert_run) request:

```bash
endly -r=test
```

@test.yaml    
 ```yaml
 pipeline:
   task1:
      workflow: assert:assert
      expected: 10
      actual: 1
 ```
 
<a name="assert_run"></a>


@run.yaml
```yaml
URL: assert.csv
tasks: assert
params:
  actual: 1
  expected: 10
```


<a name="model_vs_request"></a>
### Data model vs request namespace

In both case workflow or action node share namespace with workflow [task](../../model/task.go) or [action](../../model/action.go)
To explicitly flag key as part of model attribute use ':' prefix   


@explicit.yaml
 ```yaml
pipeline:
  task1:
    action: workflow.print
    :init:
      - va1 = $params.p1
    message: hello world
```

```bash
endly -p -r=explicit
```


@implicit.yaml
 ```yaml
pipeline:
  task1:
    action: workflow.print
    init:
      - va1 = params.p1
    message: hello world
```


```bash
endly -p -r=implicit
```


<a name="state"></a>
### State modification

#### Default values

In case a tasks share data defaults value can be used to apply the same values accross tasks

For example to avoid message attribute duplication in task1 and task2  

@test.yaml
 ```yaml
pipeline:
  task1:
    action: print
    message: hello world
    color: red
  task2:
    action: print
    message: hello world
    color: blue
```

you can use the following:

@test.yaml
 ```yaml
defaults:
  message: hello world
pipeline:
  task1:
    action: print
    color: red
  task2:
    action: print
    color: blue
```

#### State modification


State initialization can be applied on top(workflow) or task/action node level. 
Parameters can be passed from command line.


@run.yaml
 ```yaml
init: 
  var1: $params.msg
  var2: 
    k1: 1
    k2: 2
pipeline:
  task1:
    init:
      var3:
        - 1
        - 2
        
    action: print
    message: $var1 $var2 $var3
```

```bash
endly -r=run msg=hello
```


**Using udf**:

The following pipeline provide example of using WorkingDirectory and FormatTime [UDFs](../udf).

@run.yaml
 ```yaml
init:
  appPath: $WorkingDirectory(../)
  bqTimeFormatArgs:
    - now
    - yyyy-MM-dd HH:mm:ss.SSSZ
  bqTimestamp: $FormatTime($bqTimeFormatArgs)

pipeline:
  run:
    action: print
    message: upper: $appPath <-> ts: $bqTimestamp
```


```bash
endly -r=run
```


Post processing state modification.

If action or workflow returns a data then post defines a way to publish result data to context.state, 
result data can also be access with task/action name.
 
@test.yaml
 ```yaml
pipeline:
  task1:
    init:
      - var1 =  $params.p1 
    action: exec:run
    target:
      URL: ssh://127.0.0.1
      credentials: localhost
    commands:
      - ls -al /tmp/
      
    post:
      stdout: $output
  task2:
    action: print
    message: $output        
  task3:
    action: print
    message: $task1.output        
```
