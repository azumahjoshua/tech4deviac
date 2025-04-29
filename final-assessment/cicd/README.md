# CI/CD (Jenkins)
## 5. Describe the typical stages you would include in a Jenkins pipeline for a containerized application. Why is each stage important?

A **Jenkins pipeline** for a containerized application typically follows a **CI/CD** workflow, ensuring code changes are built, tested, and deployed efficiently.

Breakdown of the common stages and why each one is important:

1. **Checkout / Clone Repository**
This is usually the first stage, where the source code is pulled from the version control system (e.g., Git).
**Why it's important:** Ensures that the latest version of the codebase, including recent team changes, is used for the pipeline run.

2. **Linting / Static Analysis**
Runs static analysis tools (e.g., ESLint, Pylint, golint) to check code quality.
**Why it's important:** Prevents common errors, enforces coding standards, and maintains a clean and maintainable codebase.

3. **Unit Tests**
Runs tests on individual components of the application.
**Why it's important:** Catches issues early by validating the correctness of core logic in isolated units before integration.

4. **Build**
Compiles or prepares the application for packaging (if applicable).
**Why it's important:** Validates that the app builds successfully, ensuring there are no compilation or build-time errors.

5. **Docker Build**
Builds the Docker image for the application using a Dockerfile.
**Why it's important:** Packages the application and its dependencies into a portable container, ensuring consistent environments across all stages.

6. **Docker Push**
Pushes the Docker image to a container registry (e.g., Docker Hub, Amazon ECR).
**Why it's important:** Makes the image available for deployment in testing, staging, or production environments.

7. **Deploy to Staging / Test Environment**
Deploys the container to a staging or test Kubernetes cluster or Docker environment.
**Why it's important:** Provides a production-like environment to validate application behavior and catch bugs before going live.

8. **Integration / End-to-End (E2E) Tests**
Runs tests that simulate real-world usage and verify how different components interact.
**Why it's important:** Ensures that integrated services communicate correctly and that the application works as expected from the user’s perspective.

9. **Security Scanning**
Uses tools like Trivy or Clair to scan Docker images for vulnerabilities.
**Why it's important:** Helps identify and mitigate security risks before deployment to production.

10. **Deploy to Production**
Rolls out the application to the live production environment using tools like Helm or kubectl.
**Why it's important:** Delivers the latest working version to end-users after passing all validations.

11. **Post-Deployment Verification**
Runs basic health checks to ensure the app is running correctly in production.
**Why it's important:** Confirms that the deployment was successful and that the application is responsive and functional.

## 6. Given a sample `Jenkinsfile`, identify and explain how environment variables and credentials should be managed securely.

1. Environment Variables
- Use `environment {}` block (Declarative Pipeline) to define variables.
- **Avoid hardcoding secrets**; use Jenkins **Credentials Store** instead.
```
  pipeline {
    agent any
    environment {
        // Non-sensitive variable
        APP_VERSION = "1.0.0"
        // Sensitive variable (loaded from Jenkins credentials)
        DB_PASSWORD = credentials('db-secret') // Refers to a Jenkins-stored credential
    }
    stages {
        stage('Example') {
            steps {
                sh 'echo "Using DB password: $DB_PASSWORD"'
            }
        }
    }
}
```
- **credentials('db-secret')** securely fetches the value from Jenkins credentials.
2. Credentials Management
- **Store secrets in Jenkins Credentials Manager** (via **"Manage Jenkins" > "Credentials"**).
3. Use `withCredentials` for temporary credential binding.
```
pipeline {
    agent any
    stages {
        stage('Deploy') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'aws-creds',
                    usernameVariable: 'AWS_ACCESS_KEY',
                    passwordVariable: 'AWS_SECRET_KEY'
                )]) {
                    sh 'echo "AWS Key: $AWS_ACCESS_KEY"'
                    sh 'aws configure --profile prod'
                }
            }
        }
    }
}
```
- `withCredentials` masks secrets in logs.

## 7. What are the benefits of using declarative pipelines in Jenkins? Provide a simple example.

Jenkins provides you with two ways of developing your pipeline code: Scripted and Declarative. Scripted pipelines, also known as “traditional” pipelines, are based on Groovy as their domain-specific language.

A Jenkins declarative pipeline provides a simplified and more friendly syntax with specific statements for defining them on top of the Pipeline sub-systems without needing to learn Groovy. All valid declarative pipelines must be enclosed within a pipeline block and contain all content and instructions for executing.

### Benefits of Using Declarative Pipelines in Jenkins

- Easier to read, write, and maintain, especially for teams or beginners.
- Stages like `post`, `options`, and `when` make error handling and conditional logic cleaner.
- Follows a standardized format, reducing human error and increasing predictability.
- Supports parallel stages and reusable `steps` blocks in a clean format.

```
pipeline {
    agent any

    environment {
        IMAGE_NAME = 'myapp'
        IMAGE_TAG = "${env.BUILD_NUMBER}"
    }

    stages {
        stage('Checkout') {
            steps {
                git 'https://github.com/user/myapp.git'
            }
        }

        stage('Lint') {
            steps {
                sh 'npm run lint'
            }
        }

        stage('Test') {
            steps {
                sh 'npm test'
            }
        }

        stage('Build Docker Image') {
            steps {
                sh 'docker build -t $IMAGE_NAME:$IMAGE_TAG .'
            }
        }

        stage('Deploy') {
            when {
                branch 'main'
            }
            steps {
                sh './deploy.sh'
            }
        }
    }

    post {
        always {
            echo 'Pipeline finished!'
        }
        failure {
            echo 'Something went wrong!'
        }
    }
}

```