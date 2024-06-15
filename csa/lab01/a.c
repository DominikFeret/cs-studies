#include <stdio.h>

int main() {
    int x = 2024;
    int y = 0;
    
    asm (
        "mov eax, %1;"
        "mov ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "add ebx, eax;"
        "mov %0, ebx;"

       
        :"=r"(y) // out
        :"r"(x) // in
        :"eax", "ebx" // side effects 
    );

    printf("x = %i y = %i\n", x, y);

    return 0;
}