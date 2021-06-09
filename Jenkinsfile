library identifier: 'c3i@master', changelog: false,
  retriever: modernSCM([$class: 'GitSCMSource', remote: 'https://pagure.io/c3i-library.git'])

def artifact_glob="build/*"
// def build_image="registry.redhat.io/rhel8/go-toolset:1.15"
def build_image="quay.io/ligangty/spmm-jenkins-agent-go-centos7:latest"
// backup build image
// def build_image = "quay.io/app-sre/ubi8-go-toolset:1.15.7"
pipeline {
  agent {
    kubernetes {
      cloud params.JENKINS_AGENT_CLOUD_NAME
      label "jenkins-slave-${UUID.randomUUID().toString()}"
      serviceAccount "jenkins"
      defaultContainer 'jnlp'
      yaml """
      apiVersion: v1
      kind: Pod
      metadata:
        labels:
          app: "jenkins-${env.JOB_BASE_NAME}"
          indy-pipeline-build-number: "${env.BUILD_NUMBER}"
      spec:
        containers:
        - name: jnlp
          image: ${build_image}
          imagePullPolicy: Always
          tty: true
          env:
          - name: HOME
            value: /home/jenkins
          - name: GOROOT
            value: /usr/lib/golang
          - name: GOPATH
            value: /home/jenkins/gopath
          - name: GOPROXY
            value: https://proxy.golang.org
          resources:
            requests:
              memory: 4Gi
              cpu: 2000m
            limits:
              memory: 8Gi
              cpu: 4000m
          workingDir: /home/jenkins
      """
    }
  }
  options {
    //timestamps()
    timeout(time: 120, unit: 'MINUTES')
  }
  environment {
    PIPELINE_NAMESPACE = readFile('/run/secrets/kubernetes.io/serviceaccount/namespace').trim()
    PIPELINE_USERNAME = sh(returnStdout: true, script: 'id -un').trim()
  }
  stages {
    stage('Prepare') {
      steps {
        sh 'printenv'
      }
    }

    stage('git checkout') {
      steps{
        script{
          checkout([$class      : 'GitSCM', branches: [[name: 'main']], doGenerateSubmoduleConfigurations: false,
                    extensions  : [[$class: 'CleanCheckout']], submoduleCfg: [],
                    userRemoteConfigs: [[url: 'https://github.com/ligangty/indy-tests.git', refspec: '+refs/heads/*:refs/remotes/origin/* +refs/pull/*/head:refs/remotes/origin/pull/*/head']]])
          env.GIT_COMMIT = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()

          echo "Building main commit: ${env.GIT_COMMIT}"
        }
      }
    }

    stage('Build') {
      steps {
        sh 'make build'
      }
    }

    stage('Archive') {
      steps {
        echo "Archive"
        archiveArtifacts artifacts: "$artifact_glob", fingerprint: true
      }
    }
  }
  post {
    success {
      script {
        echo "SUCCEED"
      }
    }
    failure {
      script {
        echo "FAILED"
      }
    }
  }
}


