ffmpeg -r 30 -f x11grab -draw_mouse 0 -s 360x480 -i :0 -pix_fmt yuv420p -c:v libvpx -deadline realtime -quality realtime -f rtp rtp://127.0.0.1:5004
