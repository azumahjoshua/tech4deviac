// vars/getCommitInfo.groovy
def call() {
    // Get Git commit information
    def commitHash = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()
    def commitAuthor = sh(returnStdout: true, script: 'git log -1 --pretty=format:"%an"').trim()
    def commitMessage = sh(returnStdout: true, script: 'git log -1 --pretty=format:"%s"').trim()

    // Return as a map
    return [
        hash: commitHash,
        author: commitAuthor,
        message: commitMessage
    ]
}