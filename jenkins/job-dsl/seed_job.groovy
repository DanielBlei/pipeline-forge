def repoUrl = 'https://github.com/DanielBlei/pipeline-forge.git'
def branch  = '*/lab/jenkins-k8s'

def pipelines = [
  'ingest-test'        : 'jenkins/ci/ingest-test.Jenkinsfile',
  'operator-e2e-test'  : 'jenkins/ci/operator-e2e-test.Jenkinsfile',
  'operator-test'      : 'jenkins/ci/operator-test.Jenkinsfile',
  'trigger-test'       : 'jenkins/ci/trigger-test.Jenkinsfile',
]

pipelines.each { jobName, jenkinsfilePath ->
  pipelineJob(jobName) {
    definition {
      cpsScm {
        scm {
          git {
            remote { url(repoUrl) }
            branches(branch)
          }
        }
        scriptPath jenkinsfilePath
      }
    }
  }
}
