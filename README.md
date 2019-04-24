# Ratibor

It plows data and exposes metrics based on this data readable for Prometheus.

## Config

The example of config:

```yaml
input:
  api:
  - name: "some_api_input"
    listenAddr: 9999
    path: /some/api
template:
  type: plain
  pattern: "{{.http_code}} | {{.duration}} | {{.request}} | {{.method}}"
  delimiter: "|"
metrics:
- name: requests_total
  type: counter
  help: "Total requests"
  expr: '{{if ne .http_code "1"}}12{{else}}{{end}}'
  value: "1"
  labels:
  - name: "http_code"
    value: {{.http_code}}
  - name: "method"
    value: {{.method}}
- name: errors_total
  type: counter
  help: "Total Errors"
  expr: '{{if ne .http_code "200"}}1{{end}}'
  value: "1"
- name: response_time_milliseconds
  type: summary
  value: "{{.duration}}"
```

### Supported types

Ratibor supports the following types of templates:
 - `plain`: is a one line data, the values in which are separated by delimiter. For example: `value 1 | value 2 | value 3`. The delimiter is `|`.


### Metric value evaluation

To evaluate metric value the two fields are used: `expr` and `value`. There are 3 cases of using this fields:

1. `expr` and `value` are not empty. In this case if evaluated expr is not empty, then if `value` can be considered as float64 value, it will be as metric value, else if `value` is not empty(i.e. it is some string), it will be as `float(1)` as metric value, else `float(0)`:

Examples:

```yaml
# http_code=200

- name: errors_total
  type: counter
  help: "Total Errors"
  expr: '{{if ne .http_code "200"}}1{{end}}'
  value: "1"

# nothing will be used as metric value
```

```yaml
# http_code=100

- name: errors_total
  type: counter
  help: "Total Errors"
  expr: '{{if ne .http_code "200"}}1{{end}}'
  value: "1"

# float64(1) will be used as metric value
```


```yaml
# http_code=100, page=/search

- name: errors_total
  type: counter
  help: "Total Errors"
  expr: '{{if ne .http_code "200"}}1{{end}}'
  value: "{{if eq .page "/page"}}2{{else}}3{{end}}"

# float64(2) will be used as metric value
```

```yaml
# http_code=100, page=/product

- name: errors_total
  type: counter
  help: "Total Errors"
  expr: '{{if ne .http_code "200"}}1{{end}}'
  value: "{{if eq .page "/page"}}2{{else}}3{{end}}"

# float64(3) will be used as metric value
```

```yaml
# http_code=100, page=/product

- name: errors_total
  type: counter
  help: "Total Errors"
  expr: '{{if ne .http_code "200"}}1{{end}}'
  value: "{{if eq .page "/page"}}some_value{{else}}another_value{{end}}"

# float64(1) will be used as metric value
```

2. `expr` is not empty, but `value` is empty, then if `expr` can be considered as float64, it will be as metric value, else if `expr` is not empty(i.e. some string), then `float64(1)` value will be as metric value, else nothing will be exposed.

```yaml
# http_code=100

- name: errors_total
  type: counter
  help: "Total Errors"
  expr: '{{if ne .http_code "200"}}2{{end}}'

# float64(2) will be used as metric value
```

```yaml
# http_code=100

- name: errors_total
  type: counter
  help: "Total Errors"
  expr: '{{if ne .http_code "200"}}yes{{end}}'

# float64(1) will be used as metric value
```

```yaml
# http_code=200

- name: errors_total
  type: counter
  help: "Total Errors"
  expr: '{{if ne .http_code "200"}}yes{{end}}'

# nothing will be used as metric value
```

3. `expr` is empty, but `value` is not empty, then if `value` can be considered as float64, it will be as metric value, else if `value` is not empty, `float64(1)` will be exposed as metric value else `float64(0)` will be exposed.

```yaml
# page=/page

- name: errors_total
  type: counter
  help: "Total Errors"
  value: "{{if eq .page "/page"}}2{{else}}3{{end}}"

# float64(2) will be used as metric value
```

```yaml
# page=/product

- name: errors_total
  type: counter
  help: "Total Errors"
  value: "{{if eq .page "/page"}}3{{end}}"

# float64(0) will be used as metric value
```

```yaml
# page=/page

- name: errors_total
  type: counter
  help: "Total Errors"
  value: "{{if eq .page "/page"}}yes{{end}}"

# float64(1) will be used as metric value
```


### Inputs

Input is a source of data. There are several types of inputs.

#### Syslog

To set up syslog input add the following settings into configuration:

```yaml
input:
  syslog:
    - name: "some_syslog_input"
      listenAddr: ":9876"
```

#### API

To receive data via API it is needed to add the following configuration:

```yaml
input:
  api:
    - name: "some_api_input"
      listenAddr: ":9999"
      path: "/my_input"
```
