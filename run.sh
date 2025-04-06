LOG_LEVEL=Debug align run \
  --base 2025_ale_0002_small.jpg \
  --input data/input \
  --output data/output \
  --video output.mp4 \
  --width 1000 \
  --height 1500 \
  --model models/shape_predictor_68_face_landmarks.dat \
  --script scripts/landmark_extractor.py \
  --fps 50