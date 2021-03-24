#include <stdio.h>
#include <string.h>

int mystrlen(char *str);
char *mystradd(char *origin, char *addition);
int mystrfind(char *origin, char *substr);

int main(int argc,char **argv) {
    
    if(argc < 4){
        printf("Not enouht arguments\n");
    } else {
        if(strcmp(argv[1], "-add") == 0) {
            printf("Initial length\t: %d\n",mystrlen(argv[2]));
            char *outcome = mystradd(argv[2], argv[3]);
            printf("New String\t: %s\n" , outcome);
            printf("New length\t: %d\n",mystrlen(outcome));

        } else if(strcmp(argv[1], "-find") == 0){
            int pos = mystrfind(argv[2],argv[3]);
            if(pos < 0){
                printf("['%s'] is not in the first string\n", argv[3]);
            } else {
                printf("['%s'] was found at [%d] position\n", argv[3], pos);
            }
        } else {
            printf("Invalid parameters\n");
        }
        return 0;
    }
    
        
}
