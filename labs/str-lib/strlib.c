#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int mystrlen(char *str){
    int size = 0;

    while(*str != '\0'){
        *str++;
        size++;
    }

    return size;
}

char *mystradd(char *origin, char *addition){
    
    int originS = mystrlen(origin);
    int additionS = mystrlen(addition);
    char *outcome = malloc(originS + additionS);

    for(int i = 0; i < originS; i++){
        outcome[i] = origin[i];
    }

    int aux = 0;
    for(int j = originS; j < (originS + additionS); j++){
        outcome[j] = addition[aux];
        aux++;
    }

    return outcome;
}

int mystrfind(char *origin, char *substr){
    
    int originS = mystrlen(origin);
    int subS = mystrlen(substr);

     for(int i = 0; i < originS; i++){
        if(origin[i] == substr[0]){
            for(int j = 0; j < subS; j++){
                if(origin[i+j] != substr[j]){
                    // Keep searching 
                    break;
                } else if(j == subS-1){
                    //We only reach this condition after validating the entire string
                    return i;
                }
            }
        }
    }
    return -1;
}
