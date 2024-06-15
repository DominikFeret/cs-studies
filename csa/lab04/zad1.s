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
    ; call f
    add esp, 4
    push eax
    mov eax, OFFSET msg
    push eax
    call printf
    add esp, 8
    mov eax, 0
    ret

f:
    mov eax, [esp + 4]
    cmp eax, 1
    ja skok
    mov eax, 1
    ret

skok:
    dec eax
    push eax
    call f
    add esp, 4
    add eax, [esp + 4]
    ret

.data
msg: .asciz "Wynik = %i\n"