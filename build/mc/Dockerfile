FROM docker.io/library/openjdk:26-jdk-slim

LABEL application="Minecraft Server" version="1.21.5" org.opencontainers.image.authors="f@lex.la"

# Expose the Minecraft server port
EXPOSE 25565

# Create a directory for the server
ADD https://piston-data.mojang.com/v1/objects/e6ec2f64e6080b9b5d9b471b291c33cc7f509733/server.jar /server/server.jar

# Set permissions for the server jar file
RUN chmod 777 /server/server.jar

# Use a non-root user to run the server
RUN useradd -m -s /bin/bash mcuser
USER mcuser

# Set the working directory
WORKDIR /app

# Java memory settings
ENV JAVA_OPTS="-Xmx3G -Xms3G"

# Set the command to run the server
CMD ["java", "-jar", "/server/server.jar", "nogui"]