// def call(Map config = [:]) {
//     def language = config.language ?: error("Missing required parameter 'language'")
//     def directory = config.directory ?: '.'

//     echo "Running tests for ${language} in ${directory}"

//     switch(language.toLowerCase()) {
//         case 'go':
//             docker.image('golang:1.24.1').inside(directory) {
//                 sh 'go test ./...'
//             }
//             break
//         case 'python':
//             docker.image('python:3.12').inside(directory) {
//                 sh 'python -m venv /opt/python-tools'
//                 sh '/opt/python-tools/bin/pip install pytest'
//                 sh '/opt/python-tools/bin/pytest'
//             }
//             break
//         default:
//             error "Unsupported language: ${language}"
//         break
//     }
// }