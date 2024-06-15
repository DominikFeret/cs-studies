	.file	"a.c"
	.intel_syntax noprefix
	.text
	.section	.rodata
.LC0:
	.string	"x = %i y = %i\n"
	.text
	.globl	main
	.type	main, @function
main:
.LFB0:
	.cfi_startproc
	endbr64
	push	rbp
	.cfi_def_cfa_offset 16
	.cfi_offset 6, -16
	mov	rbp, rsp
	.cfi_def_cfa_register 6
	push	rbx
	sub	rsp, 24
	.cfi_offset 3, -24
	mov	DWORD PTR -24[rbp], 2024
	mov	DWORD PTR -20[rbp], 0
	mov	edx, DWORD PTR -24[rbp]
#APP
# 7 "a.c" 1
	mov eax, edx;mov ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;add ebx, eax;mov edx, ebx;
# 0 "" 2
#NO_APP
	mov	DWORD PTR -20[rbp], edx
	mov	edx, DWORD PTR -20[rbp]
	mov	eax, DWORD PTR -24[rbp]
	mov	esi, eax
	lea	rax, .LC0[rip]
	mov	rdi, rax
	mov	eax, 0
	call	printf@PLT
	mov	eax, 0
	mov	rbx, QWORD PTR -8[rbp]
	leave
	.cfi_def_cfa 7, 8
	ret
	.cfi_endproc
.LFE0:
	.size	main, .-main
	.ident	"GCC: (Ubuntu 11.4.0-1ubuntu1~22.04) 11.4.0"
	.section	.note.GNU-stack,"",@progbits
	.section	.note.gnu.property,"a"
	.align 8
	.long	1f - 0f
	.long	4f - 1f
	.long	5
0:
	.string	"GNU"
1:
	.align 8
	.long	0xc0000002
	.long	3f - 2f
2:
	.long	0x3
3:
	.align 8
4:
