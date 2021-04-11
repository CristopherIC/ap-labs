#include <stdio.h>
#include <stdarg.h>
#include <stdlib.h>
#include "logger.h"

//Colors
#define GREEN "\033[92m"
#define YELLOW "\033[93m"
#define RED "\033[91m"
#define PURPLE "\033[95m"
#define RESET "\033[0m"

int infof(const char *format, ...);
int warnf(const char *format, ...);
int errorf(const char *format, ...);
int panicf(const char *format, ...);

int infof(const char *format, ...){
    va_list arg;
    va_start(arg, format);

    printf(GREEN);
    vprintf (format, arg);
    printf(RESET);

    return 0;
}

int warnf(const char *format, ...){
    va_list arg;
    va_start(arg, format);

    printf(YELLOW);
    vprintf (format, arg);
    printf(RESET);

    return 0;
}

int errorf(const char *format, ...){
    va_list arg;
    va_start(arg, format);

    printf(RED);
    vprintf (format, arg);
    printf(RESET);

    return 0;
}

int panicf(const char *format, ...){
    va_list arg;
    va_start(arg, format);

    printf(PURPLE);
    vprintf (format, arg);
    printf(RESET);

    return 0;
}

