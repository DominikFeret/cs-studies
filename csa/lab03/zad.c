#include<stdio.h>

int main() {
    char s[] = "Abc xyz";

    asm (
        "mov rbx, %0;"
        "mov rcx, rbx;"

        "loop:"
            "mov al, [rcx];"
            "cmp al, 0;"
            "je endfnd;"
            "inc rcx;"
            "jmp loop;"
        "endfnd:"
            "dec rcx;"
        "loop2:"
            "cmp rcx, rbx;"
            "jbe end;"
            "mov al, [rbx];"
            "mov ah, [rcx];"
            "mov [rbx], ah;"
            "mov [rcx], al;"
            "inc rbx;"
            "dec rcx;"
            "jmp loop2;"
        "end:"

        :
        :"r"(s)
        :"rax", "rbx", "rcx"
    );

    printf("%s\n", s);
    return 0;
}