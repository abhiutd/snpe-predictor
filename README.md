# go-snpe

[![Go Report Card](https://goreportcard.com/badge/github.com/rai-project/go-mxnet)](https://goreportcard.com/report/github.com/rai-project/go-mxnet)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Go binding for Qualcomm Neural Processing Engine (SNPE) C++ API. It is also referred to as MLModelScope SNPE mobile Predictor (SNPE mPredictor). It is used to perform model inference on mobile devices. It is used by the [Qualcomm SNPE agent](https://github.com/abhiutd/snpe-agent) in [MLModelScope](mlmodelscope.org) to perform model inference in Go. More importantly, it can be used as a standalone predictor in any given Android application. Refer to [Usage Modes](Usage Modes) for further details.

## Installation

Download and install go-mxnet:

```
go get -v github.com/abhiutd/snpe-predictor
```

The binding requires Qualcomm SNPE, Gomobile and other Go packages.

### Qualcomm SNPE C++ Library

The Qualcomm SNPE C++ library is expected to be under `/opt/snpe`.

Download pre-built binaries as instructed in [SNPE documentation](https://developer.qualcomm.com/docs/snpe/setup.html).

If you get an error about not being able to write to `/opt` then perform the following

```
sudo mkdir -p /opt/snpe
sudo chown -R `whoami` /opt/snpe
```

If you are using custom path for build files, change CGO_CFLAGS, CGO_CXXFLAGS and CGO_LDFLAGS enviroment variables. Refer to [Using cgo with the go command](https://golang.org/cmd/cgo/#hdr-Using_cgo_with_the_go_command).

For example,

```
    export CGO_CFLAGS="${CGO_CFLAGS} -I/tmp/snpe/include"
    export CGO_CXXFLAGS="${CGO_CXXFLAGS} -I/tmp/snpe/include"
    export CGO_LDFLAGS="${CGO_LDFLAGS} -L/tmp/snpe/lib"
```

### Go Packages

You can install the dependency through `go get`.

```
cd $GOPATH/src/github.com/abhiutd/snpe-predictor
go get -u -v ./...
```

Or use [Dep](https://github.com/golang/dep).

```
dep ensure -v
```

This installs the dependency in `vendor/`. It is the preferred option.

Also, one needs to install `gomobile` to be able to generate Java/Objective-C bindings of the mPredictor. 

```
go get golang.org/x/mobile/cmd/gomobile
gomobile init
```

### Configure Environmental Variables

Configure the linker environmental variables since the Qualcomm SNPE C++ library is under a non-system directory. Place the following in either your `~/.bashrc` or `~/.zshrc` file

Linux
```
export LIBRARY_PATH=$LIBRARY_PATH:/opt/snpe/lib
export LD_LIBRARY_PATH=/opt/snpe/lib:$DYLD_LIBRARY_PATH

```

macOS
```
export LIBRARY_PATH=$LIBRARY_PATH:/opt/snpe/lib
export DYLD_LIBRARY_PATH=/opt/snpe/lib:$DYLD_LIBRARY_PATH
```

### Generate bindings

SNPE mPredictor is written in Go, binded with SNPE C++ API. To be able to use it in a mobile application, you would have to generate appropriate bindings (Java for Android). We provide bindings off-the-shelf in [bindings](bindings), but you can generate your own by using the following command.

```
gomobile bind -o bindings/android/snpe-predictor.aar -target=android/arm64 -v github.com/abhiutd/snpe-predictor
```

This command builds `snpe-predictor.aar` binary for Android with `arm64` ISA.

### Usage Modes

One can employ SNPE mPredictor to perform model inference in multiple ways, which are listed below.

1. Standalone Predictor (mPredictor)

There are four main API calls to be used for performing model inference in a given mobile application.

```
// create a SNPE mPredictor
New()

// perform inference on given input data
Predict()

// generate output predictions
ReadPredictedOutputFeatures()

// delete the SNPE mPredictor
Close()
```

Refer to [cbits.go](cbits.go) for details on the inputs/outputs of each API call.

2.  MLModelScope Mobile Agent

Download MLModelScope mobile agent from [agent](https://github.com/abhiutd/agent-classification-android). It has Tensorflow Lite and Qualcomm SNPE mPredictors in built. Refer to its documentation to understand its usage.

3. MLModelScope web UI

Choose Qualcomm SNPE as framework and one of the available mobile devices as hardware backend to perform model inference through web interface.
