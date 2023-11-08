	.file	"code.c"
	.text
	.def	__main;	.scl	2;	.type	32;	.endef
	.section .rdata,"dr"
	.align 8
.LC0:
	.ascii "C:\\Windows\\System32\\notepad.exe\0"
	.text
	.globl	main
	.def	main;	.scl	2;	.type	32;	.endef
	.seh_proc	main
main:
.LFB4886:
	pushq	%rbp	 #
	.seh_pushreg	%rbp
	movq	%rsp, %rbp	 #,
	.seh_setframe	%rbp, 0
	subq	$224, %rsp	 #,
	.seh_stackalloc	224
	.seh_endprologue
 # code.c:5: {
	call	__main	 #
 # code.c:8:     applicationName = "C:\\Windows\\System32\\notepad.exe";
	leaq	.LC0(%rip), %rax	 #, tmp84
	movq	%rax, -8(%rbp)	 # tmp84, applicationName
 # code.c:13:     ZeroMemory(&startupInfo, sizeof(STARTUPINFO));
	leaq	-112(%rbp), %rax	 #, tmp85
	movl	$104, %r8d	 #,
	movl	$0, %edx	 #,
	movq	%rax, %rcx	 # tmp85,
	call	memset	 #
 # code.c:14:     startupInfo.cb = sizeof(STARTUPINFO);
	movl	$104, -112(%rbp)	 #, startupInfo.cb
 # code.c:17:     CreateProcess(
	movq	-8(%rbp), %rax	 # applicationName, tmp86
	leaq	-144(%rbp), %rdx	 #, tmp87
	movq	%rdx, 72(%rsp)	 # tmp87,
	leaq	-112(%rbp), %rdx	 #, tmp88
	movq	%rdx, 64(%rsp)	 # tmp88,
	movq	$0, 56(%rsp)	 #,
	movq	$0, 48(%rsp)	 #,
	movl	$0, 40(%rsp)	 #,
	movl	$0, 32(%rsp)	 #,
	movl	$0, %r9d	 #,
	movl	$0, %r8d	 #,
	movl	$0, %edx	 #,
	movq	%rax, %rcx	 # tmp86,
	movq	__imp_CreateProcessA(%rip), %rax	 #, tmp89
	call	*%rax	 # tmp89
 # code.c:30:     return 0;
	movl	$0, %eax	 #, _6
 # code.c:31: }
	addq	$224, %rsp	 #,
	popq	%rbp	 #
	ret	
	.seh_endproc
	.ident	"GCC: (tdm64-1) 10.3.0"
	.def	memset;	.scl	2;	.type	32;	.endef
