cmake_minimum_required(VERSION 2.8)

find_package(CUDA)
include(FindCUDA)

CUDA_ADD_LIBRARY(gpucount SHARED smu.cu gpuutils.cu OPTIONS -O2 -arch=sm_20 -use_fast_math)
