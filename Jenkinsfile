pipeline {
  agent any
  stages {
    stage('build') {
      steps {
        echo 'building the application.... with automatic'
        echo "yes triggered by webhook"
      }
    }

    stage('test') {
      steps {
        echo 'testing the application ....'
      }
    }

    stage('deploy') {
      steps {
        echo 'deploying the applicatiom...'
      }
    }
  }
}
