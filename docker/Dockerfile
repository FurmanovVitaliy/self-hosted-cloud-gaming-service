# Use the base Fedora image
FROM fedora:latest

# Update the system and install required packages
RUN dnf -y update && dnf -y install \
    fuse \
    #ffmpeg \
    wget \
    glibc-locale-source \
    && dnf clean all

# Set the locale to en_US.UTF-8
RUN localedef -c -i en_US -f UTF-8 en_US.UTF-8

# Copy scripts and download files
WORKDIR /home
COPY run-stream.sh run-ffmpeg.sh
RUN chmod +x run-ffmpeg.sh
RUN wget -O yuzu.AppImage https://github.com/yuzu-emu/yuzu-mainline/releases/download/mainline-0-1606/yuzu-mainline-20231031-a0398cff0.AppImage
RUN chmod +x yuzu.AppImage

# Example command to run when the container starts
CMD ["./yuzu.AppImage"]
