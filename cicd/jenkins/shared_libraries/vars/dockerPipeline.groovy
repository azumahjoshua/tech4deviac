// def call(Map config = [:]) {

//     Map pipelineConfig = [
//         branch: config.branch ?: 'main',
//         dockerRepo: config.dockerRepo ?: 'joshua192',
//     ]

//     def dockerUtils = new org.tech4dev.DockerUtils()

//     pipeline {
//         agent any
//         environment {
//             VERSION = "v1.0.${BUILD_NUMBER}"
//         }
//         stages {
//             stage('Checkout') {
//                 steps {
//                     checkout scm
//                 }
//             }

//             stage('Lint and Format Check') {
//                 steps {
//                     dockerUtils.runLinter(language: 'go', directory: 'cicd/banking-app/backend-api')
//                     dockerUtils.runLinter(language: 'python', directory: 'cicd/banking-app/transaction-service')
//                     dockerUtils.runLinter(language: 'javascript', directory: 'cicd/banking-app/frontend')
//                 }
//             }

//             stage('Test') {
//                 steps {
//                     dockerUtils.runTests(language: 'go', directory: 'cicd/banking-app/backend-api')
//                     dockerUtils.runTests(language: 'python', directory: 'cicd/banking-app/transaction-service')
//                 }
//             }

//             stage('Build Docker Images') {
//                 steps {
//                     dockerUtils.buildDockerImages(imageName: '${config.dockerRepo}/banking-api', directory: 'cicd/banking-app/backend-api', push: true)
//                     dockerUtils.buildDockerImages(imageName: '${config.dockerRepo}/banking-processor', directory: 'cicd/banking-app/transaction-service', push: true)
//                     dockerUtils.buildDockerImages(imageName: '${config.dockerRepo}/banking-frontend', directory: 'cicd/banking-app/frontend', push: true)
//                 }
//             }
//         }
//     }
// }