
pipeline {
    options {
    buildDiscarder(logRotator(numToKeepStr: '50', artifactNumToKeepStr: '50'))
    }
    agent any 
    triggers { 
        pollSCM('H/2 * * * *') 
    } 
    stages {
        stage('Build') { 
            steps {
                echo "Build stage"
                sh 'chmod +x buildScripts/build.sh'
                sh('buildScripts/build.sh')
            }
        }
        stage('Test') { 
            steps {
                echo "Test"
            }
        }
        stage('Deploy') { 
            steps {
                echo "Deploy"
                sh 'chmod +x buildScripts/deploy.sh'
                sh('buildScripts/deploy.sh')
            }
        }
    }
}