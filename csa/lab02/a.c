#include <stdio.h>

int main() {
    int x = 0x1D39;
    int y = 0;
    
    asm (
        "mov eax, %1;"
        "xor ebx, ebx;"
        
        "petla:"
            "shl eax;"
            "jnc skok;"
            "inc ebx;"
        
        "skok:"
            "and eax, eax;"
            "jnz petla;"
            
        "mov %0, ebx;"

        :"=r"(y) // out
        :"r"(x) // in
        :"eax", "ebx" // side effects 
    );

    printf("x = %i y = %i\n", x, y);

    return 0;
}