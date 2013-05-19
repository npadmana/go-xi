#include <stdio.h>
// CUDA-C includes
#include <cuda.h>
#include <cuda_runtime.h>

extern "C" {


void checkCUDAError(char *msg)
{
    cudaError_t err = cudaGetLastError();
    if( cudaSuccess != err) 
    {
        fprintf(stderr, "Cuda error: %s: %s.\n", msg, 
                             cudaGetErrorString( err) );
        exit(EXIT_FAILURE);
    }                         
}

void allocDev(long nbytes, void **ptr) {
	cudaMalloc(ptr, nbytes);
}

void freeDev(void *ptr) {
	cudaFree(ptr);
}

int sizef4() {
	return int(sizeof(float4));
}

int sizeull() {
	return int(sizeof(unsigned long long));
}

void copyToDevice(void* dst, const void* src, long nbytes) {
	cudaMemcpy(dst, src, nbytes,  cudaMemcpyHostToDevice);
}

void copyFromDevice(void* dst, const void* src, long nbytes) {
	cudaMemcpy(dst, src, nbytes,  cudaMemcpyDeviceToHost);
}

void deviceSync() {
    cudaDeviceSynchronize();
}


}