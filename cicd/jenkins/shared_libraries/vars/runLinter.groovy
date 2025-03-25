// def call(Map config = [:]) {
//     def language = config.language ?: error("Missing required parameter 'language'")
//     def directory = config.directory ?: '.'

//     echo "Running linter for ${language} in ${directory}"

//     switch(language.toLowerCase()) {
//         case 'go':
//             docker.image('golang:1.24.1').inside(directory) {
//                 sh 'go fmt ./...'
//                 sh 'go vet ./...'
//             }
//             break
//         case 'python':
//             docker.image('python:3.12').inside(directory) {
//                 echo 'Running flake8...'
//                 sh 'python -m venv /opt/python-tools'
//                 sh '/opt/python-tools/bin/pip install flake8'
//                 sh '/opt/python-tools/bin/flake8 . --exit-zero --count --select=E9,F63,F7,F82 --show-source --statistics'
//             }
//             break
//         case 'javascript':
//             docker.image('node:22').inside(directory) {
//                 sh 'npm install eslint'
//                 sh 'npx eslint . || true' 
//             }
//             break
//         default:
//             error "Unsupported language: ${language}"
//         break
//     }
// }