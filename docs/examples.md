# 📋 Pipeline Forge Examples

This guide provides comprehensive examples of declarative YAML configurations for Pipeline Forge resources. Each example includes complete configurations with all optional fields and clear explanations of their purpose.

## 🎯 Quick Navigation

- [Staging Examples](#staging-examples)

  - [Ingestion by Reference](#ingestion-by-reference)
  - [Ingestion by Bootstrap](#ingestion-by-bootstrap)
  - [Ingestion by Trigger](#ingestion-by-trigger)
  - [Transform](#transform)

- [Trigger Examples](#trigger-examples)
  - [GCS Triggers](#gcs-triggers)
  - [Pub/Sub Triggers](#pubsub-triggers)
  - [BigQuery Triggers](#bigquery-triggers)
- [Additional Configurations](#additional-configurations)

---

## 🚀 Staging CRD Examples

Staging resources coordinate data ingestion and transformation phases with comprehensive lifecycle management and observability.

### Ingestion by Reference

This type of Ingestion is useful when you want to reference an existing CronJob or use an Event-Driven Trigger for the ingestion.

if CronJob:

- The Operator will watch for the CronJob to complete and then run the transform step.

if Trigger:

- The Operator will monitor the Trigger Resource status and run the transform step when the Trigger is in a Ready state.
- The Trigger needs to be created separately and referenced in the Staging resource.

```yaml
apiVersion: core.pipeline-forge.io/v1alpha1
kind: Staging
metadata:
  name: daily-sales-staging
  namespace: sales-analytics
spec:
  description: "Daily sales data processing pipeline"
  owner: sales-operations
# ------------------------------------------------------------
# Example 1: Ingestion by Reference (CronJob)
# ------------------------------------------------------------
  ingest:
    mode: reference
    type: cronjob
    name: daily-sales-ingest
    namespace: data-ingestion
    owner: data-engineering
    pollIntervalSeconds: 120
    suspend: false
    maxRetry: 2
# ------------------------------------------------------------
# Example 2: Ingestion by Reference (Trigger)
# ------------------------------------------------------------
  ingest:
    mode: reference
    type: trigger
    name: gcs-file-drop-trigger
    namespace: data-processing
    owner: data-engineering
    pollIntervalSeconds: 60
    suspend: false
    maxRetry: 3


  transform:
    name: sales-analytics-transform
    project: sales-analytics
    target: prod
    image: gcr.io/your-project/dbt-core:1.7.0
    models:
      - stg_daily_sales
      - fct_sales
      - dim_customers
    engine: dbt
    owner: data-analytics
    suspend: false
    full_refresh: false
    maxRetry: 3
```

### Ingestion by Bootstrap

Ingest with bootstrap mode will create a new CronJob for daily sales ingestion, managed by Pipeline Forge.

```yaml
apiVersion: core.pipeline-forge.io/v1alpha1
kind: Staging
metadata:
  name: daily-sales-staging
  namespace: sales-analytics
spec:
  description: "Daily sales data processing pipeline"
  owner: sales-operations
  ingest:
    mode: bootstrap
    type: cronjob
    name: daily-sales-ingest
    namespace: data-ingestion
    owner: data-engineering
    image: gcr.io/pipeline-forge/ingest:latest
    schedule: "0 0 * * *"
    suspend: false
    maxRetry: 3
    resources:
      requests:
        cpu: "500m"
        memory: "1Gi"
      limits:
        cpu: "2"
        memory: "4Gi"

  transform:
    name: sales-analytics-transform
    project: sales-analytics
    target: prod
    image: gcr.io/your-project/dbt-core:1.7.0
    models:
      - stg_daily_sales
    engine: dbt
    owner: data-analytics
    suspend: false
    full_refresh: false
    maxRetry: 3
```

Observations:

- By default only one Job from the CronJob will be running at a time (avoiding parallelization).
- Alterations to the Ingestion in the Staging resource will be reflected in the respective CronJob.
- The CronJob is attached to the Staging resource (by Finalizers) and will be deleted only if the Staging resource is deleted.

### Ingestion by Trigger

This type of Ingestion is useful when you want to use an Event-Driven Trigger for the ingestion.

The Operator will monitor the Trigger Resource status and run the transform step when the Trigger is in a Ready state.

The Trigger needs to be created separately and referenced in the Staging resource.

```yaml
apiVersion: core.pipeline-forge.io/v1alpha1
kind: Staging
metadata:
  name: gcs-file-staging
  namespace: file-processing
spec:
  description: "Event-driven file processing pipeline"
  owner: data-engineering

  ingest:
    mode: reference
    type: trigger
    name: gcs-file-drop-trigger
    namespace: data-processing
    owner: data-engineering
    maxRetry: 3

  transform:
    name: file-processing-transform
    project: file-processing
    target: prod
    image: gcr.io/pipeline-forge/transform:latest
    models:
      - stg_file_data
    engine: dbt
    owner: data-analytics
    suspend: false
    full_refresh: false
    maxRetry: 2
```

**What This Does**: References a GCS trigger and processes uploaded files through dbt transformations.

### Transform

Transform resources coordinate data transformation phases with comprehensive lifecycle management and observability.

````yaml
apiVersion: core.pipeline-forge.io/v1alpha1
kind: Staging
metadata:
  name: sales-analytics-staging
  namespace: sales-analytics
spec:
  description: "Daily sales data processing pipeline"
  owner: sales-operations
  ingest:
    mode: reference
    type: cronjob
    name: daily-sales-ingest

  transform:
    name: sales-analytics-transform
    project: sales-analytics
    target: prod
    image: gcr.io/pipeline-forge/transform:latest
    models:
      - stg_daily_sales
    engine: dbt
    owner: data-analytics
    suspend: false
    full_refresh: false
    maxRetry: 3
    resources:
      requests:
        cpu: "500m"
        memory: "1Gi"
      limits:
        cpu: "2"
        memory: "4Gi"

Observations:

- The Transform step will be executed only once, after the Ingest step is completed.
- The Operator will launch a Job for the Transform step, and will wait for it.
- `resources` field is optional and can be used to adjust the resources for the Transform step/job.
- `args` overrids the image entrypoint and can be used to run specific dbt commands, e.g `dbt run --select +tag:my_tag`

## 🔄 Trigger Examples

Triggers provide event-driven pipeline ingestion through various cloud services and APIs. (Currently only Google Cloud Platform is supported)

Trigger Types:

- GCS
- Pub/Sub
- BigQuery

### 📁 GCS Triggers

Monitor a GCS bucket for new data files and trigger processing when files are uploaded.

```yaml
apiVersion: core.pipeline-forge.io/v1alpha1
kind: Trigger
metadata:
  name: gcs-file-drop-trigger
  namespace: data-processing
  labels:
    app: pipeline-forge
    trigger-type: gcs
spec:
  type: gcs
  name: gcs-file-drop-trigger
  description: "Triggers pipeline when new files are uploaded to the data bucket"
  owner: data-engineering
  image: gcr.io/your-project/gcs-trigger:latest
  schedule: "*/5 * * * *" # Check every 5 minutes
  cooldownIntervalSeconds: 600 # 10 minute cooldown (Cooldown aims to prevent re-triggering)
  maxRetry: 3
  gcs:
    bucket: "my-data-bucket"
    prefix: "exports/"
````

**What This Does**: Monitors the `my-data-bucket` GCS bucket for files in the `exports/` prefix, checking every 5 minutes with a 5-minute cooldown between triggers.

### 📡 Pub/Sub Triggers

React to real-time messages from Google Cloud Pub/Sub with optional filtering. (Currently only Google Cloud Platform is supported)

```yaml
apiVersion: core.pipeline-forge.io/v1alpha1
kind: Trigger
metadata:
  name: user-events-trigger
  namespace: user-analytics
spec:
  type: pubsub
  name: user-events-trigger
  description: "Triggers pipeline on user activity events"
  owner: product-analytics
  image: gcr.io/your-project/pubsub-trigger:latest
  args:
    - "--topic=user-events"
    - "--subscription=user-events-sub"
  schedule: "*/10 * * * *" # Check every 10 minutes
  maxRetry: 3
  pubsub:
    topic: "user-events"
    messageFilter:
      attribute: "event_type"
      equals: "user_activity"
```

**What This Does**: Monitors the `user-events` Pub/Sub topic and triggers only on messages with `event_type` attribute equal to `user_activity`.

### 📊 BigQuery Trigger

Monitor a BigQuery table for new data or table updates.

```yaml
apiVersion: core.pipeline-forge.io/v1alpha1
kind: Trigger
metadata:
  name: bq-freshness-trigger
  namespace: data-quality
spec:
  type: bigquery
  name: bq-freshness-trigger
  description: "Monitors BigQuery table freshness and triggers pipeline updates"
  owner: data-engineering
  image: gcr.io/your-project/bigquery-trigger:latest
  schedule: "0 */6 * * *" # Check every 6 hours
  cooldownIntervalSeconds: 3600 # 1 hour cooldown
  runOnce: false
  maxRetry: 3
  bigquery:
    project_id: "my-gcp-project"
    dataset_id: "analytics"
    table_id: "daily_metrics"
```

**What This Does**: Monitors the `daily_metrics` table in BigQuery and triggers when data is refreshed, checking every 6 hours.

## 🔧 Additional Configurations

Complex scenarios combining multiple features and advanced configurations. Extra configurations are available for triggers, ingest and transform.

- `args`: override the default command for the image and **ignore** dedicated arguments if provided (ignores: `gcs`, `pubsub`, `bigquery`, `models`)
- `resources`: override the default resources (Memory and CPU) for the trigger, transform and ingest (bootstrap mode only)
- `schedule`: set the cadence for the trigger, transform and ingest (bootstrap mode only)
- `runOnce`: if true, the step will be executed only once
- `PollIntervalSeconds`: override the default polling interval for the Ingest step, reducing/increasing the frequency of checks for changes in the ingestion resource
- `Suspend`: if true, the step will be suspended and not processed
