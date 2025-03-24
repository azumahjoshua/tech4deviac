FROM jenkins/jenkins:lts
USER root
RUN apt-get update -qq \
    && apt-get install -qqy \
    apt-transport-https \
    ca-certificates \ 
    curl \
    wget \
    gnupg2 \
    software-properties-common \
    lsb-release \
    git \
    unzip \
    build-essential

RUN curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg \
    && echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null \
    && apt-get update \
    && apt-get install -y docker-ce-cli

RUN wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz \
    && rm go1.22.0.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

RUN apt-get install -y python3-full python3-venv pipx

# Create a virtual environment for Python tools
RUN python3 -m venv /opt/python-tools
ENV PATH="/opt/python-tools/bin:${PATH}"

# Install Python tools in the virtual environment
RUN /opt/python-tools/bin/pip install \
    flake8 \
    black \
    isort \
    pylint \
    mypy \
    pytest

RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g \
    eslint \
    prettier \
    typescript \
    ts-node \
    tslint

RUN apt-get clean && rm -rf /var/lib/apt/lists/*

RUN chown -R jenkins:jenkins /opt/python-tools

USER jenkins