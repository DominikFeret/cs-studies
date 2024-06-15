.intel_syntax noprefix
.global main
.text
main:
    mov eax, OFFSET m
    push eax
    call printf
    pop eax
    ret
    mov eax, 0

.data
m: .asciz "Hello, world\n"
