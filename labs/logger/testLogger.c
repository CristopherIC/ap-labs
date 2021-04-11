#include "logger.h"

int main(){
    infof("Test No: %d, char: %c , string: %s\n",1,'a',"info - green");
    warnf("Test No: %d, char: %c , string: %s\n",2,'b',"warn - yellow");
    errorf("Test No: %d, char: %c , string: %s\n",3,'c',"error - red");
    panicf("Test No: %d, char: %c , string: %s\n",4,'d',"panic - purple");
    return 0;
}
