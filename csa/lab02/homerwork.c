#include <stdio.h>

int main() {
    int x = 0x1D39;
    int y = 0;
    
    asm (
        "mov eax, %1;"
        "mov ecx, 0;"

        "loop:"
            "cmp ebx, ecx;"
            "jle noswap;"
            "mov ecx, ebx;"
        "noswap:"
            "cmp eax, 0;"
            "je end;"
            "mov ebx, 0;"  
        "one:"
            "shl eax;"
            "jnc loop;"
            "inc ebx;"
            "jmp one;"
        "end:"
            "mov %0, ecx;"

        :"=r"(y)
        :"r"(x) 
        :"eax", "ebx", "ecx"
    );

    printf("Najdłuższy ciąg = %i\n", y);

    return 0;
}