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
            "cmp al, 97;"
            "jb skip1;"
            "cmp al, 122;"
            "ja skip1;"
            "sub al, 32;"
        "skip1:"
            "cmp ah, 97;"
            "jb skip2;"
            "cmp ah, 122;"
            "ja skip2;"
            "sub ah, 32;"
        "skip2:"
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