#include <stdio.h>
#include <stdlib.h>


int main(int argc, char **argv)
{
    int fahr;

    if (argc == 4) {
        int totalSteps = (atoi(argv[2]) - atoi(argv[1])) / atoi(argv[3]);
        
        fahr = atoi(argv[1]);
        printf("Fahrenheit: %3d, Celcius: %6.1f\n", fahr, (5.0/9.0)*(fahr-32));

        for(int i = 0; i < totalSteps; i++){
            fahr = fahr + atoi(argv[3]);
            printf("Fahrenheit: %3d, Celcius: %6.1f\n", fahr, (5.0/9.0)*(fahr-32));
        }
    } else if (argc == 2){
        fahr = atoi(argv[1]);
	    printf("Fahrenheit: %3d, Celcius: %6.1f\n", fahr, (5.0/9.0)*(fahr-32));
    } else {
        printf("Invalid Input\n");
    } 

    return 0;
}

//Bibliography
//https://cboard.cprogramming.com/c-programming/134232-converting-command-line-char-*argv[]-int.html 