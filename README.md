# build

```
go mod vendor
go build ./...
go install ./...
```

/bin/bash --rcfile ~/myenv/bin/activate

@echo off
setlocal enabledelayedexpansion

REM Create frames directory if it doesn't exist
if not exist frames mkdir frames

REM Initialize counter
set n=1

REM Loop through each JPG file in the utput directory
for %%f in (Output\*.jpg) do (
REM Format the counter with leading zeros
set padded=0000!n!
set padded=!padded:~-4!

    REM Determine if the image is landscape or portrait by checking dimensions
    for /f "tokens=2 delims=," %%a in ('ffmpeg -i "%%f" -vf "scale=1000:-1" -frames:v 1 -f null - 2^>^&1 ^| findstr "Video"') do set /a width=%%a
    for /f "tokens=3 delims=," %%a in ('ffmpeg -i "%%f" -vf "scale=1000:-1" -frames:v 1 -f null - 2^>^&1 ^| findstr "Video"') do set /a height=%%a

    if !width! GEQ !height! (
        REM Landscape: scale height to 1500 and crop to 1000 width
        ffmpeg -y -i "%%f" -vf "scale=-2:1500,crop=1000:1500" -q:v 2 "frames/frame_!padded!.jpg"
    ) else (
        REM Portrait: scale width to 1000 and auto-adjust height to fit
        ffmpeg -y -i "%%f" -vf "scale=1000:-2" -q:v 2 "frames/frame_!padded!.jpg"
    )
    
    REM Increment the counter
    set /a n+=1
)

REM Check if there are any files in the frames folder
if not exist frames\frame_0001.jpg (
echo No images found in frames directory. Check if the conversion was successful.
pause
exit /b
)

REM Create video from the resized and cropped images
(for %%i in (frames\*) do echo file '%%i') > filelist.txt

ffmpeg -f concat -safe 0 -i filelist.txt -framerate 50 -c:v libx264 -pix_fmt yuv420p output.mp4

ffmpeg -i output.mp4 -i audio.mp3 -c:v copy -c:a aac -map 0:v:0 -map 1:a:0 -shortest photoshoot.mp4

echo Video created as output.mp4
pause
