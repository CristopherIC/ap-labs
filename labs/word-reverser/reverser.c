#include <stdio.h>
#include <time.h>
#include <stdlib.h>

#define IN   1   /* inside a word */
#define OUT  0   /* outside a word */

/* count lines, words, and characters in input */

void reverse(int lenght, char arr[]) {

    int i;
    char tmp;

    for (i = 0;  i < lenght/2; i++) {
	tmp = arr[i];
	arr[i] = arr[lenght - i - 1];
	arr[lenght - i - 1] = tmp;
    }
    
}


int main()

{
    int i, state;
    char c, word[10];
    state = OUT;

    i = 0;

    while ((c = getchar()) != EOF) {   

	    if (c == ' ' || c == '\n' || c == '\t') {
	        state = OUT;
            printf("%s \n", word);
            reverse(i, word);
            printf("%s \n", word);
            i = 0;    
        }
	    else if (state == OUT) {
	        state = IN;
            word[i] = c;
            i++;
	    } else {
            word[i] = c;
            i++;
        }
    }
    return 0;
}