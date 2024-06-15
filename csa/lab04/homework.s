.intel_syntax noprefix
.global main
.text

main:
    mov eax, [esp+8]
    mov eax, [eax+4]
    push eax
    call atoi
    add esp, 4
    
    push eax
    call fib
    add esp, 4

    push eax
    mov eax, OFFSET msg
    push eax
    call printf
    add esp, 8
    mov eax, 0

    ret

fib:
    mov eax, [esp+4]
    cmp eax, 0
    je zero
    cmp eax, 1
    ja calc_fib
    mov eax, 1
    mov ebx, 0

    ret

calc_fib:
    # get f(n-1) and f(n-2) to eax and ebx registers
    dec eax
    push eax
    dec eax
    push eax
    call fib
    add esp, 4
    call fib
    add esp, 4

    # f(n-1) + f(n-2)
    push eax
    add eax, ebx
    pop ebx
    ret

zero:
    mov eax, 0
    mov ebx, 0
    ret

.data
msg: .asciz "Wynik = %i\n"