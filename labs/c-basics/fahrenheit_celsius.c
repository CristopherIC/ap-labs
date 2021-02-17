#include <stdio.h>

int main(int argc, int **argv)
{
    if(argc > 2){
        printf("Just 1 value\n");
    } else {
        int fahr = argv[0];
        printf("%d\n", argv[0]);
	    printf("Fahrenheit: %3d, Celcius: %6.1f\n", fahr, (5.0/9.0)*(fahr-32));
    }
    return 0;
}
