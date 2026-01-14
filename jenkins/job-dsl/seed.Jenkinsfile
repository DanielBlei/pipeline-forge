pipeline {
  agent any

  options {
    timeout(time: 10, unit: 'MINUTES')
  }

  stages {
    stage('Checkout') {
      steps {
        checkout scmGit(
          branches: scm.branches,
          extensions: [cloneOption(shallow: true, depth: 1, noTags: true)],
          userRemoteConfigs: scm.userRemoteConfigs
        )
      }
    }

    stage('Generate jobs') {
      steps {
        jobDsl targets: 'jenkins/job-dsl/*.groovy',
               removedJobAction: 'DELETE',
               removedViewAction: 'DELETE',
               lookupStrategy: 'JENKINS_ROOT'
      }
    }
  }
}
