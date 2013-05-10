#include <stdio.h>
// CUDA-C includes
#include <cuda.h>
#include <cuda_runtime.h>

extern "C" {

const int BUFHIST=1024;


__global__ void shared_smu_kernel
(float4 *x1, int start1, int end1, 
    float4 *x2, int start2, int end2, 
    float scale, 
    int Nr, int Nmu, float invdr, unsigned long long *hist) {

    // Keep a shared copy of the histogram
    __shared__ long long _hist[BUFHIST];

    // Variable declarations
    int stride1, stride2, nh1, nr1, rstart, rend, ih;
    int ii, jj, ir, imu;
    float4 _x1, _x2;
    float s2, l2, sl, s1, l1;

    // Strides -- we will distribute over both x1 and x2
    stride1 = blockDim.y * gridDim.y;
    stride2 = blockDim.x * gridDim.x;


    // Compute the number of histograms we need to do
    nr1 = BUFHIST/Nmu;
    nh1 = (Nr + nr1 - 1)/nr1;


    // Do each piece of the histogram separately
    for (ih = 0; ih < nh1; ++ih) {
        rstart = ih*nr1;
        rend = rstart + nr1;

        // zero histogram
        // For simplicity, only a few threads will participate
        if (threadIdx.y == 0) {
            ii = threadIdx.x;
            while (ii < BUFHIST) {
                _hist[ii] = 0ll;
                ii += blockDim.x;
            }
        }
        __syncthreads();


        // Start loop over first set of data
        ii = threadIdx.y + blockIdx.y * blockDim.y + start1;
        while (ii < end1) {
            _x1 = x1[ii];
            jj = threadIdx.x + blockIdx.x * blockDim.x + start2;
            while (jj < end2) {
                _x2 = x2[jj];

                // X
                s1 = _x1.x - _x2.x;
                l1 = 0.5*(_x1.x + _x2.x);
                s2 = s1*s1;
                l2 = l1*l1;
                sl = s1*l1;

                // Y
                s1 = _x1.y - _x2.y;
                l1 = 0.5*(_x1.y + _x2.y);
                s2 += s1*s1;
                l2 += l1*l1;
                sl += s1*l1;

                // Z
                s1 = _x1.z - _x2.z;
                l1 = 0.5*(_x1.z + _x2.z);
                s2 += s1*s1;
                l2 += l1*l1;
                sl += s1*l1;

                // Compute s1, s2
                s1 = sqrtf(s2);
                l1 = rsqrtf(s2*l2 + 1.e-15);
                l1 = sl * l1;  // This is now mu, but save a register

                // Work out indices
                if (l1 < 0) {
                    l1 = -l1;
                }
                ir = s1 * invdr;
                imu = l1 * Nmu;
                if ((ir >= rstart) && (ir < rend)) {
                    atomicAdd( (unsigned long long*) &_hist[(ir-rstart)*Nmu + imu], _x1.w*_x2.w*scale);
                }

                // Loop over 2 ends
                jj += stride2;    
            }

            // Loop over 1 ends
            ii += stride1;
        }

        // Synchronize
        __syncthreads();

        // Copy histogram 
        // For simplicity, only a few threads will participate
        if (threadIdx.y == 0) {
            ir = Nmu*rstart;
            ii = threadIdx.x + ir;
            jj = Nmu*rend;
            while (ii < jj) {
                atomicAdd( (unsigned long long*) &hist[ii], _hist[ii-ir]);
                ii += blockDim.x;
            }
        }
        __syncthreads();

    
    // End histogram loop    
    }

}


int smu(void *p1, int start1, int end1, 
    void *p2, int start2, int end2, 
    float scale, 
    int Nr, int Nmu, float invdr, void *hist, 
    int nblocks, int dimx, int dimy) {

    dim3 Nb(nblocks);
    dim3 Nt(dimx, dimy);

    if (Nmu > BUFHIST) {
        return -1;
    }

    shared_smu_kernel<<<Nb, Nt>>>( (float4*) p1, start1, end1, 
    (float4*) p2, start2, end2, scale, 
    Nr, Nmu, invdr, (unsigned long long *)hist);

    return 0;
}




}