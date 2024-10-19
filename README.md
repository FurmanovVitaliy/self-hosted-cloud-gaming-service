# Pixel Cloud
#### Technology Stack
![Go version](https://img.shields.io/badge/Go-1.22.2-blue?logo=go)
![Docker](https://img.shields.io/badge/Docker-27.2.0-blue?logo=docker)
![MongoDB](https://img.shields.io/badge/MongoDB-7.0.14-green?logo=mongodb)
![WebRTC](https://img.shields.io/badge/WebRTC-1.0.39-yellowgreen?logo=webrtc)
![Vite](https://img.shields.io/badge/Vite-5.2-%23F7C845?logo=vite)
![TypeScript](https://img.shields.io/badge/TypeScript-5.4.5-%232b7489?logo=typescript)
![SaaS](https://img.shields.io/badge/Sass-1.77-%23c25b5d?logo=sass)
![HTML](https://img.shields.io/badge/HTML-5-%23E44D26?logo=html5)
![JSX](https://img.shields.io/badge/JSX-Template%20Syntax-%23F7DF1E?logo=javascript)




#### Packages Versions
![FFmpeg](https://img.shields.io/badge/FFmpeg-7.1-green?logo=ffmpeg)
![Arch Linux](https://img.shields.io/badge/Arch%20Linux-rolling-brightgreen?logo=archlinux)

## Introduction

I am excited to present my project, which I have been working on for the past year. It represents my personal take on cloud gaming and was inspired by projects like Google Stadia and Microsoft XCloud. The main idea of the project is to provide users with the ability to run AAA games on thin clients (such as phones, tablets, and PCs without graphics cards needed for complex graphical computations) through cloud technologies.

## Demo

[![Watch the demo](./assets/demo.gif)](https://youtu.be/q_k8pBCw4QU)


## Development environment
- **Operating System**: Arch Linux.
- **Programming Language**: Go.
- **Hardware Requirements**: AMD GPU (not tested with NVIDIA GPUs).
- **Installed Packages**:
  - `yay`: is used for installing necessary packages.
  - `xorg-xrandr`: Utility for managing screen resolution and display settings in X11.
  - `drm_info`: For gathering information about Direct Rendering Manager (DRM) on the system.
   - `mongodb`: Database used in the project.
   - `ffmpeg`: The latest version of FFmpeg for handling video encoding and streaming.
   - `js`: JavaScript-related tools or libraries.
- **AMD GPU Drivers**:
   - `mesa`, `lib32-mesa`: Open-source graphics drivers.
    - `xf86-video-amdgpu`: Driver for AMD GPUs.
    - `vulkan-radeon`, `lib32-vulkan-radeon`: Vulkan support for AMD GPUs.
    - `libva-mesa-driver`: For video acceleration (VAAPI) on AMD GPUs.
- **Certificates**: `mkcert` is used for generating local certificates.

**Setup**:
  1. Ensure that Docker and Docker Compose are installed.
  2. Install required packages using:
  3. ```bash
     ffmpeg js xorg-xrandr mesa lib32-mesa xf86-video-amdgpu vulkan-radeon  lib32-vulkan-radeon libva-mesa-driver mkcert
        ```
  4. Install required packages using from AUR via `yay` or another manager:
     ```bash
     yay -S drm_info mongodb-bin
     ```
  5. Clone the main repository:
     ```bash
     git clone https://github.com/FurmanovVitaliy/self-hosted-cloud-gaming-service.git
     cd self-hosted-cloud-gaming-service
     ```
  6. Clone the repository containing the necessary Docker images and build them :
     ```bash
     git clone https://github.com/FurmanovVitaliy/gaming-vm-in-dockers
     cd gaming-vm-in-dockers
        # Build the base image
     docker build -t arch:base ./base
        # Build the ffmpeg-audio image
     docker build -t arch:ffmpeg-audio ./ffmpeg-audio
        # Build the ffmpeg-video image
     docker build -t arch:ffmpeg-video ./ffmpeg-video
        # Build the portprotone image
     docker build -t arch:portprotone ./portprotone
        # Build the pulseaudio image
     docker build -t arch:pulseaudio ./pulseaudio
     ```
  7. Clone the repository with a script that allows you to run virtual displays:
     ```bash
     git clone https://github.com/FurmanovVitaliy/virtual-monotors-script
     cd virtual-monotors-script
     docker-compose build
     ```
  8. Ensure you create certificates using `mkcert` and build the containers.
  9.  Find `example_config.yaml` in `server/config/example_config.yaml`rename it to `local.yaml` and make sure that configaration in your config. Edit the `local.yaml` to ensure the paths to the scripts, certificates, directories, and container names are accurate.
  10. Project redy to start. 
        ```bash
        #To run serverside:
        cd self-hosted-cloud-gaming-service/server/cmd
        go run main.go
        #To run web interface:
        cd self-hosted-cloud-gaming-service/web
        npm run start
        ```
        
        
❗**For additional inquiries or issues related to building or starting the application, please refer to the Wiki section in "Technical Documents and Notes."**


## Technical documents and notes
- [Wiki](https://github.com/FurmanovVitaliy/self-hosted-cloud-gaming-service/wiki)

## Current State of the Project
- The foundation for scalability has been laid. 👌🏻
- Project in developing 🛠️
  
## Acknowledgements
A special thanks to my wife for independently creating the design for this project in Figma. This project was developed with her creative vision, and I deeply appreciate her hard work and dedication throughout the entire process.

- [Daria's Notion Portfolio](https://noiseless-giant-de5.notion.site/Daria-Furmanova-7eb639daae2a4c3a80779a1a3e47fad8)

## Plans
- [x] Add multi-controllers support  
- [ ] Add keyboard and mouce suppor
- [ ] Rewrite frontend with NEXT.js / React for SSR
- [ ] Implament PostgreSQL support instead MongoDB
- [ ] Implament QUIQ webtransport instaed of WebRTC
## From Autor 
This project has become a practical testing ground for me to master and experiment with new technologies and programming languages. Instead of writing yet another To-Do list, I chose a more ambitious idea—creating a cloud gaming platform.

The project was built from the ground up and assembled "from scratch." It has been brought to a technically functional state, but it will likely remain at this level of development. I lack the motivation to further develop this project, primarily because similar solutions already exist and offer significantly richer functionality. Some of them have already ceased to exist and have been forgotten.