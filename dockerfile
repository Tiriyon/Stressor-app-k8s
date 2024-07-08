# Use UBI base image
FROM registry.access.redhat.com/ubi8/ubi:latest

# Install necessary development tools
RUN yum install -y \
    golang \
    git \
    vim \
    sudo \
    && yum clean all

# Create a new user 'tiriyon'
RUN useradd -m -s /bin/bash tiriyon \
    && echo 'tiriyon ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# Set up the Go environment
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

# Set the correct permissions for the Go module cache directory
RUN mkdir -p /go/pkg/mod && chown -R tiriyon:tiriyon /go

# Switch to user 'tiriyon'
USER tiriyon
WORKDIR /home/tiriyon

# Expose the application port
EXPOSE 8080

# Command to start a shell
CMD ["/bin/bash"]

