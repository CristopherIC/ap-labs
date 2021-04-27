#include <stdio.h>
#include <stdarg.h>
#include <stdlib.h>
#include <syslog.h>
#include <string.h>
#include "logger.h"

//Colors
#define GREEN "\033[92m"
#define YELLOW "\033[93m"
#define RED "\033[91m"
#define PURPLE "\033[95m"
#define RESET "\033[0m"


// stdout / "" = 1
// syslog = 0;
int status = 1;

int initLogger(char *logType) {
    
    printf("Initializing Logger on: %s\n", logType);

    if(strcmp(logType,"stdout") == 0 || strcmp(logType,"") == 0){
        status = 1;
        return 0;
    } else if (strcmp(logType,"syslog") == 0) {
        status = 0;
        return 0;
    } else {
        printf("Use a valid argument");
        return 0;
    }
}

int infof(const char *format, ...){
    va_list arg;
    va_start(arg, format);
    printf(GREEN);

    if(status){
        vprintf (format, arg);
        printf("\n");
        printf(RESET);
    }else {
        vsyslog(LOG_INFO,format,arg);
        printf(RESET);
    }
    

    return 0;
}

int warnf(const char *format, ...){
    va_list arg;
    va_start(arg, format);
    printf(YELLOW);
    
    if(status){
        vprintf (format, arg);
        printf("\n");
        printf(RESET);
    }else {
        vsyslog(LOG_INFO,format,arg);
        printf(RESET);
    }

    return 0;
}

int errorf(const char *format, ...){
    va_list arg;
    va_start(arg, format);
    printf(RED);
    
    if(status){
        vprintf (format, arg);
        printf("\n");
        printf(RESET);
    }else {
        vsyslog(LOG_INFO,format,arg);
        printf(RESET);
    }

    return 0;
}

int panicf(const char *format, ...){
    va_list arg;
    va_start(arg, format);
    printf(PURPLE);
    
    if(status){
        vprintf (format, arg);
        printf("\n");
        printf(RESET);
    }else {
        vsyslog(LOG_INFO,format,arg);
        printf(RESET);
    }

    return 0;
}