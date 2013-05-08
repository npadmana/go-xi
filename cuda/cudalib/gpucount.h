#ifndef GPUCOUNT_H_


void checkCUDAError(char *msg);

void* allocDev(long nbytes);
void freeDev(void *ptr);

int sizeof_f4();
int sizeof_ull();

void copyToDevice(void* dst, const void* src, long nbytes);
void copyFromDevice(void* dst, const void* src, long nbytes);


#endif