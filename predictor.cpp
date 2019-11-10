#define _GLIBCXX_USE_CXX11_ABI 0

#include <algorithm>
#include <iosfwd>
#include <memory>
#include <string>
#include <utility>
#include <vector>
#include <iostream>
#include <iomanip>
#include <sys/time.h>

#include "SNPE/SNPE.hpp"
#include "SNPE/SNPEFactory.hpp"
#include "DlSystem/DlVersion.hpp"
#include "DlSystem/DlEnums.hpp"
#include "DlSystem/String.hpp"
#include "DlContainer/IDlContainer.hpp"
#include "DlSystem/ITensor.hpp"
#include "DlSystem/StringList.hpp"
#include "DlSystem/TensorMap.hpp"
#include "DlSystem/TensorShape.hpp"
#include "DlSystem/ITensorFactory.hpp"
#include "SNPE/SNPEBuilder.hpp"
#include "DlSystem/RuntimeList.hpp"
#include "DlSystem/UDLFunc.hpp"
#include "DlSystem/PlatformConfig.hpp"

#include "predictor.hpp"

#define LOG(x) std::cerr

using namespace snpe;
using std::string;

double get_us(struct timeval t) { return (t.tv_sec * 1000000 + t.tv_usec); }

/*
  Predictor class takes in model file (converted into .tflite from the original .pb file
  using tflite_convert CLI tool), batch size and device mode for inference
*/
class Predictor {
  public:
    Predictor(const string &model_file, int batch, int mode, bool verbose, bool profile);
    void Predict(int* inputData_quantize, float* inputData_float, bool quantize);

    std::unique_ptr<zdl::DlContainer::IDlContainer> net_;
    std::unique_ptr<zdl::SNPE::SNPE> snpe;
    int width_, height_, channels_;
    int batch_;
    int pred_len_ = 0;
    int mode_ = 0;
    TfLiteTensor* result_;
    float* result_float_;
    bool quantize_ = false;
    bool verbose_ = false; // display model details
    bool allow_fp16_ = false;
    bool profile_ = false; // operator level profiling
    bool read_outputs_ = true;
};

Predictor::Predictor(const string &model_file, int batch, int mode, bool verbose, bool profile) {
  char* model_file_char = const_cast<char*>(model_file.c_str());
  
  // set verbosity and profiling levels
  profile_ = profile;
  verbose_ = verbose;
 
  // build a runnable model from given model file
  struct timeval start_time, stop_time;
  gettimeofday(&start_time, nullptr);
  // read model file into a network 
  static zdl::DlSystem::Version_t Version = zdl::SNPE::SNPEFactory::getLibraryVersion();
  LOG(INFO) << "SNPE Version: " << Version.asString().c_str() << "\n";  
  net_ = zdl::DlContainer::IDlContainer::open(zdl::DlSystem::String(model_file_char));
  if(net_ == nullptr) {
    LOG(FATAL) << "Error while opening the container file" << "\n";
  }
  zdl::SNPE::SNPEBuilder snpeBuilder(net_.get());
  
  // set udlBundle
  zdl::DlSystem::UDLFactoryFunc udlFunc = UdlExample::MyUDLFactory;
  zdl::DlSystem::UDLBundle udlBundle;
  udlBundle.cookie = (void*)0xdeadbeaf, udlBundle.func = udlFunc;
  // set hardware backend
  zdl::DlSystem::Runtime_t runtime = zdl::DlSystem::Runtime_t::CPU;
  zdl::DlSystem::RuntimeList runtimeList;
  if(1 <= mode <= 8) {
    runtime = zdl::DlSystem::Runtime_t::CPU;
  } else if(mode == 9) {
    runtime = zdl::DlSystem::Runtime_t::GPU;
  } else if(mode == 10) {
    LOG(FATAL) << "Cannot run NNAPI through SNPE" << "\n";
  } else if (mode == 11) {
    runtime = zdl::DlSystem::Runtime_t::DSP;
  } else {
    LOG(FATAL) << "Invalid hardware mode" << "\n";
  }
  // check if chosen runtime is available on the device
  if(!zdl::SNPE::SNPEFactory::isRuntimeAvailable(runtime)) {
    LOG(INFO) << "Selected runtime not present. Falling back to CPU" << "\n";
    runtime = zdl::DlSystem::Runtime_t::CPU;
  }
  if(runtimeList.empty()) {
    runtimeList.add(runtime);
  }
  // set user supplied buffer as required
  // NOTE we do not allow user to set input output buffers
  // to make our life easier
  bool useUserSuppliedBuffers = false;
  zdl::DlSystem::PlatformConfig platformConfig;
  bool usingInitCaching = false;
  snpe = snpeBuilder.setOutputLayers({}).
      .setRuntimeProcessorOrder(runtimeList)
      .setUdlBundle(udlBundle)
      .setUseUserSuppliedBuffers(useUserSuppliedBuffers)
      .setPlatformConfig(platformConfig)
      .setInitCachedMode(usingInitCaching)
      .build();
  if(snpe == nullptr) {
    LOG(FATAL) << "Error while building SNPE object" << "\n";
  }
  gettimeofday(&stop_time, nullptr);
  // log model loading time
  if(verbose_) {
    LOG(INFO) << "Model loading (C++): " << (get_us(stop_time) - get_us(start_time))/1000 << "ms \n";
  }
  mode_ = mode;
  batch_ = batch;
  
}

void Predictor::Predict(int* inputData_quantize, float* inputData_float, bool quantize) {
  // check the batch size for the container
  zdl::DlSystem::TensorShape tensorShape;
  tensorShape = snpe->getInputDimensions()[0];
  size_t net_batchSize = tensorShape.getDimension()[0];
  if(verbose_) {
    LOG(INFO) << "Batch size for the container is " << net_batchSize << "\n";
  }
  
  std::string bufferType = ITENSOR;
  zdl::DlSystem::TensorMap outputTensorMap;

  std::unique_ptr<zdl::DlSystem::ITensor> input;
  const auto &strList_opt = snpe->getInputTensorNames();
  if(!strList_opt) {
    LOG(FATAL) << "Error obtaining Input tensor names" << "\n";
  }
  const auto &strList = *strList_opt;
  // make sure the network requires only a single input
  assert (strList.size() == 1);

  // create an input tensor that is correctly sized to hold the input of the network
  // Dimensions that have no fixed size will be represented with a value of 0
  const auto &inputDims_opt = snpe->getInputDimensions(strList.at(0));
  const auto &inputShape = *inputDims_opt;

  // calculate the total number of elements that can be stored in the tensor so that
  // we can check that the input contains the expected number of elememnts
  input = zdl::SNPE::SNPEFactory::getTensorFactory().createTensor(inputShape);

  // TODO padding the input vetcor so as to make the size of the vector to be equal
  // to an intgeret multiple of the  batch size
  // NOTE: for now we assume that the input model is going to have a batch size == 1

  // TODO set input dimensions
  width_ = 224;
  height_ = 224;
  channels_ = 3;
  if(width_ != inputShape[2]) {
    LOG(FATAL) << "width is not 224, need to resize" << "\n";
  }
  if(height != inputShape[1]) {
    LOG(FATAL) << "height is not 224, need to resize" << "\n";
  }
  if(channels_ != inputShape[3]) {
    LOG(FATAL) << "channel is not 3, can't do anything" << "\n";
  }

  // set quantization
  quantize_ = quantize;
  const int size = batch_ * width_ * height_ * channels_;
  // check if model bitwidth matches our expectation
  // input image = 224 X 224 X 3
  // resize it to what model expects if needed
  if(quantize_ == false) {
    LOG(INFO) << "Running float model" << "\n";
    float* base_pointer = &input[0];
    for(int i = 0; i < size; i++) {
      base_pointer[i] = inputData_float[i];
    }
  } else if (quantize_ == true) {
    LOG(INFO) << "Running 8-bit unsigned quantized model" << "\n";
    // TODO add quantization
  } else {
    LOG(FATAL) << "Unsupported input type: " << ", Quantize: " << quantize_ << "\n";
  }

  if(!input) {
    LOG(FATAL) << "could not read an empty tensor" << "\n";
  }

  bool execStatus = false;
  struct timeval start_time, stop_time;
  gettimeofday(&start_time, nullptr);  
  // run inference
  execStatus = snpe->execute(input.get(), outputTensorMap);
  if(execStatus == false) {
    LOG(FATAL) << "Failed to run inference" << "\n";
  }
  gettimeofday(&stop_time, nullptr);
  // log model inference
  if(verbose_) {
    LOG(INFO) << "Model computation (C++): " << (get_us(stop_time) - get_us(start_time))/1000 << "ms \n"; 
  }

  // store output
  // TODO fetch output size from output tensor map
  result_float_ = new float[output_size];
  zdl::DlSystem::StringList tensorNames = outputTensorMap.getTensorNames();
  for(auto& name : tensorNames) {
    auto tensorPtr = outputTensorMap.getTensor(name);
    for(auto it = tensorPtr->cbegin(); it != tensorPtr->cend(); it++) {
      result_float_[i] = *it;
    }
  }

  pred_len_ = output_size;
}

PredictorContext NewSnpe(char *model_file, int batch, int mode, bool verbose, bool profile) {
  try {
    const auto ctx = new Predictor(model_file, batch, mode, verbose, profile);
    return (void *) ctx;
  } catch(const std::invalid_argument &ex) {
    errno = EINVAL;
    return nullptr;
  }
}

void InitSnpe() {}

void PredictSnpe(PredictorContext pred, int* inputData_quantize, float* inputData_float, bool quantize) {
  auto predictor = (Predictor *)pred;
  if (predictor == nullptr) {
    return;
  }
  predictor->Predict(inputData_quantize, inputData_float, quantize);
  return;
}

float* GetPredictionsSnpe(PredictorContext pred) {
  auto predictor = (Predictor *)pred;
  if (predictor == nullptr) {
    return nullptr;
  }
  return predictor->result_float_;
}

void DeleteSnpe(PredictorContext pred) {
  auto predictor = (Predictor *)pred;
  if (predictor == nullptr) {
    return;
  }
  delete predictor;
}

int GetWidthSnpe(PredictorContext pred) {
  auto predictor = (Predictor *)pred;
  if (predictor == nullptr) {
    return 0;
  }
  return predictor->width_;
}

int GetHeightSnpe(PredictorContext pred) {
  auto predictor = (Predictor *)pred;
  if (predictor == nullptr) {
  return 0;
  }
  return predictor->height_;
}

int GetChannelsSnpe(PredictorContext pred) {
  auto predictor = (Predictor *)pred;
  if (predictor == nullptr) {
    return 0;
  }
  return predictor->channels_;
}

int GetPredLenSnpe(PredictorContext pred) {
  auto predictor = (Predictor *)pred;
  if (predictor == nullptr) {
    return 0;
  }
  return predictor->pred_len_;
}

