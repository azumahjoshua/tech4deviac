# Kubernetes & Orchestration
## 12. Explain the difference between Kubernetes Deployments, StatefulSets, and DaemonSets. When would you use each?

1. Deployment
A Deployment is used to manage stateless applications. It ensures that a specified number of replicas of a Pod are running at all times and can automatically handle updates and rollbacks.

- Pod Identity: Pods managed by a Deployment are interchangeable; each has a randomly generated name and no stable identity.

- Storage: Pods typically use ephemeral storage, meaning data is lost when the Pod is terminated or rescheduled.

- Scaling & Updates: Supports easy scaling and rolling updates. New Pods are rolled out gradually, and old Pods are terminated.

- Use Cases: Ideal for frontend services, REST APIs, web applications, and other stateless workloads.

    Example Use Case: Running a Node.js web server behind a load balancer.

2. StatefulSet
A StatefulSet is designed for stateful applications that require stable, persistent identity and storage.

- Pod Identity: Each Pod gets a stable hostname (e.g., app-0, app-1) and maintains its identity across restarts.

- Storage: Typically used with PersistentVolumeClaims (PVCs) to retain data even if the Pod is deleted or moved.

- Ordering: Pods are created, scaled, and deleted in a defined order, which is important for clustered or leader-based applications.

- Use Cases: Suitable for databases, distributed systems, or any service requiring persistent state and identity.

    Example Use Case: Running a MySQL or MongoDB cluster where each node needs a consistent identity and its own data volume.

3. A DaemonSet ensures that exactly one instance of a specific Pod runs on each node in the cluster.

- Pod Identity: One Pod per node, not tied to identity the same way as StatefulSets.

- Scheduling: Automatically deploys the Pod to all current and future nodes.

- Use Cases: Used for node-level background tasks such as monitoring, logging, or networking agents.

    Example Use Case: Deploying a Fluentd log collector or Prometheus Node Exporter to gather metrics from every node.

## 13. Describe the process of deploying an application using Helm. What are the advantages of using Helm charts?

Deploying an application using Helm involves using Helm charts, which are collections of YAML templates that define Kubernetes resources in a reusable and configurable way. Helm simplifies Kubernetes deployments by packaging applications into versioned, installable units.

### Process of Deploying an Application Using Helm

- Install Helm on the local machine
- Choose or create a Helm chart
    - pull from a public repo (e.g., Bitnami) with `helm repo` add & `helm search`
    - scaffold your own with `helm create <chart-name>`.
- Review chart structure
    - Charts bundle Kubernetes manifests (Deployments, Services, ConfigMaps, Secrets, etc.)
    - All templates reference values defined in `values.yaml.`
- Customize your deployment by editing `values.yaml`
    - Specify Docker image (registry, repo, tag), replicas count, ports, environment variables, resource requests/limits, and any other config.
- Deploy the chart

    ```
    helm install <release-name> <chart-path-or-repo/chart>

    ```
    Helm renders templates with your values and creates the resources via the Kubernetes API.

- Verify and inspect

    ```
    helm status <release-name>

    kubectl get all -l app.kubernetes.io/instance=<release-name>

    ```

- Roll out updates

    ```
    helm upgrade <release-name> <chart-path> --set image.tag=v2

    ```
    Helm performs a rolling update based on your new values.

- Roll back if necessary

    ```
    helm rollback <release-name> [revision]

    ```
    Instantly revert to a previous release revision.

- Uninstall when done

    ```
    helm uninstall <release-name>

    ```
    Cleans up all resources created by that release.

### Advantages of Using Helm Charts

- **Reduced YAML Sprawl**

Templates minimize repetitive manifest code, enforcing DRY principles and reducing human error.
Packaging & Abstraction
Encapsulate complex multi-resource apps into a single, versioned package.

- **Reusability & Consistency**

Share and reuse charts across teams and environments (dev, staging, prod) with predictable outcomes.

- **Declarative Configuration**

Centralize all environment-specific settings in `values.yaml`, separating code from config.

- **Release Management**

Maintain a history of releases—easily upgrade, rollback, or inspect past deployments.

- **Dependency Management**

Define chart dependencies (e.g., a web app chart depending on a database chart) for coordinated installs.

## 14. How would you securely inject secrets into a Kubernetes deployment? Provide an example using Kubernetes Secrets.

A Secret is an object that contains a small amount of sensitive data such as a password, a token, or a key. Such information might otherwise be put in a Pod specification or in a container image. Using a Secret means that you don't need to include confidential data in your application code.

Kubernetes Secrets are, by default, stored unencrypted in the API server's underlying data store (etcd).

### Create a Secret from Literal Values

- To create a Secret, we can use the imperative `kubectl create secret` command:

    ```

    kubectl create secret generic my-password \
  --from-literal=password=mysqlpassword

    ```

### Create a Secret from a Definition Manifest

We can create a Secret manually from a YAML definition manifest.

- With data maps, each value of a sensitive information field must be encoded using base64. 
    
    ```
        echo mysqlpassword | base64
        bXlzcWxwYXNzd29yZAo=
    
    ```
    and then use it in the definition manifest:

    ```
        apiVersion: v1
        kind: Secret
        metadata:
          name: my-password
        type: Opaque
        data:
          password: bXlzcWxwYXNzd29yZAo=

    ```
- With `stringData` maps, there is no need to encode the value of each sensitive information field. The value of the sensitive field will be encoded when the my-password Secret is created:

    ```
       apiVersion: v1
        kind: Secret
        metadata:
          name: my-password
        type: Opaque
        stringData:
          password: mysqlpassword

    ```

- Using the `mypass.yaml` definition file we can now create a secret with `kubectl create` command:

    ```
        kubectl create -f mypass.yaml

    ```
### Use Secrets Inside Pods: As Environment Variables

- We reference only the password key of the `my-password` Secret and assign     its value to the `DB_PASSWORD` environment variable:

    ```
        ....
        spec:
          containers:
          - image: mysql:lts
            name: mysql
            env:
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: my-password
                  key: password
        ....

    ```

### Use Secrets Inside Pods: As Volumes

- We can also mount a Secret as a Volume inside a Pod. The `secret` Volume plugin converts the Secret object into a mountable resource.

    ```
        ....
        spec:
          containers:
          - image: wordpress:4.7.3-apache
            name: wordpress
            volumeMounts:
            - name: secret-volume
              mountPath: "/etc/secret-data"
              readOnly: true
          volumes:
          - name: secret-volume
            secret:
              secretName: my-password

        ....

    ```

### Advanced Methods

-  it's best to `sync secrets from external providers` (like AWS Secrets Manager, HashiCorp Vault, Azure Key Vault) into Kubernetes, instead of manually creating them.

- `External Secrets Operator` automates this by fetching secrets from external secret stores and creating Kubernetes Secrets.

```
    apiVersion: external-secrets.io/v1beta1
    kind: ExternalSecret
    metadata:
      name: db-secret-from-aws
    spec:
      refreshInterval: 1h
      secretStoreRef:
        name: aws-secret-store
        kind: SecretStore
      target:
        name: my-db-secret
      data:
      - secretKey: DB_PASSWORD
        remoteRef:
          key: /prod/database/password

```

## 15. Given a scenario where you need to scale an application based on CPU usage, explain how you would configure Horizontal Pod Autoscaling in Kubernetes.