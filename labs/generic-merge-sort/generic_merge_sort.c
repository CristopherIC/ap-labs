#include <stdio.h>
#include <string.h>
#include <stdint.h>
#include <stdlib.h>

void mergeSort(char* [],int , int ,int (*comparator)(char *, char *));
void merge(char* [],int , int ,int (*comparator)(char *, char *));
int intComparator(char *, char *);
int strComparator(char *, char *);

int main(int argc, char **argv) {
    
    int numbersFlag;
    char *filename;
    char *data[1000];

    //Numbers
    if(argc == 3){
        numbersFlag = 1;
        filename = argv[2];
    } 
    //Strings
    else if(argc == 2){

        //Incomplete parameters
        if(strcmp(argv[1], "-n") == 0){
            printf("Incomplete parameters\n");

            return 0;
        }

        filename = argv[1];

    } else {
        printf("Use correct parameters: \n");
        printf("-n numbers.txt\n");
        printf("strings.txt\n");
        return 0;
    }

    FILE *file;

    if((file = fopen(filename,"r")) == 0){
        printf("Error: can't open %s",filename);
        return 1;
    }

    char line[50];
    int size = 0;

    while(fgets(line,50,file)){
        data[size] = (char*)malloc(strlen(line) + sizeof(char*));
        strcpy(data[size],line);
        size++;
    }

    fclose(file);

    mergeSort((void *)data,0,size-1,numbersFlag? intComparator:strComparator);
    for (int i = 0; i < size; i++){
        printf("%s", data[i]);
    }

}

void mergeSort(char* pos0[],int start, int fin,int (*comparator)(char *, char *)){
    if(fin > start){
        int mid = start + (fin - start) / 2;
        mergeSort(pos0,start,mid,comparator);
        mergeSort(pos0,mid+1,fin,comparator);
        merge(pos0,start,fin,comparator);
    }
}

void merge(char* pos0[],int start, int fin,int (*comparator)(char *, char *)){
    char *tmp[100];
    int midpoint = start + ( fin - start) / 2;
    int a = 0,
        b = 0;
    for(int i = 0; i < fin - start+1; i++){
        if(start+a <= midpoint && midpoint+1+b <= fin){
            if(comparator( pos0[start+a] , pos0[midpoint+1+b] )){
                tmp[i]=malloc(sizeof(char*)*50);
                strcpy(tmp[i],pos0[start+a]);
                a++;
            }
            else{
                tmp[i]=malloc(sizeof(char*)*50);
                strcpy(tmp[i],pos0[midpoint+1+b]);
                b++;
            }
        }
         else if(start+a <= midpoint){
            tmp[i]=malloc(sizeof(char*)*50);
            strcpy(tmp[i],pos0[start+a]);
            a++;
        }
        else{
            tmp[i]=malloc(sizeof(char*)*50);
            strcpy(tmp[i],pos0[midpoint+1+b]);
            b++;
        }
    } 

    for(int i = 0; i < fin-start+1; i++){
        pos0[i+start] = tmp[i];
    }
}

int intComparator(char *a, char *b){
    int s1, s2;
    // String -> Int 
    s1 = atoi(a);
    s2 = atoi(b);

    return s1<s2;
}

int strComparator(char *a, char *b){   
    int c = strcmp(a,b);
    // a > b
    if(c > 0){
        return 0;
    } 
    // b > a or a = b
    else {
        return 1;
    }
}
