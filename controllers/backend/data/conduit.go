package data

import "github.com/rancher/rio/pkg/apply"

const (
	conduit = `
### Namespace ###
kind: Namespace
apiVersion: v1
metadata:
  name: conduit

### Service Account Controller ###
---
kind: ServiceAccount
apiVersion: v1
metadata:
  name: conduit-controller
  namespace: conduit

### Service Account Prometheus ###
---
kind: ServiceAccount
apiVersion: v1
metadata:
  name: conduit-prometheus
  namespace: conduit

### Controller ###
---
kind: Service
apiVersion: v1
metadata:
  name: api
  namespace: conduit
  labels:
    conduit.io/control-plane-component: controller
  annotations:
    conduit.io/created-by: conduit/cli v0.4.1
spec:
  type: ClusterIP
  selector:
    conduit.io/control-plane-component: controller
  ports:
  - name: http
    port: 8085
    targetPort: 8085

---
kind: Service
apiVersion: v1
metadata:
  name: proxy-api
  namespace: conduit
  labels:
    conduit.io/control-plane-component: controller
  annotations:
    conduit.io/created-by: conduit/cli v0.4.1
spec:
  type: ClusterIP
  selector:
    conduit.io/control-plane-component: controller
  ports:
  - name: grpc
    port: 8086
    targetPort: 8086

---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    conduit.io/created-by: conduit/cli v0.4.1
  creationTimestamp: null
  labels:
    conduit.io/control-plane-component: controller
  name: controller
  namespace: conduit
spec:
  replicas: 1
  strategy: {}
  template:
    metadata:
      annotations:
        conduit.io/created-by: conduit/cli v0.4.1
        conduit.io/proxy-version: v0.4.1
      creationTimestamp: null
      labels:
        conduit.io/control-plane-component: controller
        conduit.io/control-plane-ns: conduit
        conduit.io/proxy-deployment: controller
    spec:
      containers:
      - args:
        - public-api
        - -prometheus-url=http://prometheus.conduit.svc.cluster.local:9090
        - -controller-namespace=conduit
        - -log-level=info
        - -logtostderr=true
        image: gcr.io/runconduit/controller:v0.4.1
        imagePullPolicy: IfNotPresent
        name: public-api
        ports:
        - containerPort: 8085
          name: http
        - containerPort: 9995
          name: admin-http
        resources: {}
      - args:
        - destination
        - -log-level=info
        - -logtostderr=true
        image: gcr.io/runconduit/controller:v0.4.1
        imagePullPolicy: IfNotPresent
        name: destination
        ports:
        - containerPort: 8089
          name: grpc
        - containerPort: 9999
          name: admin-http
        resources: {}
      - args:
        - proxy-api
        - -log-level=info
        - -logtostderr=true
        image: gcr.io/runconduit/controller:v0.4.1
        imagePullPolicy: IfNotPresent
        name: proxy-api
        ports:
        - containerPort: 8086
          name: grpc
        - containerPort: 9996
          name: admin-http
        resources: {}
      - args:
        - tap
        - -log-level=info
        - -logtostderr=true
        image: gcr.io/runconduit/controller:v0.4.1
        imagePullPolicy: IfNotPresent
        name: tap
        ports:
        - containerPort: 8088
          name: grpc
        - containerPort: 9998
          name: admin-http
        resources: {}
      - env:
        - name: CONDUIT_PROXY_LOG
          value: warn,conduit_proxy=info
        - name: CONDUIT_PROXY_CONTROL_URL
          value: tcp://localhost:8086
        - name: CONDUIT_PROXY_CONTROL_LISTENER
          value: tcp://0.0.0.0:4190
        - name: CONDUIT_PROXY_METRICS_LISTENER
          value: tcp://0.0.0.0:4191
        - name: CONDUIT_PROXY_PRIVATE_LISTENER
          value: tcp://127.0.0.1:4140
        - name: CONDUIT_PROXY_PUBLIC_LISTENER
          value: tcp://0.0.0.0:4143
        - name: CONDUIT_PROXY_POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: gcr.io/runconduit/proxy:v0.4.1
        imagePullPolicy: IfNotPresent
        name: conduit-proxy
        ports:
        - containerPort: 4143
          name: conduit-proxy
        - containerPort: 4191
          name: conduit-metrics
        resources: {}
        securityContext:
          runAsUser: 2102
      initContainers:
      - args:
        - --incoming-proxy-port
        - "4143"
        - --outgoing-proxy-port
        - "4140"
        - --proxy-uid
        - "2102"
        - --inbound-ports-to-ignore
        - 4190,4191
        image: gcr.io/runconduit/proxy-init:v0.4.1
        imagePullPolicy: IfNotPresent
        name: conduit-init
        resources: {}
        securityContext:
          capabilities:
            add:
            - NET_ADMIN
          privileged: false
      serviceAccount: conduit-controller
status: {}
---
kind: Service
apiVersion: v1
metadata:
  name: web
  namespace: conduit
  labels:
    conduit.io/control-plane-component: web
  annotations:
    conduit.io/created-by: conduit/cli v0.4.1
spec:
  type: ClusterIP
  selector:
    conduit.io/control-plane-component: web
  ports:
  - name: http
    port: 8084
    targetPort: 8084
  - name: admin-http
    port: 9994
    targetPort: 9994

---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    conduit.io/created-by: conduit/cli v0.4.1
  creationTimestamp: null
  labels:
    conduit.io/control-plane-component: web
  name: web
  namespace: conduit
spec:
  replicas: 1
  strategy: {}
  template:
    metadata:
      annotations:
        conduit.io/created-by: conduit/cli v0.4.1
        conduit.io/proxy-version: v0.4.1
      creationTimestamp: null
      labels:
        conduit.io/control-plane-component: web
        conduit.io/control-plane-ns: conduit
        conduit.io/proxy-deployment: web
    spec:
      containers:
      - args:
        - -api-addr=api.conduit.svc.cluster.local:8085
        - -static-dir=/dist
        - -template-dir=/templates
        - -uuid=25ef0a93-0f43-4713-947c-c936758394c1
        - -controller-namespace=conduit
        - -log-level=info
        image: gcr.io/runconduit/web:v0.4.1
        imagePullPolicy: IfNotPresent
        name: web
        ports:
        - containerPort: 8084
          name: http
        - containerPort: 9994
          name: admin-http
        resources: {}
      - env:
        - name: CONDUIT_PROXY_LOG
          value: warn,conduit_proxy=info
        - name: CONDUIT_PROXY_CONTROL_URL
          value: tcp://proxy-api.conduit.svc.cluster.local:8086
        - name: CONDUIT_PROXY_CONTROL_LISTENER
          value: tcp://0.0.0.0:4190
        - name: CONDUIT_PROXY_METRICS_LISTENER
          value: tcp://0.0.0.0:4191
        - name: CONDUIT_PROXY_PRIVATE_LISTENER
          value: tcp://127.0.0.1:4140
        - name: CONDUIT_PROXY_PUBLIC_LISTENER
          value: tcp://0.0.0.0:4143
        - name: CONDUIT_PROXY_POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: gcr.io/runconduit/proxy:v0.4.1
        imagePullPolicy: IfNotPresent
        name: conduit-proxy
        ports:
        - containerPort: 4143
          name: conduit-proxy
        - containerPort: 4191
          name: conduit-metrics
        resources: {}
        securityContext:
          runAsUser: 2102
      initContainers:
      - args:
        - --incoming-proxy-port
        - "4143"
        - --outgoing-proxy-port
        - "4140"
        - --proxy-uid
        - "2102"
        - --inbound-ports-to-ignore
        - 4190,4191
        image: gcr.io/runconduit/proxy-init:v0.4.1
        imagePullPolicy: IfNotPresent
        name: conduit-init
        resources: {}
        securityContext:
          capabilities:
            add:
            - NET_ADMIN
          privileged: false
status: {}
---
kind: Service
apiVersion: v1
metadata:
  name: prometheus
  namespace: conduit
  labels:
    conduit.io/control-plane-component: prometheus
  annotations:
    conduit.io/created-by: conduit/cli v0.4.1
spec:
  type: ClusterIP
  selector:
    conduit.io/control-plane-component: prometheus
  ports:
  - name: admin-http
    port: 9090
    targetPort: 9090

---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    conduit.io/created-by: conduit/cli v0.4.1
  creationTimestamp: null
  labels:
    conduit.io/control-plane-component: prometheus
  name: prometheus
  namespace: conduit
spec:
  replicas: 1
  strategy: {}
  template:
    metadata:
      annotations:
        conduit.io/created-by: conduit/cli v0.4.1
        conduit.io/proxy-version: v0.4.1
      creationTimestamp: null
      labels:
        conduit.io/control-plane-component: prometheus
        conduit.io/control-plane-ns: conduit
        conduit.io/proxy-deployment: prometheus
    spec:
      containers:
      - args:
        - --storage.tsdb.retention=6h
        - --config.file=/etc/prometheus/prometheus.yml
        image: prom/prometheus:v2.2.1
        imagePullPolicy: IfNotPresent
        name: prometheus
        ports:
        - containerPort: 9090
          name: admin-http
        resources: {}
        volumeMounts:
        - mountPath: /etc/prometheus
          name: prometheus-config
          readOnly: true
        securityContext:
          runAsUser: 0
      - env:
        - name: CONDUIT_PROXY_LOG
          value: warn,conduit_proxy=info
        - name: CONDUIT_PROXY_CONTROL_URL
          value: tcp://proxy-api.conduit.svc.cluster.local:8086
        - name: CONDUIT_PROXY_CONTROL_LISTENER
          value: tcp://0.0.0.0:4190
        - name: CONDUIT_PROXY_METRICS_LISTENER
          value: tcp://0.0.0.0:4191
        - name: CONDUIT_PROXY_PRIVATE_LISTENER
          value: tcp://127.0.0.1:4140
        - name: CONDUIT_PROXY_PUBLIC_LISTENER
          value: tcp://0.0.0.0:4143
        - name: CONDUIT_PROXY_POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: gcr.io/runconduit/proxy:v0.4.1
        imagePullPolicy: IfNotPresent
        name: conduit-proxy
        ports:
        - containerPort: 4143
          name: conduit-proxy
        - containerPort: 4191
          name: conduit-metrics
        resources: {}
        securityContext:
          runAsUser: 2102
      initContainers:
      - args:
        - --incoming-proxy-port
        - "4143"
        - --outgoing-proxy-port
        - "4140"
        - --proxy-uid
        - "2102"
        - --inbound-ports-to-ignore
        - 4190,4191
        image: gcr.io/runconduit/proxy-init:v0.4.1
        imagePullPolicy: IfNotPresent
        name: conduit-init
        resources: {}
        securityContext:
          capabilities:
            add:
            - NET_ADMIN
          privileged: false
      serviceAccount: conduit-prometheus
      volumes:
      - configMap:
          name: prometheus-config
        name: prometheus-config
status: {}
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: prometheus-config
  namespace: conduit
  labels:
    conduit.io/control-plane-component: prometheus
  annotations:
    conduit.io/created-by: conduit/cli v0.4.1
data:
  prometheus.yml: |-
    global:
      scrape_interval: 10s
      evaluation_interval: 10s

    scrape_configs:
    - job_name: 'prometheus'
      static_configs:
      - targets: ['localhost:9090']

    - job_name: 'conduit-controller'
      kubernetes_sd_configs:
      - role: pod
        namespaces:
          names: ['conduit']
      relabel_configs:
      - source_labels:
        - __meta_kubernetes_pod_label_conduit_io_control_plane_component
        - __meta_kubernetes_pod_container_port_name
        action: keep
        regex: (.*);admin-http$
      - source_labels: [__meta_kubernetes_pod_container_name]
        action: replace
        target_label: component

    - job_name: 'conduit-proxy'
      kubernetes_sd_configs:
      - role: pod
      relabel_configs:
      - source_labels:
        - __meta_kubernetes_pod_container_name
        - __meta_kubernetes_pod_container_port_name
        action: keep
        regex: ^conduit-proxy;conduit-metrics$
      - source_labels: [__meta_kubernetes_namespace]
        action: replace
        target_label: namespace
      - source_labels: [__meta_kubernetes_pod_name]
        action: replace
        target_label: pod
      # special case k8s' "job" label, to not interfere with prometheus' "job"
      # label
      # __meta_kubernetes_pod_label_conduit_io_proxy_job=foo =>
      # k8s_job=foo
      - source_labels: [__meta_kubernetes_pod_label_conduit_io_proxy_job]
        action: replace
        target_label: k8s_job
      # __meta_kubernetes_pod_label_conduit_io_proxy_deployment=foo =>
      # deployment=foo
      - action: labelmap
        regex: __meta_kubernetes_pod_label_conduit_io_proxy_(.+)
      # drop all labels that we just made copies of in the previous labelmap
      - action: labeldrop
        regex: __meta_kubernetes_pod_label_conduit_io_proxy_(.+)
      # __meta_kubernetes_pod_label_foo=bar => foo=bar
      - action: labelmap
        regex: __meta_kubernetes_pod_label_(.+)

### Grafana ###
---
kind: Service
apiVersion: v1
metadata:
  name: grafana
  namespace: conduit
  labels:
    conduit.io/control-plane-component: grafana
  annotations:
    conduit.io/created-by: conduit/cli v0.4.1
spec:
  type: ClusterIP
  selector:
    conduit.io/control-plane-component: grafana
  ports:
  - name: http
    port: 3000
    targetPort: 3000

---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    conduit.io/created-by: conduit/cli v0.4.1
  creationTimestamp: null
  labels:
    conduit.io/control-plane-component: grafana
  name: grafana
  namespace: conduit
spec:
  replicas: 1
  strategy: {}
  template:
    metadata:
      annotations:
        conduit.io/created-by: conduit/cli v0.4.1
        conduit.io/proxy-version: v0.4.1
      creationTimestamp: null
      labels:
        conduit.io/control-plane-component: grafana
        conduit.io/control-plane-ns: conduit
        conduit.io/proxy-deployment: grafana
    spec:
      containers:
      - image: gcr.io/runconduit/grafana:v0.4.1
        imagePullPolicy: IfNotPresent
        name: grafana
        ports:
        - containerPort: 3000
          name: http
        resources: {}
        volumeMounts:
        - mountPath: /etc/grafana
          name: grafana-config
          readOnly: true
      - env:
        - name: CONDUIT_PROXY_LOG
          value: warn,conduit_proxy=info
        - name: CONDUIT_PROXY_CONTROL_URL
          value: tcp://proxy-api.conduit.svc.cluster.local:8086
        - name: CONDUIT_PROXY_CONTROL_LISTENER
          value: tcp://0.0.0.0:4190
        - name: CONDUIT_PROXY_METRICS_LISTENER
          value: tcp://0.0.0.0:4191
        - name: CONDUIT_PROXY_PRIVATE_LISTENER
          value: tcp://127.0.0.1:4140
        - name: CONDUIT_PROXY_PUBLIC_LISTENER
          value: tcp://0.0.0.0:4143
        - name: CONDUIT_PROXY_POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: gcr.io/runconduit/proxy:v0.4.1
        imagePullPolicy: IfNotPresent
        name: conduit-proxy
        ports:
        - containerPort: 4143
          name: conduit-proxy
        - containerPort: 4191
          name: conduit-metrics
        resources: {}
        securityContext:
          runAsUser: 2102
      initContainers:
      - args:
        - --incoming-proxy-port
        - "4143"
        - --outgoing-proxy-port
        - "4140"
        - --proxy-uid
        - "2102"
        - --inbound-ports-to-ignore
        - 4190,4191
        image: gcr.io/runconduit/proxy-init:v0.4.1
        imagePullPolicy: IfNotPresent
        name: conduit-init
        resources: {}
        securityContext:
          capabilities:
            add:
            - NET_ADMIN
          privileged: false
      volumes:
      - configMap:
          items:
          - key: grafana.ini
            path: grafana.ini
          - key: datasources.yaml
            path: provisioning/datasources/datasources.yaml
          - key: dashboards.yaml
            path: provisioning/dashboards/dashboards.yaml
          name: grafana-config
        name: grafana-config
status: {}
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: grafana-config
  namespace: conduit
  labels:
    conduit.io/control-plane-component: grafana
  annotations:
    conduit.io/created-by: conduit/cli v0.4.1
data:
  grafana.ini: |-
    instance_name = conduit-grafana

    [server]
    root_url = %(protocol)s://%(domain)s:/api/v1/namespaces/conduit/services/grafana:http/proxy/

    [auth]
    disable_login_form = true

    [auth.anonymous]
    enabled = true
    org_role = Editor

    [auth.basic]
    enabled = false

    [analytics]
    check_for_updates = false

  datasources.yaml: |-
    apiVersion: 1
    datasources:
    - name: prometheus
      type: prometheus
      access: proxy
      orgId: 1
      url: http://prometheus.conduit.svc.cluster.local:9090
      isDefault: true
      jsonData:
        timeInterval: "5s"
      version: 1
      editable: true

  dashboards.yaml: |-
    apiVersion: 1
    providers:
    - name: 'default'
      orgId: 1
      folder: ''
      type: file
      disableDeletion: true
      editable: true
      options:
        path: /var/lib/grafana/dashboards
        homeDashboardId: conduit-top-line
---

`
)

func addConduit() error {
	return apply.Content([]byte(conduit))
}
