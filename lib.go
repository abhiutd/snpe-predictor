package tflite

// #cgo CXXFLAGS: -std=c++11 -I${SRCDIR}/cbits -O3 -Wall -g -Wno-sign-compare -Wno-unused-function  -I/home/as29/my_tflite/tensorflow/bazel-tensorflow/external/flatbuffers/include -I/home/as29/my_tflite/tensorflow/bazel-tensorflow/external/com_google_absl -I/home/as29/my_gles -I/home/as29/my_tflite/tensorflow
// #cgo LDFLAGS: -lstdc++ -L/home/as29/my_android_ndk/android-ndk-r19c -llog -L/home/as29/my_android_ndk/android-ndk-r19c/platforms/android-28/arch-arm64/usr/lib -lEGL -lGLESv3 -L/opt/tflite/lib -ltensorflowlite
import "C"
