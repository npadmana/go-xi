#ifndef GPUCOUNT_H_


void checkCUDAError(char *msg);

void allocDev(long nbytes, void **ptr);
void freeDev(void *ptr);

int sizef4();
int sizeull();

void copyToDevice(void* dst, const void* src, long nbytes);
void copyFromDevice(void* dst, const void* src, long nbytes);


#endif