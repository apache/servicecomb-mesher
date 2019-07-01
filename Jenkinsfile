pipeline {
    agent any

    stages {
       stage('Build proxy') {
           steps {
                echo 'Making Package'
                sh 'bash -x ci/build.sh'
            }
        }
        stage('Create proxy image') {
            steps {
                echo 'Building Docker Image'
		sh 'bash -x ci/build_image.sh'
           }
       }
    }
}

