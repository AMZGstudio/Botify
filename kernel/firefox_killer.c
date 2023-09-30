#include <linux/module.h>
#include <linux/kprobes.h>
#include <linux/sched.h>

// Define a pre-handler function that will be executed before the original function
static int my_kprobe_pre_handler(struct kprobe *p, struct pt_regs *regs)
{
    // Get the name of the program that is being executed
    struct task_struct *task = current;
    char *filename = task->comm;

    // Check if it matches "firefox"
    if (strcmp(filename, "firefox") == 0)
    {
        // Get the current task struct
        struct task_struct *task = current;

        // Send a SIGKILL signal to terminate it
        send_sig(SIGKILL, task, 0);

        // Return an error code to prevent the execution
        regs->ax = -EPERM;

        // Print a message to the kernel log
        printk(KERN_INFO "Stopped firefox from opening\n");
    }

    // Return 0 to continue the execution
    return 0;
}

// Create a kprobe instance
static struct kprobe my_kprobe = {
    .symbol_name = "__x64_sys_execve",
    .pre_handler = my_kprobe_pre_handler,
};

// A module initialization function that registers the kprobe
static int __init my_init(void)
{
    int ret;

    // Register the kprobe
    ret = register_kprobe(&my_kprobe);
    if (ret < 0)
    {
        printk(KERN_ERR "Failed to register kprobe: %d\n", ret);
        return ret;
    }

    printk(KERN_INFO "Kprobe registered successfully\n");
    return 0;
}

// A module exit function that unregisters the kprobe
static void __exit my_exit(void)
{
    unregister_kprobe(&my_kprobe);
    printk(KERN_INFO "Kprobe unregistered\n");
}

// Define the module entry and exit points
module_init(my_init);
module_exit(my_exit);

// Define some module information
MODULE_LICENSE("GPL");
MODULE_AUTHOR("Your Name");
MODULE_DESCRIPTION("A kernel module that demonstrates kprobe usage");
