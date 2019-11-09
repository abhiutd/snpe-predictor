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
#include "SNPE/SNPEBuilder.hpp"

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
  // TODO read model file into a network 
  static zdl::DlSystem::Version_t Version = zdl::SNPE::SNPEFactory::getLibraryVersion();
  LOG(INFO) << "SNPE Version: " << Version.asString().c_str() << "\n";  
  net_ = zdl::DlContainer::IDlContainer::open(zdl::DlSystem::String(model_file_char));
  zdl::SNPE::SNPEBuilder snpeBuilder(net_.get());
  
  gettimeofday(&stop_time, nullptr);
  // log model loading time
  if(verbose_) {
    LOG(INFO) << "Model loading (C++): " << (get_us(stop_time) - get_us(start_time))/1000 << "ms \n";
  }
  mode_ = mode;
  batch_ = batch;
  
}

void Predictor::Predict(int* inputData_quantize, float* inputData_float, bool quantize) {
  int input = interpreter->inputs()[0];
  if(verbose_)
    LOG(INFO) << "input: " << input << "\n";
  const std::vector<int> inputs = interpreter->inputs();
  const std::vector<int> outputs = interpreter->outputs();
  if(verbose_) {
    LOG(INFO) << "number of inputs: " << inputs.size() << "\n";
    LOG(INFO) << "number of outputs: " << outputs.size() << "\n";
  }

  // set appropriate hardware backend
  switch(mode_) {
    case 7: {
      const TfLiteGpuDelegateOptions options = {
        .metadata = NULL,
        .compile_options = {
          .precision_loss_allowed = 1, // FP16
          .preferred_gl_object_type = TFLITE_GL_OBJECT_TYPE_FASTEST,
          .dynamic_batch_enabled = 0, // Not fully functional yet
        },
      };
      auto* delegate = TfLiteGpuDelegateCreate(&options);
      if(!delegate) {
        LOG(FATAL) << "Unable to create GPU delegate" << "\n";
      }
      if(interpreter->ModifyGraphWithDelegate(delegate) != kTfLiteOk) {
         LOG(FATAL) << "Failed to apply " << "GPU delegate" << "\n";
      } else {
         LOG(INFO) << "Applied " << "GPU delegate" << "\n";
      }
      break; }
    case 8: {
      auto delegate = tflite::evaluation::CreateNNAPIDelegate();
      if(!delegate) {
        LOG(INFO) << "NNAPI acceleration is unsupported on this platform" << "\n";
      }
      interpreter->UseNNAPI(true);
      break; }
    case 1: {
      interpreter->SetNumThreads(1); 
      break; }
    case 2: {
      interpreter->SetNumThreads(2); 
      break; }
    case 3: {
      interpreter->SetNumThreads(3); 
      break; }
    case 4: {
      interpreter->SetNumThreads(4); 
      break; }
    case 5: {
      interpreter->SetNumThreads(5); 
      break; }
    case 6: {
      interpreter->SetNumThreads(6);
      break; }
    default: {
      interpreter->SetNumThreads(4); }
  }
  
  if(interpreter->AllocateTensors() != kTfLiteOk) {
    LOG(FATAL) << "Failed to allocate tensors!";
  }

  // fill input buffers
  TfLiteTensor* input_tensor = interpreter->tensor(input);
  TfLiteIntArray* input_dims = input_tensor->dims;
  height_ = input_dims->data[1];
  width_ = input_dims->data[2];
  channels_ = input_dims->data[3];
  if(verbose_) {
    LOG(INFO) << "Model input height is " << height_ << "\n";
    LOG(INFO) << "Model input width is " << width_ << "\n";
    LOG(INFO) << "Model input channel is " << channels_ << "\n";
  }

  // set quantization
  quantize_ = quantize;
  const int size = batch_ * width_ * height_ * channels_;
  // check if model bitwidth matches our expectation
  // input image = 224 X 224 X 3
  // resize it to what model expects if needed
  if(interpreter->tensor(input)->type == kTfLiteFloat32 && quantize_ == false) {
    LOG(INFO) << "Running float model" << "\n";
    if(height_ != 224 || width_ != 224 | channels_ != 3) {
      SetInputTflite_float(interpreter->typed_tensor<float>(input), inputData_float, 224, 224, 3, height_, width_, channels_);
    } else {
      memcpy(interpreter->typed_tensor<float>(input), &inputData_float[0], size);
    }
  } else if (interpreter->tensor(input)->type == kTfLiteUInt8 && quantize_ == true) {
    LOG(INFO) << "Running 8-bit unsigned quantized model" << "\n";
    if(height_ != 224 || width_ != 224 || channels_ != 3) {
      SetInputTflite_quantize_8_unsigned(interpreter->typed_tensor<uint8_t>(input), inputData_quantize, 224, 224, 3, height_, width_, channels_);
    } else {
      uint8_t* base_pointer = interpreter->typed_tensor<uint8_t>(input);
      for(int i = 0; i < size; i++) {
        base_pointer[i] = (uint8_t)inputData_quantize[i];
      }
    }
  } else if(interpreter->tensor(input)->type == kTfLiteInt8 && quantize_ == true) {
    LOG(INFO) << "Running 8-bit signed quantized model" << "\n";
    if(height_ != 224 || width_ != 224 || channels_ != 3) {
      SetInputTflite_quantize_8_signed(interpreter->typed_tensor<int8_t>(input), inputData_quantize, 224, 224, 3, height_, width_, channels_);
    } else {
      int8_t* base_pointer = interpreter->typed_tensor<int8_t>(input);
      for(int i = 0; i < size; i++) {
        base_pointer[i] = (int8_t)inputData_quantize[i];
      }
    }
  } else {
    LOG(FATAL) << "Unsupported input type: " << interpreter->tensor(input)->type << ", Quantize: " << quantize_ << "\n";
  }

  // TODO interpreter profiler not fetching any information
  auto profiler = absl::make_unique<profiling::Profiler>(1024);
  interpreter->SetProfiler(profiler.get());
  if(profile_ == true) {
    LOG(INFO) << "Starting profiler" << "\n";
    profiler->StartProfiling();
  }

  struct timeval start_time, stop_time;
  gettimeofday(&start_time, nullptr);  
  // run inference
  if(interpreter->Invoke() != kTfLiteOk) {
    LOG(FATAL) << "Failed to invoke tflite" << "\n";
  }
  gettimeofday(&stop_time, nullptr);
  // log model inference
  if(verbose_) {
    LOG(INFO) << "Model computation (C++): " << (get_us(stop_time) - get_us(start_time))/1000 << "ms \n"; 
  }

  if(profile_ == true) {
    LOG(INFO) << "Stopping profiler" << "\n";
    profiler->StopProfiling();
    auto profile_events = profiler->GetProfileEvents();
    for(int i = 0; i < profile_events.size(); i++) {
      LOG(INFO) << "Inside profiler loop" << "\n";
      auto op_index = profile_events[i]->event_metadata;
      const auto node_and_registration = interpreter->node_and_registration(op_index);
      const TfLiteRegistration registration = node_and_registration->second;
      LOG(INFO) << std::fixed << std::setw(10) << std::setprecision(3)
                << (profile_events[i]->end_timestamp_us - profile_events[i]->begin_timestamp_us) / 1000.0
                << ", Node" << std::setw(3) << std::setprecision(3) << op_index
                << ", OpCode" << std::setw(3) << std::setprecision(3)
                << registration.builtin_code << ", "
                << EnumNameBuiltinOperator(static_cast<BuiltinOperator>(registration.builtin_code))
                << "\n";
    }
    LOG(INFO) << "Displayed layer wise profiling information" << "\n" ;
  }

  // read and store model predictions
  int output = interpreter->outputs()[0];
  TfLiteIntArray* output_dims = interpreter->tensor(output)->dims;
  auto output_size = output_dims->data[output_dims->size-1];
  pred_len_ = output_size;
  
  result_float_ = new float[output_size];
  if(interpreter->tensor(output)->type == kTfLiteFloat32) {
    float* prediction = interpreter->typed_output_tensor<float>(0);
    for(int i = 0; i < output_size; i++)
      result_float_[i] = prediction[i]; 
  }	else if(interpreter->tensor(output)->type == kTfLiteUInt8) {
    uint8_t* prediction = interpreter->typed_output_tensor<uint8_t>(0);
    for(int i = 0; i < output_size; i++)
      result_float_[i] = prediction[i] / 255.0; 
  } else if(interpreter->tensor(output)->type == kTfLiteInt8) {
    int8_t* prediction = interpreter->typed_output_tensor<int8_t>(0);
    for(int i = 0; i < output_size; i++)
      result_float_[i] = prediction[i] / 255.0;
  } else {
    LOG(FATAL) << "Unsupported output type: " << interpreter->tensor(output)->type << "\n";
  }
}

PredictorContext NewTflite(char *model_file, int batch, int mode, bool verbose, bool profile) {
  try {
    const auto ctx = new Predictor(model_file, batch, mode, verbose, profile);
    return (void *) ctx;
  } catch(const std::invalid_argument &ex) {
    errno = EINVAL;
    return nullptr;
  }
}

void InitTflite() {}

void PredictTflite(PredictorContext pred, int* inputData_quantize, float* inputData_float, bool quantize) {
  auto predictor = (Predictor *)pred;
  if (predictor == nullptr) {
    return;
  }
  predictor->Predict(inputData_quantize, inputData_float, quantize);
  return;
}

float* GetPredictionsTflite(PredictorContext pred) {
  auto predictor = (Predictor *)pred;
  if (predictor == nullptr) {
    return nullptr;
  }
  return predictor->result_float_;
}

void DeleteTflite(PredictorContext pred) {
  auto predictor = (Predictor *)pred;
  if (predictor == nullptr) {
    return;
  }
  delete predictor;
}

int GetWidthTflite(PredictorContext pred) {
  auto predictor = (Predictor *)pred;
  if (predictor == nullptr) {
    return 0;
  }
  return predictor->width_;
}

int GetHeightTflite(PredictorContext pred) {
  auto predictor = (Predictor *)pred;
  if (predictor == nullptr) {
  return 0;
  }
  return predictor->height_;
}

int GetChannelsTflite(PredictorContext pred) {
  auto predictor = (Predictor *)pred;
  if (predictor == nullptr) {
    return 0;
  }
  return predictor->channels_;
}

int GetPredLenTflite(PredictorContext pred) {
  auto predictor = (Predictor *)pred;
  if (predictor == nullptr) {
    return 0;
  }
  return predictor->pred_len_;
}

void SetInputTflite_float(float* out, float* in, int image_height, int image_width, int image_channels, int model_height, int model_width, int model_channels) {
  
  int number_of_pixels = image_height * image_width * image_channels;
  
  // create  a new interpreter to resize input image into model's desired dimensions
  std::unique_ptr<Interpreter> interpreter(new Interpreter);
  int base_index = 0;
  // two inputs: input and new_sizes
  interpreter->AddTensors(2, &base_index);
  // one output
  interpreter->AddTensors(1, &base_index);
  // set input and output tensors
  interpreter->SetInputs({0, 1});
  interpreter->SetOutputs({2});

  // set parameters of tensors
  TfLiteQuantizationParams quant;
  interpreter->SetTensorParametersReadWrite(
      0, kTfLiteFloat32, "input",
      {1, image_height, image_width, image_channels}, quant);
  interpreter->SetTensorParametersReadWrite(1, kTfLiteInt32, "new_size", {2},
                                            quant);
  interpreter->SetTensorParametersReadWrite(
      2, kTfLiteFloat32, "output",
      {1, model_height, model_width, model_channels}, quant);

  ops::builtin::BuiltinOpResolver resolver;
  const TfLiteRegistration* resize_op =
      resolver.FindOp(BuiltinOperator_RESIZE_BILINEAR, 1);
  auto* params = reinterpret_cast<TfLiteResizeBilinearParams*>(
      malloc(sizeof(TfLiteResizeBilinearParams)));
  params->align_corners = false;
  interpreter->AddNodeWithParameters({0, 1}, {2}, nullptr, 0, params, resize_op,
                                     nullptr);

  interpreter->AllocateTensors();

  // fill input image
  auto input = interpreter->typed_tensor<float>(0);
  for (int i = 0; i < number_of_pixels; i++) {
    input[i] = in[i];
  }

  // fill new_sizes
  interpreter->typed_tensor<int>(1)[0] = model_height;
  interpreter->typed_tensor<int>(1)[1] = model_width;

  interpreter->Invoke();

  auto output = interpreter->typed_tensor<float>(2);
  auto output_number_of_pixels = model_height * model_width * model_channels;

  for (int i = 0; i < output_number_of_pixels; i++) {
      out[i] = output[i];
  }
}

void SetInputTflite_quantize_8_unsigned(uint8_t* out, int* in, int image_height, int image_width, int image_channels, int model_height, int model_width, int model_channels) {
  
  int number_of_pixels = image_height * image_width * image_channels;
  
  // create  a new interpreter to resize input image into model's desired dimensions
  std::unique_ptr<Interpreter> interpreter(new Interpreter);
  int base_index = 0;
  // two inputs: input and new_sizes
  interpreter->AddTensors(2, &base_index);
  // one output
  interpreter->AddTensors(1, &base_index);
  // set input and output tensors
  interpreter->SetInputs({0, 1});
  interpreter->SetOutputs({2});

  // set parameters of tensors
  TfLiteQuantizationParams quant;
  interpreter->SetTensorParametersReadWrite(
      0, kTfLiteFloat32, "input",
      {1, image_height, image_width, image_channels}, quant);
  interpreter->SetTensorParametersReadWrite(1, kTfLiteInt32, "new_size", {2},
                                            quant);
  interpreter->SetTensorParametersReadWrite(
      2, kTfLiteFloat32, "output",
      {1, model_height, model_width, model_channels}, quant);

  ops::builtin::BuiltinOpResolver resolver;
  const TfLiteRegistration* resize_op =
      resolver.FindOp(BuiltinOperator_RESIZE_BILINEAR, 1);
  auto* params = reinterpret_cast<TfLiteResizeBilinearParams*>(
      malloc(sizeof(TfLiteResizeBilinearParams)));
  params->align_corners = false;
  interpreter->AddNodeWithParameters({0, 1}, {2}, nullptr, 0, params, resize_op,
                                     nullptr);

  interpreter->AllocateTensors();

  // fill input image
  auto input = interpreter->typed_tensor<float>(0);
  for (int i = 0; i < number_of_pixels; i++) {
    input[i] = in[i];
  }

  // fill new_sizes
  interpreter->typed_tensor<int>(1)[0] = model_height;
  interpreter->typed_tensor<int>(1)[1] = model_width;

  interpreter->Invoke();

  auto output = interpreter->typed_tensor<float>(2);
  auto output_number_of_pixels = model_height * model_width * model_channels;

  for (int i = 0; i < output_number_of_pixels; i++) {
      out[i] = (uint8_t)output[i];
  }
}


void SetInputTflite_quantize_8_signed(int8_t* out, int* in, int image_height, int image_width, int image_channels, int model_height, int model_width, int model_channels) {
  
  int number_of_pixels = image_height * image_width * image_channels;
  
  // create  a new interpreter to resize input image into model's desired dimensions
  std::unique_ptr<Interpreter> interpreter(new Interpreter);
  int base_index = 0;
  // two inputs: input and new_sizes
  interpreter->AddTensors(2, &base_index);
  // one output
  interpreter->AddTensors(1, &base_index);
  // set input and output tensors
  interpreter->SetInputs({0, 1});
  interpreter->SetOutputs({2});

  // set parameters of tensors
  TfLiteQuantizationParams quant;
  interpreter->SetTensorParametersReadWrite(
      0, kTfLiteFloat32, "input",
      {1, image_height, image_width, image_channels}, quant);
  interpreter->SetTensorParametersReadWrite(1, kTfLiteInt32, "new_size", {2},
                                            quant);
  interpreter->SetTensorParametersReadWrite(
      2, kTfLiteFloat32, "output",
      {1, model_height, model_width, model_channels}, quant);

  ops::builtin::BuiltinOpResolver resolver;
  const TfLiteRegistration* resize_op =
      resolver.FindOp(BuiltinOperator_RESIZE_BILINEAR, 1);
  auto* params = reinterpret_cast<TfLiteResizeBilinearParams*>(
      malloc(sizeof(TfLiteResizeBilinearParams)));
  params->align_corners = false;
  interpreter->AddNodeWithParameters({0, 1}, {2}, nullptr, 0, params, resize_op,
                                     nullptr);

  interpreter->AllocateTensors();

  // fill input image
  auto input = interpreter->typed_tensor<float>(0);
  for (int i = 0; i < number_of_pixels; i++) {
    input[i] = in[i];
  }

  // fill new_sizes
  interpreter->typed_tensor<int>(1)[0] = model_height;
  interpreter->typed_tensor<int>(1)[1] = model_width;

  interpreter->Invoke();

  auto output = interpreter->typed_tensor<float>(2);
  auto output_number_of_pixels = model_height * model_width * model_channels;

  for (int i = 0; i < output_number_of_pixels; i++) {
      out[i] = (int8_t)output[i];
  }
}

