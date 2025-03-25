def call(Map config = [:]) {
    def directory = config.directory ?: '.'
    def imageName = config.imageName ?: error("Missing required parameter 'imageName'")
    def tag = config.tag ?: 'latest'
    def dockerfile = config.dockerfile ?: 'Dockerfile'
    def buildArgs = config.buildArgs ?: ''
    def push = config.push ?: false

    echo "Building Docker image ${imageName}:${tag} from ${dockerfile} in ${directory}"

    docker.build("${imageName}:${tag}", dockerfile, buildArgs)

    if(push) {
        echo "Pushing Docker image ${imageName}:${tag}"
        docker.withRegistry('https://registry.hub.docker.com', 'docker-hub-credentials') {
            docker.image("${imageName}:${tag}").push()
        }
    }
}