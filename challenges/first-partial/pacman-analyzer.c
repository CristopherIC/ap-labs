#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>  

// Estructura donde guardaremos la información de los distintos paquetes
struct package{
    char name[100];
    char installDate[100];
    char lastUpdate[100];
    int noUpdates;
    char removalDate[100];
};

//Array donde guardaremos los paquetes
struct package packages[2000];

void analizeLog(char *logFile, char *report);
int findPackageIndex(char *name);

int installed = 0,
    removed = 0,
    upgraded = 0,
    current = 0;

int main(int argc, char **argv) {

    // [0] = pacman-analizer [1] = -input [2] = pacman.txt [3] = -report [4] = packages_report.txt
    
    if (argc != 5) {
	printf("Invalid Input, please use the following format: \n");
    printf("-input log.txt -report reportName.txt \n");
    } else {
        analizeLog(argv[2], argv[4]);
    }
    return 0;
    
}

void analizeLog(char *logFile, char *report) {
    printf("Generating Report from: [%s] log file\n", logFile);

    char *line_buf = NULL;
    size_t line_buf_size = 0;
    int line_count = 0;
    ssize_t line_size;

    FILE *fp = fopen(logFile, "r");

    if (!fp) {
        fprintf(stderr, "Error opening file '%s'\n", logFile);
        return;
    }

    //1era Linea
    line_size = getline(&line_buf, &line_buf_size, fp);
    

    while (line_size >= 0) {
  
        int i = 0;
        char* split = strtok(line_buf, " ");
        char* elements[200];

        while(split != NULL){
            elements[i] = split;
            i++;
            split = strtok(NULL, " ");  //Evita que se quede en bucle
        }

        elements[i] = '\0';
        
        // Para evitar acceder a elementos inexistentes en el array
        if(i >= 4) {
            
            // Para el formato en pacman2.txt

            if(strcmp(elements[2], "installed") == 0){
                strcpy(packages[current].name, elements[3]);
                strcpy(packages[current].installDate, elements[0]);
                strcpy(packages[current].removalDate, "-");
                packages[current].noUpdates = 0;
                current++;
                installed++;
            } else if(strcmp(elements[2], "upgraded") == 0){
                int index = findPackageIndex(elements[3]);
                if(index != -1){
                   strcpy(packages[index].lastUpdate, elements[0]);         
                    if(packages[index].noUpdates == 0){
                        upgraded++;
                    }
                    packages[index].noUpdates++;
                }
            } else if(strcmp(elements[2], "removed") == 0){
                int index = findPackageIndex(elements[3]);
                if(index != -1){
                    strcpy(packages[index].removalDate, elements[0]);
                }
                removed++;
            }

            // Para el formato en pacman.txt

            else if(strcmp(elements[3], "installed") == 0){
                strcpy(packages[current].name, elements[4]);
                strcpy(packages[current].installDate, elements[0]);
                strcat(packages[current].installDate, " ");
                strcat(packages[current].installDate, elements[1]);
                strcpy(packages[current].removalDate, "-");
                packages[current].noUpdates = 0;
                current++;
                installed++;
            } else if(strcmp(elements[3], "upgraded") == 0){
                int index = findPackageIndex(elements[4]);
                if(index != -1){
                    strcpy(packages[index].lastUpdate, elements[0]);
                    strcat(packages[index].lastUpdate, " ");
                    strcat(packages[index].lastUpdate, elements[1]);         
                    if(packages[index].noUpdates == 0){
                        upgraded++;
                    }
                    packages[index].noUpdates++;
                }
            } else if(strcmp(elements[3], "removed") == 0){
                int index = findPackageIndex(elements[4]);
                if(index != -1){
                    strcpy(packages[index].removalDate, elements[0]);
                    strcat(packages[index].removalDate, " ");
                    strcat(packages[index].removalDate, elements[1]);
                }
                removed++;
            }

        }

        //Pasa a la siguiente linea
        line_size = getline(&line_buf, &line_buf_size, fp);
    }

    free(line_buf);
    line_buf = NULL;

    fclose(fp);

    // Crear reporte

    int reportFile = open(report, O_WRONLY|O_CREAT|O_TRUNC, 0644);

    if (reportFile < 0) {
        printf("An error occurred during report generation"); 
    } else {
        char aux[100];

        sprintf(aux,"Pacman Packages Report\n----------------------\n");
        write(reportFile,aux,strlen(aux));
        sprintf(aux,"- Installed packages\t: %d\n",installed);
        write(reportFile,aux,strlen(aux));
        sprintf(aux,"- Removed packages\t: %d\n",removed);
        write(reportFile,aux,strlen(aux));
        sprintf(aux,"- Upgraded packages\t: %d\n",upgraded);
        write(reportFile,aux,strlen(aux));
        sprintf(aux,"- Current installed\t: %d\n",(installed-removed));
        write(reportFile,aux,strlen(aux));

        sprintf(aux,"\nList of packages\n----------------\n");
        write(reportFile,aux,strlen(aux));

        for(int j = 0; j < current; j++){
            sprintf(aux,"- Package Name\t: %s\n",packages[j].name);
            write(reportFile,aux,strlen(aux));
            sprintf(aux,"  - Install date\t: %s\n",packages[j].installDate);
            write(reportFile,aux,strlen(aux));
            sprintf(aux,"  - Last update date\t: %s\n",packages[j].lastUpdate);
            write(reportFile,aux,strlen(aux));
            sprintf(aux,"  - How many updates\t: %d\n",packages[j].noUpdates);
            write(reportFile,aux,strlen(aux));
            sprintf(aux,"  - Removal date\t: %s\n",packages[j].removalDate);
            write(reportFile,aux,strlen(aux));
        }

        close(reportFile);
    }

    printf("Report is generated at: [%s]\n", report);

}

int findPackageIndex(char *name){
    for(int j = 0; j <= current; j++){
        if(strcmp(name,packages[j].name)==0){
            return j;
        }
    }
    return -1;
}
