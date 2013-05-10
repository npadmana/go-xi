#ifndef GPUCOUNT_H_


void checkCUDAError(char *msg);

void allocDev(long nbytes, void **ptr);
void freeDev(void *ptr);

int sizef4();
int sizeull();

void copyToDevice(void* dst, const void* src, long nbytes);
void copyFromDevice(void* dst, const void* src, long nbytes);

int smu(void *p1, int start1, int end1, 
    void *p2, int start2, int end2, 
    float scale, 
    int Nr, int Nmu, float invdr, void *hist, 
    int nblocks, int dimx, int dimy);

#endif