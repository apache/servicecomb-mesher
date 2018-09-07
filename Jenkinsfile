pipeline {
    agent any

    stages {
       stage('Build Mesher') {
           steps {
                echo 'Making Package'
                sh 'bash -x scripts/build/build.sh'
            }
        }
        stage('Create Docker Image') {
            steps {
                echo 'Building Docker Image'
		sh 'bash -x scripts/build/build_image.sh'
           }
       }
    }
}

