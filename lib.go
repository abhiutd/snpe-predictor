package snpe

// #cgo CXXFLAGS: -std=c++11 -I${SRCDIR}/cbits -O3 -Wall -g -Wno-sign-compare -Wno-unused-function  -I/home/as29/my_snpe/snpe-1.32.0.555/include/zdl -I/home/as29/my_gles
// #cgo LDFLAGS: -lstdc++ -L/home/as29/my_android_ndk/android-ndk-r19c -L/home/as29/my_android_ndk/android-ndk-r19c/platforms/android-28/arch-arm64/usr/lib -llog -lEGL -lGLESv3 -L/home/as29/my_snpe/snpe-1.32.0.555/lib/aarch64-android-clang6.0 -lc++_shared -lPlatformValidatorShared -lPSNPE -lSNPE_G -lSNPE -lsymphony-cpu
import "C"
