#include <math.h>
#include <stdio.h>

typedef struct {
	double x, y, z, w;
} Particle; 

void smu(void *v1, void *v2, int n1, int n2, double *data, int nmu, double maxs2, double invdmu, double invds, double scale) {
	Particle *p1, *p2;
	double x1, y1, z1, w1, s2, l2, sl, s1, l1, mu;
	int ip1, ip2, imu, is;

	p1 = (Particle*) v1;
	p2 = (Particle*) v2;

	for (ip1=0; ip1 < n1; ++ip1) {
		x1 = p1[ip1].x; y1 = p1[ip1].y; z1 = p1[ip1].z; 
		w1 = p1[ip1].w * scale;

		for (ip2=0;ip2 < n2; ++ip2) {

			// x
			s1 = x1 - p2[ip2].x;
			l1 = 0.5*(x1 + p2[ip2].x);
			s2 = s1*s1;
			l2 = l1*l1;
			sl = s1*l1;

			// y
			s1 = y1 - p2[ip2].y;
			l1 = 0.5*(y1 + p2[ip2].y);
			s2 += s1*s1;
			l2 += l1*l1;
			sl += s1*l1;

			// z 
			s1 = z1 - p2[ip2].z;
			l1 = 0.5*(z1 + p2[ip2].z);
			s2 += s1*s1;
			l2 += l1*l1;
			sl += s1*l1;

			if (s2 >= maxs2) continue;
			
			s1 = sqrt(s2);
			l1 = 1./sqrt(s2*l2+1.e-15);
			mu = sl*l1;
			if (mu < 0) mu = -mu;

			imu = (int)(mu*invdmu);
			is = (int)(s1*invds);
			data[is*nmu + imu] += w1*p2[ip2].w;
	
		}


	}


}