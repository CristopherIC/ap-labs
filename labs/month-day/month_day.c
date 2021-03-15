#include <stdio.h>
#include <stdlib.h>

void month_day(int year, int yearday, int *pmonth, int *pday);
int isLeap(int year);


int isLeap(int year){
    if(year % 400 == 0 || year % 4 == 0){
        return 1;
    } else if (year % 100 == 0){
        return 0;
    } else {
        return 0;
    }
}


void month_day(int year, int yearday, int *pmonth, int *pday){

    int leap = isLeap(year);

    // Entradas invalidas
    if(yearday < 1 || (yearday>365 && !leap) || (yearday>366 && leap) ){
        *pmonth = 0;
        *pday = yearday;
        return;
    }
    
    int days[12] = {31,28,31,30,31,30,31,31,30,31,30,31};

    if(leap) {
        days[1] = 29;
    }


    for(int i = 0; i <= 12; i++){

        if(yearday > days[i]){
            yearday -= days[i];

        } else{
            *pmonth = i + 1;
            *pday = yearday;
            break;
        }
    }
    
}

int main(int argc, char ** argv) {
    int *pmonth, *pday, day, month;
    day = 0;
    month = 0;
    pday = &day;
    pmonth = &month;
    char * months[13] = {"Invalid month","January","February","March","April",
                        "May","June","July","August","September","October",
                        "November","December"};

    int year = atoi(argv[1]);
    int yday = atoi(argv[2]);
    month_day(year, yday, pmonth, pday);
    printf("%s %d, %d\n", months[*pmonth], *pday, year);

}