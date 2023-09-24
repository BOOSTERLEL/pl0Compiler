#include <stdio.h>
#include <stdlib.h>

int pl_0_builtin_println(int x){
    return printf("%d\n",x);
}

int pl_0_builtin_write(){
    int x;
	scanf_s("%d",&x);
    return x;
}

int pl_0_builtin_exit(int x){
    exit(x);
    return 0;
}